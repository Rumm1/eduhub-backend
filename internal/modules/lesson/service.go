package lesson

import (
	"context"
	"errors"
	"strings"
	"time"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/google/uuid"
)

var (
	ErrTenantRequired          = errors.New("tenant organization is required")
	ErrGroupIDRequired         = errors.New("group id is required")
	ErrGroupIDInvalid          = errors.New("group id is invalid")
	ErrGroupNotFound           = errors.New("group not found in organization")
	ErrLessonIDInvalid         = errors.New("lesson id is invalid")
	ErrLessonNotFound          = errors.New("lesson not found in organization")
	ErrLessonDateRequired      = errors.New("lesson date is required")
	ErrLessonDateInvalid       = errors.New("lesson date is invalid")
	ErrStartTimeRequired       = errors.New("start time is required")
	ErrStartTimeInvalid        = errors.New("start time is invalid")
	ErrEndTimeRequired         = errors.New("end time is required")
	ErrEndTimeInvalid          = errors.New("end time is invalid")
	ErrLessonTimeInvalid       = errors.New("end time must be after start time")
	ErrActualTeacherIDRequired = errors.New("actual teacher id is required")
	ErrActualTeacherIDInvalid  = errors.New("actual teacher id is invalid")
	ErrActualTeacherNotFound   = errors.New("actual teacher not found in organization")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateLessonRequest) (LessonResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return LessonResponse{}, ErrTenantRequired
	}

	groupIDRaw := strings.TrimSpace(req.GroupID)
	if groupIDRaw == "" {
		return LessonResponse{}, ErrGroupIDRequired
	}

	groupID, err := uuid.Parse(groupIDRaw)
	if err != nil {
		return LessonResponse{}, ErrGroupIDInvalid
	}

	lessonDate := strings.TrimSpace(req.LessonDate)
	if lessonDate == "" {
		return LessonResponse{}, ErrLessonDateRequired
	}

	if _, err := time.Parse("2006-01-02", lessonDate); err != nil {
		return LessonResponse{}, ErrLessonDateInvalid
	}

	startTime := strings.TrimSpace(req.StartTime)
	if startTime == "" {
		return LessonResponse{}, ErrStartTimeRequired
	}

	parsedStartTime, err := time.Parse("15:04", startTime)
	if err != nil {
		return LessonResponse{}, ErrStartTimeInvalid
	}

	endTime := strings.TrimSpace(req.EndTime)
	if endTime == "" {
		return LessonResponse{}, ErrEndTimeRequired
	}

	parsedEndTime, err := time.Parse("15:04", endTime)
	if err != nil {
		return LessonResponse{}, ErrEndTimeInvalid
	}

	if !parsedEndTime.After(parsedStartTime) {
		return LessonResponse{}, ErrLessonTimeInvalid
	}

	newLesson := Lesson{
		ID:             uuid.New(),
		OrganizationID: *currentUser.OrganizationID,
		GroupID:        groupID,
		LessonDate:     lessonDate,
		StartTime:      startTime,
		EndTime:        endTime,
		Topic:          strings.TrimSpace(req.Topic),
		Status:         "planned",
	}

	createdLesson, err := s.repo.Create(ctx, newLesson)
	if err != nil {
		return LessonResponse{}, err
	}

	return buildLessonResponse(createdLesson), nil
}

func (s *Service) List(ctx context.Context) (ListLessonsResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return ListLessonsResponse{}, ErrTenantRequired
	}

	lessons, err := s.repo.ListByOrganizationID(ctx, *currentUser.OrganizationID)
	if err != nil {
		return ListLessonsResponse{}, err
	}

	items := make([]LessonResponse, 0, len(lessons))

	for _, item := range lessons {
		items = append(items, buildLessonResponse(item))
	}

	return ListLessonsResponse{
		Items: items,
		Total: len(items),
	}, nil
}

func (s *Service) UpdateActualTeacher(
	ctx context.Context,
	lessonIDRaw string,
	req UpdateLessonTeacherRequest,
) (LessonResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return LessonResponse{}, ErrTenantRequired
	}

	lessonID, err := uuid.Parse(strings.TrimSpace(lessonIDRaw))
	if err != nil {
		return LessonResponse{}, ErrLessonIDInvalid
	}

	actualTeacherIDRaw := strings.TrimSpace(req.ActualTeacherID)
	if actualTeacherIDRaw == "" {
		return LessonResponse{}, ErrActualTeacherIDRequired
	}

	actualTeacherID, err := uuid.Parse(actualTeacherIDRaw)
	if err != nil {
		return LessonResponse{}, ErrActualTeacherIDInvalid
	}

	updatedLesson, err := s.repo.UpdateActualTeacher(
		ctx,
		*currentUser.OrganizationID,
		lessonID,
		actualTeacherID,
		strings.TrimSpace(req.SubstitutionReason),
	)
	if err != nil {
		return LessonResponse{}, err
	}

	return buildLessonResponse(updatedLesson), nil
}

func buildLessonResponse(lesson Lesson) LessonResponse {
	return LessonResponse{
		ID:                 lesson.ID.String(),
		OrganizationID:     lesson.OrganizationID.String(),
		BranchID:           lesson.BranchID.String(),
		GroupID:            lesson.GroupID.String(),
		TeacherID:          lesson.TeacherID,
		PlannedTeacherID:   lesson.PlannedTeacherID,
		ActualTeacherID:    lesson.ActualTeacherID,
		SubjectID:          lesson.SubjectID.String(),
		ScheduleID:         lesson.ScheduleID,
		LessonDate:         lesson.LessonDate,
		StartTime:          lesson.StartTime,
		EndTime:            lesson.EndTime,
		Topic:              lesson.Topic,
		Status:             lesson.Status,
		SubstitutionReason: lesson.SubstitutionReason,
	}
}
