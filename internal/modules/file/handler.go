package file

import (
	"errors"
	"net/http"

	"github.com/Rumm1/eduhub-backend/internal/shared/response"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.List(r.Context())
	if err != nil {
		writeFileError(w, err)
		return
	}

	response.Success(w, http.StatusOK, ListFilesResponse{
		Items: items,
		Total: len(items),
	})
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	fileID := chi.URLParam(r, "fileID")

	result, err := h.service.GetByID(r.Context(), fileID)
	if err != nil {
		writeFileError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes+1024)

	if err := r.ParseMultipartForm(maxUploadBytes); err != nil {
		response.Error(w, http.StatusBadRequest, "FILE_UPLOAD_INVALID", "Invalid multipart form or file is too large")
		return
	}

	uploadedFile, header, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "FILE_REQUIRED", "File is required")
		return
	}
	defer uploadedFile.Close()

	result, err := h.service.Upload(
		r.Context(),
		r.FormValue("folder"),
		header.Filename,
		header.Header.Get("Content-Type"),
		header.Size,
		uploadedFile,
	)
	if err != nil {
		writeFileError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, result)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	fileID := chi.URLParam(r, "fileID")

	if err := h.service.Delete(r.Context(), fileID); err != nil {
		writeFileError(w, err)
		return
	}

	response.Message(w, http.StatusOK, "File deleted")
}

func writeFileError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTenantRequired):
		response.Error(w, http.StatusForbidden, "TENANT_REQUIRED", "Tenant organization is required")
	case errors.Is(err, ErrFileRequired):
		response.Error(w, http.StatusBadRequest, "FILE_REQUIRED", "File is required")
	case errors.Is(err, ErrFileTooLarge):
		response.Error(w, http.StatusBadRequest, "FILE_TOO_LARGE", "File is too large")
	case errors.Is(err, ErrFileIDInvalid):
		response.Error(w, http.StatusBadRequest, "FILE_ID_INVALID", "File id is invalid")
	case errors.Is(err, ErrFileNotFound):
		response.Error(w, http.StatusNotFound, "FILE_NOT_FOUND", "File not found")
	default:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	}
}
