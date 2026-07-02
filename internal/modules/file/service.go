package file

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/google/uuid"
)

const (
	uploadRoot     = "uploads"
	defaultFolder  = "general"
	maxUploadBytes = 20 * 1024 * 1024
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context) ([]FileResponse, error) {
	currentUser, err := getCurrentFileUser(ctx)
	if err != nil {
		return nil, err
	}

	items, err := s.repository.List(ctx, *currentUser.OrganizationID)
	if err != nil {
		return nil, err
	}

	return mapFiles(items), nil
}

func (s *Service) GetByID(ctx context.Context, fileIDRaw string) (FileResponse, error) {
	currentUser, err := getCurrentFileUser(ctx)
	if err != nil {
		return FileResponse{}, err
	}

	fileID, err := uuid.Parse(fileIDRaw)
	if err != nil {
		return FileResponse{}, ErrFileIDInvalid
	}

	item, err := s.repository.GetByID(ctx, *currentUser.OrganizationID, fileID)
	if err != nil {
		return FileResponse{}, err
	}

	return mapFile(item), nil
}

func (s *Service) Upload(
	ctx context.Context,
	folder string,
	fileName string,
	mimeType string,
	sizeBytes int64,
	reader io.Reader,
) (FileResponse, error) {
	currentUser, err := getCurrentFileUser(ctx)
	if err != nil {
		return FileResponse{}, err
	}

	if reader == nil || strings.TrimSpace(fileName) == "" {
		return FileResponse{}, ErrFileRequired
	}

	if sizeBytes > maxUploadBytes {
		return FileResponse{}, ErrFileTooLarge
	}

	safeFolder := sanitizeFolder(folder)
	safeFileName := sanitizeFileName(fileName)

	fileID := uuid.New()
	storageFileName := fileID.String() + "_" + safeFileName

	relativePath := filepath.Join(uploadRoot, safeFolder, storageFileName)
	relativePath = filepath.ToSlash(relativePath)

	fullDir := filepath.Join(uploadRoot, safeFolder)

	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return FileResponse{}, err
	}

	destination, err := os.Create(filepath.FromSlash(relativePath))
	if err != nil {
		return FileResponse{}, err
	}
	defer destination.Close()

	writtenBytes, err := io.Copy(destination, reader)
	if err != nil {
		return FileResponse{}, err
	}

	if sizeBytes <= 0 {
		sizeBytes = writtenBytes
	}

	item := File{
		ID:             fileID,
		OrganizationID: *currentUser.OrganizationID,
		UploadedBy:     currentUser.UserID,
		Folder:         safeFolder,
		FileName:       fileName,
		FilePath:       relativePath,
		MimeType:       strings.TrimSpace(mimeType),
		SizeBytes:      sizeBytes,
	}

	result, err := s.repository.Create(ctx, item)
	if err != nil {
		_ = os.Remove(filepath.FromSlash(relativePath))
		return FileResponse{}, err
	}

	return mapFile(result), nil
}

func (s *Service) Delete(ctx context.Context, fileIDRaw string) error {
	currentUser, err := getCurrentFileUser(ctx)
	if err != nil {
		return err
	}

	fileID, err := uuid.Parse(fileIDRaw)
	if err != nil {
		return ErrFileIDInvalid
	}

	filePath, err := s.repository.Delete(ctx, *currentUser.OrganizationID, fileID)
	if err != nil {
		return err
	}

	if strings.TrimSpace(filePath) != "" {
		_ = os.Remove(filepath.FromSlash(filePath))
	}

	return nil
}

func getCurrentFileUser(ctx context.Context) (usercontext.UserContext, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return usercontext.UserContext{}, ErrTenantRequired
	}

	return currentUser, nil
}

func sanitizeFolder(folder string) string {
	folder = strings.ToLower(strings.TrimSpace(folder))
	if folder == "" {
		return defaultFolder
	}

	re := regexp.MustCompile(`[^a-zA-Z0-9_-]+`)
	folder = re.ReplaceAllString(folder, "-")
	folder = strings.Trim(folder, "-_")

	if folder == "" {
		return defaultFolder
	}

	return folder
}

func sanitizeFileName(fileName string) string {
	fileName = filepath.Base(strings.TrimSpace(fileName))
	if fileName == "" || fileName == "." {
		return "file"
	}

	re := regexp.MustCompile(`[^a-zA-Z0-9._-]+`)
	fileName = re.ReplaceAllString(fileName, "-")
	fileName = strings.Trim(fileName, "-_.")

	if fileName == "" {
		return "file"
	}

	return fileName
}

func mapFiles(items []File) []FileResponse {
	result := make([]FileResponse, 0, len(items))

	for _, item := range items {
		result = append(result, mapFile(item))
	}

	return result
}

func mapFile(item File) FileResponse {
	organizationID := ""
	if item.OrganizationID != uuid.Nil {
		organizationID = item.OrganizationID.String()
	}

	uploadedBy := ""
	if item.UploadedBy != uuid.Nil {
		uploadedBy = item.UploadedBy.String()
	}

	fileURL := "/" + strings.TrimLeft(item.FilePath, "/")

	return FileResponse{
		ID:             item.ID.String(),
		OrganizationID: organizationID,
		UploadedBy:     uploadedBy,
		Folder:         item.Folder,
		FileName:       item.FileName,
		FilePath:       item.FilePath,
		FileURL:        fileURL,
		MimeType:       item.MimeType,
		SizeBytes:      item.SizeBytes,
		CreatedAt:      item.CreatedAt,
	}
}
