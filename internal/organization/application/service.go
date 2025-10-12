package application

import (
	"context"
	"strconv"

	appDto "github.com/kiin21/go-rest/internal/organization/application/dto"
	"github.com/kiin21/go-rest/internal/organization/domain"
	"github.com/kiin21/go-rest/pkg/response"
)

// OrganizationApplicationService handles application-level business logic for organization entities.
type OrganizationApplicationService struct {
	departmentRepo domain.DepartmentRepository
}

// NewOrganizationApplicationService creates a new application service.
func NewOrganizationApplicationService(departmentRepo domain.DepartmentRepository) *OrganizationApplicationService {
	return &OrganizationApplicationService{
		departmentRepo: departmentRepo,
	}
}

// GetAllDepartments returns a paginated list of departments with full details.
func (s *OrganizationApplicationService) GetAllDepartments(ctx context.Context, query appDto.ListDepartmentsQuery) (*response.PaginatedResult[*domain.DepartmentWithDetails], error) {
	filter := domain.DepartmentListFilter{
		BusinessUnitID:        query.BusinessUnitID,
		IncludeSubdepartments: query.IncludeSubdepartments,
	}

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

// CreateDepartment creates a new department
func (s *OrganizationApplicationService) CreateDepartment(ctx context.Context, cmd appDto.CreateDepartmentCommand) (*domain.Department, error) {
	// Create domain entity
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

	// Convert to DTO and return
	return department, nil
}

// UpdateDepartment updates an existing department with partial update support
func (s *OrganizationApplicationService) UpdateDepartment(ctx context.Context, id int64, cmd appDto.UpdateDepartmentCommand) (*domain.Department, error) {
	// Find existing department
	department, err := s.departmentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates using domain methods (only update fields that are provided)
	department.UpdateInfo(cmd.FullName, cmd.Shortname)
	department.AssignToBusinessUnit(cmd.BusinessUnitID)
	department.AssignToGroupDepartment(cmd.GroupDepartmentID)
	department.AssignLeader(cmd.LeaderID)

	// Persist changes
	if err := s.departmentRepo.Update(ctx, department); err != nil {
		return nil, err
	}

	// Convert to DTO and return
	return department, nil
}
