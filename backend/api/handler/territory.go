package handler

import (
	"errors"
	"github.com/Rastaiha/bermudia/internal/domain"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) GetTerritory(w http.ResponseWriter, r *http.Request) {
	territoryID := chi.URLParam(r, "territoryID")

	if territoryID == "" {
		sendError(w, http.StatusBadRequest, "Territory ID is required")
		return
	}

	territory, err := h.territoryService.GetTerritory(r.Context(), territoryID)
	if err != nil {
		// Check for specific error types using errors.Is
		if errors.Is(err, domain.ErrTerritoryNotFound) {
			sendError(w, http.StatusNotFound, "Territory not found")
			return
		}

		handleError(w, err)
		return
	}

	sendResult(w, territory)
}
