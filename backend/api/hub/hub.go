package hub

import (
	"errors"
	"github.com/gorilla/websocket"
	"log/slog"
	"sync"
	"time"
)

var (
	errClosed = errors.New("connection closed")
)

type connection struct {
	lock   sync.Mutex
	conn   *websocket.Conn
	closed bool
}

func (c *connection) Write(data any, timeout time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.closed {
		return errClosed
	}
	err := c.conn.SetWriteDeadline(time.Now().Add(timeout))
	if err != nil {
		return err
	}
	return c.conn.WriteJSON(data)
}

func (c *connection) Close() {
	c.lock.Lock()
	c.closed = true
	c.lock.Unlock()
	if err := c.conn.Close(); err != nil {
		slog.Error("failed to close websocket connection",
			slog.String("reason", err.Error()),
		)
	}
}

type Hub struct {
	lock        sync.RWMutex
	connections map[int32]*connection
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[int32]*connection),
	}
}

func (h *Hub) Send(userId int32, data any, timeout time.Duration) {
	h.lock.RLock()
	c, ok := h.connections[userId]
	h.lock.RUnlock()
	if !ok {
		return
	}
	err := c.Write(data, timeout)
	if err == nil {
		return
	}

	slog.Error("failed to write to websocket: closing connection",
		slog.Int("user_id", int(userId)),
		slog.String("reason", err.Error()),
	)

	// received error; remove and close the connection
	h.removeConnection(userId, c)
}

func (h *Hub) Register(userId int32, conn *websocket.Conn) {
	newConn := &connection{conn: conn}
	h.lock.Lock()
	old := h.connections[userId]
	h.connections[userId] = newConn
	go h.readMessages(userId, newConn)
	h.lock.Unlock()
	if old != nil {
		old.Close()
	}
}

func (h *Hub) removeConnection(userId int32, c *connection) {
	h.lock.Lock()
	n, ok := h.connections[userId]
	if !ok || c != n {
		h.lock.Unlock()
		return
	}
	delete(h.connections, userId)
	h.lock.Unlock()
	c.Close()
}

func (h *Hub) readMessages(userId int32, c *connection) {
	c.conn.SetReadLimit(1024)
	for {
		_, _, err := c.conn.NextReader()
		if err != nil {
			slog.Error("connection read error",
				slog.Int("user_id", int(userId)),
				slog.String("error", err.Error()),
			)
			h.removeConnection(userId, c)
			return
		}
	}
}
