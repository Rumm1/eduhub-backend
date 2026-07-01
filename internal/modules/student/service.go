package student

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
	ErrBranchIDRequired    = errors.New("branch id is required")
	ErrBranchIDInvalid     = errors.New("branch id is invalid")
	ErrBranchNotFound      = errors.New("branch not found in organization")
	ErrStudentNameRequired = errors.New("student name is required")
	ErrBirthDateInvalid    = errors.New("birth date is invalid")
	ErrParentNameRequired  = errors.New("parent name is required")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateStudentRequest) (StudentResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return StudentResponse{}, ErrTenantRequired
	}

	rawBranchID := strings.TrimSpace(req.BranchID)
	if rawBranchID == "" {
		return StudentResponse{}, ErrBranchIDRequired
	}

	branchID, err := uuid.Parse(rawBranchID)
	if err != nil {
		return StudentResponse{}, ErrBranchIDInvalid
	}

	fullName := strings.TrimSpace(req.FullName)
	if fullName == "" {
		return StudentResponse{}, ErrStudentNameRequired
	}

	birthDate := strings.TrimSpace(req.BirthDate)
	if birthDate != "" {
		if _, err := time.Parse("2006-01-02", birthDate); err != nil {
			return StudentResponse{}, ErrBirthDateInvalid
		}
	}

	newStudent := Student{
		ID:             uuid.New(),
		OrganizationID: *currentUser.OrganizationID,
		BranchID:       branchID,
		FullName:       fullName,
		Phone:          strings.TrimSpace(req.Phone),
		BirthDate:      birthDate,
		Gender:         strings.TrimSpace(req.Gender),
		Status:         "active",
		Source:         strings.TrimSpace(req.Source),
		Notes:          strings.TrimSpace(req.Notes),
	}

	var parent *Parent
	if req.Parent != nil {
		parentFullName := strings.TrimSpace(req.Parent.FullName)
		if parentFullName == "" {
			return StudentResponse{}, ErrParentNameRequired
		}

		parent = &Parent{
			ID:             uuid.New(),
			OrganizationID: *currentUser.OrganizationID,
			FullName:       parentFullName,
			Phone:          strings.TrimSpace(req.Parent.Phone),
			Email:          strings.ToLower(strings.TrimSpace(req.Parent.Email)),
			Relation:       strings.TrimSpace(req.Parent.Relation),
		}
	}

	createdStudent, err := s.repo.CreateWithParent(ctx, newStudent, parent)
	if err != nil {
		return StudentResponse{}, err
	}

	parents, err := s.repo.GetParentsByStudentID(ctx, createdStudent.ID)
	if err != nil {
		return StudentResponse{}, err
	}

	return buildStudentResponse(createdStudent, parents), nil
}

func (s *Service) List(ctx context.Context) (ListStudentsResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return ListStudentsResponse{}, ErrTenantRequired
	}

	students, err := s.repo.ListByOrganizationID(ctx, *currentUser.OrganizationID)
	if err != nil {
		return ListStudentsResponse{}, err
	}

	items := make([]StudentResponse, 0, len(students))

	for _, item := range students {
		parents, err := s.repo.GetParentsByStudentID(ctx, item.ID)
		if err != nil {
			return ListStudentsResponse{}, err
		}

		items = append(items, buildStudentResponse(item, parents))
	}

	return ListStudentsResponse{
		Items: items,
		Total: len(items),
	}, nil
}

func buildStudentResponse(student Student, parents []Parent) StudentResponse {
	return StudentResponse{
		ID:             student.ID.String(),
		OrganizationID: student.OrganizationID.String(),
		BranchID:       student.BranchID.String(),
		FullName:       student.FullName,
		Phone:          student.Phone,
		BirthDate:      student.BirthDate,
		Gender:         student.Gender,
		Status:         student.Status,
		Source:         student.Source,
		Notes:          student.Notes,
		Parents:        buildParentResponses(parents),
	}
}

func buildParentResponses(parents []Parent) []ParentResponse {
	result := make([]ParentResponse, 0, len(parents))

	for _, parent := range parents {
		result = append(result, ParentResponse{
			ID:       parent.ID.String(),
			FullName: parent.FullName,
			Phone:    parent.Phone,
			Email:    parent.Email,
			Relation: parent.Relation,
		})
	}

	return result
}
