package main

import (
	"encoding/json"
	"fmt"
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
			curX, _ := player.currentPosition.Load("x")
			curY, _ := player.currentPosition.Load("y")
			player.send <- g.toJSON()
			g.board.fields[curX.(int)][curY.(int)].setUsed(player)
		case player := <-g.endGame:
			for id, p := range g.players {
				if player.id != id {
					g.winner = p
					g.sendToAll(g.toJSON())
				}
			}
			for id, p := range g.players {
				if p.rotationTicker != nil {
					p.stopRotation()
				}
				delete(g.players, id)
				close(p.send)
			}
			delete(activeGames, g.id)
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
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	startTime := ""

	for t := range ticker.C {
		if g.winner != nil {
			break
		}

		if startTime == "" {
			startTime = t.String()
			g.broadcast <- []byte(fmt.Sprintf("game started at %s", startTime))
		}

		if len(g.players) < 2 {
			return
		}
		movePlayers(g)
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

func movePlayers(g *Game) {
	for _, p := range g.players {
		p.move()
	}

	g.broadcast <- g.toJSON()
}

func (g *Game) toJSON() []byte {
	temp := make(map[string]interface{})
	playersStatusMap := make(map[int]interface{})
	for id, p := range g.players {
		playersStatusMap[id] = syncMapToMap(p.currentPosition)
	}
	temp["players"] = playersStatusMap

	if g.winner != nil {
		temp["winner"] = g.winner.id
	}

	res, err := json.Marshal(&temp)
	if err != nil {
		panic(err)
	}

	return res
}
