package service

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	businessunitquery "github.com/kiin21/go-rest/internal/organization/application/dto/businessunit/query"
	departmentcommand "github.com/kiin21/go-rest/internal/organization/application/dto/department/command"
	departmentquery "github.com/kiin21/go-rest/internal/organization/application/dto/department/query"
	model "github.com/kiin21/go-rest/internal/organization/domain/model"
	"github.com/kiin21/go-rest/internal/organization/domain/repository"
	sharedDomain "github.com/kiin21/go-rest/internal/shared/domain"
	"github.com/kiin21/go-rest/pkg/response"
)

type OrganizationApplicationService struct {
	departmentRepo   repository.DepartmentRepository
	businessUnitRepo repository.BusinessUnitRepository
	starterRepo      repository.StarterRepository
}

func NewOrganizationApplicationService(
	departmentRepo repository.DepartmentRepository,
	businessUnitRepo repository.BusinessUnitRepository,
	starterRepo repository.StarterRepository,
) *OrganizationApplicationService {
	return &OrganizationApplicationService{
		departmentRepo:   departmentRepo,
		businessUnitRepo: businessUnitRepo,
		starterRepo:      starterRepo,
	}
}

func (s *OrganizationApplicationService) GetAllDepartments(
	ctx context.Context,
	query departmentquery.ListDepartmentsQuery,
) (*response.PaginatedResult[*model.DepartmentWithDetails], error) {
	filter := model.DepartmentListFilter{BusinessUnitID: query.BusinessUnitID}

	departments, total, err := s.departmentRepo.ListWithDetails(ctx, filter, query.Pagination)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / query.Pagination.Limit
	if int(total)%query.Pagination.Limit > 0 {
		totalPages++
	}

	var prev *string = nil
	var next *string = nil

	if query.Pagination.Page > 1 {
		value := strconv.Itoa(query.Pagination.Page - 1)
		prev = &value
	}
	if query.Pagination.Page < totalPages {
		value := strconv.Itoa(query.Pagination.Page + 1)
		next = &value
	}

	return &response.PaginatedResult[*model.DepartmentWithDetails]{
		Data: departments,
		Pagination: response.RespPagination{
			Limit:      query.Pagination.Limit,
			TotalItems: total,
			Prev:       prev,
			Next:       next,
		},
	}, nil
}

func (s *OrganizationApplicationService) GetOneDepartment(
	ctx context.Context,
	query departmentquery.GetDepartmentQuery,
) (*model.DepartmentWithDetails, error) {
	ids := make([]int64, 0, 1)
	ids = append(ids, query.ID)

	departments, err := s.departmentRepo.FindByIDsWithDetails(ctx, ids)

	if err != nil {
		return nil, err
	}
	if len(departments) == 0 || len(departments) > 1 {
		return nil, response.NewAPIError(400, "Bad request", sharedDomain.ErrInvalidInput)
	}

	return departments[0], nil
}

func (s *OrganizationApplicationService) CreateDepartment(
	ctx context.Context,
	cmd departmentcommand.CreateDepartmentCommand,
) (*model.DepartmentWithDetails, error) {
	department := &model.Department{
		FullName:          cmd.FullName,
		Shortname:         cmd.Shortname,
		BusinessUnitID:    cmd.BusinessUnitID,
		GroupDepartmentID: cmd.GroupDepartmentID,
		LeaderID:          cmd.LeaderID,
	}

	if err := s.departmentRepo.Create(ctx, department); err != nil {
		return nil, err
	}

	detailedDepartments, err := s.departmentRepo.FindByIDsWithDetails(ctx, []int64{department.ID})
	if err != nil {
		return nil, err
	}
	if len(detailedDepartments) == 0 {
		return nil, errors.New("failed to fetch department details after creation")
	}

	return detailedDepartments[0], nil
}

func (s *OrganizationApplicationService) UpdateDepartment(ctx context.Context, id int64, cmd departmentcommand.UpdateDepartmentCommand) (*model.DepartmentWithDetails, error) {

	ids := make([]int64, 0, 1)
	ids = append(ids, id)
	departments, err := s.departmentRepo.FindByIDsWithDetails(ctx, ids)
	if err != nil {
		return nil, err
	}

	department := departments[0]

	department.FullName = *cmd.FullName
	department.Shortname = *cmd.Shortname
	department.GroupDepartmentID = cmd.GroupDepartmentID
	department.LeaderID = cmd.LeaderID
	department.BusinessUnitID = cmd.BusinessUnitID

	if err := s.departmentRepo.Update(ctx, &model.Department{
		ID:                department.ID,
		GroupDepartmentID: department.GroupDepartmentID,
		FullName:          department.FullName,
		Shortname:         department.Shortname,
		BusinessUnitID:    department.BusinessUnitID,
		LeaderID:          department.LeaderID,
		CreatedAt:         department.CreatedAt,
		UpdatedAt:         department.UpdatedAt,
		DeletedAt:         department.DeletedAt,
	}); err != nil {
		return nil, err
	}

	return department, nil
}

func (s *OrganizationApplicationService) AssignLeader(ctx context.Context, cmd departmentcommand.AssignLeaderCommand) (*model.DepartmentWithDetails, error) {
	deptId := make([]int64, 0, 1)
	deptId = append(deptId, cmd.DepartmentID)
	departments, err := s.departmentRepo.FindByIDsWithDetails(ctx, deptId)
	if err != nil {
		return nil, err
	}

	if len(departments) == 0 {
		return nil, sharedDomain.ErrNotFound
	}

	department := departments[0]

	switch cmd.LeaderIdentifierType {
	case "id":
		if cmd.LeaderID != nil {
			department.LeaderID = cmd.LeaderID
		}
	case "domain":
		if cmd.LeaderDomain != nil {
			if s.starterRepo == nil {
				return nil, response.NewAPIError(http.StatusNotImplemented, "starter repository not configured", nil)
			}

			starter, err := s.starterRepo.FindByDomain(ctx, *cmd.LeaderDomain)
			if err != nil {
				if errors.Is(err, sharedDomain.ErrNotFound) {
					return nil, sharedDomain.ErrNotFound
				}
				return nil, err
			}

			leaderID := starter.ID()
			department.LeaderID = &leaderID
		}
	default:
		return nil, response.NewAPIError(400, "Invalid leader identifier type", nil)
	}

	if err := s.departmentRepo.Update(ctx, &model.Department{
		ID:                department.ID,
		GroupDepartmentID: department.GroupDepartmentID,
		FullName:          department.FullName,
		Shortname:         department.Shortname,
		BusinessUnitID:    department.BusinessUnitID,
		LeaderID:          department.LeaderID,
		CreatedAt:         department.CreatedAt,
		UpdatedAt:         department.UpdatedAt,
		DeletedAt:         department.DeletedAt,
	}); err != nil {
		return nil, err
	}

	return department, nil
}

func (s *OrganizationApplicationService) ListBusinessUnits(ctx context.Context, query businessunitquery.ListBusinessUnitsQuery) (*response.PaginatedResult[*model.BusinessUnit], error) {
	units, total, err := s.businessUnitRepo.List(ctx, query.Pagination)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / query.Pagination.Limit
	if int(total)%query.Pagination.Limit > 0 {
		totalPages++
	}

	var prev, next *string
	if query.Pagination.Page > 1 {
		value := strconv.Itoa(query.Pagination.Page - 1)
		prev = &value
	}
	if query.Pagination.Page < totalPages {
		value := strconv.Itoa(query.Pagination.Page + 1)
		next = &value
	}

	return &response.PaginatedResult[*model.BusinessUnit]{
		Data: units,
		Pagination: response.RespPagination{
			Limit:      query.Pagination.Limit,
			TotalItems: total,
			Prev:       prev,
			Next:       next,
		},
	}, nil
}

func (s *OrganizationApplicationService) GetBusinessUnit(ctx context.Context, id int64) (*model.BusinessUnit, error) {
	ids := make([]int64, 0, 1)
	ids = append(ids, id)
	bus, err := s.businessUnitRepo.FindByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	if len(bus) == 0 || len(bus) > 1 {
		return nil, sharedDomain.ErrNotFound
	}
	return bus[0], nil
}

func (s *OrganizationApplicationService) ListBusinessUnitsWithDetails(ctx context.Context, query businessunitquery.ListBusinessUnitsQuery) (*response.PaginatedResult[*model.BusinessUnitWithDetails], error) {
	units, total, err := s.businessUnitRepo.ListWithDetails(ctx, query.Pagination)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / query.Pagination.Limit
	if int(total)%query.Pagination.Limit > 0 {
		totalPages++
	}

	var prev, next *string
	if query.Pagination.Page > 1 {
		value := strconv.Itoa(query.Pagination.Page - 1)
		prev = &value
	}
	if query.Pagination.Page < totalPages {
		value := strconv.Itoa(query.Pagination.Page + 1)
		next = &value
	}

	return &response.PaginatedResult[*model.BusinessUnitWithDetails]{
		Data: units,
		Pagination: response.RespPagination{
			Limit:      query.Pagination.Limit,
			TotalItems: total,
			Prev:       prev,
			Next:       next,
		},
	}, nil
}

func (s *OrganizationApplicationService) GetBusinessUnitWithDetails(ctx context.Context, id int64) (*model.BusinessUnitWithDetails, error) {
	return s.businessUnitRepo.FindByIDWithDetails(ctx, id)
}
