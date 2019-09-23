package main

import (
	"encoding/json"
	"log"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Human represents the websocket player
type Human struct {
	PlayerData
	conn *websocket.Conn
}

// ID returns the players Id
func (h *Human) ID() int {
	return h.id
}

// ClientID returns the id obtained from the WS client
func (h *Human) ClientID() int {
	return h.clientID
}

func (h *Human) setClientID(id int) {
	h.clientID = id
}

// Game returns the pointer to Game
func (h *Human) Game() *Game {
	return h.game
}

// CurrentPosition returns a map that contains
// info about the human players current position
func (h *Human) CurrentPosition() *sync.Map {
	return h.currentPosition
}

// IsAlive returns the players alive status
func (h *Human) IsAlive() bool {
	return h.alive
}

// SetAlive sets the players alive status
func (h *Human) SetAlive(alive bool) {
	h.alive = alive
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (h *Human) readPump() {
	h.conn.SetReadLimit(maxMessageSize)
	h.conn.SetReadDeadline(time.Now().Add(pongWait))
	h.conn.SetPongHandler(func(string) error { h.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := h.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		var event map[string]string
		if err := json.Unmarshal(message, &event); err != nil {
			log.Printf("unmarshal error: %v", err)
		}
		if event["dir"] == directionDown {
			h.rotationChannel <- RotationData{dir: event["dir"], key: event["key"]}
		} else if event["dir"] == directionUp {
			h.rotationChannel <- RotationData{dir: event["dir"]}
		} else if event["clientId"] != "" {
			clientID, err := strconv.Atoi(event["clientId"])
			if err != nil {
				log.Printf("Cannot convert %s to int", event["clientId"])
			}
			h.setClientID(clientID)
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (h *Human) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		h.conn.Close()
	}()
	for {
		select {
		case message, ok := <-h.send:
			h.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				h.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := h.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			h.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := h.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Broadcast sends the message to writePump, which eventually sends it
// to the websocket client
func (h *Human) Broadcast(message []byte) {
	h.send <- message
}

// Destroy closes all channels and removes player from the game
func (h *Human) Destroy() {
	h.StopRotation()
	close(h.send)
	delete(h.game.players, h.id)
}

// InitPlayer initializes the players position and opens read/write channels
func (h *Human) InitPlayer() {
	go h.writePump()
	go h.readPump()

	startX := width / 2
	startY := height/5 + h.id*3*height/5
	h.currentPosition.Store("x", startX)
	h.currentPosition.Store("y", startY)
	h.currentPosition.Store("rotation", getStartRotation())
	h.currentPosition.Store("trace", false)
	h.game.board.fields[startX][startY].setUsed(h)
}

// StartRotation creates a new ticker and updates
// the players rotation angle on each tick
func (h *Human) StartRotation(direction string) {
	h.StopRotation()
	h.rotationTicker = time.NewTicker(30 * time.Millisecond)
	go func() {
		for {
			select {
			case <-h.rotationTicker.C:
				currRotation, _ := h.currentPosition.Load("rotation")

				if direction == directionRight {
					h.currentPosition.Store("rotation", (currRotation.(int)+5)%360)
				} else {
					h.currentPosition.Store("rotation", (currRotation.(int)+355)%360)
				}
				h.currentPosition.Store("rotationDir", direction)
			}
		}
	}()
}

// StopRotation stops the rotation ticker
func (h *Human) StopRotation() {
	if h.rotationTicker != nil {
		h.rotationTicker.Stop()
		h.rotationTicker = nil
		h.currentPosition.Store("rotationDir", nil)
	}
}

// BroadcastCurrentPosition sends the players current position
// to all clients
func (h *Human) BroadcastCurrentPosition() {
	temp := make(map[string]interface{})
	playerPositionMap := syncMapToMap(h.currentPosition)
	playerPositionMap["clientId"] = h.ClientID()
	playersStatusMap := make(map[int]interface{})
	playersStatusMap[h.id] = playerPositionMap
	temp["players"] = playersStatusMap

	res, err := json.Marshal(&temp)
	if err != nil {
		log.Printf("Could not convert to JSON, %v", err)
		return
	}

	h.game.sendToAll(res)
}

// Move sets the new position for the player calculated by
// the rotation angle on each tick.
func (h *Human) Move() {
	mainTicker := time.NewTicker(1000 / fps * time.Millisecond)
	defer mainTicker.Stop()

	visitedTicker := createRandomIntervalTicker(1000, 2000)

	for {
		select {
		case <-mainTicker.C:
			if !h.game.started {
				h.BroadcastCurrentPosition()
			} else {
				if !h.alive || h.game.winner != nil {
					visitedTicker.Stop()
					return
				}
				curX, _ := h.currentPosition.Load("x")
				curY, _ := h.currentPosition.Load("y")
				curRotation, _ := h.currentPosition.Load("rotation")
				rotationRad := float64(curRotation.(int)) * math.Pi / 180
				moveBresenham(h, curX.(int), curY.(int), rotationRad)
			}
		case rotationData := <-h.rotationChannel:
			if rotationData.dir == directionDown {
				h.StartRotation(rotationData.key)
			} else if rotationData.dir == directionUp {
				h.StopRotation()
			}
		case <-visitedTicker.C:
			trace, _ := h.currentPosition.Load("trace")
			h.currentPosition.Store("trace", !trace.(bool))
			visitedTicker = createRandomIntervalTicker(1000, 2000)
		}
	}
}
