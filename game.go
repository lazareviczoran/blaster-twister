package main

import (
	"fmt"
	"sort"
	"time"
)

// Game holds the connections to the players
type Game struct {
	id        string
	players   map[*Player]bool
	register  chan *Player
	endGame   chan *Player
	broadcast chan []byte
	board     *Board
}

func (g *Game) run() {
	for {
		select {
		case player := <-g.register:
			g.players[player] = true
			currPos := player.currentPosition
			player.game.board.fields[currPos["x"]][currPos["y"]].isUsed = true
		case player := <-g.endGame:
			for p := range g.players {
				if player.id != p.id {
					g.sendToAll([]byte(fmt.Sprintf("{winner: %d}", p.id)))
				}
			}
			for p := range g.players {
				if _, ok := g.players[p]; ok {
					if p.rotationTicker != nil {
						p.stopRotation()
					}
					delete(g.players, p)
					close(p.send)
				}
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
		players:   make(map[*Player]bool),
		board:     initBoard(height, width),
	}
}

func (g *Game) startGame() {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	startTime := ""

	for t := range ticker.C {
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
	for player := range g.players {
		select {
		case player.send <- message:
		default:
			close(player.send)
			delete(g.players, player)
		}
	}
}

func movePlayers(g *Game) {
	players := make([]*Player, 0)
	for p := range g.players {
		players = append(players, p)
	}

	sort.SliceStable(players, func(i, j int) bool { return players[i].id < players[j].id })
	for _, p := range players {
		p.move()
	}

	g.broadcast <- []byte(
		fmt.Sprintf("{p0:{x:%d, y:%d}, p1:{x:%d, y:%d}}",
			players[0].currentPosition["x"],
			players[0].currentPosition["y"],
			players[1].currentPosition["x"],
			players[1].currentPosition["y"],
		),
	)
}
