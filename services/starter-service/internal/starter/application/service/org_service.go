package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/kiin21/go-rest/pkg/events"
	"github.com/kiin21/go-rest/pkg/httputil"
	businessunitquery "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/businessunit/query"
	departmentcommand "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/department/command"
	departmentquery "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/department/query"
	sharedDomain "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/error"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
	messagebroker "github.com/kiin21/go-rest/services/starter-service/internal/starter/infrastructure/messagebroker"
)

type OrganizationApplicationService struct {
	departmentRepo   repository.DepartmentRepository
	businessUnitRepo repository.BusinessUnitRepository
	starterRepo      repository.StarterRepository
	notificationPub  messagebroker.NotificationPublisher
}

func NewOrganizationApplicationService(
	departmentRepo repository.DepartmentRepository,
	businessUnitRepo repository.BusinessUnitRepository,
	starterRepo repository.StarterRepository,
	notificationPublisher messagebroker.NotificationPublisher,
) *OrganizationApplicationService {
	return &OrganizationApplicationService{
		departmentRepo:   departmentRepo,
		businessUnitRepo: businessUnitRepo,
		starterRepo:      starterRepo,
		notificationPub:  notificationPublisher,
	}
}

func (s *OrganizationApplicationService) GetAllDepartments(ctx context.Context, query departmentquery.ListDepartmentsQuery) (*httputil.PaginatedResult[*model.DepartmentWithDetails], error) {
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

	return &httputil.PaginatedResult[*model.DepartmentWithDetails]{
		Data: departments,
		Pagination: httputil.RespPagination{
			Limit:      query.Pagination.Limit,
			TotalItems: total,
			Prev:       prev,
			Next:       next,
		},
	}, nil
}

func (s *OrganizationApplicationService) GetOneDepartment(ctx context.Context, query departmentquery.GetDepartmentQuery) (*model.DepartmentWithDetails, error) {
	ids := make([]int64, 0, 1)
	ids = append(ids, query.ID)

	departments, err := s.departmentRepo.FindByIDsWithDetails(ctx, ids)

	if err != nil {
		return nil, err
	}
	if len(departments) == 0 {
		return nil, sharedDomain.ErrNotFound
	}

	return departments[0], nil
}

func (s *OrganizationApplicationService) CreateDepartment(ctx context.Context, cmd departmentcommand.CreateDepartmentCommand) (*model.DepartmentWithDetails, error) {
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

func (s *OrganizationApplicationService) UpdateDepartment(ctx context.Context, cmd departmentcommand.UpdateDepartmentCommand) (*model.DepartmentWithDetails, error) {

	ids := make([]int64, 0, 1)
	ids = append(ids, cmd.ID)
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

func (s *OrganizationApplicationService) DeleteDepartment(ctx context.Context, cmd departmentcommand.DeleteDepartmentCommand) error {
	return s.departmentRepo.Delete(ctx, cmd.ID)
}

func (s *OrganizationApplicationService) AssignLeader(ctx context.Context, cmd departmentcommand.AssignLeaderCommand) (*model.DepartmentWithDetails, error) {
	deptIDList := []int64{cmd.DepartmentID}
	departments, err := s.departmentRepo.FindByIDsWithDetails(ctx, deptIDList)
	if err != nil {
		return nil, err
	}

	if len(departments) == 0 {
		return nil, sharedDomain.ErrNotFound
	}

	department := departments[0]
	previousLeaderDomain := ""
	if department.Leader != nil {
		previousLeaderDomain = department.Leader.Domain
	}

	switch cmd.LeaderIdentifierType {
	case "id":
		if cmd.LeaderID != nil {
			department.LeaderID = cmd.LeaderID
		}
	case "domain":
		if cmd.LeaderDomain != nil {
			if s.starterRepo == nil {
				return nil, httputil.NewAPIError(http.StatusNotImplemented, "starter repository not configured", nil)
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
		return nil, sharedDomain.ErrInvalidInput
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

	updated, err := s.departmentRepo.FindByIDsWithDetails(ctx, []int64{department.ID})
	if err != nil {
		return nil, err
	}
	if len(updated) == 0 {
		return nil, sharedDomain.ErrNotFound
	}
	department = updated[0]

	s.publishLeaderAssignmentNotification(ctx, department, previousLeaderDomain)

	return department, nil
}

func (s *OrganizationApplicationService) publishLeaderAssignmentNotification(ctx context.Context, department *model.DepartmentWithDetails, previousLeaderDomain string) {
	if s.notificationPub == nil || department == nil || department.Leader == nil {
		return
	}

	toDomain := department.Leader.Domain
	if toDomain == "" {
		return
	}

	fromDomain := previousLeaderDomain
	if fromDomain == "" {
		fromDomain = "system"
	}

	message := fmt.Sprintf("You have been assigned as leader of %s", department.FullName)
	event := &events.LeaderAssignmentNotification{
		FromStarter: fromDomain,
		ToStarter:   toDomain,
		Message:     message,
		Type:        "leader_assignment",
		Timestamp:   time.Now().UTC(),
	}

	if err := s.notificationPub.PublishLeaderAssignment(ctx, event); err != nil {
		log.Printf("failed to publish leader assignment notification: %v", err)
	}
}

func (s *OrganizationApplicationService) ListBusinessUnits(ctx context.Context, query businessunitquery.ListBusinessUnitsQuery) (*httputil.PaginatedResult[*model.BusinessUnit], error) {
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

	return &httputil.PaginatedResult[*model.BusinessUnit]{
		Data: units,
		Pagination: httputil.RespPagination{
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

func (s *OrganizationApplicationService) ListBusinessUnitsWithDetails(ctx context.Context, query businessunitquery.ListBusinessUnitsQuery) (*httputil.PaginatedResult[*model.BusinessUnitWithDetails], error) {
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

	return &httputil.PaginatedResult[*model.BusinessUnitWithDetails]{
		Data: units,
		Pagination: httputil.RespPagination{
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
