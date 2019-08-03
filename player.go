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

	fps = 10
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
	alive           bool
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (p *Player) readPump() {
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
			log.Printf("unmarshal error: %v", err)
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
		player := &Player{playerCount, game, conn, send, &currentPosition, nil, true}
		go player.writePump()
		go player.readPump()

		player.game.register <- player
	}
}

func (p *Player) startRotation(direction string) {
	p.rotationTicker = time.NewTicker(30 * time.Millisecond)

	for range p.rotationTicker.C {
		currRotation, _ := p.currentPosition.Load("rotation")
		if direction == "right" {
			p.currentPosition.Store("rotation", (currRotation.(int)+5)%360)
		} else {
			p.currentPosition.Store("rotation", (currRotation.(int)+355)%360)
		}
	}
}

func (p *Player) stopRotation() {
	p.rotationTicker.Stop()
	p.rotationTicker = nil
}

func moveBresenham(p *Player, x0 int, y0 int, rotationRad float64) {
	x1 := int(float64(x0) + math.Cos(rotationRad)*1000)
	y1 := int(float64(y0) + math.Sin(rotationRad)*1000)

	dx := int(math.Abs(float64(x1 - x0)))
	dy := -(int(math.Abs(float64(y1 - y0))))
	sx := 1
	if x1 < x0 {
		sx = -1
	}
	sy := 1
	if y1 < y0 {
		sy = -1
	}
	err := dx + dy
	for i := 0; i < 3; i++ {
		e2 := 2 * err
		if e2 >= dy {
			if x0 != x1 {
				err += dy
				x0 += sx
			}
		}
		if e2 <= dx {
			if y0 != y1 {
				err += dx
				y0 += sy
			}
		}
		if p.game.board.isValidMove(x0, y0) {
			p.currentPosition.Store("x", x0)
			p.currentPosition.Store("y", y0)
			trace, ok := p.currentPosition.Load("trace")
			if !ok {
				log.Printf("Could not load value from map")
				continue
			}
			if trace.(bool) {
				p.game.board.fields[x0][y0].setUsed(p)
			}
			p.broadcastCurrentPosition()
		} else {
			p.alive = false
			p.game.endGame <- p
			break
		}
	}
}

func (p *Player) broadcastCurrentPosition() {
	temp := make(map[string]interface{})
	playersStatusMap := make(map[int]interface{})
	playersStatusMap[p.id] = syncMapToMap(p.currentPosition)
	temp["players"] = playersStatusMap

	res, err := json.Marshal(&temp)
	if err != nil {
		log.Printf("Could not convert to JSON")
		return
	}

	p.game.sendToAll(res)
}

func (p *Player) move() {
	mainTicker := time.NewTicker(1000 / fps * time.Millisecond)
	defer mainTicker.Stop()

	visitedTicker := createRandomIntervalTicker(1000, 2000)

	for {
		select {
		case <-mainTicker.C:
			if !p.alive || p.game.winner != nil {
				visitedTicker.Stop()
				break
			}
			curX, _ := p.currentPosition.Load("x")
			curY, _ := p.currentPosition.Load("y")
			curRotation, _ := p.currentPosition.Load("rotation")
			rotationRad := float64(curRotation.(int)) * math.Pi / 180
			moveBresenham(p, curX.(int), curY.(int), rotationRad)
		case <-visitedTicker.C:
			trace, ok := p.currentPosition.Load("trace")
			if !ok {
				log.Printf("Could not load value from map")
				continue
			}
			p.currentPosition.Store("trace", !trace.(bool))
			visitedTicker = createRandomIntervalTicker(1000, 2000)
		}
	}
}
