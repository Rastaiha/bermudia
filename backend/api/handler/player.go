package handler

import (
	"encoding/json"
	"errors"
	"github.com/Rastaiha/bermudia/internal/domain"
	"net/http"
	"strconv"
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

type anchorCheckRequest struct {
	Island string `json:"island"`
}

func (h *Handler) AnchorCheck(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	var req anchorCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendDecodeError(w)
		return
	}

	result, err := h.playerService.AnchorCheck(r.Context(), user.ID, req.Island)
	if err != nil {
		handleError(w, err)
		return
	}

	sendResult(w, result)
}

type anchorRequest struct {
	Island string `json:"island"`
}

func (h *Handler) Anchor(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	var req anchorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendDecodeError(w)
		return
	}

	err = h.playerService.Anchor(r.Context(), user.ID, req.Island)
	if err != nil {
		handleError(w, err)
		return
	}

	sendResult(w, struct{}{})
}

func (h *Handler) MigrateCheck(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	result, err := h.playerService.MigrateCheck(r.Context(), user.ID)
	if err != nil {
		handleError(w, err)
		return
	}

	sendResult(w, result)
}

type migrateRequest struct {
	ToTerritory string `json:"toTerritory"`
}

func (h *Handler) Migrate(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	var req migrateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendDecodeError(w)
		return
	}

	err = h.playerService.Migrate(r.Context(), user.ID, req.ToTerritory)
	if err != nil {
		handleError(w, err)
		return
	}

	sendResult(w, struct{}{})
}

type unlockTreasureCheckRequest struct {
	TreasureID string `json:"treasureID"`
}

func (h *Handler) UnlockTreasureCheck(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	var req unlockTreasureCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendDecodeError(w)
		return
	}

	result, err := h.playerService.UnlockTreasureCheck(r.Context(), user.ID, req.TreasureID)
	if err != nil {
		handleError(w, err)
		return
	}

	sendResult(w, result)
}

type unlockTreasureRequest struct {
	TreasureID string `json:"treasureID"`
}

func (h *Handler) UnlockTreasure(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	var req unlockTreasureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendDecodeError(w)
		return
	}

	result, err := h.playerService.UnlockTreasure(r.Context(), user.ID, req.TreasureID)
	if err != nil {
		handleError(w, err)
		return
	}

	sendResult(w, result)
}

func (h *Handler) MakeOfferCheck(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	result, err := h.playerService.MakeOfferCheck(r.Context(), user.ID)
	if err != nil {
		handleError(w, err)
		return
	}
	sendResult(w, result)
}

type makeOfferRequest struct {
	Offered   domain.Cost `json:"offered"`
	Requested domain.Cost `json:"requested"`
}

func (h *Handler) MakeOffer(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	var req makeOfferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendDecodeError(w)
		return
	}

	result, err := h.playerService.MakeOffer(r.Context(), user, req.Offered, req.Requested)
	if err != nil {
		handleError(w, err)
		return
	}

	sendResult(w, result)
}

type acceptOfferRequest struct {
	OfferID string `json:"offerID"`
}

func (h *Handler) AcceptOffer(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	var req acceptOfferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendDecodeError(w)
		return
	}

	err = h.playerService.AcceptOffer(r.Context(), user.ID, req.OfferID)
	if err != nil {
		handleError(w, err)
		return
	}

	sendResult(w, struct{}{})
}

type deleteOfferRequest struct {
	OfferID string `json:"offerID"`
}

func (h *Handler) DeleteOffer(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	var req deleteOfferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendDecodeError(w)
		return
	}

	err = h.playerService.DeleteOffer(r.Context(), user.ID, req.OfferID)
	if err != nil {
		handleError(w, err)
		return
	}

	sendResult(w, struct{}{})
}

func (h *Handler) GetTradeOffers(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	offsetStr := r.URL.Query().Get("offset")
	offset := int64(0)
	if offsetStr != "" {
		offset, err = strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			sendDecodeError(w)
			return
		}
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	filter := r.URL.Query().Get("by")

	offers, err := h.playerService.GetTradeOffers(r.Context(), user.ID, domain.GetOffersByFilterType(filter), offset, limit)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidFilter) {
			sendError(w, http.StatusBadRequest, "invalid filter")
			return
		}
		handleError(w, err)
		return
	}

	sendResult(w, offers)
}
