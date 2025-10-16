package application

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	appDto "github.com/kiin21/go-rest/internal/organization/application/dto"
	"github.com/kiin21/go-rest/internal/organization/domain"
	"github.com/kiin21/go-rest/pkg/response"
)

// DepartmentApplicationService handles application-level business logic for departments.
type DepartmentApplicationService struct {
	departmentRepo domain.DepartmentRepository
}

// NewDepartmentApplicationService creates a new application service.
func NewDepartmentApplicationService(departmentRepo domain.DepartmentRepository) *DepartmentApplicationService {
	return &DepartmentApplicationService{
		departmentRepo: departmentRepo,
	}
}

// GetAllDepartments returns a paginated list of departments with full details.
func (s *DepartmentApplicationService) GetAllDepartments(ctx context.Context, query appDto.ListDepartmentsQuery) (*response.PaginatedResult[*domain.DepartmentWithDetails], error) {
	filter := domain.DepartmentListFilter{
		BusinessUnitID: query.BusinessUnitID,
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

// GetOneDepartment retrieves departments (by IDs) with their relations loaded.
// It accepts a query that contains one or more department IDs and returns
// a slice of DepartmentWithDetails from the repository.
func (s *DepartmentApplicationService) GetOneDepartment(ctx context.Context, query appDto.GetDepartmentQuery) ([]*domain.DepartmentWithDetails, error) {
	// Convert []int to []int64 for repository
	ids := make([]int64, 0, 1)
	ids = append(ids, *query.ID)

	departments, err := s.departmentRepo.FindByIDsWithRelations(ctx, ids)
	if err != nil {
		return nil, err
	}

	return departments, nil
}

// CreateDepartment creates a new department
func (s *DepartmentApplicationService) CreateDepartment(ctx context.Context, cmd appDto.CreateDepartmentCommand) (*domain.DepartmentWithDetails, error) {
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
func (s *DepartmentApplicationService) UpdateDepartment(ctx context.Context, id int64, cmd appDto.UpdateDepartmentCommand) (*domain.DepartmentWithDetails, error) {
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
func (s *DepartmentApplicationService) AssignLeader(ctx context.Context, cmd appDto.AssignLeaderCommand) (*domain.DepartmentWithDetails, error) {
	// Find existing department (with details)
	ids := make([]int64, 0, 1)
	ids = append(ids, cmd.DepartmentID)
	departments, err := s.departmentRepo.FindByIDsWithRelations(ctx, ids)
	if err != nil {
		return nil, err
	}

	if len(departments) == 0 {
		return nil, response.NewAPIError(404, "Department not found", nil)
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
		// For domain-based assignment, you might need to:
		// 1. Look up user by domain to get their ID
		// 2. Then assign the ID to the department
		// For now, this is a placeholder - you'll need to implement user lookup logic
		if cmd.LeaderDomain != nil {
			// TODO: Implement user lookup by domain
			// userID, err := s.userService.FindUserIDByDomain(ctx, *cmd.LeaderDomain)
			// if err != nil {
			//     return nil, err
			// }
			// department.LeaderID = &userID
			return nil, response.NewAPIError(400, "Leader assignment by domain not yet implemented", nil)
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
