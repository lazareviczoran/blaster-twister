package main

import (
	"fmt"
	"log"
	"sort"
	"time"
)

// Game holds the connections to the players
type Game struct {
	players    map[*Player]bool
	register   chan *Player
	unregister chan *Player
	broadcast  chan []byte
	board      *Board
}

func (g *Game) run() {
	for {
		select {
		case player := <-g.register:
			g.players[player] = true
			currPos := player.currentPosition
			player.game.board.fields[currPos["x"]][currPos["y"]].isUsed = true
		case player := <-g.unregister:
			log.Printf("unregister")
			if _, ok := g.players[player]; ok {
				delete(g.players, player)
				close(player.send)
			}
		case message := <-g.broadcast:
			for player := range g.players {
				select {
				case player.send <- message:
				default:
					close(player.send)
					delete(g.players, player)
				}
			}
		}
	}
}

func newGame(height, width int) *Game {
	return &Game{
		broadcast:  make(chan []byte),
		register:   make(chan *Player),
		unregister: make(chan *Player),
		players:    make(map[*Player]bool),
		board:      initBoard(height, width),
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
		movePlayers(g)
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
		fmt.Sprintf("new positions p1{%d:%d}, p2{%d:%d}",
			players[0].currentPosition["x"],
			players[0].currentPosition["y"],
			players[1].currentPosition["x"],
			players[1].currentPosition["y"],
		),
	)
}
