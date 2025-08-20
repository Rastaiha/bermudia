package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Rastaiha/bermudia/internal/domain"
	"log/slog"
	"net/http"
)

type APIResponse struct {
	OK     bool   `json:"ok"`
	Error  string `json:"error,omitempty"`
	Result any    `json:"result,omitempty"`
}

func handleError(w http.ResponseWriter, err error) {
	if errors.Is(err, context.Canceled) {
		sendError(w, http.StatusRequestTimeout, "Request cancelled")
		return
	}
	var domainError domain.Error
	if errors.As(err, &domainError) {
		switch domainError.Reason() {
		case domain.ErrorReasonResourceNotFound:
			sendError(w, http.StatusNotFound, err.Error())
		case domain.ErrorReasonRuleViolation:
			sendError(w, http.StatusConflict, err.Error())
		default:
			sendError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	slog.Error("internal error", err)
	sendError(w, http.StatusInternalServerError, "Internal server error")
}

func sendResult(w http.ResponseWriter, result any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := APIResponse{
		OK:     true,
		Result: result,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If encoding fails, send a basic error response
		http.Error(w, "Internal server error: could not successful response", http.StatusInternalServerError)
	}
}

func sendDecodeError(w http.ResponseWriter) {
	sendError(w, http.StatusBadRequest, "Invalid request payload")
}

func sendError(w http.ResponseWriter, statusCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		OK:    false,
		Error: errorMsg,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If encoding fails, send a basic error response
		http.Error(w, "Internal server error: could not error response", http.StatusInternalServerError)
	}
}
