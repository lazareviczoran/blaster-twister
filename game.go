package main

import (
	"fmt"
	"log"
	"time"
)

// Game holds the connections to the players
type Game struct {
	players    map[*Player]bool
	register   chan *Player
	unregister chan *Player
	broadcast  chan []byte
}

func (g *Game) run() {
	for {
		select {
		case player := <-g.register:
			g.players[player] = true
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

func newGame() *Game {
	return &Game{
		broadcast:  make(chan []byte),
		register:   make(chan *Player),
		unregister: make(chan *Player),
		players:    make(map[*Player]bool),
	}
}

func (g *Game) startGame() {
	log.Printf("startGame")
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	startTime := ""

	for t := range ticker.C {
		if startTime == "" {
			startTime = t.String()
			g.broadcast <- []byte(fmt.Sprintf("game started at %s", startTime))
		}
		random := getStartRotation()
		g.broadcast <- []byte(fmt.Sprintf("%d sent at %s", random, t.String()))
	}
}
