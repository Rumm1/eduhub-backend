package importer

import "errors"

var (
	ErrTenantRequired      = errors.New("tenant organization is required")
	ErrFileRequired        = errors.New("file is required")
	ErrFileTypeUnsupported = errors.New("file type is unsupported")
	ErrEmptyImportFile     = errors.New("import file is empty")
	ErrGroupNotFound       = errors.New("group not found")
)
