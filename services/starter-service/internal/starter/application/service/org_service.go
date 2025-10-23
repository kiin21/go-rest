package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/kiin21/go-rest/pkg/events"
	"github.com/kiin21/go-rest/pkg/httputil"
	businessunitquery "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/businessunit/query"
	departmentcommand "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/department/command"
	departmentquery "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/department/query"
	sharedDomain "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/error"
	domainmessaging "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/messaging"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
)

type OrganizationApplicationService struct {
	departmentRepo   repository.DepartmentRepository
	businessUnitRepo repository.BusinessUnitRepository
	starterRepo      repository.StarterRepository
	notificationPub  domainmessaging.NotificationProducer
}

func NewOrganizationApplicationService(
	departmentRepo repository.DepartmentRepository,
	businessUnitRepo repository.BusinessUnitRepository,
	starterRepo repository.StarterRepository,
	notificationPublisher domainmessaging.NotificationProducer,
) *OrganizationApplicationService {
	return &OrganizationApplicationService{
		departmentRepo:   departmentRepo,
		businessUnitRepo: businessUnitRepo,
		starterRepo:      starterRepo,
		notificationPub:  notificationPublisher,
	}
}

func (s *OrganizationApplicationService) GetAllDepartments(ctx context.Context, query *departmentquery.ListDepartmentsQuery) (*httputil.PaginatedResult[*model.DepartmentWithDetails], error) {
	filter := &model.DepartmentListFilter{BusinessUnitID: query.BusinessUnitID}

	departments, total, err := s.departmentRepo.ListWithDetails(ctx, filter, &query.Pagination)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / query.Pagination.GetLimit()
	if int(total)%(query.Pagination.GetLimit()) > 0 {
		totalPages++
	}

	var prev *string = nil
	var next *string = nil

	if query.Pagination.GetPage() > 1 {
		value := strconv.Itoa(query.Pagination.GetPage() - 1)
		prev = &value
	}
	if query.Pagination.GetPage() < totalPages {
		value := strconv.Itoa(query.Pagination.GetPage() + 1)
		next = &value
	}

	return &httputil.PaginatedResult[*model.DepartmentWithDetails]{
		Data: departments,
		Pagination: httputil.RespPagination{
			Limit:      query.Pagination.GetLimit(),
			TotalItems: total,
			Prev:       prev,
			Next:       next,
		},
	}, nil
}

func (s *OrganizationApplicationService) GetOneDepartment(ctx context.Context, ID int64) (*model.DepartmentWithDetails, error) {
	ids := make([]int64, 0, 1)
	ids = append(ids, ID)

	departments, err := s.departmentRepo.FindByIDsWithDetails(ctx, ids)

	if err != nil {
		return nil, err
	}
	if len(departments) == 0 {
		return nil, sharedDomain.ErrNotFound
	}

	return departments[0], nil
}

func (s *OrganizationApplicationService) CreateDepartment(ctx context.Context, cmd *departmentcommand.CreateDepartmentCommand) (*model.DepartmentWithDetails, error) {
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

func (s *OrganizationApplicationService) UpdateDepartment(ctx context.Context, cmd *departmentcommand.UpdateDepartmentCommand) (*model.DepartmentWithDetails, error) {

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

func (s *OrganizationApplicationService) DeleteDepartment(ctx context.Context, ID int64) error {
	return s.departmentRepo.Delete(ctx, ID)
}

func (s *OrganizationApplicationService) AssignLeader(ctx context.Context, cmd *departmentcommand.AssignLeaderCommand) (*model.DepartmentWithDetails, error) {
	deptIDList := []int64{cmd.DepartmentID}
	departments, err := s.departmentRepo.FindByIDsWithDetails(ctx, deptIDList)
	if err != nil {
		return nil, err
	}
	if len(departments) == 0 {
		return nil, sharedDomain.ErrNotFound
	}

	oldDept := departments[0]
	previousLeaderDomain := ""
	if oldDept.Leader != nil {
		previousLeaderDomain = oldDept.Leader.Domain
	}

	if cmd.LeaderID != nil && cmd.LeaderDomain == nil { // assign by starter ID
		oldDept.LeaderID = cmd.LeaderID
	} else if cmd.LeaderID == nil && cmd.LeaderDomain != nil { // assign by starter domain
		starter, err := s.starterRepo.FindByDomain(ctx, *cmd.LeaderDomain)

		if err != nil {
			return nil, err
		}

		leaderID := starter.ID()
		oldDept.LeaderID = &leaderID
	} else {
		return nil, sharedDomain.ErrInvalidInput
	}

	if err := s.departmentRepo.Update(ctx, &model.Department{
		ID:                oldDept.ID,
		GroupDepartmentID: oldDept.GroupDepartmentID,
		FullName:          oldDept.FullName,
		Shortname:         oldDept.Shortname,
		BusinessUnitID:    oldDept.BusinessUnitID,
		LeaderID:          oldDept.LeaderID,
		CreatedAt:         oldDept.CreatedAt,
		UpdatedAt:         oldDept.UpdatedAt,
		DeletedAt:         oldDept.DeletedAt,
	}); err != nil {
		return nil, err
	}

	updated, err := s.departmentRepo.FindByIDsWithDetails(ctx, []int64{oldDept.ID})
	if err != nil {
		return nil, err
	}
	if len(updated) == 0 {
		return nil, sharedDomain.ErrNotFound
	}
	oldDept = updated[0]

	// Send msg to a Kafka topic
	s.publishLeaderAssignmentNotification(oldDept, previousLeaderDomain)

	return oldDept, nil
}

func (s *OrganizationApplicationService) publishLeaderAssignmentNotification(department *model.DepartmentWithDetails, previousLeaderDomain string) {
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
	payload := events.LeaderAssignmentEventPayload{
		FromStarter: fromDomain,
		ToStarter:   toDomain,
		Message:     message,
	}

	// Create event
	event, err := events.NewEvent(events.EventTypeNotificationLeaderAssignment, payload)
	if err != nil {
		log.Printf("failed to create leader assignment notification event: %v", err)
		return
	}

	if err := s.notificationPub.SendNotification(event); err != nil {
		log.Printf("failed to publish leader assignment notification: %v", err)
	}
}

func (s *OrganizationApplicationService) ListBusinessUnits(
	ctx context.Context, query *businessunitquery.ListBusinessUnitsQuery,
) (*httputil.PaginatedResult[*model.BusinessUnit], error) {
	units, total, err := s.businessUnitRepo.List(ctx, query.Pagination)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / (query.Pagination.GetLimit())
	if int(total)%(query.Pagination.GetLimit()) > 0 {
		totalPages++
	}

	var prev, next *string
	if query.Pagination.GetPage() > 1 {
		value := strconv.Itoa(query.Pagination.GetPage() - 1)
		prev = &value
	}
	if query.Pagination.GetPage() < totalPages {
		value := strconv.Itoa(query.Pagination.GetPage() + 1)
		next = &value
	}

	return &httputil.PaginatedResult[*model.BusinessUnit]{
		Data: units,
		Pagination: httputil.RespPagination{
			Limit:      query.Pagination.GetLimit(),
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

func (s *OrganizationApplicationService) ListBusinessUnitsWithDetails(ctx context.Context, query *businessunitquery.ListBusinessUnitsQuery) (*httputil.PaginatedResult[*model.BusinessUnitWithDetails], error) {
	units, total, err := s.businessUnitRepo.ListWithDetails(ctx, query.Pagination)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / query.Pagination.GetLimit()
	if int(total)%(query.Pagination.GetLimit()) > 0 {
		totalPages++
	}

	// TODO: refactor this kind of handling pagination into a common utility function
	var prev, next *string
	if query.Pagination.GetPage() > 1 {
		value := strconv.Itoa(query.Pagination.GetPage() - 1)
		prev = &value
	}
	if query.Pagination.GetPage() < totalPages {
		value := strconv.Itoa(query.Pagination.GetPage() + 1)
		next = &value
	}

	return &httputil.PaginatedResult[*model.BusinessUnitWithDetails]{
		Data: units,
		Pagination: httputil.RespPagination{
			Limit:      query.Pagination.GetLimit(),
			TotalItems: total,
			Prev:       prev,
			Next:       next,
		},
	}, nil
}

func (s *OrganizationApplicationService) GetBusinessUnitWithDetails(ctx context.Context, id int64) (*model.BusinessUnitWithDetails, error) {
	return s.businessUnitRepo.FindByIDWithDetails(ctx, id)
}
