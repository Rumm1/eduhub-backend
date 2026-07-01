package group

import (
	"context"
	"errors"
	"strings"
	"time"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/google/uuid"
)

var (
	ErrTenantRequired        = errors.New("tenant organization is required")
	ErrBranchIDRequired      = errors.New("branch id is required")
	ErrBranchIDInvalid       = errors.New("branch id is invalid")
	ErrSubjectIDRequired     = errors.New("subject id is required")
	ErrSubjectIDInvalid      = errors.New("subject id is invalid")
	ErrTeacherIDInvalid      = errors.New("teacher id is invalid")
	ErrGroupIDInvalid        = errors.New("group id is invalid")
	ErrStudentIDRequired     = errors.New("student id is required")
	ErrStudentIDInvalid      = errors.New("student id is invalid")
	ErrGroupNameRequired     = errors.New("group name is required")
	ErrStartDateInvalid      = errors.New("start date is invalid")
	ErrEndDateInvalid        = errors.New("end date is invalid")
	ErrBranchNotFound        = errors.New("branch not found in organization")
	ErrSubjectNotFound       = errors.New("subject not found in organization")
	ErrTeacherNotFound       = errors.New("teacher not found in organization")
	ErrGroupNotFound         = errors.New("group not found in organization")
	ErrStudentNotFound       = errors.New("student not found in organization")
	ErrStudentBranchMismatch = errors.New("student branch does not match group branch")
	ErrGroupIsFull           = errors.New("group is full")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateGroupRequest) (GroupResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return GroupResponse{}, ErrTenantRequired
	}

	branchID, err := parseRequiredUUID(req.BranchID, ErrBranchIDRequired, ErrBranchIDInvalid)
	if err != nil {
		return GroupResponse{}, err
	}

	subjectID, err := parseRequiredUUID(req.SubjectID, ErrSubjectIDRequired, ErrSubjectIDInvalid)
	if err != nil {
		return GroupResponse{}, err
	}

	teacherID := strings.TrimSpace(req.TeacherID)
	if teacherID != "" {
		if _, err := uuid.Parse(teacherID); err != nil {
			return GroupResponse{}, ErrTeacherIDInvalid
		}
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return GroupResponse{}, ErrGroupNameRequired
	}

	startDate := strings.TrimSpace(req.StartDate)
	if startDate != "" {
		if _, err := time.Parse("2006-01-02", startDate); err != nil {
			return GroupResponse{}, ErrStartDateInvalid
		}
	}

	endDate := strings.TrimSpace(req.EndDate)
	if endDate != "" {
		if _, err := time.Parse("2006-01-02", endDate); err != nil {
			return GroupResponse{}, ErrEndDateInvalid
		}
	}

	maxStudents := req.MaxStudents
	if maxStudents <= 0 {
		maxStudents = 15
	}

	newGroup := Group{
		ID:             uuid.New(),
		OrganizationID: *currentUser.OrganizationID,
		BranchID:       branchID,
		SubjectID:      subjectID,
		TeacherID:      teacherID,
		Name:           name,
		Level:          strings.TrimSpace(req.Level),
		Status:         "active",
		MaxStudents:    maxStudents,
		StartDate:      startDate,
		EndDate:        endDate,
	}

	createdGroup, err := s.repo.Create(ctx, newGroup)
	if err != nil {
		return GroupResponse{}, err
	}

	return buildGroupResponse(createdGroup), nil
}

func (s *Service) List(ctx context.Context) (ListGroupsResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return ListGroupsResponse{}, ErrTenantRequired
	}

	groups, err := s.repo.ListByOrganizationID(ctx, *currentUser.OrganizationID)
	if err != nil {
		return ListGroupsResponse{}, err
	}

	items := make([]GroupResponse, 0, len(groups))

	for _, item := range groups {
		items = append(items, buildGroupResponse(item))
	}

	return ListGroupsResponse{
		Items: items,
		Total: len(items),
	}, nil
}

func (s *Service) AddStudent(ctx context.Context, groupIDRaw string, req AddStudentToGroupRequest) error {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return ErrTenantRequired
	}

	groupID, err := uuid.Parse(strings.TrimSpace(groupIDRaw))
	if err != nil {
		return ErrGroupIDInvalid
	}

	studentID, err := parseRequiredUUID(req.StudentID, ErrStudentIDRequired, ErrStudentIDInvalid)
	if err != nil {
		return err
	}

	return s.repo.AddStudent(ctx, *currentUser.OrganizationID, groupID, studentID)
}

func (s *Service) ListStudents(ctx context.Context, groupIDRaw string) (ListGroupStudentsResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return ListGroupStudentsResponse{}, ErrTenantRequired
	}

	groupID, err := uuid.Parse(strings.TrimSpace(groupIDRaw))
	if err != nil {
		return ListGroupStudentsResponse{}, ErrGroupIDInvalid
	}

	students, err := s.repo.ListStudents(ctx, *currentUser.OrganizationID, groupID)
	if err != nil {
		return ListGroupStudentsResponse{}, err
	}

	items := make([]GroupStudentResponse, 0, len(students))

	for _, item := range students {
		items = append(items, GroupStudentResponse{
			StudentID: item.StudentID.String(),
			FullName:  item.FullName,
			Phone:     item.Phone,
			Status:    item.Status,
			JoinedAt:  item.JoinedAt,
		})
	}

	return ListGroupStudentsResponse{
		Items: items,
		Total: len(items),
	}, nil
}

func parseRequiredUUID(value string, requiredErr error, invalidErr error) (uuid.UUID, error) {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return uuid.Nil, requiredErr
	}

	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, invalidErr
	}

	return id, nil
}

func buildGroupResponse(group Group) GroupResponse {
	return GroupResponse{
		ID:             group.ID.String(),
		OrganizationID: group.OrganizationID.String(),
		BranchID:       group.BranchID.String(),
		SubjectID:      group.SubjectID.String(),
		TeacherID:      group.TeacherID,
		Name:           group.Name,
		Level:          group.Level,
		Status:         group.Status,
		MaxStudents:    group.MaxStudents,
		StartDate:      group.StartDate,
		EndDate:        group.EndDate,
		StudentsCount:  group.StudentsCount,
	}
}
