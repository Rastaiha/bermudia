package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Rastaiha/bermudia/internal/models"
	"log/slog"
	"net/http"
)

func handleError(w http.ResponseWriter, err error) {
	if errors.Is(err, context.Canceled) {
		sendError(w, http.StatusRequestTimeout, "Request cancelled")
		return
	}

	slog.Error("internal error", err)
	sendError(w, http.StatusInternalServerError, "Internal server error")
}

func sendResult(w http.ResponseWriter, result any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := models.APIResponse{
		OK:     true,
		Result: result,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If encoding fails, send a basic error response
		http.Error(w, "Internal server error: could not successful response", http.StatusInternalServerError)
	}
}

func sendError(w http.ResponseWriter, statusCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := models.APIResponse{
		OK:    false,
		Error: errorMsg,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If encoding fails, send a basic error response
		http.Error(w, "Internal server error: could not error response", http.StatusInternalServerError)
	}
}
