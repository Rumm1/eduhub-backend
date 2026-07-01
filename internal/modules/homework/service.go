package homework

import (
	"context"
	"errors"
	"strings"
	"time"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/google/uuid"
)

var (
	ErrTenantRequired   = errors.New("tenant organization is required")
	ErrLessonIDRequired = errors.New("lesson id is required")
	ErrLessonIDInvalid  = errors.New("lesson id is invalid")
	ErrLessonNotFound   = errors.New("lesson not found in organization")
	ErrTeacherNotFound  = errors.New("teacher not found for lesson")
	ErrTitleRequired    = errors.New("homework title is required")
	ErrDueDateInvalid   = errors.New("due date is invalid")
	ErrHomeworkDisabled = errors.New("homework is disabled for this subject or group")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateHomeworkRequest) (HomeworkResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return HomeworkResponse{}, ErrTenantRequired
	}

	lessonIDRaw := strings.TrimSpace(req.LessonID)
	if lessonIDRaw == "" {
		return HomeworkResponse{}, ErrLessonIDRequired
	}

	lessonID, err := uuid.Parse(lessonIDRaw)
	if err != nil {
		return HomeworkResponse{}, ErrLessonIDInvalid
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return HomeworkResponse{}, ErrTitleRequired
	}

	dueDate := strings.TrimSpace(req.DueDate)
	if dueDate != "" {
		if _, err := time.Parse("2006-01-02", dueDate); err != nil {
			return HomeworkResponse{}, ErrDueDateInvalid
		}
	}

	newHomework := Homework{
		ID:             uuid.New(),
		OrganizationID: *currentUser.OrganizationID,
		LessonID:       lessonID,
		Title:          title,
		Description:    strings.TrimSpace(req.Description),
		DueDate:        dueDate,
	}

	createdHomework, err := s.repo.Create(ctx, newHomework)
	if err != nil {
		return HomeworkResponse{}, err
	}

	return buildHomeworkResponse(createdHomework), nil
}

func (s *Service) List(ctx context.Context) (ListHomeworksResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return ListHomeworksResponse{}, ErrTenantRequired
	}

	homeworks, err := s.repo.ListByOrganizationID(ctx, *currentUser.OrganizationID)
	if err != nil {
		return ListHomeworksResponse{}, err
	}

	items := make([]HomeworkResponse, 0, len(homeworks))

	for _, item := range homeworks {
		items = append(items, buildHomeworkResponse(item))
	}

	return ListHomeworksResponse{
		Items: items,
		Total: len(items),
	}, nil
}

func (s *Service) ListByLessonID(ctx context.Context, lessonIDRaw string) (ListHomeworksResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return ListHomeworksResponse{}, ErrTenantRequired
	}

	lessonID, err := uuid.Parse(strings.TrimSpace(lessonIDRaw))
	if err != nil {
		return ListHomeworksResponse{}, ErrLessonIDInvalid
	}

	homeworks, err := s.repo.ListByLessonID(ctx, *currentUser.OrganizationID, lessonID)
	if err != nil {
		return ListHomeworksResponse{}, err
	}

	items := make([]HomeworkResponse, 0, len(homeworks))

	for _, item := range homeworks {
		items = append(items, buildHomeworkResponse(item))
	}

	return ListHomeworksResponse{
		Items: items,
		Total: len(items),
	}, nil
}

func buildHomeworkResponse(homework Homework) HomeworkResponse {
	return HomeworkResponse{
		ID:             homework.ID.String(),
		OrganizationID: homework.OrganizationID.String(),
		GroupID:        homework.GroupID.String(),
		LessonID:       homework.LessonID.String(),
		TeacherID:      homework.TeacherID,
		Title:          homework.Title,
		Description:    homework.Description,
		DueDate:        homework.DueDate,
	}
}
