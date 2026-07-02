package branding

import (
	"context"
	"strings"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetCurrentBranding(ctx context.Context) (BrandingResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok {
		return BrandingResponse{}, ErrTenantRequired
	}

	data, err := s.repository.GetCurrentBranding(ctx, currentUser.UserID, currentUser.OrganizationID)
	if err != nil {
		return BrandingResponse{}, err
	}

	return mapBrandingResponse(data), nil
}

func (s *Service) UpdateMyAvatar(ctx context.Context, request UpdateAvatarRequest) (BrandingResponse, error) {
	avatarPath := strings.TrimSpace(request.AvatarPath)
	if avatarPath == "" {
		return BrandingResponse{}, ErrAvatarRequired
	}

	currentUser, ok := usercontext.GetUser(ctx)
	if !ok {
		return BrandingResponse{}, ErrTenantRequired
	}

	if err := s.repository.UpdateUserAvatar(ctx, currentUser.UserID, avatarPath); err != nil {
		return BrandingResponse{}, err
	}

	return s.GetCurrentBranding(ctx)
}

func (s *Service) ClearMyAvatar(ctx context.Context) (BrandingResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok {
		return BrandingResponse{}, ErrTenantRequired
	}

	if err := s.repository.ClearUserAvatar(ctx, currentUser.UserID); err != nil {
		return BrandingResponse{}, err
	}

	return s.GetCurrentBranding(ctx)
}

func (s *Service) UpdateOrganizationLogo(ctx context.Context, request UpdateLogoRequest) (BrandingResponse, error) {
	logoPath := strings.TrimSpace(request.LogoPath)
	if logoPath == "" {
		return BrandingResponse{}, ErrLogoRequired
	}

	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return BrandingResponse{}, ErrTenantRequired
	}

	if err := s.repository.UpdateOrganizationLogo(ctx, *currentUser.OrganizationID, logoPath); err != nil {
		return BrandingResponse{}, err
	}

	return s.GetCurrentBranding(ctx)
}

func (s *Service) ClearOrganizationLogo(ctx context.Context) (BrandingResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return BrandingResponse{}, ErrTenantRequired
	}

	if err := s.repository.ClearOrganizationLogo(ctx, *currentUser.OrganizationID); err != nil {
		return BrandingResponse{}, err
	}

	return s.GetCurrentBranding(ctx)
}

func mapBrandingResponse(data BrandingData) BrandingResponse {
	response := BrandingResponse{
		UserAvatarPath:       data.UserAvatarPath,
		UserAvatarURL:        toPublicURL(data.UserAvatarPath),
		OrganizationLogoPath: data.OrganizationLogoPath,
		OrganizationLogoURL:  toPublicURL(data.OrganizationLogoPath),
		DefaultLogo:          data.OrganizationLogoPath == "",
		DefaultName:          "EduHub",
	}

	return response
}

func toPublicURL(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}

	if strings.HasPrefix(path, "/") {
		return path
	}

	return "/" + path
}
