package response

import (
	"encoding/json"
	"net/http"
)

type Envelope map[string]interface{}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

func OK(w http.ResponseWriter, data interface{}) {
	Success(w, http.StatusOK, data)
}

func Success(w http.ResponseWriter, status int, data interface{}) {
	JSON(w, status, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

func Message(w http.ResponseWriter, status int, message string) {
	JSON(w, status, SuccessResponse{
		Success: true,
		Message: message,
	})
}

func Error(w http.ResponseWriter, status int, values ...string) {
	code := "ERROR"
	message := "Something went wrong"

	if len(values) == 1 {
		message = values[0]
	}

	if len(values) >= 2 {
		code = values[0]
		message = values[1]
	}

	JSON(w, status, ErrorResponse{
		Success: false,
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}
