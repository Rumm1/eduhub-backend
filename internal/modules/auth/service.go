package auth

import (
	"context"
	"errors"
	"strings"

	platformjwt "github.com/Rumm1/eduhub-backend/internal/platform/jwt"
	"github.com/Rumm1/eduhub-backend/internal/platform/password"
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

	var organizationID *string
	if accessData.User.OrganizationID != nil {
		value := accessData.User.OrganizationID.String()
		organizationID = &value
	}

	branchIDs := make([]string, 0, len(accessData.BranchIDs))
	for _, id := range accessData.BranchIDs {
		branchIDs = append(branchIDs, id.String())
	}

	return LoginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		User: UserResponse{
			ID:             accessData.User.ID.String(),
			OrganizationID: organizationID,
			Email:          accessData.User.Email,
			FullName:       accessData.User.FullName,
			Roles:          accessData.Roles,
			Permissions:    accessData.Permissions,
			BranchIDs:      branchIDs,
		},
	}, nil
}
