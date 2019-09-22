package main

import (
	"math"
	"sync"
	"time"
)

const (
	fps = 15

	directionLeft  = "left"
	directionRight = "right"
	directionUp    = "up"
	directionDown  = "down"
)

// Player contains the info about the current position and actions for moving
type Player interface {
	InitPlayer()
	Move()
	StartRotation(direction string)
	StopRotation()
	ID() int
	ClientID() int
	Game() *Game
	CurrentPosition() *sync.Map
	IsAlive() bool
	SetAlive(alive bool)
	BroadcastCurrentPosition()
	Broadcast(message []byte)
	Destroy()
}

// RotationData struct is used to send rotation data through a channel
type RotationData struct {
	dir string
	key string
}

// PlayerData contains the info about player and the players position
type PlayerData struct {
	id              int
	clientID        int
	game            *Game
	send            chan []byte
	currentPosition *sync.Map
	rotationChannel chan RotationData
	rotationTicker  *time.Ticker
	alive           bool
}

func moveBresenham(p Player, x0 int, y0 int, rotationRad float64) {
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
		fromX := x0
		fromY := y0
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
		if p.Game().board.isValidMove(fromX, fromY, x0, y0) {
			p.CurrentPosition().Store("x", x0)
			p.CurrentPosition().Store("y", y0)
			trace, _ := p.CurrentPosition().Load("trace")
			if trace.(bool) {
				p.Game().board.fields[x0][y0].setUsed(p)
			}
			if p.IsAlive() && p.Game().winner == nil {
				p.BroadcastCurrentPosition()
			}
		} else {
			p.Game().endGame <- p
			return
		}
	}
}
