package importer

import (
	"errors"
	"net/http"

	"github.com/Rumm1/eduhub-backend/internal/shared/response"
)

const maxImportFileBytes = 10 * 1024 * 1024

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) PreviewStudents(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxImportFileBytes+1024)

	if err := r.ParseMultipartForm(maxImportFileBytes); err != nil {
		response.Error(w, http.StatusBadRequest, "IMPORT_FILE_INVALID", "Invalid multipart form or file is too large")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "FILE_REQUIRED", "File is required")
		return
	}
	defer file.Close()

	result, err := h.service.PreviewStudents(r.Context(), header.Filename, file)
	if err != nil {
		writeImportError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) ConfirmStudents(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxImportFileBytes+1024)

	if err := r.ParseMultipartForm(maxImportFileBytes); err != nil {
		response.Error(w, http.StatusBadRequest, "IMPORT_FILE_INVALID", "Invalid multipart form or file is too large")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "FILE_REQUIRED", "File is required")
		return
	}
	defer file.Close()

	result, err := h.service.ConfirmStudents(r.Context(), header.Filename, file)
	if err != nil {
		writeImportError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func writeImportError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTenantRequired):
		response.Error(w, http.StatusForbidden, "TENANT_REQUIRED", "Tenant organization is required")
	case errors.Is(err, ErrFileRequired):
		response.Error(w, http.StatusBadRequest, "FILE_REQUIRED", "File is required")
	case errors.Is(err, ErrFileTypeUnsupported):
		response.Error(w, http.StatusBadRequest, "FILE_TYPE_UNSUPPORTED", "Only .xlsx and .csv files are supported")
	case errors.Is(err, ErrEmptyImportFile):
		response.Error(w, http.StatusBadRequest, "EMPTY_IMPORT_FILE", "Import file is empty")
	default:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
}
