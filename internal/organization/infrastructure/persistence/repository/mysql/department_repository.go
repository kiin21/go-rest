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

type DepartmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) repo.DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) ListWithDetails(
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

	var departmentWithCounts []DeptWithCounts
	var total int64

	baseQuery := r.db.WithContext(ctx).Table("v_departments_with_counts")

	if filter.BusinessUnitID != nil {
		baseQuery = baseQuery.Where("business_unit_id = ?", *filter.BusinessUnitID)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Main query with pagination
	offset := (pg.Page - 1) * pg.Limit
	if err := baseQuery.Order("full_name ASC").Offset(offset).Limit(pg.Limit).Find(&departmentWithCounts).Error; err != nil {
		return nil, 0, err
	}

	leaderIDs := make(map[int64]bool)
	businessUnitIDs := make(map[int64]bool)
	parentDeptIDs := make(map[int64]bool)
	deptIDs := make([]int64, len(departmentWithCounts))

	for i, d := range departmentWithCounts {
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
		leaderIDList := idsToSlice(leaderIDs)

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
		var businessUnitEntities []entity.BusinessUnitEntity
		buIDList := idsToSlice(businessUnitIDs)

		err := r.db.WithContext(ctx).
			Where("id IN ?", buIDList).
			Find(&businessUnitEntities).Error
		if err == nil {
			for i := range businessUnitEntities {
				buMapById[businessUnitEntities[i].ID] = &model.BusinessUnit{
					ID:        businessUnitEntities[i].ID,
					Name:      businessUnitEntities[i].Name,
					Shortname: businessUnitEntities[i].Shortname,
					CompanyID: businessUnitEntities[i].CompanyID,
					LeaderID:  businessUnitEntities[i].LeaderID,
					CreatedAt: businessUnitEntities[i].CreatedAt,
					UpdatedAt: businessUnitEntities[i].UpdatedAt,
				}
			}
		}
	}

	// Batch load parent departments
	parentDeptMap := make(map[int64]*sharedModel.OrgDepartmentNested)
	if len(parentDeptIDs) > 0 {
		var parents []sharedModel.OrgDepartmentNested
		parentIDList := idsToSlice(parentDeptIDs)

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


	subDepartmentMap := make(map[int64][]*sharedModel.OrgDepartmentNested)
	if len(deptIDs) > 0 {
		type SubDepartment struct {
			ID                int64  `gorm:"column:id"`
			GroupDepartmentID int64  `gorm:"column:group_department_id"`
			FullName          string `gorm:"column:full_name"`
			Shortname         string `gorm:"column:shortname"`
		}

		var subDepartments []SubDepartment

		if err := r.db.WithContext(ctx).Table("v_departments_with_counts").
			Select("id, group_department_id, full_name, shortname").
			Where("group_department_id IN ?", deptIDs).
			Find(&subDepartments).Error; err == nil {
			for _, sd := range subDepartments {
				subDepartmentMap[sd.GroupDepartmentID] = append(subDepartmentMap[sd.GroupDepartmentID], &sharedModel.OrgDepartmentNested{
					ID:        sd.ID,
					FullName:  sd.FullName,
					Shortname: sd.Shortname,
				})
			}
		}
	}

	results := make([]*model.DepartmentWithDetails, len(departmentWithCounts))
	for i, d := range departmentWithCounts {
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
			Subdepartments: subDepartmentMap[d.ID],
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

func (r *DepartmentRepository) FindByIDsWithDetails(
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

	// Collect group_department_id, business_unit_id, leader_id
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
		var entities []entity.DepartmentEntity
		gdIDs := idsToSlice(groupDeptIDs)

		if err := r.db.WithContext(ctx).
			Where("id IN ? AND deleted_at IS NULL", gdIDs).
			Find(&entities).Error; err == nil {
			for _, gd := range entities {
				groupDeptMapByDeptId[gd.ID] = r.toModel(&gd)
			}
		}
	}

	buMapByBuId := make(map[int64]*model.BusinessUnit)
	if len(buIDs) > 0 {
		var businessUnitEntities []entity.BusinessUnitEntity
		buIDList := idsToSlice(buIDs)

		if err := r.db.WithContext(ctx).
			Where("id IN ?", buIDList).
			Find(&businessUnitEntities).Error; err == nil {
			for _, bu := range businessUnitEntities {
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
		var entities []entity.StarterEntity
		leaderIDList := idsToSlice(leaderIDs)

		if err := r.db.WithContext(ctx).
			Where("id IN ?", leaderIDList).
			Find(&entities).Error; err == nil {
			for _, l := range entities {
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

	resultMap := make(map[int64]*model.DepartmentWithDetails)
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

		resultMap[vr.ID] = rel
	}

	relations := make([]*model.DepartmentWithDetails, 0, len(ids))
	for _, id := range ids {
		if rel, ok := resultMap[id]; ok {
			relations = append(relations, rel)
		}
	}

	return relations, nil
}

func (r *DepartmentRepository) Create(ctx context.Context, department *model.Department) error {
	_model := &entity.DepartmentEntity{
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

func (r *DepartmentRepository) Update(ctx context.Context, department *model.Department) error {
	deptEntity := &entity.DepartmentEntity{
		ID:                department.ID,
		GroupDepartmentID: department.GroupDepartmentID,
		FullName:          department.FullName,
		Shortname:         department.Shortname,
		BusinessUnitID:    department.BusinessUnitID,
		LeaderID:          department.LeaderID,
	}

	if err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", department.ID).
		Updates(deptEntity).Error; err != nil {
		return err
	}

	return nil
}

func (r *DepartmentRepository) toModel(dm *entity.DepartmentEntity) *model.Department {
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

// Helper function
func idsToSlice(idMap map[int64]bool) []int64 {
	ids := make([]int64, len(idMap))
	i := 0
	for id := range idMap {
		ids[i] = id
		i++
	}
	return ids
}