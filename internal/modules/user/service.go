package user

import (
	"context"
	"errors"
	"strings"

	"github.com/Rumm1/eduhub-backend/internal/platform/password"
	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/google/uuid"
)

var (
	ErrTenantRequired   = errors.New("tenant organization is required")
	ErrEmailRequired    = errors.New("email is required")
	ErrPasswordRequired = errors.New("password is required")
	ErrFullNameRequired = errors.New("full name is required")
	ErrRoleRequired     = errors.New("role is required")
	ErrRoleInvalid      = errors.New("role is invalid")
	ErrBranchIDInvalid  = errors.New("branch id is invalid")
	ErrBranchNotFound   = errors.New("branch not found in organization")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateUserRequest) (UserResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return UserResponse{}, ErrTenantRequired
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" {
		return UserResponse{}, ErrEmailRequired
	}

	plainPassword := strings.TrimSpace(req.Password)
	if plainPassword == "" {
		return UserResponse{}, ErrPasswordRequired
	}

	fullName := strings.TrimSpace(req.FullName)
	if fullName == "" {
		return UserResponse{}, ErrFullNameRequired
	}

	roleCode := normalizeRoleCode(req.RoleCode)
	if roleCode == "" {
		return UserResponse{}, ErrRoleRequired
	}

	role, ok := roleTemplateByCode(roleCode)
	if !ok {
		return UserResponse{}, ErrRoleInvalid
	}

	branchIDs, err := parseBranchIDs(req.BranchIDs)
	if err != nil {
		return UserResponse{}, err
	}

	hashedPassword, err := password.Hash(plainPassword)
	if err != nil {
		return UserResponse{}, err
	}

	newUser := User{
		ID:             uuid.New(),
		OrganizationID: *currentUser.OrganizationID,
		Email:          email,
		PasswordHash:   hashedPassword,
		FullName:       fullName,
		Phone:          strings.TrimSpace(req.Phone),
		AvatarPath:     strings.TrimSpace(req.AvatarPath),
		Status:         "active",
	}

	createdUser, err := s.repo.CreateUserWithRoleAndBranches(ctx, newUser, role, branchIDs)
	if err != nil {
		return UserResponse{}, err
	}

	roles, savedBranchIDs, err := s.repo.BuildUserResponseData(ctx, createdUser.ID)
	if err != nil {
		return UserResponse{}, err
	}

	return buildUserResponse(createdUser, roles, savedBranchIDs), nil
}

func (s *Service) List(ctx context.Context) (ListUsersResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return ListUsersResponse{}, ErrTenantRequired
	}

	users, err := s.repo.ListByOrganizationID(ctx, *currentUser.OrganizationID)
	if err != nil {
		return ListUsersResponse{}, err
	}

	items := make([]UserResponse, 0, len(users))

	for _, item := range users {
		roles, branchIDs, err := s.repo.BuildUserResponseData(ctx, item.ID)
		if err != nil {
			return ListUsersResponse{}, err
		}

		items = append(items, buildUserResponse(item, roles, branchIDs))
	}

	return ListUsersResponse{
		Items: items,
		Total: len(items),
	}, nil
}

func normalizeRoleCode(roleCode string) string {
	roleCode = strings.TrimSpace(roleCode)
	roleCode = strings.ReplaceAll(roleCode, "-", "_")
	roleCode = strings.ToUpper(roleCode)

	return roleCode
}

func parseBranchIDs(rawBranchIDs []string) ([]uuid.UUID, error) {
	branchIDs := make([]uuid.UUID, 0, len(rawBranchIDs))

	for _, rawID := range rawBranchIDs {
		rawID = strings.TrimSpace(rawID)
		if rawID == "" {
			continue
		}

		branchID, err := uuid.Parse(rawID)
		if err != nil {
			return nil, ErrBranchIDInvalid
		}

		branchIDs = append(branchIDs, branchID)
	}

	return branchIDs, nil
}

func buildUserResponse(item User, roles []string, branchIDs []uuid.UUID) UserResponse {
	return UserResponse{
		ID:             item.ID.String(),
		OrganizationID: item.OrganizationID.String(),
		Email:          item.Email,
		FullName:       item.FullName,
		Phone:          item.Phone,
		AvatarPath:     item.AvatarPath,
		Status:         item.Status,
		Roles:          roles,
		BranchIDs:      uuidSliceToStringSlice(branchIDs),
	}
}

func uuidSliceToStringSlice(ids []uuid.UUID) []string {
	result := make([]string, 0, len(ids))

	for _, id := range ids {
		result = append(result, id.String())
	}

	return result
}

func roleTemplateByCode(code string) (RoleTemplate, bool) {
	templates := map[string]RoleTemplate{
		"TEACHER": {
			Code:        "TEACHER",
			Name:        "Teacher",
			Description: "Teacher role",
			Permissions: []string{
				"students.read",
				"groups.read",
				"lessons.read",
				"lessons.create",
				"lessons.update",
				"attendance.read",
				"attendance.manage",
				"homeworks.read",
				"homeworks.manage",
				"files.upload",
				"files.read",
				"notifications.read",
			},
		},
		"MANAGER": {
			Code:        "MANAGER",
			Name:        "Manager",
			Description: "Manager role",
			Permissions: []string{
				"branches.read",
				"users.read",
				"subjects.read",
				"teachers.read",
				"students.read",
				"students.create",
				"students.update",
				"groups.read",
				"groups.create",
				"groups.update",
				"lessons.read",
				"lessons.create",
				"lessons.update",
				"attendance.read",
				"homeworks.read",
				"payments.read",
				"files.upload",
				"files.read",
				"notifications.read",
				"notifications.manage",
			},
		},
		"ACCOUNTANT": {
			Code:        "ACCOUNTANT",
			Name:        "Accountant",
			Description: "Accountant role",
			Permissions: []string{
				"payments.read",
				"payments.manage",
				"payroll.read",
				"payroll.manage",
				"files.upload",
				"files.read",
				"notifications.read",
			},
		},
		"BRANCH_ADMIN": {
			Code:        "BRANCH_ADMIN",
			Name:        "Branch Admin",
			Description: "Branch administrator role",
			Permissions: []string{
				"branches.read",
				"users.read",
				"users.create",
				"users.update",
				"subjects.read",
				"teachers.read",
				"teachers.create",
				"teachers.update",
				"students.read",
				"students.create",
				"students.update",
				"groups.read",
				"groups.create",
				"groups.update",
				"lessons.read",
				"lessons.create",
				"lessons.update",
				"attendance.read",
				"attendance.manage",
				"homeworks.read",
				"homeworks.manage",
				"payments.read",
				"files.upload",
				"files.read",
				"notifications.read",
				"notifications.manage",
			},
		},
		"RECEPTIONIST": {
			Code:        "RECEPTIONIST",
			Name:        "Receptionist",
			Description: "Receptionist role",
			Permissions: []string{
				"branches.read",
				"students.read",
				"students.create",
				"students.update",
				"groups.read",
				"lessons.read",
				"payments.read",
				"payments.manage",
				"files.upload",
				"files.read",
				"notifications.read",
			},
		},
	}

	template, ok := templates[code]
	return template, ok
}
