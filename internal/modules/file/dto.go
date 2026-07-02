package file

type FileResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id,omitempty"`
	UploadedBy     string `json:"uploaded_by,omitempty"`
	Folder         string `json:"folder,omitempty"`
	FileName       string `json:"file_name"`
	FilePath       string `json:"file_path"`
	FileURL        string `json:"file_url"`
	MimeType       string `json:"mime_type,omitempty"`
	SizeBytes      int64  `json:"size_bytes"`
	CreatedAt      string `json:"created_at"`
}

type ListFilesResponse struct {
	Items []FileResponse `json:"items"`
	Total int            `json:"total"`
}

type UploadInput struct {
	Folder    string
	FileName  string
	MimeType  string
	SizeBytes int64
	Reader    interface{}
}
