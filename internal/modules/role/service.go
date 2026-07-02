package role

import (
	"context"
	"errors"
	"strings"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/google/uuid"
)

var (
	ErrTenantRequired     = errors.New("tenant organization is required")
	ErrRoleIDInvalid      = errors.New("role id is invalid")
	ErrRoleNotFound       = errors.New("role not found")
	ErrRoleNameRequired   = errors.New("role name is required")
	ErrRoleCodeRequired   = errors.New("role code is required")
	ErrPermissionRequired = errors.New("permission is required")
	ErrPermissionNotFound = errors.New("permission not found")
	ErrSystemRoleReadonly = errors.New("system role is readonly")
	ErrRoleInUse          = errors.New("role is in use")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context) (ListRolesResponse, error) {
	organizationID, err := getTenantOrganizationID(ctx)
	if err != nil {
		return ListRolesResponse{}, err
	}

	items, err := s.repo.List(ctx, organizationID)
	if err != nil {
		return ListRolesResponse{}, err
	}

	return ListRolesResponse{
		Items: buildRoleResponses(items),
		Total: len(items),
	}, nil
}

func (s *Service) GetByID(ctx context.Context, rawRoleID string) (RoleResponse, error) {
	organizationID, err := getTenantOrganizationID(ctx)
	if err != nil {
		return RoleResponse{}, err
	}

	roleID, err := parseUUID(rawRoleID, ErrRoleIDInvalid)
	if err != nil {
		return RoleResponse{}, err
	}

	item, err := s.repo.GetByID(ctx, organizationID, roleID)
	if err != nil {
		return RoleResponse{}, err
	}

	return buildRoleResponse(item), nil
}

func (s *Service) Create(ctx context.Context, req CreateRoleRequest) (RoleResponse, error) {
	organizationID, err := getTenantOrganizationID(ctx)
	if err != nil {
		return RoleResponse{}, err
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return RoleResponse{}, ErrRoleNameRequired
	}

	code := normalizeRoleCode(req.Code)
	if code == "" {
		return RoleResponse{}, ErrRoleCodeRequired
	}

	permissionCodes := normalizePermissionCodes(req.PermissionCodes)

	item := Role{
		ID:              uuid.New(),
		OrganizationID:  &organizationID,
		Name:            name,
		Code:            code,
		Description:     strings.TrimSpace(req.Description),
		IsSystem:        false,
		PermissionCodes: permissionCodes,
	}

	created, err := s.repo.Create(ctx, organizationID, item)
	if err != nil {
		return RoleResponse{}, err
	}

	return buildRoleResponse(created), nil
}

func (s *Service) Update(ctx context.Context, rawRoleID string, req UpdateRoleRequest) (RoleResponse, error) {
	organizationID, err := getTenantOrganizationID(ctx)
	if err != nil {
		return RoleResponse{}, err
	}

	roleID, err := parseUUID(rawRoleID, ErrRoleIDInvalid)
	if err != nil {
		return RoleResponse{}, err
	}

	input := UpdateRoleInput{}

	if req.Name != nil {
		value := strings.TrimSpace(*req.Name)
		if value == "" {
			return RoleResponse{}, ErrRoleNameRequired
		}
		input.Name = &value
	}

	if req.Code != nil {
		value := normalizeRoleCode(*req.Code)
		if value == "" {
			return RoleResponse{}, ErrRoleCodeRequired
		}
		input.Code = &value
	}

	if req.Description != nil {
		value := strings.TrimSpace(*req.Description)
		input.Description = &value
	}

	updated, err := s.repo.Update(ctx, organizationID, roleID, input)
	if err != nil {
		return RoleResponse{}, err
	}

	return buildRoleResponse(updated), nil
}

func (s *Service) Delete(ctx context.Context, rawRoleID string) error {
	organizationID, err := getTenantOrganizationID(ctx)
	if err != nil {
		return err
	}

	roleID, err := parseUUID(rawRoleID, ErrRoleIDInvalid)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, organizationID, roleID)
}

func (s *Service) AddPermission(ctx context.Context, rawRoleID string, req AddPermissionRequest) (RoleResponse, error) {
	organizationID, err := getTenantOrganizationID(ctx)
	if err != nil {
		return RoleResponse{}, err
	}

	roleID, err := parseUUID(rawRoleID, ErrRoleIDInvalid)
	if err != nil {
		return RoleResponse{}, err
	}

	permissionCode := normalizePermissionCode(req.PermissionCode)
	if permissionCode == "" {
		return RoleResponse{}, ErrPermissionRequired
	}

	updated, err := s.repo.AddPermission(ctx, organizationID, roleID, permissionCode)
	if err != nil {
		return RoleResponse{}, err
	}

	return buildRoleResponse(updated), nil
}

func (s *Service) RemovePermission(ctx context.Context, rawRoleID string, rawPermissionCode string) (RoleResponse, error) {
	organizationID, err := getTenantOrganizationID(ctx)
	if err != nil {
		return RoleResponse{}, err
	}

	roleID, err := parseUUID(rawRoleID, ErrRoleIDInvalid)
	if err != nil {
		return RoleResponse{}, err
	}

	permissionCode := normalizePermissionCode(rawPermissionCode)
	if permissionCode == "" {
		return RoleResponse{}, ErrPermissionRequired
	}

	updated, err := s.repo.RemovePermission(ctx, organizationID, roleID, permissionCode)
	if err != nil {
		return RoleResponse{}, err
	}

	return buildRoleResponse(updated), nil
}

func getTenantOrganizationID(ctx context.Context) (uuid.UUID, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return uuid.Nil, ErrTenantRequired
	}

	return *currentUser.OrganizationID, nil
}

func normalizeRoleCode(code string) string {
	code = strings.TrimSpace(code)
	code = strings.ReplaceAll(code, "-", "_")
	code = strings.ReplaceAll(code, " ", "_")
	code = strings.ToUpper(code)

	return code
}

func normalizePermissionCode(code string) string {
	return strings.TrimSpace(code)
}

func normalizePermissionCodes(rawCodes []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(rawCodes))

	for _, rawCode := range rawCodes {
		code := normalizePermissionCode(rawCode)
		if code == "" {
			continue
		}

		if _, exists := seen[code]; exists {
			continue
		}

		seen[code] = struct{}{}
		result = append(result, code)
	}

	return result
}

func parseUUID(rawID string, invalidErr error) (uuid.UUID, error) {
	rawID = strings.TrimSpace(rawID)
	if rawID == "" {
		return uuid.Nil, invalidErr
	}

	id, err := uuid.Parse(rawID)
	if err != nil {
		return uuid.Nil, invalidErr
	}

	return id, nil
}

func buildRoleResponses(items []Role) []RoleResponse {
	responses := make([]RoleResponse, 0, len(items))

	for _, item := range items {
		responses = append(responses, buildRoleResponse(item))
	}

	return responses
}

func buildRoleResponse(item Role) RoleResponse {
	return RoleResponse{
		ID:              item.ID.String(),
		OrganizationID:  uuidPointerToString(item.OrganizationID),
		Name:            item.Name,
		Code:            item.Code,
		Description:     item.Description,
		IsSystem:        item.IsSystem,
		PermissionCodes: item.PermissionCodes,
	}
}

func uuidPointerToString(id *uuid.UUID) string {
	if id == nil {
		return ""
	}

	return id.String()
}
