package hub

import (
	"errors"
	"github.com/gorilla/websocket"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

var (
	errClosed = errors.New("connection closed")
)

type Connection struct {
	lock   sync.Mutex
	conn   *websocket.Conn
	closed atomic.Bool
}

func (c *Connection) write(data any, timeout time.Duration) error {
	if c.closed.Load() {
		return errClosed
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	err := c.conn.SetWriteDeadline(time.Now().Add(timeout))
	if err != nil {
		return err
	}
	return c.conn.WriteJSON(data)
}

func (c *Connection) close(closeMessage []byte) {
	c.closed.Store(true)
	c.lock.Lock()
	defer c.lock.Unlock()

	if closeMessage != nil {
		err := c.conn.WriteControl(
			websocket.CloseMessage,
			closeMessage,
			time.Now().Add(5*time.Second),
		)
		if err != nil {
			slog.Error("failed to send close message in websocket connection",
				slog.String("reason", err.Error()),
			)
		}
	}

	if err := c.conn.Close(); err != nil {
		slog.Error("failed to close websocket connection",
			slog.String("reason", err.Error()),
		)
	}
}

type Hub struct {
	lock        sync.RWMutex
	connections map[int32]*Connection
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[int32]*Connection),
	}
}

func (h *Hub) Send(userId int32, data any, timeout time.Duration) {
	h.lock.RLock()
	c, ok := h.connections[userId]
	h.lock.RUnlock()
	if !ok {
		return
	}
	h.SendOnConn(c, userId, data, timeout)
}

func (h *Hub) SendOnConn(c *Connection, userId int32, data any, timeout time.Duration) {
	err := c.write(data, timeout)
	if err == nil {
		return
	}

	slog.Error("failed to write to websocket: closing connection",
		slog.Int("user_id", int(userId)),
		slog.String("reason", err.Error()),
	)

	// received error; remove and close the connection
	h.removeConnection(userId, c, nil)
}

func (h *Hub) Register(userId int32, conn *websocket.Conn) *Connection {
	newConn := &Connection{conn: conn}
	h.lock.Lock()
	old := h.connections[userId]
	h.connections[userId] = newConn
	go h.readMessages(userId, newConn)
	h.lock.Unlock()
	if old != nil {
		old.close(
			websocket.FormatCloseMessage(
				websocket.CloseNormalClosure,
				"received a newer connection from your user id; closing the old one",
			),
		)
	}
	return newConn
}

func (h *Hub) RemoveConnection(userId int32, c *Connection, err error) {
	var closeMessage []byte
	if err != nil {
		closeMessage = websocket.FormatCloseMessage(
			websocket.CloseInternalServerErr,
			err.Error(),
		)
	}
	h.removeConnection(userId, c, closeMessage)
}

func (h *Hub) removeConnection(userId int32, c *Connection, closeMessage []byte) {
	h.lock.Lock()
	n, ok := h.connections[userId]
	if ok && c == n {
		delete(h.connections, userId)
	}
	h.lock.Unlock()
	c.close(closeMessage)
}

func (h *Hub) readMessages(userId int32, c *Connection) {
	c.conn.SetReadLimit(1024)
	for {
		_, _, err := c.conn.NextReader()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseNormalClosure,
				websocket.CloseNoStatusReceived,
				websocket.CloseGoingAway,
			) {
				slog.Error("connection read error",
					slog.Int("user_id", int(userId)),
					slog.String("error", err.Error()),
				)
			}
			h.removeConnection(userId, c, nil)
			return
		}
	}
}

func (h *Hub) Broadcast(callback func(userId int32, c *Connection)) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	for userId, conn := range h.connections {
		callback(userId, conn)
	}
}
