package ai

import "errors"

var (
	ErrTenantRequired   = errors.New("tenant organization is required")
	ErrChatMessageEmpty = errors.New("chat message is required")
)
