package parent

import (
	"encoding/json"
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
		writeParentError(w, err)
		return
	}

	response.Success(w, http.StatusOK, ListParentsResponse{
		Items: items,
		Total: len(items),
	})
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	parentID := chi.URLParam(r, "parentID")

	result, err := h.service.GetByID(r.Context(), parentID)
	if err != nil {
		writeParentError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateParentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.Create(r.Context(), req)
	if err != nil {
		writeParentError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, result)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	parentID := chi.URLParam(r, "parentID")

	var req UpdateParentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.Update(r.Context(), parentID, req)
	if err != nil {
		writeParentError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	parentID := chi.URLParam(r, "parentID")

	if err := h.service.Delete(r.Context(), parentID); err != nil {
		writeParentError(w, err)
		return
	}

	response.Message(w, http.StatusOK, "Parent deleted")
}

func (h *Handler) AttachStudent(w http.ResponseWriter, r *http.Request) {
	parentID := chi.URLParam(r, "parentID")
	studentID := chi.URLParam(r, "studentID")

	var req AttachStudentRequest

	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}

	if err := h.service.AttachStudent(r.Context(), parentID, studentID, req); err != nil {
		writeParentError(w, err)
		return
	}

	response.Message(w, http.StatusOK, "Student attached to parent")
}

func (h *Handler) DetachStudent(w http.ResponseWriter, r *http.Request) {
	parentID := chi.URLParam(r, "parentID")
	studentID := chi.URLParam(r, "studentID")

	if err := h.service.DetachStudent(r.Context(), parentID, studentID); err != nil {
		writeParentError(w, err)
		return
	}

	response.Message(w, http.StatusOK, "Student detached from parent")
}

func (h *Handler) ListStudents(w http.ResponseWriter, r *http.Request) {
	parentID := chi.URLParam(r, "parentID")

	items, err := h.service.ListStudents(r.Context(), parentID)
	if err != nil {
		writeParentError(w, err)
		return
	}

	response.Success(w, http.StatusOK, ListParentStudentsResponse{
		Items: items,
		Total: len(items),
	})
}

func writeParentError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTenantRequired):
		response.Error(w, http.StatusForbidden, "TENANT_REQUIRED", "Tenant organization is required")
	case errors.Is(err, ErrParentIDInvalid):
		response.Error(w, http.StatusBadRequest, "PARENT_ID_INVALID", "Parent id is invalid")
	case errors.Is(err, ErrStudentIDInvalid):
		response.Error(w, http.StatusBadRequest, "STUDENT_ID_INVALID", "Student id is invalid")
	case errors.Is(err, ErrFullNameRequired):
		response.Error(w, http.StatusBadRequest, "FULL_NAME_REQUIRED", "Full name is required")
	case errors.Is(err, ErrParentNotFound):
		response.Error(w, http.StatusNotFound, "PARENT_NOT_FOUND", "Parent not found")
	case errors.Is(err, ErrStudentNotFound):
		response.Error(w, http.StatusNotFound, "STUDENT_NOT_FOUND", "Student not found")
	default:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	}
}
