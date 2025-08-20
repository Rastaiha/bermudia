package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	player, err := h.playerService.GetPlayer(r.Context(), user)
	if err != nil {
		handleError(w, err)
		return
	}
	sendResult(w, player)
}

type travelRequest struct {
	FromIsland string `json:"fromIsland"`
	ToIsland   string `json:"toIsland"`
}

func (h *Handler) Travel(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	var req travelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendDecodeError(w)
		return
	}

	err = h.playerService.Travel(r.Context(), user, req.FromIsland, req.ToIsland)
	if err != nil {
		handleError(w, err)
		return
	}
	sendResult(w, struct{}{})
}
