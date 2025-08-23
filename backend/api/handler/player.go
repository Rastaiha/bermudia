package handler

import (
	"encoding/json"
	"github.com/Rastaiha/bermudia/internal/domain"
	"net/http"
)

func (h *Handler) HandlePlayerUpdateEvent(e *domain.FullPlayerUpdateEvent) {
	h.sendEvent(e.Player.UserId, event{PlayerUpdate: e})
}

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

type travelCheckRequest struct {
	FromIsland string `json:"fromIsland"`
	ToIsland   string `json:"toIsland"`
}

func (h *Handler) TravelCheck(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	var req travelCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendDecodeError(w)
		return
	}

	checkResult, err := h.playerService.TravelCheck(r.Context(), user, req.FromIsland, req.ToIsland)
	if err != nil {
		handleError(w, err)
		return
	}
	sendResult(w, checkResult)
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

func (h *Handler) RefuelCheck(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	check, err := h.playerService.RefuelCheck(r.Context(), user.ID)
	if err != nil {
		handleError(w, err)
		return
	}
	sendResult(w, check)
}

type refuelRequest struct {
	Amount int32 `json:"amount"`
}

func (h *Handler) Refuel(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	var req refuelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendDecodeError(w)
		return
	}

	err = h.playerService.Refuel(r.Context(), user.ID, req.Amount)
	if err != nil {
		handleError(w, err)
		return
	}

	sendResult(w, struct{}{})
}
