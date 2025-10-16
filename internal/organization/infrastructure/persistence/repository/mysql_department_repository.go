package repository

import (
	"context"

	"github.com/kiin21/go-rest/internal/organization/domain"
	"github.com/kiin21/go-rest/internal/organization/infrastructure/persistence/model"
	"github.com/kiin21/go-rest/pkg/response"
	"gorm.io/gorm"
)

// MySQLDepartmentRepository implements DepartmentRepository using MySQL/GORM
type MySQLDepartmentRepository struct {
	db *gorm.DB
}

// NewMySQLDepartmentRepository creates a new MySQL-based department repository
func NewMySQLDepartmentRepository(db *gorm.DB) domain.DepartmentRepository {
	return &MySQLDepartmentRepository{db: db}
}

// List retrieves departments with optional filtering and pagination
func (r *MySQLDepartmentRepository) List(ctx context.Context, filter domain.DepartmentListFilter, pg response.ReqPagination) ([]*domain.Department, int64, error) {
	var models []model.DepartmentModel
	var total int64

	// Build base query - exclude soft deleted
	query := r.db.WithContext(ctx).
		Model(&model.DepartmentModel{}).
		Where("deleted_at IS NULL")

	// Apply filters
	if filter.BusinessUnitID != nil {
		// Only direct departments
		query = query.Where("business_unit_id = ?", *filter.BusinessUnitID)

	}

	// Count total matching records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	offset := (pg.Page - 1) * pg.Limit
	query = query.
		Order("full_name ASC").
		Offset(offset).
		Limit(pg.Limit)

	// Execute query
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	// Convert to domain entities
	departments := make([]*domain.Department, len(models))
	for i, _model := range models {
		departments[i] = r.toDomain(&_model)
	}

	return departments, total, nil
}

// ListWithDetails retrieves departments with all related data using optimized queries
func (r *MySQLDepartmentRepository) ListWithDetails(ctx context.Context, filter domain.DepartmentListFilter, pg response.ReqPagination) ([]*domain.DepartmentWithDetails, int64, error) {
	// Step 1: Get departments with counts from view
	type DeptWithCounts struct {
		ID                  int64  `gorm:"column:id"`
		GroupDepartmentID   *int64 `gorm:"column:group_department_id"`
		FullName            string `gorm:"column:full_name"`
		Shortname           string `gorm:"column:shortname"`
		LeaderID            *int64 `gorm:"column:leader_id"`
		BusinessUnitID      *int64 `gorm:"column:business_unit_id"`
		CreatedAt           string `gorm:"column:created_at"`
		UpdatedAt           string `gorm:"column:updated_at"`
		MembersCount        int    `gorm:"column:members_count"`
		SubdepartmentsCount int    `gorm:"column:subdepartments_count"`
	}

	var deptsWithCounts []DeptWithCounts
	var total int64

	// Build query from view
	query := r.db.WithContext(ctx).
		Table("v_departments_with_counts").
		Where("deleted_at IS NULL")

	// Apply filters
	if filter.BusinessUnitID != nil {
		query = query.Where("business_unit_id = ?", *filter.BusinessUnitID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (pg.Page - 1) * pg.Limit
	query = query.Order("full_name ASC").Offset(offset).Limit(pg.Limit)

	if err := query.Find(&deptsWithCounts).Error; err != nil {
		return nil, 0, err
	}

	// Step 2: Collect IDs for batch loading
	leaderIDs := make(map[int64]bool)
	businessUnitIDs := make(map[int64]bool)
	parentDeptIDs := make(map[int64]bool)
	deptIDs := make([]int64, len(deptsWithCounts))

	for i, d := range deptsWithCounts {
		deptIDs[i] = d.ID
		if d.LeaderID != nil {
			leaderIDs[*d.LeaderID] = true
		}
		if d.BusinessUnitID != nil {
			businessUnitIDs[*d.BusinessUnitID] = true
		}
		if d.GroupDepartmentID != nil {
			parentDeptIDs[*d.GroupDepartmentID] = true
		}
	}

	// Step 3: Batch load leaders
	leaderMap := make(map[int64]*domain.Leader)
	if len(leaderIDs) > 0 {
		type LeaderResult struct {
			ID       int64  `gorm:"column:id"`
			Domain   string `gorm:"column:domain"`
			Name     string `gorm:"column:name"`
			Email    string `gorm:"column:email"`
			JobTitle string `gorm:"column:job_title"`
		}
		var leaders []LeaderResult
		leaderIDList := make([]int64, 0, len(leaderIDs))
		for id := range leaderIDs {
			leaderIDList = append(leaderIDList, id)
		}
		if err := r.db.WithContext(ctx).Table("starters").
			Select("id, domain, name, email, job_title").
			Where("id IN ? AND deleted_at IS NULL", leaderIDList).
			Find(&leaders).Error; err == nil {
			for _, l := range leaders {
				leaderMap[l.ID] = &domain.Leader{
					ID:       l.ID,
					Domain:   l.Domain,
					Name:     l.Name,
					Email:    l.Email,
					JobTitle: l.JobTitle,
				}
			}
		}
	}

	// Step 4: Batch load business units
	buMap := make(map[int64]*domain.BusinessUnit)
	if len(businessUnitIDs) > 0 {
		var buModels []model.BusinessUnitModel
		buIDList := make([]int64, 0, len(businessUnitIDs))
		for id := range businessUnitIDs {
			buIDList = append(buIDList, id)
		}
		if err := r.db.WithContext(ctx).
			Where("id IN ?", buIDList).
			Find(&buModels).Error; err == nil {
			for _, bu := range buModels {
				buMap[bu.ID] = &domain.BusinessUnit{
					ID:        bu.ID,
					Name:      bu.Name,
					Shortname: bu.Shortname,
				}
			}
		}
	}

	// Step 5: Batch load parent departments
	parentDeptMap := make(map[int64]*domain.DepartmentNested)
	if len(parentDeptIDs) > 0 {
		type ParentDept struct {
			ID           int64  `gorm:"column:id"`
			FullName     string `gorm:"column:full_name"`
			Shortname    string `gorm:"column:shortname"`
			MembersCount int    `gorm:"column:members_count"`
		}
		var parents []ParentDept
		parentIDList := make([]int64, 0, len(parentDeptIDs))
		for id := range parentDeptIDs {
			parentIDList = append(parentIDList, id)
		}
		if err := r.db.WithContext(ctx).Table("v_departments_with_counts").
			Select("id, full_name, shortname, members_count").
			Where("id IN ?", parentIDList).
			Find(&parents).Error; err == nil {
			for _, p := range parents {
				parentDeptMap[p.ID] = &domain.DepartmentNested{
					ID:           p.ID,
					FullName:     p.FullName,
					Shortname:    p.Shortname,
					MembersCount: p.MembersCount,
				}
			}
		}
	}

	// Step 6: Batch load subdepartments
	subdeptsMap := make(map[int64][]*domain.DepartmentNested)
	if len(deptIDs) > 0 {
		type Subdept struct {
			ID                int64  `gorm:"column:id"`
			GroupDepartmentID int64  `gorm:"column:group_department_id"`
			FullName          string `gorm:"column:full_name"`
			Shortname         string `gorm:"column:shortname"`
			MembersCount      int    `gorm:"column:members_count"`
		}
		var subdepts []Subdept
		if err := r.db.WithContext(ctx).Table("v_departments_with_counts").
			Select("id, group_department_id, full_name, shortname, members_count").
			Where("group_department_id IN ? AND deleted_at IS NULL", deptIDs).
			Find(&subdepts).Error; err == nil {
			for _, sd := range subdepts {
				subdeptsMap[sd.GroupDepartmentID] = append(subdeptsMap[sd.GroupDepartmentID], &domain.DepartmentNested{
					ID:           sd.ID,
					FullName:     sd.FullName,
					Shortname:    sd.Shortname,
					MembersCount: sd.MembersCount,
				})
			}
		}
	}

	// Step 7: Assemble final result
	results := make([]*domain.DepartmentWithDetails, len(deptsWithCounts))
	for i, d := range deptsWithCounts {
		dept := &domain.Department{
			ID:                d.ID,
			GroupDepartmentID: d.GroupDepartmentID,
			FullName:          d.FullName,
			Shortname:         d.Shortname,
			BusinessUnitID:    d.BusinessUnitID,
			LeaderID:          d.LeaderID,
		}

		// Build department with details, handling nil pointers safely
		details := &domain.DepartmentWithDetails{
			Department:          dept,
			MembersCount:        d.MembersCount,
			SubdepartmentsCount: d.SubdepartmentsCount,
			Subdepartments:      subdeptsMap[d.ID],
		}

		// Safely assign business unit if exists
		if d.BusinessUnitID != nil {
			details.BusinessUnit = buMap[*d.BusinessUnitID]
		}

		// Safely assign leader if exists
		if d.LeaderID != nil {
			details.Leader = leaderMap[*d.LeaderID]
		}

		// Safely assign parent department if exists
		if d.GroupDepartmentID != nil {
			details.ParentDepartment = parentDeptMap[*d.GroupDepartmentID]
		}

		results[i] = details
	}

	return results, total, nil
}

// FindByIDsWithRelations retrieves departments with group_department and business_unit using view
//
//goland:noinspection DuplicatedCode
func (r *MySQLDepartmentRepository) FindByIDsWithRelations(ctx context.Context, ids []int64) ([]*domain.DepartmentWithDetails, error) {
	if len(ids) == 0 {
		return []*domain.DepartmentWithDetails{}, nil
	}

	// Query from view with resolved business_unit_id
	type ViewResult struct {
		ID                int64  `gorm:"column:id"`
		GroupDepartmentID *int64 `gorm:"column:group_department_id"`
		FullName          string `gorm:"column:full_name"`
		Shortname         string `gorm:"column:shortname"`
		BusinessUnitID    *int64 `gorm:"column:business_unit_id"` // Auto-resolved from view
		LeaderID          *int64 `gorm:"column:leader_id"`
		CreatedAt         string `gorm:"column:created_at"`
		UpdatedAt         string `gorm:"column:updated_at"`
	}

	var viewResults []ViewResult
	if err := r.db.WithContext(ctx).
		Table("v_departments_with_bu").
		Where("id IN ?", ids).
		Find(&viewResults).Error; err != nil {
		return nil, err
	}

	// Map to store departments for quick lookup
	deptMap := make(map[int64]*ViewResult)
	for i := range viewResults {
		deptMap[viewResults[i].ID] = &viewResults[i]
	}

	// Collect group_department_id và business_unit_id cần load
	groupDeptIDs := make(map[int64]bool)
	buIDs := make(map[int64]bool)

	for _, vr := range viewResults {
		if vr.GroupDepartmentID != nil {
			groupDeptIDs[*vr.GroupDepartmentID] = true
		}
		if vr.BusinessUnitID != nil {
			buIDs[*vr.BusinessUnitID] = true
		}
	}

	// Load group departments nếu có
	groupDeptMap := make(map[int64]*domain.Department)
	if len(groupDeptIDs) > 0 {
		var groupDeptModels []model.DepartmentModel
		gdIDs := make([]int64, 0, len(groupDeptIDs))
		for id := range groupDeptIDs {
			gdIDs = append(gdIDs, id)
		}

		if err := r.db.WithContext(ctx).
			Where("id IN ? AND deleted_at IS NULL", gdIDs).
			Find(&groupDeptModels).Error; err == nil {
			for _, gd := range groupDeptModels {
				groupDeptMap[gd.ID] = r.toDomain(&gd)
			}
		}
	}

	// Load business units nếu có
	buMap := make(map[int64]*domain.BusinessUnit)
	if len(buIDs) > 0 {
		var buModels []model.BusinessUnitModel
		buIDList := make([]int64, 0, len(buIDs))
		for id := range buIDs {
			buIDList = append(buIDList, id)
		}

		if err := r.db.WithContext(ctx).
			Where("id IN ?", buIDList).
			Find(&buModels).Error; err == nil {
			for _, bu := range buModels {
				buMap[bu.ID] = &domain.BusinessUnit{
					ID:        bu.ID,
					Name:      bu.Name,
					Shortname: bu.Shortname,
				}
			}
		}
	}

	// Build final result
	relations := make([]*domain.DepartmentWithDetails, 0, len(viewResults))
	for _, vr := range viewResults {
		rel := &domain.DepartmentWithDetails{
			Department: &domain.Department{
				ID:                vr.ID,
				FullName:          vr.FullName,
				Shortname:         vr.Shortname,
				GroupDepartmentID: vr.GroupDepartmentID,
				BusinessUnitID:    vr.BusinessUnitID,
			},
		}

		// Add group department if exists
		if vr.GroupDepartmentID != nil {
			if gd, ok := groupDeptMap[*vr.GroupDepartmentID]; ok {
				rel.ParentDepartment = &domain.DepartmentNested{
					ID:           gd.ID,
					FullName:     gd.FullName,
					Shortname:    gd.Shortname,
					MembersCount: 2107,
				}
			}
		}

		// Add business unit if exists (từ view đã resolve)
		if vr.BusinessUnitID != nil {
			if bu, ok := buMap[*vr.BusinessUnitID]; ok {
				rel.BusinessUnit = bu
			}
		}

		relations = append(relations, rel)
	}

	return relations, nil
}

// Create inserts a new department
func (r *MySQLDepartmentRepository) Create(ctx context.Context, department *domain.Department) error {
	_model := &model.DepartmentModel{
		GroupDepartmentID: department.GroupDepartmentID,
		FullName:          department.FullName,
		Shortname:         department.Shortname,
		BusinessUnitID:    department.BusinessUnitID,
		LeaderID:          department.LeaderID,
	}

	if err := r.db.WithContext(ctx).Create(_model).Error; err != nil {
		return err
	}

	// Update the domain entity with the generated ID and timestamps
	department.ID = _model.ID
	department.CreatedAt = _model.CreatedAt
	department.UpdatedAt = _model.UpdatedAt

	return nil
}

// Update updates an existing department
func (r *MySQLDepartmentRepository) Update(ctx context.Context, department *domain.Department) error {
	model := &model.DepartmentModel{
		ID:                department.ID,
		GroupDepartmentID: department.GroupDepartmentID,
		FullName:          department.FullName,
		Shortname:         department.Shortname,
		BusinessUnitID:    department.BusinessUnitID,
		LeaderID:          department.LeaderID,
	}

	if err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", department.ID).
		Updates(model).Error; err != nil {
		return err
	}

	return nil
}

// toDomain converts model.DepartmentModel to domain.Department
func (r *MySQLDepartmentRepository) toDomain(model *model.DepartmentModel) *domain.Department {
	return &domain.Department{
		ID:                model.ID,
		GroupDepartmentID: model.GroupDepartmentID,
		FullName:          model.FullName,
		Shortname:         model.Shortname,
		BusinessUnitID:    model.BusinessUnitID,
		LeaderID:          model.LeaderID,
		CreatedAt:         model.CreatedAt,
		UpdatedAt:         model.UpdatedAt,
		DeletedAt:         model.DeletedAt,
	}
}
