package handler

import (
	"errors"
	"github.com/Rastaiha/bermudia/internal/domain"
	"log/slog"
	"net/http"
	"time"
)

type event struct {
	PlayerUpdate *domain.FullPlayerUpdateEvent `json:"playerUpdate"`
	Timestamp    int64                         `json:"timestamp,string"`
}

func (h *Handler) sendEvent(userId int32, e event) {
	e.Timestamp = time.Now().UTC().UnixMilli()
	go func() {
		h.connectionHub.Send(userId, e, 15*time.Second)
	}()
}

func (h *Handler) StreamEvents(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	user, ok := h.authService.ValidateToken(r.Context(), token)
	if !ok {
		sendError(w, http.StatusUnauthorized, "Invalid auth token")
		return
	}
	conn, err := h.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("websocket upgrade failed", err)
		return
	}
	connection := h.connectionHub.Register(user.ID, conn)
	if err := h.playerService.SendInitialEvents(r.Context(), user.ID); err != nil {
		slog.Error("send initial events failed", slog.String("error", err.Error()))
		h.connectionHub.RemoveConnection(user.ID, connection, errors.New("failed to send initial events"))
	}
}
