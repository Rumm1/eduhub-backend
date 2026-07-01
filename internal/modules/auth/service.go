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

	if !password.Compare(accessData.User.PasswordHash, req.Password) {
		return LoginResponse{}, ErrInvalidCredentials
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(platformjwt.AccessTokenPayload{
		UserID:         accessData.User.ID,
		OrganizationID: accessData.User.OrganizationID,
		Roles:          accessData.Roles,
		Permissions:    accessData.Permissions,
		BranchIDs:      accessData.BranchIDs,
	})
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
		return UserResponse{}, errors.New("user context not found")
	}

	accessData, err := s.repo.GetUserAccessDataByID(ctx, currentUser.UserID)
	if err != nil {
		return UserResponse{}, err
	}

	return buildUserResponse(accessData), nil
}

func buildUserResponse(accessData UserAccessData) UserResponse {
	return UserResponse{
		ID:             accessData.User.ID.String(),
		OrganizationID: uuidToStringPointer(accessData.User.OrganizationID),
		Email:          accessData.User.Email,
		FullName:       accessData.User.FullName,
		Roles:          accessData.Roles,
		Permissions:    accessData.Permissions,
		BranchIDs:      uuidSliceToStringSlice(accessData.BranchIDs),
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
