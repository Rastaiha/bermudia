package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/Rastaiha/rasta-1404-contest/internal/models"
	"github.com/Rastaiha/rasta-1404-contest/internal/repository"
	"github.com/Rastaiha/rasta-1404-contest/internal/service"
)

// TerritoryHandler handles HTTP requests for territories
type TerritoryHandler struct {
	service *service.TerritoryService
}

// NewTerritoryHandler creates a new territory handler
func NewTerritoryHandler(service *service.TerritoryService) *TerritoryHandler {
	return &TerritoryHandler{
		service: service,
	}
}

// GetTerritory handles GET /api/v1/territories/{territoryID}
func (h *TerritoryHandler) GetTerritory(w http.ResponseWriter, r *http.Request) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	territoryID := chi.URLParam(r, "territoryID")

	if territoryID == "" {
		h.sendResponse(w, http.StatusBadRequest, false, "Territory ID is required", nil)
		return
	}

	territory, err := h.service.GetTerritory(ctx, territoryID)
	if err != nil {
		// Check for specific error types using errors.Is
		if errors.Is(err, repository.ErrTerritoryNotFound) {
			h.sendResponse(w, http.StatusNotFound, false, "Territory not found", nil)
			return
		}

		if errors.Is(err, repository.ErrInvalidData) {
			h.sendResponse(w, http.StatusInternalServerError, false, "Invalid territory data", nil)
			return
		}

		// Check for context cancellation
		if errors.Is(err, context.Canceled) {
			h.sendResponse(w, http.StatusRequestTimeout, false, "Request cancelled", nil)
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			h.sendResponse(w, http.StatusRequestTimeout, false, "Request timeout", nil)
			return
		}

		h.sendResponse(w, http.StatusInternalServerError, false, "Failed to retrieve territory", nil)
		return
	}

	h.sendResponse(w, http.StatusOK, true, "", territory)
}

// Helper method to send generic API responses
func (h *TerritoryHandler) sendResponse(w http.ResponseWriter, statusCode int, ok bool, errorMsg string, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := models.APIResponse{
		OK: ok,
	}

	if !ok && errorMsg != "" {
		response.Error = errorMsg
	}

	if ok && result != nil {
		response.Result = result
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If encoding fails, send a basic error response
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
