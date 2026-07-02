package file

import "github.com/google/uuid"

type File struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	UploadedBy     uuid.UUID
	Folder         string
	FileName       string
	FilePath       string
	MimeType       string
	SizeBytes      int64
	CreatedAt      string
}
