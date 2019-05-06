package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
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

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Player is a middleman between the connection and the Game
type Player struct {
	id              int
	game            *Game
	conn            *websocket.Conn
	send            chan []byte
	currentPosition *sync.Map
	rotationTicker  *time.Ticker
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (p *Player) readPump() {
	// defer func() {
	// 	p.game.unregister <- p
	// 	p.conn.Close()
	// }()
	p.conn.SetReadLimit(maxMessageSize)
	p.conn.SetReadDeadline(time.Now().Add(pongWait))
	p.conn.SetPongHandler(func(string) error { p.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := p.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		var event map[string]string
		if err := json.Unmarshal(message, &event); err != nil {
			panic(err)
		}
		if event["dir"] == "down" {
			go p.startRotation(event["key"])
		} else {
			p.stopRotation()
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (p *Player) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		p.conn.Close()
	}()
	for {
		select {
		case message, ok := <-p.send:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				p.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := p.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(p.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-p.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func connect(game *Game, w http.ResponseWriter, r *http.Request) {
	playerCount := len(game.players)
	if playerCount < 2 {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		// defer conn.Close()

		send := make(chan []byte, 256)
		currentPosition := sync.Map{}
		currentPosition.Store("x", width/2)
		currentPosition.Store("y", height/5+playerCount*3*height/5)
		currentPosition.Store("rotation", getStartRotation())
		player := &Player{playerCount, game, conn, send, &currentPosition, nil}
		player.game.register <- player

		go player.writePump()
		go player.readPump()

		if playerCount == 1 {
			log.Printf("starting game")
			game.startGame()
		}
	} else {
		log.Printf("game has started, cannot join :(")
	}
}

func (p *Player) startRotation(direction string) {
	p.rotationTicker = time.NewTicker(50 * time.Millisecond)

	for range p.rotationTicker.C {
		currRotation, _ := p.currentPosition.Load("rotation")
		if direction == "left" {
			p.currentPosition.Store("rotation", (currRotation.(int)+5)%360)
		} else {
			p.currentPosition.Store("rotation", (currRotation.(int)+355)%360)
		}
	}
}

func (p *Player) stopRotation() {
	p.rotationTicker.Stop()
}

func (p *Player) move() {
	curX, _ := p.currentPosition.Load("x")
	curY, _ := p.currentPosition.Load("y")
	curRotation, _ := p.currentPosition.Load("rotation")
	rotationRad := float64(curRotation.(int)) * math.Pi / 180
	sinValue := math.Sin(rotationRad)
	cosValue := math.Cos(rotationRad)
	newX := curX.(int)
	newY := curY.(int)

	if math.Abs(cosValue) >= 0.5 {
		newX = curX.(int) + int(cosValue/math.Abs(cosValue))
	}
	if math.Abs(sinValue) >= 0.5 {
		newY = curY.(int) - int(sinValue/math.Abs(sinValue))
	}
	if p.game.board.isValidMove(newX, newY) {
		p.currentPosition.Store("x", newX)
		p.currentPosition.Store("y", newY)
		p.game.board.fields[newX][newY].isUsed = true
	} else {
		p.game.endGame <- p
	}
}
