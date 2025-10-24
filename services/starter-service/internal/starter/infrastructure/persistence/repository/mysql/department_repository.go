package mysql

import (
	"context"

	"github.com/kiin21/go-rest/pkg/httputil"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
	repo "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/infrastructure/persistence/entity"
	"gorm.io/gorm"
)

type DepartmentRepository struct {
	db *gorm.DB
}

func (r *DepartmentRepository) SearchByKeyword(ctx context.Context, keyword string) ([]*model.Department, int64, error) {
	var entities []entity.DepartmentEntity
	var total int64

	baseQuery := r.db.WithContext(ctx).
		Model(&entity.DepartmentEntity{}).
		Where("deleted_at IS NULL")

	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		baseQuery = baseQuery.Where(
			"full_name LIKE ? OR shortname LIKE ?",
			searchPattern,
			searchPattern,
		)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := baseQuery.
		Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return r.entitiesToModels(entities), total, nil
}

func NewDepartmentRepository(db *gorm.DB) repo.DepartmentRepository {
	return &DepartmentRepository{db: db}
}

// ============================================================================
// ============================================================================

func (r *DepartmentRepository) FindByIDs(ctx context.Context, ids []int64) ([]*model.Department, error) {
	if len(ids) == 0 {
		return []*model.Department{}, nil
	}

	var entities []entity.DepartmentEntity
	if err := r.db.WithContext(ctx).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Find(&entities).Error; err != nil {
		return nil, err
	}

	return r.entitiesToModels(entities), nil
}

func (r *DepartmentRepository) ListWithDetails(
	ctx context.Context,
	filter *model.DepartmentListFilter,
	pg *httputil.ReqPagination,
) ([]*model.DepartmentWithDetails, int64, error) {
	// Fetch main data with count
	viewResults, total, err := r.fetchDepartmentsWithCounts(ctx, filter, pg)
	if err != nil {
		return nil, 0, err
	}
	// Extract related IDs
	relatedIDs := r.extractRelatedIDsFromCounts(viewResults)
	// Fetch related data
	relatedData := r.fetchRelatedData(ctx, relatedIDs)
	results := r.buildDepartmentDetailsFromCounts(viewResults, relatedData)

	return results, total, nil
}

func (r *DepartmentRepository) FindByIDsWithDetails(
	ctx context.Context,
	ids []int64,
) ([]*model.DepartmentWithDetails, error) {
	if len(ids) == 0 {
		return []*model.DepartmentWithDetails{}, nil
	}

	viewResults, err := r.fetchDepartmentViewResults(ctx, ids)
	if err != nil {
		return nil, err
	}
	
	relatedIDs := r.extractRelatedIDs(viewResults)
	relatedData := r.fetchRelatedData(ctx, relatedIDs)
	
	return r.buildDepartmentDetailsPreserveOrder(ids, viewResults, relatedData), nil
}

func (r *DepartmentRepository) Create(ctx context.Context, department *model.Department) error {
	newEntity := &entity.DepartmentEntity{
		GroupDepartmentID: department.GroupDepartmentID,
		FullName:          department.FullName,
		Shortname:         department.Shortname,
		BusinessUnitID:    department.BusinessUnitID,
		LeaderID:          department.LeaderID,
	}

	if err := r.db.WithContext(ctx).Create(newEntity).Error; err != nil {
		return err
	}

	department.ID = newEntity.ID
	department.CreatedAt = newEntity.CreatedAt
	department.UpdatedAt = newEntity.UpdatedAt

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

	return r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", department.ID).
		Updates(deptEntity).
		Error
}

func (r *DepartmentRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Exec("CALL sp_delete_department(?)", id).Error
}

// ============================================================================
// ============================================================================

type deptWithCounts struct {
	ID                int64  `gorm:"column:id"`
	FullName          string `gorm:"column:full_name"`
	Shortname         string `gorm:"column:shortname"`
	LeaderID          *int64 `gorm:"column:leader_id"`
	GroupDepartmentID *int64 `gorm:"column:group_department_id"`
	BusinessUnitID    *int64 `gorm:"column:business_unit_id"`
	CreatedAt         string `gorm:"column:created_at"`
	UpdatedAt         string `gorm:"column:updated_at"`
}

func (r *DepartmentRepository) fetchDepartmentsWithCounts(
	ctx context.Context,
	filter *model.DepartmentListFilter,
	pg *httputil.ReqPagination,
) ([]deptWithCounts, int64, error) {
	var results []deptWithCounts
	var total int64

	baseQuery := r.db.WithContext(ctx).
		Table("v_departments_with_counts").
		Where("deleted_at IS NULL")

	if filter.BusinessUnitID != nil {
		baseQuery = baseQuery.Where("business_unit_id = ?", *filter.BusinessUnitID)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := baseQuery.
		Offset(pg.GetOffset()).
		Limit(pg.GetLimit()).
		Find(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

type departmentViewResult struct {
	ID                int64  `gorm:"column:id"`
	FullName          string `gorm:"column:full_name"`
	Shortname         string `gorm:"column:shortname"`
	GroupDepartmentID *int64 `gorm:"column:group_department_id"`
	BusinessUnitID    *int64 `gorm:"column:business_unit_id"`
	LeaderID          *int64 `gorm:"column:leader_id"`
	CreatedAt         string `gorm:"column:created_at"`
	UpdatedAt         string `gorm:"column:updated_at"`
}

func (r *DepartmentRepository) fetchDepartmentViewResults(ctx context.Context, ids []int64) ([]departmentViewResult, error) {
	var results []departmentViewResult

	if err := r.db.WithContext(ctx).
		Table("v_departments_with_bu").
		Where("id IN ? AND deleted_at IS NULL", ids).
		Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

// ============================================================================
// ============================================================================

type relatedIDs struct {
	deptIDs      []int64
	groupDeptIDs []int64
	buIDs        []int64
	leaderIDs    []int64
}

type idMapper struct {
	deptIDs      []int64
	groupDeptIDs map[int64]bool
	buIDs        map[int64]bool
	leaderIDs    map[int64]bool
}

func newIDMapper(capacity int) *idMapper {
	return &idMapper{
		deptIDs:      make([]int64, 0, capacity),
		groupDeptIDs: make(map[int64]bool),
		buIDs:        make(map[int64]bool),
		leaderIDs:    make(map[int64]bool),
	}
}

func (e *idMapper) add(deptID int64, groupDeptID, buID, leaderID *int64) {
	e.deptIDs = append(e.deptIDs, deptID)

	if groupDeptID != nil {
		e.groupDeptIDs[*groupDeptID] = true
	}

	if buID != nil {
		e.buIDs[*buID] = true
	}

	if leaderID != nil {
		e.leaderIDs[*leaderID] = true
	}
}

func (e *idMapper) build() relatedIDs {
	return relatedIDs{
		deptIDs:      e.deptIDs,
		groupDeptIDs: mapKeysToSlice(e.groupDeptIDs),
		buIDs:        mapKeysToSlice(e.buIDs),
		leaderIDs:    mapKeysToSlice(e.leaderIDs),
	}
}

func (r *DepartmentRepository) extractRelatedIDs(
	viewResults []departmentViewResult,
) relatedIDs {
	mapper := newIDMapper(len(viewResults))

	for _, vr := range viewResults {
		mapper.add(vr.ID, vr.GroupDepartmentID, vr.BusinessUnitID, vr.LeaderID)
	}

	return mapper.build()
}

func (r *DepartmentRepository) extractRelatedIDsFromCounts(
	counts []deptWithCounts,
) relatedIDs {
	extractor := newIDMapper(len(counts))

	for _, d := range counts {
		extractor.add(d.ID, d.GroupDepartmentID, d.BusinessUnitID, d.LeaderID)
	}

	return extractor.build()
}

// ============================================================================
// ============================================================================

type relatedData struct {
	groupDepts    map[int64]*model.Department
	businessUnits map[int64]*model.BusinessUnit
	leaders       map[int64]*model.LineManagerNested
	subdepts      map[int64][]*model.OrgDepartmentNested
}

func (r *DepartmentRepository) fetchRelatedData(
	ctx context.Context,
	ids relatedIDs,
) relatedData {
	return relatedData{
		groupDepts:    r.fetchGroupDepartments(ctx, ids.groupDeptIDs),
		businessUnits: r.fetchBusinessUnits(ctx, ids.buIDs),
		leaders:       r.fetchLeaders(ctx, ids.leaderIDs),
		subdepts:      r.fetchSubdepartments(ctx, ids.deptIDs),
	}
}

func (r *DepartmentRepository) fetchGroupDepartments(
	ctx context.Context,
	ids []int64,
) map[int64]*model.Department {
	result := make(map[int64]*model.Department)
	if len(ids) == 0 {
		return result
	}

	var entities []entity.DepartmentEntity
	if err := r.db.WithContext(ctx).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Find(&entities).Error; err != nil {
		return result
	}

	for i := range entities {
		result[entities[i].ID] = r.toModel(&entities[i])
	}

	return result
}

func (r *DepartmentRepository) fetchBusinessUnits(
	ctx context.Context,
	ids []int64,
) map[int64]*model.BusinessUnit {
	result := make(map[int64]*model.BusinessUnit)
	if len(ids) == 0 {
		return result
	}

	var entities []entity.BusinessUnitEntity
	if err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&entities).Error; err != nil {
		return result
	}

	for i := range entities {
		e := &entities[i]
		result[e.ID] = &model.BusinessUnit{
			ID:        e.ID,
			Name:      e.Name,
			Shortname: e.Shortname,
			CompanyID: e.CompanyID,
			LeaderID:  e.LeaderID,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		}
	}

	return result
}

func (r *DepartmentRepository) fetchLeaders(
	ctx context.Context,
	ids []int64,
) map[int64]*model.LineManagerNested {
	result := make(map[int64]*model.LineManagerNested)
	if len(ids) == 0 {
		return result
	}

	var leaders []model.LineManagerNested
	if err := r.db.WithContext(ctx).
		Table("starters").
		Select("id, domain, name, email, job_title").
		Where("id IN ? AND deleted_at IS NULL", ids).
		Find(&leaders).Error; err != nil {
		return result
	}

	for i := range leaders {
		result[leaders[i].ID] = &leaders[i]
	}

	return result
}

func (r *DepartmentRepository) fetchSubdepartments(
	ctx context.Context,
	parentIDs []int64,
) map[int64][]*model.OrgDepartmentNested {
	result := make(map[int64][]*model.OrgDepartmentNested)
	if len(parentIDs) == 0 {
		return result
	}

	type subDeptRow struct {
		ID                int64  `gorm:"column:id"`
		GroupDepartmentID *int64 `gorm:"column:group_department_id"`
		FullName          string `gorm:"column:full_name"`
		Shortname         string `gorm:"column:shortname"`
	}

	var rows []subDeptRow
	if err := r.db.WithContext(ctx).
		Table("departments").
		Select("id, group_department_id, full_name, shortname").
		Where("group_department_id IN ? AND deleted_at IS NULL", parentIDs).
		Find(&rows).Error; err != nil {
		return result
	}

	for _, row := range rows {
		if row.GroupDepartmentID == nil {
			continue
		}

		result[*row.GroupDepartmentID] = append(
			result[*row.GroupDepartmentID],
			&model.OrgDepartmentNested{
				ID:        row.ID,
				FullName:  row.FullName,
				Shortname: row.Shortname,
			},
		)
	}

	return result
}

// ============================================================================
// ============================================================================

type departmentDetailsBuilder struct {
	related relatedData
}

func (b *departmentDetailsBuilder) attachRelations(
	dept *model.DepartmentWithDetails,
	groupDeptID, buID, leaderID *int64,
	deptID int64,
) {
	if groupDeptID != nil {
		if gd, ok := b.related.groupDepts[*groupDeptID]; ok {
			dept.ParentDepartment = &model.OrgDepartmentNested{
				ID:        gd.ID,
				FullName:  gd.FullName,
				Shortname: gd.Shortname,
			}
		}
	}

	if buID != nil {
		if bu, ok := b.related.businessUnits[*buID]; ok {
			dept.BusinessUnit = bu
		}
	}

	if leaderID != nil {
		if leader, ok := b.related.leaders[*leaderID]; ok {
			dept.Leader = leader
		}
	}

	if subs, ok := b.related.subdepts[deptID]; ok {
		dept.Subdepartments = subs
	}
}

func (r *DepartmentRepository) buildDepartmentDetailsFromCounts(
	counts []deptWithCounts,
	related relatedData,
) []*model.DepartmentWithDetails {
	results := make([]*model.DepartmentWithDetails, len(counts))
	builder := &departmentDetailsBuilder{related: related}

	for i, d := range counts {
		dept := &model.DepartmentWithDetails{
			Department: &model.Department{
				ID:                d.ID,
				GroupDepartmentID: d.GroupDepartmentID,
				FullName:          d.FullName,
				Shortname:         d.Shortname,
				BusinessUnitID:    d.BusinessUnitID,
				LeaderID:          d.LeaderID,
			},
		}

		builder.attachRelations(
			dept,
			d.GroupDepartmentID,
			d.BusinessUnitID,
			d.LeaderID,
			d.ID,
		)

		results[i] = dept
	}

	return results
}

func (r *DepartmentRepository) buildDepartmentDetailsPreserveOrder(
	originalIDs []int64,
	viewResults []departmentViewResult,
	related relatedData,
) []*model.DepartmentWithDetails {
	resultMap := make(map[int64]*model.DepartmentWithDetails, len(viewResults))
	builder := &departmentDetailsBuilder{related: related}

	for _, vr := range viewResults {
		dept := &model.DepartmentWithDetails{
			Department: &model.Department{
				ID:                vr.ID,
				FullName:          vr.FullName,
				Shortname:         vr.Shortname,
				GroupDepartmentID: vr.GroupDepartmentID,
				BusinessUnitID:    vr.BusinessUnitID,
				LeaderID:          vr.LeaderID,
			},
		}

		builder.attachRelations(
			dept,
			vr.GroupDepartmentID,
			vr.BusinessUnitID,
			vr.LeaderID,
			vr.ID,
		)

		resultMap[vr.ID] = dept
	}

	// Preserve original order
	result := make([]*model.DepartmentWithDetails, 0, len(originalIDs))
	for _, id := range originalIDs {
		if dept, ok := resultMap[id]; ok {
			result = append(result, dept)
		}
	}

	return result
}

// ============================================================================
// ============================================================================

func (r *DepartmentRepository) toModel(
	dm *entity.DepartmentEntity,
) *model.Department {
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

func (r *DepartmentRepository) entitiesToModels(
	entities []entity.DepartmentEntity,
) []*model.Department {
	result := make([]*model.Department, len(entities))
	for i := range entities {
		result[i] = r.toModel(&entities[i])
	}
	return result
}

func mapKeysToSlice(m map[int64]bool) []int64 {
	if len(m) == 0 {
		return nil
	}

	result := make([]int64, 0, len(m))
	for id := range m {
		result = append(result, id)
	}
	return result
}
