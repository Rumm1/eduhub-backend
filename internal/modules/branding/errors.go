package branding

import "errors"

var (
	ErrTenantRequired = errors.New("tenant organization is required")
	ErrAvatarRequired = errors.New("avatar path is required")
	ErrLogoRequired   = errors.New("logo path is required")
)
