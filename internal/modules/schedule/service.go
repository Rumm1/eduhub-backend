package schedule

import (
	"context"
	"errors"
	"strings"
	"time"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/google/uuid"
)

var (
	ErrTenantRequired      = errors.New("tenant organization is required")
	ErrGroupIDRequired     = errors.New("group id is required")
	ErrGroupIDInvalid      = errors.New("group id is invalid")
	ErrGroupNotFound       = errors.New("group not found in organization")
	ErrWeekdayInvalid      = errors.New("weekday is invalid")
	ErrStartTimeRequired   = errors.New("start time is required")
	ErrStartTimeInvalid    = errors.New("start time is invalid")
	ErrEndTimeRequired     = errors.New("end time is required")
	ErrEndTimeInvalid      = errors.New("end time is invalid")
	ErrScheduleTimeInvalid = errors.New("end time must be after start time")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateScheduleRequest) (ScheduleResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return ScheduleResponse{}, ErrTenantRequired
	}

	groupIDRaw := strings.TrimSpace(req.GroupID)
	if groupIDRaw == "" {
		return ScheduleResponse{}, ErrGroupIDRequired
	}

	groupID, err := uuid.Parse(groupIDRaw)
	if err != nil {
		return ScheduleResponse{}, ErrGroupIDInvalid
	}

	if req.Weekday < 1 || req.Weekday > 7 {
		return ScheduleResponse{}, ErrWeekdayInvalid
	}

	startTime := strings.TrimSpace(req.StartTime)
	if startTime == "" {
		return ScheduleResponse{}, ErrStartTimeRequired
	}

	parsedStartTime, err := time.Parse("15:04", startTime)
	if err != nil {
		return ScheduleResponse{}, ErrStartTimeInvalid
	}

	endTime := strings.TrimSpace(req.EndTime)
	if endTime == "" {
		return ScheduleResponse{}, ErrEndTimeRequired
	}

	parsedEndTime, err := time.Parse("15:04", endTime)
	if err != nil {
		return ScheduleResponse{}, ErrEndTimeInvalid
	}

	if !parsedEndTime.After(parsedStartTime) {
		return ScheduleResponse{}, ErrScheduleTimeInvalid
	}

	newSchedule := Schedule{
		ID:             uuid.New(),
		OrganizationID: *currentUser.OrganizationID,
		GroupID:        groupID,
		Weekday:        req.Weekday,
		StartTime:      startTime,
		EndTime:        endTime,
		Room:           strings.TrimSpace(req.Room),
	}

	createdSchedule, err := s.repo.Create(ctx, newSchedule)
	if err != nil {
		return ScheduleResponse{}, err
	}

	return buildScheduleResponse(createdSchedule), nil
}

func (s *Service) List(ctx context.Context) (ListSchedulesResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return ListSchedulesResponse{}, ErrTenantRequired
	}

	schedules, err := s.repo.ListByOrganizationID(ctx, *currentUser.OrganizationID)
	if err != nil {
		return ListSchedulesResponse{}, err
	}

	items := make([]ScheduleResponse, 0, len(schedules))

	for _, item := range schedules {
		items = append(items, buildScheduleResponse(item))
	}

	return ListSchedulesResponse{
		Items: items,
		Total: len(items),
	}, nil
}

func buildScheduleResponse(schedule Schedule) ScheduleResponse {
	return ScheduleResponse{
		ID:             schedule.ID.String(),
		OrganizationID: schedule.OrganizationID.String(),
		BranchID:       schedule.BranchID.String(),
		GroupID:        schedule.GroupID.String(),
		Weekday:        schedule.Weekday,
		StartTime:      schedule.StartTime,
		EndTime:        schedule.EndTime,
		Room:           schedule.Room,
	}
}
