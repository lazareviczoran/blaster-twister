package main

import (
	"encoding/json"
	"log"
	"time"
)

// Game holds the connections to the players
type Game struct {
	id        string
	players   map[int]*Player
	register  chan *Player
	endGame   chan *Player
	broadcast chan []byte
	board     *Board
	winner    *Player
}

func (g *Game) run() {
	for {
		select {
		case player := <-g.register:
			g.players[player.id] = player
			startX := width / 2
			startY := height/5 + player.id*3*height/5
			player.currentPosition.Store("x", startX)
			player.currentPosition.Store("y", startY)
			player.currentPosition.Store("rotation", getStartRotation())
			player.currentPosition.Store("trace", false)
			player.game.board.fields[startX][startY].setUsed(player)
			player.broadcastCurrentPosition()
			if len(g.players) == 2 {
				g.startGame()
			}
		case player := <-g.endGame:
			var winner *Player
			for id, p := range g.players {
				if p.alive && player.id != id {
					if winner != nil {
						// we have more than one players alive
						return
					}
					winner = p
				}
			}
			if winner != nil {
				g.winner = winner
				temp := make(map[string]interface{})
				temp["winner"] = winner.id
				res, err := json.Marshal(&temp)
				if err != nil {
					log.Printf("Could not convert to JSON")
				}
				for id, p := range g.players {
					p.send <- res
					if p.rotationTicker != nil {
						p.stopRotation()
					}
					delete(g.players, id)
					close(p.send)
				}
				delete(activeGames, g.id)
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
		register:  make(chan *Player),
		endGame:   make(chan *Player),
		players:   make(map[int]*Player),
		board:     initBoard(height, width),
		winner:    nil,
	}
}

func (g *Game) startGame() {
	startTime := time.Now()
	// g.broadcast <- []byte(fmt.Sprintf("game started at %s", startTime))
	log.Printf("game started at %v", startTime)

	for _, p := range g.players {
		go p.move()
	}
}

func (g *Game) sendToAll(message []byte) {
	for id, player := range g.players {
		select {
		case player.send <- message:
		default:
			close(player.send)
			delete(g.players, id)
		}
	}
}
