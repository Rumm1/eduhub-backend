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
	ErrTenantRequired          = errors.New("tenant organization is required")
	ErrEmailRequired           = errors.New("email is required")
	ErrPasswordRequired        = errors.New("password is required")
	ErrFullNameRequired        = errors.New("full name is required")
	ErrProfilesRequired        = errors.New("profiles are required")
	ErrDefaultProfileRequired  = errors.New("default profile is required")
	ErrDefaultProfileDuplicate = errors.New("only one default profile is allowed")
	ErrRoleRequired            = errors.New("role is required")
	ErrRoleInvalid             = errors.New("role is invalid")
	ErrBranchIDInvalid         = errors.New("branch id is invalid")
	ErrBranchNotFound          = errors.New("branch not found in organization")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateUserRequest) (CreateUserResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return CreateUserResponse{}, ErrTenantRequired
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" {
		return CreateUserResponse{}, ErrEmailRequired
	}

	plainPassword := strings.TrimSpace(req.Password)
	if plainPassword == "" {
		generatedPassword, err := generateTemporaryPassword()
		if err != nil {
			return CreateUserResponse{}, err
		}

		plainPassword = generatedPassword
	}

	fullName := strings.TrimSpace(req.FullName)
	if fullName == "" {
		return CreateUserResponse{}, ErrFullNameRequired
	}

	profiles, err := buildProfilesFromRequest(req.Profiles, *currentUser.OrganizationID, fullName)
	if err != nil {
		return CreateUserResponse{}, err
	}

	hashedPassword, err := password.Hash(plainPassword)
	if err != nil {
		return CreateUserResponse{}, err
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

	for index := range profiles {
		profiles[index].UserID = newUser.ID
	}

	createdUser, err := s.repo.CreateUserWithProfiles(ctx, newUser, profiles)
	if err != nil {
		return CreateUserResponse{}, err
	}

	roles, savedBranchIDs, savedProfiles, err := s.repo.BuildUserResponseData(ctx, createdUser.ID)
	if err != nil {
		return CreateUserResponse{}, err
	}

	return CreateUserResponse{
		User: buildUserResponse(createdUser, roles, savedBranchIDs, savedProfiles),
		TemporaryCredentials: TemporaryCredentialsResponse{
			Login:              createdUser.Email,
			Password:           plainPassword,
			MustChangePassword: true,
		},
	}, nil
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
		roles, branchIDs, profiles, err := s.repo.BuildUserResponseData(ctx, item.ID)
		if err != nil {
			return ListUsersResponse{}, err
		}

		items = append(items, buildUserResponse(item, roles, branchIDs, profiles))
	}

	return ListUsersResponse{
		Items: items,
		Total: len(items),
	}, nil
}

func buildProfilesFromRequest(
	requestProfiles []CreateUserProfileRequest,
	organizationID uuid.UUID,
	fullName string,
) ([]UserProfile, error) {
	if len(requestProfiles) == 0 {
		return nil, ErrProfilesRequired
	}

	profiles := make([]UserProfile, 0, len(requestProfiles))
	defaultCount := 0

	for index, requestProfile := range requestProfiles {
		isDefault := requestProfile.IsDefault

		if len(requestProfiles) == 1 {
			isDefault = true
		}

		if isDefault {
			defaultCount++
		}

		roleCodes := normalizeRoleCodes(requestProfile.RoleCodes)
		if len(roleCodes) == 0 {
			return nil, ErrRoleRequired
		}

		branchIDs, err := parseBranchIDs(requestProfile.BranchIDs)
		if err != nil {
			return nil, err
		}

		branchID, err := parseOptionalBranchID(requestProfile.BranchID)
		if err != nil {
			return nil, err
		}

		if branchID == nil && len(branchIDs) > 0 {
			firstBranchID := branchIDs[0]
			branchID = &firstBranchID
		}

		displayName := strings.TrimSpace(requestProfile.DisplayName)
		if displayName == "" {
			displayName = fullName
		}

		position := strings.TrimSpace(requestProfile.Position)
		if position == "" {
			position = roleCodes[0]
		}

		profileType := normalizeProfileType(requestProfile.ProfileType)
		if profileType == "" {
			profileType = strings.ToLower(roleCodes[0])
		}

		profiles = append(profiles, UserProfile{
			ID:             uuid.New(),
			OrganizationID: organizationID,
			BranchID:       branchID,
			DisplayName:    displayName,
			Position:       position,
			ProfileType:    profileType,
			Status:         "active",
			IsDefault:      isDefault,
			RoleCodes:      roleCodes,
			BranchIDs:      branchIDs,
		})

		_ = index
	}

	if defaultCount == 0 {
		return nil, ErrDefaultProfileRequired
	}

	if defaultCount > 1 {
		return nil, ErrDefaultProfileDuplicate
	}

	return profiles, nil
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

func parseOptionalBranchID(rawBranchID string) (*uuid.UUID, error) {
	rawBranchID = strings.TrimSpace(rawBranchID)
	if rawBranchID == "" {
		return nil, nil
	}

	branchID, err := uuid.Parse(rawBranchID)
	if err != nil {
		return nil, ErrBranchIDInvalid
	}

	return &branchID, nil
}

func parseBranchIDs(rawBranchIDs []string) ([]uuid.UUID, error) {
	branchIDs := make([]uuid.UUID, 0, len(rawBranchIDs))
	seen := make(map[uuid.UUID]struct{})

	for _, rawID := range rawBranchIDs {
		rawID = strings.TrimSpace(rawID)
		if rawID == "" {
			continue
		}

		branchID, err := uuid.Parse(rawID)
		if err != nil {
			return nil, ErrBranchIDInvalid
		}

		if _, exists := seen[branchID]; exists {
			continue
		}

		seen[branchID] = struct{}{}
		branchIDs = append(branchIDs, branchID)
	}

	return branchIDs, nil
}

func buildUserResponse(item User, roles []string, branchIDs []uuid.UUID, profiles []UserProfile) UserResponse {
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
		Profiles:       buildProfileResponses(profiles),
	}
}

func buildProfileResponses(profiles []UserProfile) []UserProfileResponse {
	result := make([]UserProfileResponse, 0, len(profiles))

	for _, profile := range profiles {
		result = append(result, UserProfileResponse{
			ID:             profile.ID.String(),
			OrganizationID: profile.OrganizationID.String(),
			BranchID:       uuidPointerToString(profile.BranchID),
			DisplayName:    profile.DisplayName,
			Position:       profile.Position,
			ProfileType:    profile.ProfileType,
			Status:         profile.Status,
			IsDefault:      profile.IsDefault,
			RoleCodes:      profile.RoleCodes,
			BranchIDs:      uuidSliceToStringSlice(profile.BranchIDs),
		})
	}

	return result
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
