package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Candidate represents a potential player which is waiting in the lobby
type Candidate struct {
	send    chan []byte
	receive chan []byte
	conn    *websocket.Conn
}

func newCandidate(conn *websocket.Conn) Candidate {
	c := Candidate{make(chan []byte), make(chan []byte), conn}
	go c.writePump()
	go c.readPump()
	return c
}

func (c *Candidate) IsConnected() bool {
	timeoutTicker := time.NewTicker(200 * time.Millisecond)
	defer func() {
		timeoutTicker.Stop()
	}()
	c.send <- []byte("ping")
	for {
		select {
		case <-c.receive:
			return true
		case <-timeoutTicker.C:
			return false
		}
	}
}

func (c *Candidate) Redirect(message []byte) {
	c.send <- message
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Candidate) readPump() {
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.receive <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Candidate) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
