package organization

import (
	"context"
	"errors"
	"strings"

	"github.com/Rumm1/eduhub-backend/internal/platform/password"
	"github.com/google/uuid"
)

var (
	ErrOrganizationNameRequired    = errors.New("organization name is required")
	ErrOrganizationLanguageInvalid = errors.New("organization language is invalid")
	ErrAdminEmailRequired          = errors.New("admin email is required")
	ErrAdminPasswordRequired       = errors.New("admin password is required")
	ErrAdminFullNameRequired       = errors.New("admin full name is required")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateOrganization(ctx context.Context, req CreateOrganizationRequest) (CreateOrganizationResponse, error) {
	defaultLanguage := normalizeLanguage(req.DefaultLanguage)
	if !isValidLanguage(defaultLanguage) {
		return CreateOrganizationResponse{}, ErrOrganizationLanguageInvalid
	}

	nameRU := strings.TrimSpace(req.NameRU)
	nameKK := strings.TrimSpace(req.NameKK)
	nameEN := strings.TrimSpace(req.NameEN)

	oldName := strings.TrimSpace(req.Name)
	if nameRU == "" && nameKK == "" && nameEN == "" && oldName != "" {
		switch defaultLanguage {
		case "kk":
			nameKK = oldName
		case "en":
			nameEN = oldName
		default:
			nameRU = oldName
		}
	}

	if nameRU == "" && nameKK == "" && nameEN == "" {
		return CreateOrganizationResponse{}, ErrOrganizationNameRequired
	}

	orgName := chooseOrganizationName(defaultLanguage, nameRU, nameKK, nameEN)

	adminEmail := strings.ToLower(strings.TrimSpace(req.AdminEmail))
	adminPassword := strings.TrimSpace(req.AdminPassword)
	adminFullName := strings.TrimSpace(req.AdminFullName)

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
		nameRU,
		nameKK,
		nameEN,
		defaultLanguage,
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
			ID:              createdOrg.ID.String(),
			Name:            createdOrg.Name,
			NameRU:          createdOrg.NameRU,
			NameKK:          createdOrg.NameKK,
			NameEN:          createdOrg.NameEN,
			DefaultLanguage: createdOrg.DefaultLanguage,
			BIN:             createdOrg.BIN,
			Phone:           createdOrg.Phone,
			Email:           createdOrg.Email,
			Status:          createdOrg.Status,
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

func normalizeLanguage(language string) string {
	language = strings.ToLower(strings.TrimSpace(language))
	if language == "" {
		return "ru"
	}

	return language
}

func isValidLanguage(language string) bool {
	return language == "ru" || language == "kk" || language == "en"
}

func chooseOrganizationName(defaultLanguage string, nameRU string, nameKK string, nameEN string) string {
	switch defaultLanguage {
	case "kk":
		if nameKK != "" {
			return nameKK
		}
	case "en":
		if nameEN != "" {
			return nameEN
		}
	default:
		if nameRU != "" {
			return nameRU
		}
	}

	if nameRU != "" {
		return nameRU
	}

	if nameKK != "" {
		return nameKK
	}

	return nameEN
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
		"payments.create",
		"payments.update_group_price",

		"payroll.read",
		"payroll.manage",
		"payroll.approve",
		"payroll.rules.manage",
		"payroll.read_all",
		"payroll.generate",
		"payroll.adjustments.manage",
		"payroll.send_to_teacher",
		"payroll.mark_paid",
		"payroll.export",

		"reports.teacher_schedule.read",
		"reports.payments.read",
		"reports.payroll.read",
		"reports.student_balance.read",
		"reports.export",

		"files.upload",
		"files.read",
		"files.delete",

		"notifications.read",
		"notifications.manage",

		"audit.read",
		"audit_logs.read",
		"audit_logs.export",
	}
}
