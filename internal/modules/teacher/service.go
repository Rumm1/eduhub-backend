package teacher

import (
	"context"
	"errors"
	"strings"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/google/uuid"
)

var (
	ErrTenantRequired     = errors.New("tenant organization is required")
	ErrUserIDRequired     = errors.New("user id is required")
	ErrUserIDInvalid      = errors.New("user id is invalid")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserIsNotTeacher   = errors.New("user is not teacher")
	ErrSubjectIDInvalid   = errors.New("subject id is invalid")
	ErrSubjectNotFound    = errors.New("subject not found")
	ErrExperienceInvalid  = errors.New("experience years is invalid")
	ErrHourlyRateInvalid  = errors.New("hourly rate is invalid")
	ErrFixedSalaryInvalid = errors.New("fixed salary is invalid")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateTeacherRequest) (TeacherResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return TeacherResponse{}, ErrTenantRequired
	}

	rawUserID := strings.TrimSpace(req.UserID)
	if rawUserID == "" {
		return TeacherResponse{}, ErrUserIDRequired
	}

	userID, err := uuid.Parse(rawUserID)
	if err != nil {
		return TeacherResponse{}, ErrUserIDInvalid
	}

	if req.ExperienceYears < 0 {
		return TeacherResponse{}, ErrExperienceInvalid
	}

	if req.HourlyRate < 0 {
		return TeacherResponse{}, ErrHourlyRateInvalid
	}

	if req.FixedSalary < 0 {
		return TeacherResponse{}, ErrFixedSalaryInvalid
	}

	subjectIDs, err := parseSubjectIDs(req.SubjectIDs)
	if err != nil {
		return TeacherResponse{}, err
	}

	newTeacher := Teacher{
		UserID:          userID,
		OrganizationID:  *currentUser.OrganizationID,
		Bio:             strings.TrimSpace(req.Bio),
		ExperienceYears: req.ExperienceYears,
		EmploymentType:  strings.TrimSpace(req.EmploymentType),
		HourlyRate:      req.HourlyRate,
		FixedSalary:     req.FixedSalary,
	}

	createdTeacher, err := s.repo.Create(ctx, *currentUser.OrganizationID, newTeacher, subjectIDs)
	if err != nil {
		return TeacherResponse{}, err
	}

	subjects, err := s.repo.GetSubjectsByTeacherID(ctx, createdTeacher.UserID)
	if err != nil {
		return TeacherResponse{}, err
	}

	return buildTeacherResponse(createdTeacher, subjects), nil
}

func (s *Service) List(ctx context.Context) (ListTeachersResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return ListTeachersResponse{}, ErrTenantRequired
	}

	teachers, err := s.repo.ListByOrganizationID(ctx, *currentUser.OrganizationID)
	if err != nil {
		return ListTeachersResponse{}, err
	}

	items := make([]TeacherResponse, 0, len(teachers))

	for _, teacher := range teachers {
		subjects, err := s.repo.GetSubjectsByTeacherID(ctx, teacher.UserID)
		if err != nil {
			return ListTeachersResponse{}, err
		}

		items = append(items, buildTeacherResponse(teacher, subjects))
	}

	return ListTeachersResponse{
		Items: items,
		Total: len(items),
	}, nil
}

func parseSubjectIDs(rawSubjectIDs []string) ([]uuid.UUID, error) {
	subjectIDs := make([]uuid.UUID, 0, len(rawSubjectIDs))

	for _, rawID := range rawSubjectIDs {
		rawID = strings.TrimSpace(rawID)
		if rawID == "" {
			continue
		}

		subjectID, err := uuid.Parse(rawID)
		if err != nil {
			return nil, ErrSubjectIDInvalid
		}

		subjectIDs = append(subjectIDs, subjectID)
	}

	return subjectIDs, nil
}

func buildTeacherResponse(teacher Teacher, subjects []Subject) TeacherResponse {
	return TeacherResponse{
		UserID:          teacher.UserID.String(),
		OrganizationID:  teacher.OrganizationID.String(),
		Email:           teacher.Email,
		FullName:        teacher.FullName,
		Phone:           teacher.Phone,
		Bio:             teacher.Bio,
		ExperienceYears: teacher.ExperienceYears,
		EmploymentType:  teacher.EmploymentType,
		HourlyRate:      teacher.HourlyRate,
		FixedSalary:     teacher.FixedSalary,
		Subjects:        buildSubjectResponses(subjects),
	}
}

func buildSubjectResponses(subjects []Subject) []SubjectResponse {
	result := make([]SubjectResponse, 0, len(subjects))

	for _, subject := range subjects {
		result = append(result, SubjectResponse{
			ID:   subject.ID.String(),
			Name: subject.Name,
		})
	}

	return result
}
