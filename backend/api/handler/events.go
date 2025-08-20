package handler

import (
	"github.com/Rastaiha/bermudia/internal/domain"
	"log/slog"
	"net/http"
	"time"
)

type event struct {
	PlayerUpdate *domain.PlayerUpdateEvent `json:"playerUpdate"`
	Timestamp    int64                     `json:"timestamp,string"`
}

func (h *Handler) sendEvent(userId int32, e event) {
	e.Timestamp = time.Now().UTC().UnixMilli()
	go func() {
		h.connectionHub.Send(userId, e, 15*time.Second)
	}()
}

func (h *Handler) StreamEvents(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	conn, err := h.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("websocket upgrade failed", err)
		return
	}
	h.connectionHub.Register(user.ID, conn)
}
