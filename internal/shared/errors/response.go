package apperrors

import (
	"net/http"

	"github.com/Rumm1/eduhub-backend/internal/shared/response"
)

func WriteSuccess(w http.ResponseWriter, status int, data interface{}) {
	response.Success(w, status, data)
}

func WriteError(w http.ResponseWriter, status int, code string, message string) {
	response.Error(w, status, code, message)
}
