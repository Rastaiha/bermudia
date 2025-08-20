package hub

import (
	"errors"
	"github.com/gorilla/websocket"
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
	_ = c.conn.Close()
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

	// receive error; remove and close the connection
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

func (h *Hub) Register(userId int32, conn *websocket.Conn) {
	newConn := &connection{conn: conn}
	h.lock.Lock()
	old := h.connections[userId]
	h.connections[userId] = newConn
	h.lock.Unlock()
	if old != nil {
		old.Close()
	}
}
