package parent

import (
	"context"
	"strings"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/google/uuid"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context) ([]ParentResponse, error) {
	organizationID, err := getOrganizationID(ctx)
	if err != nil {
		return nil, err
	}

	items, err := s.repository.List(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	return mapParents(items), nil
}

func (s *Service) GetByID(ctx context.Context, id string) (ParentResponse, error) {
	organizationID, err := getOrganizationID(ctx)
	if err != nil {
		return ParentResponse{}, err
	}

	parentID, err := uuid.Parse(id)
	if err != nil {
		return ParentResponse{}, ErrParentIDInvalid
	}

	item, err := s.repository.GetByID(ctx, organizationID, parentID)
	if err != nil {
		return ParentResponse{}, err
	}

	return mapParent(item), nil
}

func (s *Service) Create(ctx context.Context, request CreateParentRequest) (ParentResponse, error) {
	organizationID, err := getOrganizationID(ctx)
	if err != nil {
		return ParentResponse{}, err
	}

	fullName := strings.TrimSpace(request.FullName)
	if fullName == "" {
		return ParentResponse{}, ErrFullNameRequired
	}

	item := Parent{
		ID:             uuid.New(),
		OrganizationID: organizationID,
		FullName:       fullName,
		Phone:          strings.TrimSpace(request.Phone),
		Email:          strings.TrimSpace(request.Email),
	}

	result, err := s.repository.Create(ctx, item)
	if err != nil {
		return ParentResponse{}, err
	}

	return mapParent(result), nil
}

func (s *Service) Update(ctx context.Context, id string, request UpdateParentRequest) (ParentResponse, error) {
	organizationID, err := getOrganizationID(ctx)
	if err != nil {
		return ParentResponse{}, err
	}

	parentID, err := uuid.Parse(id)
	if err != nil {
		return ParentResponse{}, ErrParentIDInvalid
	}

	input := Parent{
		FullName: strings.TrimSpace(request.FullName),
		Phone:    strings.TrimSpace(request.Phone),
		Email:    strings.TrimSpace(request.Email),
	}

	result, err := s.repository.Update(ctx, organizationID, parentID, input)
	if err != nil {
		return ParentResponse{}, err
	}

	return mapParent(result), nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	organizationID, err := getOrganizationID(ctx)
	if err != nil {
		return err
	}

	parentID, err := uuid.Parse(id)
	if err != nil {
		return ErrParentIDInvalid
	}

	return s.repository.Delete(ctx, organizationID, parentID)
}

func (s *Service) AttachStudent(ctx context.Context, parentIDRaw string, studentIDRaw string, request AttachStudentRequest) error {
	organizationID, err := getOrganizationID(ctx)
	if err != nil {
		return err
	}

	parentID, err := uuid.Parse(parentIDRaw)
	if err != nil {
		return ErrParentIDInvalid
	}

	studentID, err := uuid.Parse(studentIDRaw)
	if err != nil {
		return ErrStudentIDInvalid
	}

	return s.repository.AttachStudent(
		ctx,
		organizationID,
		parentID,
		studentID,
		strings.TrimSpace(request.Relation),
	)
}

func (s *Service) DetachStudent(ctx context.Context, parentIDRaw string, studentIDRaw string) error {
	organizationID, err := getOrganizationID(ctx)
	if err != nil {
		return err
	}

	parentID, err := uuid.Parse(parentIDRaw)
	if err != nil {
		return ErrParentIDInvalid
	}

	studentID, err := uuid.Parse(studentIDRaw)
	if err != nil {
		return ErrStudentIDInvalid
	}

	return s.repository.DetachStudent(ctx, organizationID, parentID, studentID)
}

func (s *Service) ListStudents(ctx context.Context, parentIDRaw string) ([]StudentResponse, error) {
	organizationID, err := getOrganizationID(ctx)
	if err != nil {
		return nil, err
	}

	parentID, err := uuid.Parse(parentIDRaw)
	if err != nil {
		return nil, ErrParentIDInvalid
	}

	items, err := s.repository.ListStudents(ctx, organizationID, parentID)
	if err != nil {
		return nil, err
	}

	return mapStudents(items), nil
}

func getOrganizationID(ctx context.Context) (uuid.UUID, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return uuid.Nil, ErrTenantRequired
	}

	return *currentUser.OrganizationID, nil
}

func mapParents(items []Parent) []ParentResponse {
	result := make([]ParentResponse, 0, len(items))

	for _, item := range items {
		result = append(result, mapParent(item))
	}

	return result
}

func mapParent(item Parent) ParentResponse {
	return ParentResponse{
		ID:             item.ID.String(),
		OrganizationID: item.OrganizationID.String(),
		FullName:       item.FullName,
		Phone:          item.Phone,
		Email:          item.Email,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}
}

func mapStudents(items []Student) []StudentResponse {
	result := make([]StudentResponse, 0, len(items))

	for _, item := range items {
		result = append(result, StudentResponse{
			ID:       item.ID.String(),
			BranchID: item.BranchID.String(),
			FullName: item.FullName,
			Phone:    item.Phone,
			Status:   item.Status,
			Relation: item.Relation,
		})
	}

	return result
}
