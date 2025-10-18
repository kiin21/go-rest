package mysql

import (
	"context"

	"github.com/kiin21/go-rest/internal/organization/domain/model"
	repo "github.com/kiin21/go-rest/internal/organization/domain/repository"
	"github.com/kiin21/go-rest/internal/organization/infrastructure/persistence/entity"
	sharedModel "github.com/kiin21/go-rest/internal/shared/domain/model"
	"github.com/kiin21/go-rest/pkg/response"
	"gorm.io/gorm"
)

type MySQLDepartmentRepository struct {
	db *gorm.DB
}

func NewMySQLDepartmentRepository(db *gorm.DB) repo.DepartmentRepository {
	return &MySQLDepartmentRepository{db: db}
}

func (r *MySQLDepartmentRepository) ListWithDetails(
	ctx context.Context,
	filter model.DepartmentListFilter,
	pg response.ReqPagination,
) ([]*model.DepartmentWithDetails, int64, error) {
	// View struct
	type DeptWithCounts struct {
		ID                int64  `gorm:"column:id"`
		FullName          string `gorm:"column:full_name"`
		Shortname         string `gorm:"column:shortname"`
		LeaderID          *int64 `gorm:"column:leader_id"`
		GroupDepartmentID *int64 `gorm:"column:group_department_id"`
		BusinessUnitID    *int64 `gorm:"column:business_unit_id"`
		CreatedAt         string `gorm:"column:created_at"`
		UpdatedAt         string `gorm:"column:updated_at"`
	}

	var deptsWithCounts []DeptWithCounts
	var total int64

	// Build query from view
	query := r.db.WithContext(ctx).Table("v_departments_with_counts")

	if filter.BusinessUnitID != nil {
		query = query.Where("business_unit_id = ?", *filter.BusinessUnitID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (pg.Page - 1) * pg.Limit
	query = query.Order("full_name ASC").Offset(offset).Limit(pg.Limit)
	if err := query.Find(&deptsWithCounts).Error; err != nil {
		return nil, 0, err
	}

	// Collect IDs for batch loading, use map to avoid duplicated IDs
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

	// Batch load leaders
	leaderMapById := make(map[int64]*sharedModel.LineManagerNested)
	if len(leaderIDs) > 0 {
		var leaders []sharedModel.LineManagerNested
		leaderIDList := make([]int64, 0, len(leaderIDs))
		for id := range leaderIDs {
			leaderIDList = append(leaderIDList, id)
		}

		err := r.db.WithContext(ctx).Table("starters").
			Select("id, domain, name, email, job_title").
			Where("id IN ? AND deleted_at IS NULL", leaderIDList).
			Find(&leaders).Error

		if err == nil {
			for i := range leaders {
				leaderMapById[leaders[i].ID] = &leaders[i]
			}
		}
	}

	// Batch load business units
	buMapById := make(map[int64]*model.BusinessUnit)
	if len(businessUnitIDs) > 0 {
		var buModels []entity.BusinessUnitModel
		buIDList := make([]int64, 0, len(businessUnitIDs))
		for id := range businessUnitIDs {
			buIDList = append(buIDList, id)
		}
		err := r.db.WithContext(ctx).
			Where("id IN ?", buIDList).
			Find(&buModels).Error
		if err == nil {
			for i := range buModels {
				buMapById[buModels[i].ID] = &model.BusinessUnit{
					ID:        buModels[i].ID,
					Name:      buModels[i].Name,
					Shortname: buModels[i].Shortname,
					CompanyID: buModels[i].CompanyID,
					LeaderID:  buModels[i].LeaderID,
					CreatedAt: buModels[i].CreatedAt,
					UpdatedAt: buModels[i].UpdatedAt,
				}
			}
		}
	}

	// Batch load parent departments
	parentDeptMap := make(map[int64]*sharedModel.OrgDepartmentNested)
	if len(parentDeptIDs) > 0 {
		var parents []sharedModel.OrgDepartmentNested
		parentIDList := make([]int64, 0, len(parentDeptIDs))
		for id := range parentDeptIDs {
			parentIDList = append(parentIDList, id)
		}

		err := r.db.WithContext(ctx).Table("v_departments_with_counts").
			Select("id, full_name, shortname").
			Where("id IN ?", parentIDList).
			Find(&parents).Error

		if err == nil {
			for i := range parents {
				parentDeptMap[parents[i].ID] = &parents[i]
			}
		}
	}

	// Batch load subdepartments
	subdeptsMap := make(map[int64][]*sharedModel.OrgDepartmentNested)
	if len(deptIDs) > 0 {

		type Subdept struct {
			ID                int64  `gorm:"column:id"`
			GroupDepartmentID int64  `gorm:"column:group_department_id"`
			FullName          string `gorm:"column:full_name"`
			Shortname         string `gorm:"column:shortname"`
		}

		var subdepts []Subdept

		if err := r.db.WithContext(ctx).Table("v_departments_with_counts").
			Select("id, group_department_id, full_name, shortname").
			Where("group_department_id IN ?", deptIDs).
			Find(&subdepts).Error; err == nil {
			for _, sd := range subdepts {
				subdeptsMap[sd.GroupDepartmentID] = append(subdeptsMap[sd.GroupDepartmentID], &sharedModel.OrgDepartmentNested{
					ID:        sd.ID,
					FullName:  sd.FullName,
					Shortname: sd.Shortname,
				})
			}
		}
	}

	// Assemble result
	results := make([]*model.DepartmentWithDetails, len(deptsWithCounts))
	for i, d := range deptsWithCounts {
		dept := &model.Department{
			ID:                d.ID,
			GroupDepartmentID: d.GroupDepartmentID,
			FullName:          d.FullName,
			Shortname:         d.Shortname,
			BusinessUnitID:    d.BusinessUnitID,
			LeaderID:          d.LeaderID,
		}

		details := &model.DepartmentWithDetails{
			Department:     dept,
			Subdepartments: subdeptsMap[d.ID],
		}

		if d.BusinessUnitID != nil {
			details.BusinessUnit = buMapById[*d.BusinessUnitID]
		}

		if d.LeaderID != nil {
			details.Leader = leaderMapById[*d.LeaderID]
		}

		if d.GroupDepartmentID != nil {
			details.ParentDepartment = parentDeptMap[*d.GroupDepartmentID]
		}

		results[i] = details
	}

	return results, total, nil
}

func (r *MySQLDepartmentRepository) FindByIDsWithDetails(
	ctx context.Context,
	ids []int64,
) ([]*model.DepartmentWithDetails, error) {
	if len(ids) == 0 {
		return []*model.DepartmentWithDetails{}, nil
	}

	// View struct
	type ViewResult struct {
		ID                int64  `gorm:"column:id"`
		FullName          string `gorm:"column:full_name"`
		Shortname         string `gorm:"column:shortname"`
		GroupDepartmentID *int64 `gorm:"column:group_department_id"`
		BusinessUnitID    *int64 `gorm:"column:business_unit_id"`
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

	deptMapByDeptId := make(map[int64]*ViewResult)
	for i := range viewResults {
		deptMapByDeptId[viewResults[i].ID] = &viewResults[i]
	}

	// Collect group_department_id và business_unit_id
	groupDeptIDs := make(map[int64]bool)
	buIDs := make(map[int64]bool)
	leaderIDs := make(map[int64]bool)

	for _, vr := range viewResults {
		if vr.GroupDepartmentID != nil {
			groupDeptIDs[*vr.GroupDepartmentID] = true
		}
		if vr.BusinessUnitID != nil {
			buIDs[*vr.BusinessUnitID] = true
		}
		if vr.LeaderID != nil {
			leaderIDs[*vr.LeaderID] = true
		}
	}

	groupDeptMapByDeptId := make(map[int64]*model.Department)
	if len(groupDeptIDs) > 0 {
		var groupDeptModels []entity.DepartmentModel
		gdIDs := make([]int64, 0, len(groupDeptIDs))
		for id := range groupDeptIDs {
			gdIDs = append(gdIDs, id)
		}

		if err := r.db.WithContext(ctx).
			Where("id IN ? AND deleted_at IS NULL", gdIDs).
			Find(&groupDeptModels).Error; err == nil {
			for _, gd := range groupDeptModels {
				groupDeptMapByDeptId[gd.ID] = r.toDomain(&gd)
			}
		}
	}

	buMapByBuId := make(map[int64]*model.BusinessUnit)
	if len(buIDs) > 0 {
		var buModels []entity.BusinessUnitModel
		buIDList := make([]int64, 0, len(buIDs))
		for id := range buIDs {
			buIDList = append(buIDList, id)
		}

		if err := r.db.WithContext(ctx).
			Where("id IN ?", buIDList).
			Find(&buModels).Error; err == nil {
			for _, bu := range buModels {
				buMapByBuId[bu.ID] = &model.BusinessUnit{
					ID:        bu.ID,
					Name:      bu.Name,
					Shortname: bu.Shortname,
				}
			}
		}
	}

	leaderMapByLeaderId := make(map[int64]*sharedModel.LineManagerNested)
	if len(leaderIDs) > 0 {
		var leaderModels []entity.StarterModel
		leaderIDList := make([]int64, 0, len(leaderIDs))
		for id := range leaderIDs {
			leaderIDList = append(leaderIDList, id)
		}

		if err := r.db.WithContext(ctx).
			Where("id IN ?", leaderIDList).
			Find(&leaderModels).Error; err == nil {
			for _, l := range leaderModels {
				leaderMapByLeaderId[l.ID] = &sharedModel.LineManagerNested{
					ID:       l.ID,
					Domain:   l.Domain,
					Name:     l.Name,
					Email:    l.Email,
					JobTitle: l.JobTitle,
				}
			}
		}
	}

	// Assemble final result
	relations := make([]*model.DepartmentWithDetails, 0, len(viewResults))
	for _, vr := range viewResults {
		rel := &model.DepartmentWithDetails{
			Department: &model.Department{
				ID:                vr.ID,
				FullName:          vr.FullName,
				Shortname:         vr.Shortname,
				GroupDepartmentID: vr.GroupDepartmentID,
				BusinessUnitID:    vr.BusinessUnitID,
				LeaderID:          vr.LeaderID,
			},
		}

		if vr.GroupDepartmentID != nil {
			if gd, ok := groupDeptMapByDeptId[*vr.GroupDepartmentID]; ok {
				rel.ParentDepartment = &sharedModel.OrgDepartmentNested{
					ID:        gd.ID,
					FullName:  gd.FullName,
					Shortname: gd.Shortname,
				}
			}
		}

		if vr.BusinessUnitID != nil {
			if bu, ok := buMapByBuId[*vr.BusinessUnitID]; ok {
				rel.BusinessUnit = bu
			}
		}

		if vr.LeaderID != nil {
			if lm, ok := leaderMapByLeaderId[*vr.LeaderID]; ok {
				rel.Leader = lm
			}
		}

		relations = append(relations, rel)
	}

	return relations, nil
}

func (r *MySQLDepartmentRepository) Create(ctx context.Context, department *model.Department) error {
	_model := &entity.DepartmentModel{
		GroupDepartmentID: department.GroupDepartmentID,
		FullName:          department.FullName,
		Shortname:         department.Shortname,
		BusinessUnitID:    department.BusinessUnitID,
		LeaderID:          department.LeaderID,
	}

	if err := r.db.WithContext(ctx).Create(_model).Error; err != nil {
		return err
	}

	department.ID = _model.ID
	department.CreatedAt = _model.CreatedAt
	department.UpdatedAt = _model.UpdatedAt

	return nil
}

func (r *MySQLDepartmentRepository) Update(ctx context.Context, department *model.Department) error {
	model := &entity.DepartmentModel{
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

func (r *MySQLDepartmentRepository) toDomain(dm *entity.DepartmentModel) *model.Department {
	return &model.Department{
		ID:                dm.ID,
		GroupDepartmentID: dm.GroupDepartmentID,
		FullName:          dm.FullName,
		Shortname:         dm.Shortname,
		BusinessUnitID:    dm.BusinessUnitID,
		LeaderID:          dm.LeaderID,
		CreatedAt:         dm.CreatedAt,
		UpdatedAt:         dm.UpdatedAt,
		DeletedAt:         dm.DeletedAt,
	}
}
