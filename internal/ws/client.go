package ws

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pingPeriod = (writeWait * 9) / 10
)

// Client represents a single WebSocket connection to the server, allowing for sending and receiving messages
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// constructor function NewClient initializes a new Client instance with the provided WebSocket connection
func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}
}

// ReadPump listens for incoming messages from the WebSocket connection and processes them using the provided Handler
func (c *Client) ReadPump(h *Handler) {
	defer func() {
		h.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		h.handleClaimBlock(message)
	}
}

// WritePump listens for outgoing messages on the send channel and writes them to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			// set a write deadline to ensure that the connection is closed if the client is unresponsive
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// write the message to the WebSocket connection
			writer, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err := writer.Write(message); err != nil {
				return
			}

			// batch messages from the send channel to minimize the number of writes to the WebSocket connection
			n := len(c.send)
			for i := 0; i < n; i++ {
				if _, err := writer.Write(<-c.send); err != nil {
					return
				}
			}
			if err := writer.Close(); err != nil {
				return
			}
			// send a ping message to the client to keep the connection alive
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
