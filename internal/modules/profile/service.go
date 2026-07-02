package profile

import (
	"context"
	"errors"
	"strings"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/google/uuid"
)

var (
	ErrTenantRequired              = errors.New("tenant organization is required")
	ErrUserIDInvalid               = errors.New("user id is invalid")
	ErrProfileIDInvalid            = errors.New("profile id is invalid")
	ErrBranchIDInvalid             = errors.New("branch id is invalid")
	ErrRoleRequired                = errors.New("role is required")
	ErrRoleInvalid                 = errors.New("role is invalid")
	ErrUserNotFound                = errors.New("user not found")
	ErrProfileNotFound             = errors.New("profile not found")
	ErrBranchNotFound              = errors.New("branch not found in organization")
	ErrProfileInactive             = errors.New("profile is inactive")
	ErrDefaultProfileRequired      = errors.New("default profile is required")
	ErrCannotDisableDefaultProfile = errors.New("cannot disable default profile")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListByUserID(ctx context.Context, rawUserID string) (ListProfilesResponse, error) {
	organizationID, err := getTenantOrganizationID(ctx)
	if err != nil {
		return ListProfilesResponse{}, err
	}

	userID, err := parseUUID(rawUserID, ErrUserIDInvalid)
	if err != nil {
		return ListProfilesResponse{}, err
	}

	items, err := s.repo.ListByUserID(ctx, organizationID, userID)
	if err != nil {
		return ListProfilesResponse{}, err
	}

	return ListProfilesResponse{
		Items: buildProfileResponses(items),
		Total: len(items),
	}, nil
}

func (s *Service) GetByID(ctx context.Context, rawProfileID string) (ProfileResponse, error) {
	organizationID, err := getTenantOrganizationID(ctx)
	if err != nil {
		return ProfileResponse{}, err
	}

	profileID, err := parseUUID(rawProfileID, ErrProfileIDInvalid)
	if err != nil {
		return ProfileResponse{}, err
	}

	item, err := s.repo.GetByID(ctx, organizationID, profileID)
	if err != nil {
		return ProfileResponse{}, err
	}

	return buildProfileResponse(item), nil
}

func (s *Service) Create(ctx context.Context, rawUserID string, req CreateProfileRequest) (ProfileResponse, error) {
	organizationID, err := getTenantOrganizationID(ctx)
	if err != nil {
		return ProfileResponse{}, err
	}

	userID, err := parseUUID(rawUserID, ErrUserIDInvalid)
	if err != nil {
		return ProfileResponse{}, err
	}

	item, err := buildProfileFromCreateRequest(req, organizationID, userID)
	if err != nil {
		return ProfileResponse{}, err
	}

	created, err := s.repo.Create(ctx, organizationID, userID, item)
	if err != nil {
		return ProfileResponse{}, err
	}

	return buildProfileResponse(created), nil
}

func (s *Service) Update(ctx context.Context, rawProfileID string, req UpdateProfileRequest) (ProfileResponse, error) {
	organizationID, err := getTenantOrganizationID(ctx)
	if err != nil {
		return ProfileResponse{}, err
	}

	profileID, err := parseUUID(rawProfileID, ErrProfileIDInvalid)
	if err != nil {
		return ProfileResponse{}, err
	}

	input, err := buildUpdateProfileInput(req)
	if err != nil {
		return ProfileResponse{}, err
	}

	updated, err := s.repo.Update(ctx, organizationID, profileID, input)
	if err != nil {
		return ProfileResponse{}, err
	}

	return buildProfileResponse(updated), nil
}

func (s *Service) Disable(ctx context.Context, rawProfileID string) error {
	organizationID, err := getTenantOrganizationID(ctx)
	if err != nil {
		return err
	}

	profileID, err := parseUUID(rawProfileID, ErrProfileIDInvalid)
	if err != nil {
		return err
	}

	return s.repo.Disable(ctx, organizationID, profileID)
}

func (s *Service) SetDefault(ctx context.Context, rawProfileID string) (ProfileResponse, error) {
	organizationID, err := getTenantOrganizationID(ctx)
	if err != nil {
		return ProfileResponse{}, err
	}

	profileID, err := parseUUID(rawProfileID, ErrProfileIDInvalid)
	if err != nil {
		return ProfileResponse{}, err
	}

	updated, err := s.repo.SetDefault(ctx, organizationID, profileID)
	if err != nil {
		return ProfileResponse{}, err
	}

	return buildProfileResponse(updated), nil
}

func (s *Service) AddRole(ctx context.Context, rawProfileID string, req AddRoleRequest) (ProfileResponse, error) {
	organizationID, err := getTenantOrganizationID(ctx)
	if err != nil {
		return ProfileResponse{}, err
	}

	profileID, err := parseUUID(rawProfileID, ErrProfileIDInvalid)
	if err != nil {
		return ProfileResponse{}, err
	}

	roleCode := normalizeRoleCode(req.RoleCode)
	if roleCode == "" {
		return ProfileResponse{}, ErrRoleRequired
	}

	updated, err := s.repo.AddRole(ctx, organizationID, profileID, roleCode)
	if err != nil {
		return ProfileResponse{}, err
	}

	return buildProfileResponse(updated), nil
}

func (s *Service) RemoveRole(ctx context.Context, rawProfileID string, rawRoleCode string) (ProfileResponse, error) {
	organizationID, err := getTenantOrganizationID(ctx)
	if err != nil {
		return ProfileResponse{}, err
	}

	profileID, err := parseUUID(rawProfileID, ErrProfileIDInvalid)
	if err != nil {
		return ProfileResponse{}, err
	}

	roleCode := normalizeRoleCode(rawRoleCode)
	if roleCode == "" {
		return ProfileResponse{}, ErrRoleRequired
	}

	updated, err := s.repo.RemoveRole(ctx, organizationID, profileID, roleCode)
	if err != nil {
		return ProfileResponse{}, err
	}

	return buildProfileResponse(updated), nil
}

func (s *Service) AddBranch(ctx context.Context, rawProfileID string, req AddBranchRequest) (ProfileResponse, error) {
	organizationID, err := getTenantOrganizationID(ctx)
	if err != nil {
		return ProfileResponse{}, err
	}

	profileID, err := parseUUID(rawProfileID, ErrProfileIDInvalid)
	if err != nil {
		return ProfileResponse{}, err
	}

	branchID, err := parseUUID(req.BranchID, ErrBranchIDInvalid)
	if err != nil {
		return ProfileResponse{}, err
	}

	updated, err := s.repo.AddBranch(ctx, organizationID, profileID, branchID)
	if err != nil {
		return ProfileResponse{}, err
	}

	return buildProfileResponse(updated), nil
}

func (s *Service) RemoveBranch(ctx context.Context, rawProfileID string, rawBranchID string) (ProfileResponse, error) {
	organizationID, err := getTenantOrganizationID(ctx)
	if err != nil {
		return ProfileResponse{}, err
	}

	profileID, err := parseUUID(rawProfileID, ErrProfileIDInvalid)
	if err != nil {
		return ProfileResponse{}, err
	}

	branchID, err := parseUUID(rawBranchID, ErrBranchIDInvalid)
	if err != nil {
		return ProfileResponse{}, err
	}

	updated, err := s.repo.RemoveBranch(ctx, organizationID, profileID, branchID)
	if err != nil {
		return ProfileResponse{}, err
	}

	return buildProfileResponse(updated), nil
}

func getTenantOrganizationID(ctx context.Context) (uuid.UUID, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return uuid.Nil, ErrTenantRequired
	}

	return *currentUser.OrganizationID, nil
}

func buildProfileFromCreateRequest(req CreateProfileRequest, organizationID uuid.UUID, userID uuid.UUID) (Profile, error) {
	roleCodes := normalizeRoleCodes(req.RoleCodes)
	if len(roleCodes) == 0 {
		return Profile{}, ErrRoleRequired
	}

	branchIDs, err := parseUUIDSlice(req.BranchIDs, ErrBranchIDInvalid)
	if err != nil {
		return Profile{}, err
	}

	branchID, err := parseOptionalUUID(req.BranchID, ErrBranchIDInvalid)
	if err != nil {
		return Profile{}, err
	}

	if branchID == nil && len(branchIDs) > 0 {
		firstBranchID := branchIDs[0]
		branchID = &firstBranchID
	}

	displayName := strings.TrimSpace(req.DisplayName)
	position := strings.TrimSpace(req.Position)
	profileType := normalizeProfileType(req.ProfileType)

	if position == "" {
		position = roleCodes[0]
	}

	if profileType == "" {
		profileType = strings.ToLower(roleCodes[0])
	}

	return Profile{
		ID:             uuid.New(),
		UserID:         userID,
		OrganizationID: organizationID,
		BranchID:       branchID,
		DisplayName:    displayName,
		Position:       position,
		ProfileType:    profileType,
		Status:         "active",
		IsDefault:      req.IsDefault,
		RoleCodes:      roleCodes,
		BranchIDs:      branchIDs,
	}, nil
}

func buildUpdateProfileInput(req UpdateProfileRequest) (UpdateProfileInput, error) {
	input := UpdateProfileInput{}

	if req.BranchID != nil {
		input.HasBranchID = true

		trimmedBranchID := strings.TrimSpace(*req.BranchID)
		if trimmedBranchID == "" {
			input.BranchID = nil
		} else {
			branchID, err := uuid.Parse(trimmedBranchID)
			if err != nil {
				return UpdateProfileInput{}, ErrBranchIDInvalid
			}

			input.BranchID = &branchID
		}
	}

	if req.DisplayName != nil {
		value := strings.TrimSpace(*req.DisplayName)
		input.DisplayName = &value
	}

	if req.Position != nil {
		value := strings.TrimSpace(*req.Position)
		input.Position = &value
	}

	if req.ProfileType != nil {
		value := normalizeProfileType(*req.ProfileType)
		input.ProfileType = &value
	}

	if req.Status != nil {
		value := normalizeStatus(*req.Status)
		input.Status = &value
	}

	if req.IsDefault != nil {
		input.IsDefault = req.IsDefault
	}

	return input, nil
}

func normalizeRoleCodes(rawRoleCodes []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(rawRoleCodes))

	for _, rawRoleCode := range rawRoleCodes {
		roleCode := normalizeRoleCode(rawRoleCode)
		if roleCode == "" {
			continue
		}

		if _, exists := seen[roleCode]; exists {
			continue
		}

		seen[roleCode] = struct{}{}
		result = append(result, roleCode)
	}

	return result
}

func normalizeRoleCode(roleCode string) string {
	roleCode = strings.TrimSpace(roleCode)
	roleCode = strings.ReplaceAll(roleCode, "-", "_")
	roleCode = strings.ReplaceAll(roleCode, " ", "_")
	roleCode = strings.ToUpper(roleCode)

	return roleCode
}

func normalizeProfileType(profileType string) string {
	profileType = strings.TrimSpace(profileType)
	profileType = strings.ReplaceAll(profileType, "-", "_")
	profileType = strings.ReplaceAll(profileType, " ", "_")
	profileType = strings.ToLower(profileType)

	return profileType
}

func normalizeStatus(status string) string {
	status = strings.TrimSpace(status)
	status = strings.ToLower(status)

	if status == "" {
		return "active"
	}

	return status
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

func parseOptionalUUID(rawID string, invalidErr error) (*uuid.UUID, error) {
	rawID = strings.TrimSpace(rawID)
	if rawID == "" {
		return nil, nil
	}

	id, err := uuid.Parse(rawID)
	if err != nil {
		return nil, invalidErr
	}

	return &id, nil
}

func parseUUIDSlice(rawIDs []string, invalidErr error) ([]uuid.UUID, error) {
	items := make([]uuid.UUID, 0, len(rawIDs))
	seen := make(map[uuid.UUID]struct{})

	for _, rawID := range rawIDs {
		rawID = strings.TrimSpace(rawID)
		if rawID == "" {
			continue
		}

		id, err := uuid.Parse(rawID)
		if err != nil {
			return nil, invalidErr
		}

		if _, exists := seen[id]; exists {
			continue
		}

		seen[id] = struct{}{}
		items = append(items, id)
	}

	return items, nil
}

func buildProfileResponses(items []Profile) []ProfileResponse {
	responses := make([]ProfileResponse, 0, len(items))

	for _, item := range items {
		responses = append(responses, buildProfileResponse(item))
	}

	return responses
}

func buildProfileResponse(item Profile) ProfileResponse {
	return ProfileResponse{
		ID:             item.ID.String(),
		UserID:         item.UserID.String(),
		OrganizationID: item.OrganizationID.String(),
		BranchID:       uuidPointerToString(item.BranchID),
		DisplayName:    item.DisplayName,
		Position:       item.Position,
		ProfileType:    item.ProfileType,
		Status:         item.Status,
		IsDefault:      item.IsDefault,
		RoleCodes:      item.RoleCodes,
		BranchIDs:      uuidSliceToStringSlice(item.BranchIDs),
	}
}

func uuidPointerToString(id *uuid.UUID) string {
	if id == nil {
		return ""
	}

	return id.String()
}

func uuidSliceToStringSlice(ids []uuid.UUID) []string {
	result := make([]string, 0, len(ids))

	for _, id := range ids {
		result = append(result, id.String())
	}

	return result
}
