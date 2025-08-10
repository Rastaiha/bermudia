package handler

import (
	"errors"
	"github.com/Rastaiha/bermudia/internal/repository"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) GetIsland(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "islandID")
	if id == "" {
		sendError(w, http.StatusBadRequest, "Island ID is required")
		return
	}

	island, err := h.islandService.GetIsland(r.Context(), id)
	if err != nil {
		// Check for specific error types using errors.Is
		if errors.Is(err, repository.ErrIslandNotFound) {
			sendError(w, http.StatusNotFound, "Island not found")
			return
		}

		handleError(w, err)
		return
	}

	sendResult(w, island)
}

func (h *Handler) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "inputID")
	if id == "" {
		sendError(w, http.StatusBadRequest, "input ID is required")
		return
	}

	sendResult(w, map[string]any{})
}
