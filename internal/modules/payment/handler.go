package payment

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
	var req CreatePaymentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.Create(r.Context(), req)
	if err != nil {
		writePaymentError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, result)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.List(r.Context())
	if err != nil {
		writePaymentError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) ListByStudentID(w http.ResponseWriter, r *http.Request) {
	studentID := chi.URLParam(r, "studentID")

	result, err := h.service.ListByStudentID(r.Context(), studentID)
	if err != nil {
		writePaymentError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) UpdateGroupPrice(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupID")

	var req UpdateGroupPriceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.UpdateGroupPrice(r.Context(), groupID, req)
	if err != nil {
		writePaymentError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) GetStudentBalance(w http.ResponseWriter, r *http.Request) {
	studentID := chi.URLParam(r, "studentID")
	groupID := r.URL.Query().Get("group_id")
	period := r.URL.Query().Get("period")

	result, err := h.service.GetStudentBalance(r.Context(), studentID, groupID, period)
	if err != nil {
		writePaymentError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func writePaymentError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTenantRequired):
		response.Error(w, http.StatusForbidden, "TENANT_REQUIRED", "Tenant organization is required")
	case errors.Is(err, ErrStudentIDRequired):
		response.Error(w, http.StatusBadRequest, "STUDENT_ID_REQUIRED", "Student id is required")
	case errors.Is(err, ErrStudentIDInvalid):
		response.Error(w, http.StatusBadRequest, "STUDENT_ID_INVALID", "Student id is invalid")
	case errors.Is(err, ErrStudentNotFound):
		response.Error(w, http.StatusBadRequest, "STUDENT_NOT_FOUND", "Student not found in organization")
	case errors.Is(err, ErrGroupIDRequired):
		response.Error(w, http.StatusBadRequest, "GROUP_ID_REQUIRED", "Group id is required")
	case errors.Is(err, ErrGroupIDInvalid):
		response.Error(w, http.StatusBadRequest, "GROUP_ID_INVALID", "Group id is invalid")
	case errors.Is(err, ErrGroupNotFound):
		response.Error(w, http.StatusBadRequest, "GROUP_NOT_FOUND", "Group not found in organization")
	case errors.Is(err, ErrStudentNotInGroup):
		response.Error(w, http.StatusBadRequest, "STUDENT_NOT_IN_GROUP", "Student is not in group")
	case errors.Is(err, ErrAmountRequired):
		response.Error(w, http.StatusBadRequest, "AMOUNT_REQUIRED", "Amount is required")
	case errors.Is(err, ErrAmountInvalid):
		response.Error(w, http.StatusBadRequest, "AMOUNT_INVALID", "Amount must be greater than zero")
	case errors.Is(err, ErrMonthlyPriceRequired):
		response.Error(w, http.StatusBadRequest, "MONTHLY_PRICE_REQUIRED", "Monthly price is required")
	case errors.Is(err, ErrMonthlyPriceInvalid):
		response.Error(w, http.StatusBadRequest, "MONTHLY_PRICE_INVALID", "Monthly price is invalid")
	case errors.Is(err, ErrPaymentDateRequired):
		response.Error(w, http.StatusBadRequest, "PAYMENT_DATE_REQUIRED", "Payment date is required")
	case errors.Is(err, ErrPaymentDateInvalid):
		response.Error(w, http.StatusBadRequest, "PAYMENT_DATE_INVALID", "Payment date must be in YYYY-MM-DD format")
	case errors.Is(err, ErrPaymentPeriodInvalid):
		response.Error(w, http.StatusBadRequest, "PAYMENT_PERIOD_INVALID", "Payment period must be in YYYY-MM format")
	case errors.Is(err, ErrPaymentStatusInvalid):
		response.Error(w, http.StatusBadRequest, "PAYMENT_STATUS_INVALID", "Payment status is invalid")
	default:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	}
}
