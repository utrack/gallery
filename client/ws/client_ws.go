/*Package clientWs provides the client.Connection implementation
that uses websockets as a transport.
*/
package clientWs

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	clientIface "github.com/utrack/gallery/client"
	"sync"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// client wraps the websocket connection.
type client struct {
	ws *websocket.Conn

	queueMsgsOut       chan []byte
	queueMsgsOutClosed bool
	queueMsgsOutMu     sync.Mutex

	disconChan chan error
}

// NewClient returns a client.Connection that uses passed
// WebSocket connection as a transport.
func NewClient(ws *websocket.Conn) clientIface.Connection {
	ret := &client{
		ws: ws,

		queueMsgsOut: make(chan []byte, 10),
		disconChan:   make(chan error, 2),
	}
	go ret.readPump()
	go ret.writePump()
	return ret
}

func (c *client) Send(m json.RawMessage) {
	c.queueMsgsOut <- m
}

func (c *client) DisconChan() <-chan error {
	return c.disconChan
}

func (c *client) Disconnect() {
	c.discon(errors.New("Spite"))
	c.ws.Close()
}

func (c *client) discon(reason error) {
	c.queueMsgsOutMu.Lock()
	defer c.queueMsgsOutMu.Unlock()

	if c.queueMsgsOutClosed {
		return
	}

	c.queueMsgsOutClosed = true
	close(c.queueMsgsOut)

	c.disconChan <- reason
}

func (c *client) readPump() {
	var err error
	// Forward the discon info on disconnect
	defer func() {
		c.discon(err)
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		if _, _, err = c.ws.NextReader(); err != nil {
			break
		}
	}
}

func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	var err error
	// Forward the discon info on disconnect
	defer func() {
		c.discon(err)
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.queueMsgsOut:
			// Input chan was closed, shutdown the sock
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			if err = c.write(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			if err = c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// write writes a message with the given message type and payload.
func (c *client) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}
