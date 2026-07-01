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
	ErrTenantRequired       = errors.New("tenant organization is required")
	ErrGroupIDRequired      = errors.New("group id is required")
	ErrGroupIDInvalid       = errors.New("group id is invalid")
	ErrGroupNotFound        = errors.New("group not found in organization")
	ErrScheduleIDInvalid    = errors.New("schedule id is invalid")
	ErrScheduleNotFound     = errors.New("schedule not found in organization")
	ErrWeekdayInvalid       = errors.New("weekday is invalid")
	ErrStartTimeRequired    = errors.New("start time is required")
	ErrStartTimeInvalid     = errors.New("start time is invalid")
	ErrEndTimeRequired      = errors.New("end time is required")
	ErrEndTimeInvalid       = errors.New("end time is invalid")
	ErrScheduleTimeInvalid  = errors.New("end time must be after start time")
	ErrFromDateRequired     = errors.New("from date is required")
	ErrToDateRequired       = errors.New("to date is required")
	ErrFromDateInvalid      = errors.New("from date is invalid")
	ErrToDateInvalid        = errors.New("to date is invalid")
	ErrGenerateRangeInvalid = errors.New("to date must be after or equal from date")
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

func (s *Service) GenerateLessons(
	ctx context.Context,
	scheduleIDRaw string,
	req GenerateLessonsRequest,
) (GenerateLessonsResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return GenerateLessonsResponse{}, ErrTenantRequired
	}

	scheduleID, err := uuid.Parse(strings.TrimSpace(scheduleIDRaw))
	if err != nil {
		return GenerateLessonsResponse{}, ErrScheduleIDInvalid
	}

	fromDateRaw := strings.TrimSpace(req.FromDate)
	if fromDateRaw == "" {
		return GenerateLessonsResponse{}, ErrFromDateRequired
	}

	fromDate, err := time.Parse("2006-01-02", fromDateRaw)
	if err != nil {
		return GenerateLessonsResponse{}, ErrFromDateInvalid
	}

	toDateRaw := strings.TrimSpace(req.ToDate)
	if toDateRaw == "" {
		return GenerateLessonsResponse{}, ErrToDateRequired
	}

	toDate, err := time.Parse("2006-01-02", toDateRaw)
	if err != nil {
		return GenerateLessonsResponse{}, ErrToDateInvalid
	}

	if toDate.Before(fromDate) {
		return GenerateLessonsResponse{}, ErrGenerateRangeInvalid
	}

	topic := strings.TrimSpace(req.Topic)
	if topic == "" {
		topic = "Auto-generated lesson"
	}

	createdLessons, skippedDates, err := s.repo.GenerateLessons(
		ctx,
		*currentUser.OrganizationID,
		scheduleID,
		fromDate,
		toDate,
		topic,
	)
	if err != nil {
		return GenerateLessonsResponse{}, err
	}

	createdResponses := make([]GeneratedLessonResponse, 0, len(createdLessons))

	for _, item := range createdLessons {
		createdResponses = append(createdResponses, GeneratedLessonResponse{
			ID:         item.ID.String(),
			ScheduleID: item.ScheduleID.String(),
			GroupID:    item.GroupID.String(),
			LessonDate: item.LessonDate,
			StartTime:  item.StartTime,
			EndTime:    item.EndTime,
			Topic:      item.Topic,
			Status:     item.Status,
		})
	}

	return GenerateLessonsResponse{
		Created:      createdResponses,
		SkippedDates: skippedDates,
		CreatedCount: len(createdResponses),
		SkippedCount: len(skippedDates),
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
