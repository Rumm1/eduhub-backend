package subject

import (
	"context"
	"errors"
	"strings"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/google/uuid"
)

var (
	ErrTenantRequired      = errors.New("tenant organization is required")
	ErrSubjectNameRequired = errors.New("subject name is required")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateSubjectRequest) (SubjectResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return SubjectResponse{}, ErrTenantRequired
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return SubjectResponse{}, ErrSubjectNameRequired
	}

	newSubject := Subject{
		ID:             uuid.New(),
		OrganizationID: *currentUser.OrganizationID,
		Name:           name,
		Description:    strings.TrimSpace(req.Description),
		Status:         "active",
	}

	createdSubject, err := s.repo.Create(ctx, newSubject)
	if err != nil {
		return SubjectResponse{}, err
	}

	return buildSubjectResponse(createdSubject), nil
}

func (s *Service) List(ctx context.Context) (ListSubjectsResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return ListSubjectsResponse{}, ErrTenantRequired
	}

	subjects, err := s.repo.ListByOrganizationID(ctx, *currentUser.OrganizationID)
	if err != nil {
		return ListSubjectsResponse{}, err
	}

	items := make([]SubjectResponse, 0, len(subjects))
	for _, item := range subjects {
		items = append(items, buildSubjectResponse(item))
	}

	return ListSubjectsResponse{
		Items: items,
		Total: len(items),
	}, nil
}

func buildSubjectResponse(item Subject) SubjectResponse {
	return SubjectResponse{
		ID:             item.ID.String(),
		OrganizationID: item.OrganizationID.String(),
		Name:           item.Name,
		Description:    item.Description,
		Status:         item.Status,
	}
}
