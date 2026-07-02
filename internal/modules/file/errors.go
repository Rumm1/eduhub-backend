package file

import "errors"

var (
	ErrTenantRequired = errors.New("tenant organization is required")
	ErrFileRequired   = errors.New("file is required")
	ErrFileTooLarge   = errors.New("file is too large")
	ErrFileIDInvalid  = errors.New("file id is invalid")
	ErrFileNotFound   = errors.New("file not found")
)
