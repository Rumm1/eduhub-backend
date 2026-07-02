package parent

import "errors"

var (
	ErrTenantRequired   = errors.New("tenant organization is required")
	ErrParentIDInvalid  = errors.New("parent id is invalid")
	ErrStudentIDInvalid = errors.New("student id is invalid")
	ErrParentNotFound   = errors.New("parent not found")
	ErrStudentNotFound  = errors.New("student not found")
	ErrFullNameRequired = errors.New("full name is required")
)
