package branding

type UpdateAvatarRequest struct {
	AvatarPath string `json:"avatar_path"`
}

type UpdateLogoRequest struct {
	LogoPath string `json:"logo_path"`
}

type BrandingResponse struct {
	UserAvatarPath string `json:"user_avatar_path"`
	UserAvatarURL  string `json:"user_avatar_url"`

	OrganizationLogoPath string `json:"organization_logo_path"`
	OrganizationLogoURL  string `json:"organization_logo_url"`

	DefaultLogo bool   `json:"default_logo"`
	DefaultName string `json:"default_name"`
}
