package notification

import "errors"

var (
	ErrTenantRequired        = errors.New("tenant organization is required")
	ErrNotificationIDInvalid = errors.New("notification id is invalid")
	ErrNotificationNotFound  = errors.New("notification not found")
	ErrTitleRequired         = errors.New("title is required")
	ErrTypeInvalid           = errors.New("notification type is invalid")
	ErrTargetUserInvalid     = errors.New("target user id is invalid")
	ErrTargetUserNotFound    = errors.New("target user not found in organization")
)
