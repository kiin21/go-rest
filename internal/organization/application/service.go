package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	appDto "github.com/kiin21/go-rest/internal/organization/application/dto"
	"github.com/kiin21/go-rest/internal/organization/domain"
	sharedDomain "github.com/kiin21/go-rest/internal/shared/domain"
	"github.com/kiin21/go-rest/pkg/response"
)

type OrganizationApplicationService struct {
	departmentRepo   domain.DepartmentRepository
	businessUnitRepo domain.BusinessUnitRepository
	leaderLookup     LeaderLookup
}

func NewOrganizationApplicationService(
	departmentRepo domain.DepartmentRepository,
	businessUnitRepo domain.BusinessUnitRepository,
	leaderLookup LeaderLookup,
) *OrganizationApplicationService {
	return &OrganizationApplicationService{
		departmentRepo:   departmentRepo,
		businessUnitRepo: businessUnitRepo,
		leaderLookup:     leaderLookup,
	}
}

func (s *OrganizationApplicationService) GetAllDepartments(
	ctx context.Context,
	query appDto.ListDepartmentsQuery,
) (*response.PaginatedResult[*domain.DepartmentWithDetails], error) {
	filter := domain.DepartmentListFilter{BusinessUnitID: query.BusinessUnitID}

	departments, total, err := s.departmentRepo.ListWithDetails(ctx, filter, query.Pagination)
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
	} else {
		prev = nil
	}
	if query.Pagination.Page < totalPages {
		value := strconv.Itoa(query.Pagination.Page + 1)
		next = &value
	} else {
		next = nil
	}

	return &response.PaginatedResult[*domain.DepartmentWithDetails]{
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
	query appDto.GetDepartmentQuery,
) (*domain.DepartmentWithDetails, error) {
	ids := make([]int64, 0, 1)
	ids = append(ids, query.ID)

	departments, err := s.departmentRepo.FindByIDsWithRelations(ctx, ids)
	if err != nil {
		return nil, err
	}

	if len(departments) > 1 {
		return nil, response.NewAPIError(400, "Bad request", sharedDomain.ErrInvalidInput)
	}

	leader, err := s.leaderLookup.FindStarterById(ctx, *departments[0].LeaderID)
	if err != nil {
		return nil, response.NewAPIError(400, "Bad request", sharedDomain.ErrInvalidInput)
	}
	departments[0].Leader = &domain.LineManager{
		ID:       leader.ID(),
		Domain:   leader.Domain(),
		Email:    leader.Email(),
		JobTitle: leader.JobTitle(),
		Name:     leader.Name(),
	}

	return departments[0], nil
}

// CreateDepartment creates a new department
func (s *OrganizationApplicationService) CreateDepartment(ctx context.Context, cmd appDto.CreateDepartmentCommand) (*domain.DepartmentWithDetails, error) {
	department := &domain.Department{
		FullName:          cmd.FullName,
		Shortname:         cmd.Shortname,
		BusinessUnitID:    cmd.BusinessUnitID,
		GroupDepartmentID: cmd.GroupDepartmentID,
		LeaderID:          cmd.LeaderID,
	}

	// Persist to repository
	if err := s.departmentRepo.Create(ctx, department); err != nil {
		return nil, err
	}

	// Fetch the newly created department with all its details
	detailedDepartments, err := s.departmentRepo.FindByIDsWithRelations(ctx, []int64{department.ID})
	if err != nil {
		// Log the error, but proceed with the basic info if details fail
		return nil, err
	}
	if len(detailedDepartments) == 0 {
		return nil, errors.New("failed to fetch department details after creation")
	}

	return detailedDepartments[0], nil
}

// UpdateDepartment updates an existing department with partial update support
func (s *OrganizationApplicationService) UpdateDepartment(ctx context.Context, id int64, cmd appDto.UpdateDepartmentCommand) (*domain.DepartmentWithDetails, error) {
	// Find existing department (with details)
	ids := make([]int64, 0, 1)
	ids = append(ids, id)
	departments, err := s.departmentRepo.FindByIDsWithRelations(ctx, ids)
	if err != nil {
		return nil, err
	}

	// Extract the underlying domain.Department pointer
	department := departments[0]

	// Apply updates using domain methods (only update fields that are provided)
	department.FullName = *cmd.FullName
	department.Shortname = *cmd.Shortname
	department.GroupDepartmentID = cmd.GroupDepartmentID
	department.LeaderID = cmd.LeaderID
	department.BusinessUnitID = cmd.BusinessUnitID

	fmt.Println("BYHBUYJBKJH: ", department.BusinessUnitID)
	// Persist changes
	if err := s.departmentRepo.Update(ctx, &domain.Department{
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

	// Return the updated domain.Department
	return department, nil
}

// AssignLeader assigns a leader to a department using either ID or domain
func (s *OrganizationApplicationService) AssignLeader(ctx context.Context, cmd appDto.AssignLeaderCommand) (*domain.DepartmentWithDetails, error) {
	// Find existing department (with details)
	deptId := make([]int64, 0, 1)
	deptId = append(deptId, cmd.DepartmentID)
	departments, err := s.departmentRepo.FindByIDsWithRelations(ctx, deptId)
	if err != nil {
		return nil, err
	}

	if len(departments) == 0 {
		return nil, sharedDomain.ErrNotFound
	}

	// Extract the underlying domain.Department pointer
	department := departments[0]

	// Handle leader assignment based on identifier type
	switch cmd.LeaderIdentifierType {
	case "id":
		if cmd.LeaderID != nil {
			department.LeaderID = cmd.LeaderID
		}
	case "domain":
		if cmd.LeaderDomain != nil {
			if s.leaderLookup == nil {
				return nil, response.NewAPIError(http.StatusNotImplemented, "starter lookup not configured", nil)
			}

			leaderID, err := s.leaderLookup.FindStarterIDByDomain(ctx, *cmd.LeaderDomain)
			if err != nil {
				if errors.Is(err, sharedDomain.ErrNotFound) {
					return nil, sharedDomain.ErrNotFound
				}
				return nil, err
			}

			department.LeaderID = &leaderID
		}
	default:
		return nil, response.NewAPIError(400, "Invalid leader identifier type", nil)
	}

	// Persist changes
	if err := s.departmentRepo.Update(ctx, &domain.Department{
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

	// Return the updated department
	return department, nil
}

// ListBusinessUnits returns a paginated list of business units.
func (s *OrganizationApplicationService) ListBusinessUnits(ctx context.Context, query appDto.ListBusinessUnitsQuery) (*response.PaginatedResult[*domain.BusinessUnit], error) {
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

	return &response.PaginatedResult[*domain.BusinessUnit]{
		Data: units,
		Pagination: response.RespPagination{
			Limit:      query.Pagination.Limit,
			TotalItems: total,
			Prev:       prev,
			Next:       next,
		},
	}, nil
}

// GetBusinessUnit retrieves a single business unit by ID.
func (s *OrganizationApplicationService) GetBusinessUnit(ctx context.Context, id int64) (*domain.BusinessUnit, error) {
	return s.businessUnitRepo.FindByID(ctx, id)
}

// ListBusinessUnitsWithDetails returns a paginated list of business units with company and leader details.
func (s *OrganizationApplicationService) ListBusinessUnitsWithDetails(ctx context.Context, query appDto.ListBusinessUnitsQuery) (*response.PaginatedResult[*domain.BusinessUnitWithDetails], error) {
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

	return &response.PaginatedResult[*domain.BusinessUnitWithDetails]{
		Data: units,
		Pagination: response.RespPagination{
			Limit:      query.Pagination.Limit,
			TotalItems: total,
			Prev:       prev,
			Next:       next,
		},
	}, nil
}

// GetBusinessUnitWithDetails retrieves a single business unit with company and leader details by ID.
func (s *OrganizationApplicationService) GetBusinessUnitWithDetails(ctx context.Context, id int64) (*domain.BusinessUnitWithDetails, error) {
	return s.businessUnitRepo.FindByIDWithDetails(ctx, id)
}
