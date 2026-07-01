package organization

import (
	"context"
	"errors"
	"strings"

	"github.com/Rumm1/eduhub-backend/internal/platform/password"
	"github.com/google/uuid"
)

var (
	ErrOrganizationNameRequired = errors.New("organization name is required")
	ErrAdminEmailRequired       = errors.New("admin email is required")
	ErrAdminPasswordRequired    = errors.New("admin password is required")
	ErrAdminFullNameRequired    = errors.New("admin full name is required")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateOrganization(ctx context.Context, req CreateOrganizationRequest) (CreateOrganizationResponse, error) {
	orgName := strings.TrimSpace(req.Name)
	adminEmail := strings.ToLower(strings.TrimSpace(req.AdminEmail))
	adminPassword := strings.TrimSpace(req.AdminPassword)
	adminFullName := strings.TrimSpace(req.AdminFullName)

	if orgName == "" {
		return CreateOrganizationResponse{}, ErrOrganizationNameRequired
	}

	if adminEmail == "" {
		return CreateOrganizationResponse{}, ErrAdminEmailRequired
	}

	if adminPassword == "" {
		return CreateOrganizationResponse{}, ErrAdminPasswordRequired
	}

	if adminFullName == "" {
		return CreateOrganizationResponse{}, ErrAdminFullNameRequired
	}

	hashedPassword, err := password.Hash(adminPassword)
	if err != nil {
		return CreateOrganizationResponse{}, err
	}

	org := NewOrganization(
		orgName,
		strings.TrimSpace(req.BIN),
		strings.TrimSpace(req.Phone),
		strings.ToLower(strings.TrimSpace(req.Email)),
	)

	admin := User{
		ID:             uuid.New(),
		OrganizationID: org.ID,
		Email:          adminEmail,
		PasswordHash:   hashedPassword,
		FullName:       adminFullName,
		Phone:          strings.TrimSpace(req.AdminPhone),
		Status:         "active",
	}

	role := Role{
		ID:             uuid.New(),
		OrganizationID: org.ID,
		Name:           "Organization Admin",
		Code:           "ORG_ADMIN",
		Description:    "Organization administrator",
		IsSystem:       true,
	}

	createdOrg, createdAdmin, err := s.repo.CreateOrganizationWithAdmin(
		ctx,
		org,
		admin,
		role,
		orgAdminPermissions(),
	)
	if err != nil {
		return CreateOrganizationResponse{}, err
	}

	return CreateOrganizationResponse{
		Organization: OrganizationResponse{
			ID:     createdOrg.ID.String(),
			Name:   createdOrg.Name,
			BIN:    createdOrg.BIN,
			Phone:  createdOrg.Phone,
			Email:  createdOrg.Email,
			Status: createdOrg.Status,
		},
		Admin: AdminResponse{
			ID:       createdAdmin.ID.String(),
			Email:    createdAdmin.Email,
			FullName: createdAdmin.FullName,
			Phone:    createdAdmin.Phone,
			Roles:    []string{"ORG_ADMIN"},
		},
	}, nil
}

func orgAdminPermissions() []string {
	return []string{
		"branches.read",
		"branches.create",
		"branches.update",
		"branches.delete",

		"users.read",
		"users.create",
		"users.update",
		"users.delete",

		"roles.read",
		"roles.manage",

		"subjects.read",
		"subjects.create",
		"subjects.update",
		"subjects.delete",

		"teachers.read",
		"teachers.create",
		"teachers.update",
		"teachers.delete",

		"students.read",
		"students.create",
		"students.update",
		"students.delete",

		"groups.read",
		"groups.create",
		"groups.update",
		"groups.delete",

		"lessons.read",
		"lessons.create",
		"lessons.update",
		"lessons.delete",

		"attendance.read",
		"attendance.manage",

		"homeworks.read",
		"homeworks.manage",

		"payments.read",
		"payments.manage",

		"payroll.read",
		"payroll.manage",
		"payroll.approve",
		"payroll.rules.manage",

		"files.upload",
		"files.read",
		"files.delete",

		"notifications.read",
		"notifications.manage",

		"audit.read",
	}
}
