package handler

import (
	"errors"
	"github.com/Rastaiha/bermudia/api/hub"
	"github.com/Rastaiha/bermudia/internal/domain"
	"log/slog"
	"net/http"
	"time"
)

const (
	eventSendTimeout = 15 * time.Second
)

func (h *Handler) createConnection(w http.ResponseWriter, r *http.Request, hb *hub.Hub) (*domain.User, *hub.Connection) {
	token := r.URL.Query().Get("token")
	user, ok := h.authService.ValidateToken(r.Context(), token)
	if !ok {
		sendError(w, http.StatusUnauthorized, "Invalid auth token")
		return nil, nil
	}
	conn, err := h.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("websocket upgrade failed", err)
		return user, nil
	}
	connection := hb.Register(user.ID, conn)
	return user, connection
}

type playerEvent struct {
	PlayerUpdate *domain.FullPlayerUpdateEvent `json:"playerUpdate"`
	Timestamp    int64                         `json:"timestamp,string"`
}

func (h *Handler) HandlePlayerUpdateEvent(e *domain.FullPlayerUpdateEvent) {
	event := &playerEvent{
		PlayerUpdate: e,
		Timestamp:    time.Now().UTC().UnixMilli(),
	}
	go h.playerHub.Send(e.Player.UserId, event, eventSendTimeout)
}

func (h *Handler) StreamPlayerEvents(w http.ResponseWriter, r *http.Request) {
	user, c := h.createConnection(w, r, h.playerHub)
	if c == nil {
		return
	}
	if err := h.playerService.SendInitialEvents(r.Context(), user.ID); err != nil {
		slog.Error("send initial events failed", slog.String("error", err.Error()))
		h.playerHub.RemoveConnection(user.ID, c, errors.New("failed to send initial events"))
	}
}

func (h *Handler) HandleTradeEventBroadcast(eventProvider func(userId int32) *domain.TradeEvent) {
	h.tradeHub.Broadcast(func(userId int32, c *hub.Connection) {
		event := eventProvider(userId)
		go h.tradeHub.SendOnConn(c, userId, event, eventSendTimeout)
	})
}

func (h *Handler) StreamTradeEvents(w http.ResponseWriter, r *http.Request) {
	user, c := h.createConnection(w, r, h.tradeHub)
	if c == nil {
		return
	}

	event, err := h.playerService.GetInitialTradeEvent(r.Context())
	if err != nil {
		slog.Error("get initial trade event failed", slog.String("error", err.Error()))
		h.tradeHub.RemoveConnection(user.ID, c, errors.New("failed to get initial trade events"))
		return
	}
	h.tradeHub.SendOnConn(c, user.ID, event, eventSendTimeout)
}

func (h *Handler) HandleInboxEvent(e *domain.InboxEvent) {
	go h.inboxHub.Send(e.UserId, e, eventSendTimeout)
}

func (h *Handler) StreamInboxEvents(w http.ResponseWriter, r *http.Request) {
	user, c := h.createConnection(w, r, h.inboxHub)
	if c == nil {
		return
	}
	event, err := h.playerService.GetInitialInboxEvent(r.Context(), user.ID)
	if err != nil {
		slog.Error("get initial inbox event failed", slog.String("error", err.Error()))
		h.inboxHub.RemoveConnection(user.ID, c, errors.New("failed to get initial inbox events"))
		return
	}
	h.inboxHub.SendOnConn(c, user.ID, event, eventSendTimeout)
}

func (h *Handler) HandleBroadcastMessage(eventProvider func(userId int32) *domain.InboxMessageView) {
	h.inboxHub.Broadcast(func(userId int32, c *hub.Connection) {
		event := eventProvider(userId)
		if event != nil {
			go h.inboxHub.SendOnConn(c, userId, event, eventSendTimeout)
		}
	})
}
