package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Game holds the connections to the players
type Game struct {
	id        string
	players   map[int]Player
	register  chan Player
	endGame   chan Player
	broadcast chan []byte
	board     *Board
	winner    Player
	createdAt time.Time
	available bool
}

func (g *Game) run() {
	for {
		select {
		case player := <-g.register:
			g.players[player.ID()] = player
			for _, p := range g.players {
				p.BroadcastCurrentPosition()
			}
			if len(g.players) == 2 {
				g.startCountdown()
				g.startGame()
			}
		case player := <-g.endGame:
			player.SetAlive(false)
			var winner Player
			alivePlayers := 0
			for _, p := range g.players {
				if p.IsAlive() {
					alivePlayers++
					winner = p
				}
			}
			if alivePlayers == 1 {
				g.winner = winner
				temp := make(map[string]interface{})
				temp["winner"] = winner.ID()
				res, err := json.Marshal(&temp)
				if err != nil {
					log.Printf("Could not convert to JSON, %v", err)
				}
				g.sendToAll(res)
				g.destroyPlayers()
				delete(activeGames, g.id)
				return
			}
		case message := <-g.broadcast:
			g.sendToAll(message)
		}
	}
}

func newGame(id string, height, width int) *Game {
	return &Game{
		id:        id,
		broadcast: make(chan []byte),
		register:  make(chan Player),
		endGame:   make(chan Player),
		players:   make(map[int]Player),
		board:     initBoard(height, width),
		winner:    nil,
		createdAt: time.Now(),
		available: true,
	}
}

func (g *Game) startGame() {
	g.available = false
	startTime := time.Now()
	// g.broadcast <- []byte(fmt.Sprintf("game started at %s", startTime))
	log.Printf("game started at %v", startTime)

	for _, p := range g.players {
		go p.Move()
	}
}

func (g *Game) startCountdown() {
	countdownTicker := time.NewTicker(1000 * time.Millisecond)
	defer countdownTicker.Stop()

	counter := 3
	for {
		select {
		case <-countdownTicker.C:
			temp := make(map[string]interface{})
			temp["countdown"] = counter
			res, err := json.Marshal(&temp)
			if err != nil {
				log.Printf("Could not convert to JSON, %v", err)
			}
			g.sendToAll(res)
			counter--
			if counter < 0 {
				return
			}
		}
	}
}

func (g *Game) sendToAll(message []byte) {
	for _, p := range g.players {
		p.Broadcast(message)
	}
}

func connectPlayer(game *Game, w http.ResponseWriter, r *http.Request) {
	id := len(game.players)
	if id < 2 {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		createPlayer(game, id, conn)
	}
}

func connectBot(game *Game) {
	id := len(game.players)
	if id < 2 {
		createPlayer(game, id, nil)
	}
}

func createPlayer(game *Game, id int, conn *websocket.Conn) {
	currentPosition := sync.Map{}
	send := make(chan []byte, 256)
	var player Player
	if conn != nil {
		human := &Human{PlayerData{id, -1, game, send, &currentPosition, nil, true}, conn}
		player = human
	} else {
		bot := &Bot{PlayerData{id, -1, game, send, &currentPosition, nil, true}}
		player = bot
	}
	player.InitPlayer()
	game.register <- player
}

func (g *Game) destroyPlayers() {
	for _, p := range g.players {
		p.Destroy()
	}
}
