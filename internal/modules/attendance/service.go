package attendance

import (
	"context"
	"errors"
	"strings"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/google/uuid"
)

var (
	ErrTenantRequired           = errors.New("tenant organization is required")
	ErrLessonIDInvalid          = errors.New("lesson id is invalid")
	ErrLessonNotFound           = errors.New("lesson not found in organization")
	ErrAttendanceItemsRequired  = errors.New("attendance items are required")
	ErrStudentIDRequired        = errors.New("student id is required")
	ErrStudentIDInvalid         = errors.New("student id is invalid")
	ErrAttendanceStatusRequired = errors.New("attendance status is required")
	ErrAttendanceStatusInvalid  = errors.New("attendance status is invalid")
	ErrStudentNotInLessonGroup  = errors.New("student is not in lesson group")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) MarkLessonAttendance(
	ctx context.Context,
	lessonIDRaw string,
	req MarkLessonAttendanceRequest,
) (ListLessonAttendanceResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return ListLessonAttendanceResponse{}, ErrTenantRequired
	}

	lessonID, err := uuid.Parse(strings.TrimSpace(lessonIDRaw))
	if err != nil {
		return ListLessonAttendanceResponse{}, ErrLessonIDInvalid
	}

	if len(req.Items) == 0 {
		return ListLessonAttendanceResponse{}, ErrAttendanceItemsRequired
	}

	items := make([]Attendance, 0, len(req.Items))

	for _, reqItem := range req.Items {
		studentIDRaw := strings.TrimSpace(reqItem.StudentID)
		if studentIDRaw == "" {
			return ListLessonAttendanceResponse{}, ErrStudentIDRequired
		}

		studentID, err := uuid.Parse(studentIDRaw)
		if err != nil {
			return ListLessonAttendanceResponse{}, ErrStudentIDInvalid
		}

		status := strings.ToLower(strings.TrimSpace(reqItem.Status))
		if status == "" {
			return ListLessonAttendanceResponse{}, ErrAttendanceStatusRequired
		}

		if !isValidStatus(status) {
			return ListLessonAttendanceResponse{}, ErrAttendanceStatusInvalid
		}

		items = append(items, Attendance{
			LessonID:  lessonID,
			StudentID: studentID,
			Status:    status,
			Reason:    strings.ToLower(strings.TrimSpace(reqItem.Reason)),
			Comment:   strings.TrimSpace(reqItem.Comment),
		})
	}

	if err := s.repo.MarkLessonAttendance(
		ctx,
		*currentUser.OrganizationID,
		lessonID,
		currentUser.UserID,
		items,
	); err != nil {
		return ListLessonAttendanceResponse{}, err
	}

	return s.ListByLessonID(ctx, lessonIDRaw)
}

func (s *Service) ListByLessonID(ctx context.Context, lessonIDRaw string) (ListLessonAttendanceResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return ListLessonAttendanceResponse{}, ErrTenantRequired
	}

	lessonID, err := uuid.Parse(strings.TrimSpace(lessonIDRaw))
	if err != nil {
		return ListLessonAttendanceResponse{}, ErrLessonIDInvalid
	}

	items, err := s.repo.ListByLessonID(ctx, *currentUser.OrganizationID, lessonID)
	if err != nil {
		return ListLessonAttendanceResponse{}, err
	}

	responses := make([]AttendanceResponse, 0, len(items))

	for _, item := range items {
		responses = append(responses, AttendanceResponse{
			ID:              item.ID,
			LessonID:        item.LessonID.String(),
			StudentID:       item.StudentID.String(),
			StudentFullName: item.StudentFullName,
			Status:          item.Status,
			Reason:          item.Reason,
			Comment:         item.Comment,
			MarkedBy:        item.MarkedBy,
			MarkedAt:        item.MarkedAt,
		})
	}

	return ListLessonAttendanceResponse{
		Items: responses,
		Total: len(responses),
	}, nil
}

func isValidStatus(status string) bool {
	switch status {
	case "present", "absent", "late", "excused", "online":
		return true
	default:
		return false
	}
}
