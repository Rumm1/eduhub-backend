package platformuser

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) ResetPassword(ctx context.Context, userID uuid.UUID) (ResetPasswordResponse, error) {
	user, err := s.repository.GetUserByID(ctx, userID)
	if err != nil {
		return ResetPasswordResponse{}, err
	}

	temporaryPassword, err := GenerateTemporaryPassword()
	if err != nil {
		return ResetPasswordResponse{}, err
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(temporaryPassword), bcrypt.DefaultCost)
	if err != nil {
		return ResetPasswordResponse{}, err
	}

	if err := s.repository.UpdatePassword(ctx, user.ID, string(hashBytes)); err != nil {
		return ResetPasswordResponse{}, err
	}

	return ResetPasswordResponse{
		UserID:             user.ID.String(),
		Login:              user.Email,
		TemporaryPassword:  temporaryPassword,
		MustChangePassword: true,
	}, nil
}
