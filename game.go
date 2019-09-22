package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const width int = 500
const height int = 600

// Game holds the connections to the players
type Game struct {
	id        string
	players   map[int]Player
	lobby     chan int
	register  chan Player
	endGame   chan Player
	broadcast chan []byte
	board     *Board
	winner    Player
	createdAt time.Time
	available bool
	started   bool
}

func (g *Game) run() {
	for {
		select {
		case player := <-g.register:
			g.players[player.ID()] = player
			if len(g.players) == 2 {
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
		lobby:     make(chan int),
		register:  make(chan Player),
		endGame:   make(chan Player),
		players:   make(map[int]Player),
		board:     initBoard(height, width),
		winner:    nil,
		createdAt: time.Now(),
		available: true,
		started:   false,
	}
}

func createGameAndWait() (string, error) {
	gameID := randToken()
	game := newGame(gameID, height, width)
	activeGames[gameID] = game
	go game.run()

	timeoutTicker := time.NewTicker(time.Minute)
	defer timeoutTicker.Stop()
	joinedPlayers := 1

	for {
		select {
		case <-game.lobby:
			joinedPlayers++
			if joinedPlayers == 2 {
				game.available = false
				return gameID, nil
			}
		case <-timeoutTicker.C:
			delete(activeGames, gameID)
			return "", errors.New("There are no active players to join")
		}
	}
}

func findAvailableGameAndJoin() (string, error) {
	timeoutTicker := time.NewTicker(time.Minute)
	countdownTicker := time.NewTicker(time.Second)
	defer func() {
		countdownTicker.Stop()
		timeoutTicker.Stop()
	}()
	var oldestGame *Game

	for {
		select {
		case <-countdownTicker.C:
			for _, game := range activeGames {
				if game.available && (oldestGame == nil || oldestGame.createdAt.After(game.createdAt)) {
					oldestGame = game
				}
			}
			if oldestGame != nil {
				oldestGame.lobby <- 1
				return oldestGame.id, nil
			}
		case <-timeoutTicker.C:
			return "", errors.New("There are no active games")
		}
	}
}

func (g *Game) startGame() {
	startTime := time.Now()
	log.Printf("game started at %v", startTime)

	for _, p := range g.players {
		p.BroadcastCurrentPosition()
		go p.Move()
	}

	g.startCountdown()
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
				g.started = true
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
	rotationChannel := make(chan RotationData)
	var player Player
	if conn != nil {
		human := &Human{PlayerData{id, -1, game, send, &currentPosition, rotationChannel, nil, true}, conn}
		player = human
	} else {
		bot := &Bot{PlayerData{id, -1, game, send, &currentPosition, rotationChannel, nil, true}}
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
