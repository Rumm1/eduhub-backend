package auth

import (
	"context"
	"errors"
	"strings"

	platformjwt "github.com/Rumm1/eduhub-backend/internal/platform/jwt"
	"github.com/Rumm1/eduhub-backend/internal/platform/password"
	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserInactive       = errors.New("user is inactive")
	ErrProfileInactive    = errors.New("profile is inactive")
	ErrProfileIDInvalid   = errors.New("profile id is invalid")
	ErrUserContextMissing = errors.New("user context missing")
)

type Service struct {
	repo       *Repository
	jwtManager *platformjwt.Manager
}

func NewService(repo *Repository, jwtManager *platformjwt.Manager) *Service {
	return &Service{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || req.Password == "" {
		return LoginResponse{}, ErrInvalidCredentials
	}

	accessData, err := s.repo.GetUserAccessData(ctx, email)
	if err != nil {
		return LoginResponse{}, ErrInvalidCredentials
	}

	if accessData.User.Status != "active" {
		return LoginResponse{}, ErrUserInactive
	}

	if accessData.Profile.Status != "active" {
		return LoginResponse{}, ErrProfileInactive
	}

	if !password.Compare(accessData.User.PasswordHash, req.Password) {
		return LoginResponse{}, ErrInvalidCredentials
	}

	accessToken, err := s.buildAccessToken(accessData)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		User:        buildUserResponse(accessData),
	}, nil
}

func (s *Service) Me(ctx context.Context) (UserResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok {
		return UserResponse{}, ErrUserContextMissing
	}

	var accessData UserAccessData
	var err error

	if currentUser.ProfileID != nil {
		accessData, err = s.repo.GetUserAccessDataByProfileID(ctx, currentUser.UserID, *currentUser.ProfileID)
	} else {
		accessData, err = s.repo.GetUserAccessDataByID(ctx, currentUser.UserID)
	}

	if err != nil {
		return UserResponse{}, err
	}

	return buildUserResponse(accessData), nil
}

func (s *Service) SwitchProfile(ctx context.Context, req SwitchProfileRequest) (LoginResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok {
		return LoginResponse{}, ErrUserContextMissing
	}

	profileID, err := uuid.Parse(strings.TrimSpace(req.ProfileID))
	if err != nil {
		return LoginResponse{}, ErrProfileIDInvalid
	}

	accessData, err := s.repo.GetUserAccessDataByProfileID(ctx, currentUser.UserID, profileID)
	if err != nil {
		return LoginResponse{}, ErrProfileIDInvalid
	}

	if accessData.User.Status != "active" {
		return LoginResponse{}, ErrUserInactive
	}

	if accessData.Profile.Status != "active" {
		return LoginResponse{}, ErrProfileInactive
	}

	accessToken, err := s.buildAccessToken(accessData)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		User:        buildUserResponse(accessData),
	}, nil
}

func (s *Service) buildAccessToken(accessData UserAccessData) (string, error) {
	profileID := accessData.Profile.ID

	return s.jwtManager.GenerateAccessToken(platformjwt.AccessTokenPayload{
		UserID:         accessData.User.ID,
		ProfileID:      &profileID,
		OrganizationID: accessData.Profile.OrganizationID,
		Roles:          accessData.Roles,
		Permissions:    accessData.Permissions,
		BranchIDs:      accessData.BranchIDs,
	})
}

func buildUserResponse(accessData UserAccessData) UserResponse {
	profileID := accessData.Profile.ID.String()

	return UserResponse{
		ID:                accessData.User.ID.String(),
		ProfileID:         &profileID,
		OrganizationID:    uuidToStringPointer(accessData.Profile.OrganizationID),
		Email:             accessData.User.Email,
		FullName:          accessData.User.FullName,
		Roles:             accessData.Roles,
		Permissions:       accessData.Permissions,
		BranchIDs:         uuidSliceToStringSlice(accessData.BranchIDs),
		CurrentProfile:    buildProfileResponsePointer(accessData.Profile),
		AvailableProfiles: buildProfileResponses(accessData.AvailableProfiles),
	}
}

func buildProfileResponses(profiles []UserProfile) []ProfileResponse {
	result := make([]ProfileResponse, 0, len(profiles))

	for _, profile := range profiles {
		result = append(result, buildProfileResponse(profile))
	}

	return result
}

func buildProfileResponsePointer(profile UserProfile) *ProfileResponse {
	response := buildProfileResponse(profile)
	return &response
}

func buildProfileResponse(profile UserProfile) ProfileResponse {
	return ProfileResponse{
		ID:             profile.ID.String(),
		OrganizationID: uuidToStringPointer(profile.OrganizationID),
		BranchID:       uuidToStringPointer(profile.BranchID),
		DisplayName:    profile.DisplayName,
		Position:       profile.Position,
		ProfileType:    profile.ProfileType,
		Status:         profile.Status,
		IsDefault:      profile.IsDefault,
		Roles:          profile.Roles,
		BranchIDs:      uuidSliceToStringSlice(profile.BranchIDs),
	}
}

func uuidToStringPointer(id *uuid.UUID) *string {
	if id == nil {
		return nil
	}

	value := id.String()
	return &value
}

func uuidSliceToStringSlice(ids []uuid.UUID) []string {
	result := make([]string, 0, len(ids))

	for _, id := range ids {
		result = append(result, id.String())
	}

	return result
}
