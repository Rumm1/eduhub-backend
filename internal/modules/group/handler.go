package group

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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateGroupRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.Create(r.Context(), req)
	if err != nil {
		writeGroupError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, result)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.List(r.Context())
	if err != nil {
		writeGroupError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) AddStudent(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupID")

	var req AddStudentToGroupRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if err := h.service.AddStudent(r.Context(), groupID, req); err != nil {
		writeGroupError(w, err)
		return
	}

	response.Message(w, http.StatusCreated, "Student added to group")
}

func (h *Handler) ListStudents(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupID")

	result, err := h.service.ListStudents(r.Context(), groupID)
	if err != nil {
		writeGroupError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func writeGroupError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTenantRequired):
		response.Error(w, http.StatusForbidden, "TENANT_REQUIRED", "Tenant organization is required")
	case errors.Is(err, ErrBranchIDRequired):
		response.Error(w, http.StatusBadRequest, "BRANCH_ID_REQUIRED", "Branch id is required")
	case errors.Is(err, ErrBranchIDInvalid):
		response.Error(w, http.StatusBadRequest, "BRANCH_ID_INVALID", "Branch id is invalid")
	case errors.Is(err, ErrSubjectIDRequired):
		response.Error(w, http.StatusBadRequest, "SUBJECT_ID_REQUIRED", "Subject id is required")
	case errors.Is(err, ErrSubjectIDInvalid):
		response.Error(w, http.StatusBadRequest, "SUBJECT_ID_INVALID", "Subject id is invalid")
	case errors.Is(err, ErrTeacherIDInvalid):
		response.Error(w, http.StatusBadRequest, "TEACHER_ID_INVALID", "Teacher id is invalid")
	case errors.Is(err, ErrGroupIDInvalid):
		response.Error(w, http.StatusBadRequest, "GROUP_ID_INVALID", "Group id is invalid")
	case errors.Is(err, ErrStudentIDRequired):
		response.Error(w, http.StatusBadRequest, "STUDENT_ID_REQUIRED", "Student id is required")
	case errors.Is(err, ErrStudentIDInvalid):
		response.Error(w, http.StatusBadRequest, "STUDENT_ID_INVALID", "Student id is invalid")
	case errors.Is(err, ErrGroupNameRequired):
		response.Error(w, http.StatusBadRequest, "GROUP_NAME_REQUIRED", "Group name is required")
	case errors.Is(err, ErrStartDateInvalid):
		response.Error(w, http.StatusBadRequest, "START_DATE_INVALID", "Start date must be in YYYY-MM-DD format")
	case errors.Is(err, ErrEndDateInvalid):
		response.Error(w, http.StatusBadRequest, "END_DATE_INVALID", "End date must be in YYYY-MM-DD format")
	case errors.Is(err, ErrBranchNotFound):
		response.Error(w, http.StatusBadRequest, "BRANCH_NOT_FOUND", "Branch not found in organization")
	case errors.Is(err, ErrSubjectNotFound):
		response.Error(w, http.StatusBadRequest, "SUBJECT_NOT_FOUND", "Subject not found in organization")
	case errors.Is(err, ErrTeacherNotFound):
		response.Error(w, http.StatusBadRequest, "TEACHER_NOT_FOUND", "Teacher not found in organization")
	case errors.Is(err, ErrGroupNotFound):
		response.Error(w, http.StatusNotFound, "GROUP_NOT_FOUND", "Group not found in organization")
	case errors.Is(err, ErrStudentNotFound):
		response.Error(w, http.StatusBadRequest, "STUDENT_NOT_FOUND", "Student not found in organization")
	case errors.Is(err, ErrStudentBranchMismatch):
		response.Error(w, http.StatusBadRequest, "STUDENT_BRANCH_MISMATCH", "Student branch does not match group branch")
	case errors.Is(err, ErrGroupIsFull):
		response.Error(w, http.StatusBadRequest, "GROUP_IS_FULL", "Group is full")
	default:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	}
}
