package main

import (
	"encoding/json"
	"log"
	"math"
	"sync"
	"time"
)

// Bot represents the computer player
type Bot struct {
	PlayerData
}

type intersection struct {
	distance int
	angle    int
}

// ID returns the players Id
func (b *Bot) ID() int {
	return b.id
}

// ClientID returns the id obtained from the WS client,
// in this case we return -1 since the bot is not a WS client
func (b *Bot) ClientID() int {
	return b.clientID
}

// Game returns the pointer to Game
func (b *Bot) Game() *Game {
	return b.game
}

// CurrentPosition returns a map that contains
// info about the human players current position
func (b *Bot) CurrentPosition() *sync.Map {
	return b.currentPosition
}

// IsAlive returns the players alive status
func (b *Bot) IsAlive() bool {
	return b.alive
}

// SetAlive sets the players alive status
func (b *Bot) SetAlive(alive bool) {
	b.alive = alive
}

// Broadcast sends the message to writePump, which eventually sends it
// to the websocket client
func (b *Bot) Broadcast(message []byte) {
	b.send <- message
}

// Destroy closes all channels and removes player from the game
func (b *Bot) Destroy() {
	b.StopRotation()
	close(b.send)
	delete(b.game.players, b.id)
}

// InitPlayer initializes the players position
func (b *Bot) InitPlayer() {
	go func() {
		for {
			select {
			case <-b.send:
				// listen to bots send
				if !b.IsAlive() || b.game.winner != nil {
					return
				}
			}
		}
	}()

	startX := width / 2
	startY := height/5 + b.id*3*height/5
	b.currentPosition.Store("x", startX)
	b.currentPosition.Store("y", startY)
	b.currentPosition.Store("rotation", getStartRotation())
	b.currentPosition.Store("trace", false)
	b.game.board.fields[startX][startY].setUsed(b)
}

// StartRotation creates a new ticker and updates
// the players rotation angle on each tick
func (b *Bot) StartRotation(direction string) {
	go func() {
		b.rotationTicker = time.NewTicker(30 * time.Millisecond)
		for {
			select {
			case <-b.rotationTicker.C:
				currRotation, _ := b.currentPosition.Load("rotation")

				if direction == directionRight {
					b.currentPosition.Store("rotation", (currRotation.(int)+5)%360)
				} else {
					b.currentPosition.Store("rotation", (currRotation.(int)+355)%360)
				}
				b.currentPosition.Store("rotationDir", direction)
			}
		}
	}()
}

// StopRotation stops the rotation ticker
func (b *Bot) StopRotation() {
	if b.rotationTicker != nil {
		b.rotationTicker.Stop()
		b.rotationTicker = nil
		b.currentPosition.Store("rotationDir", nil)
	}
}

// BroadcastCurrentPosition sends the players current position
// to all clients
func (b *Bot) BroadcastCurrentPosition() {
	temp := make(map[string]interface{})
	playersStatusMap := make(map[int]interface{})
	playersStatusMap[b.id] = syncMapToMap(b.currentPosition)
	temp["players"] = playersStatusMap

	res, err := json.Marshal(&temp)
	if err != nil {
		log.Printf("Could not convert to JSON, %v", err)
		return
	}

	b.game.sendToAll(res)
}

// Move sets the new position for the player calculated by
// the rotation angle on each tick.
func (b *Bot) Move() {
	mainTicker := time.NewTicker(1000 / fps * time.Millisecond)
	defer mainTicker.Stop()

	visitedTicker := createRandomIntervalTicker(1000, 2000)

	for {
		select {
		case <-mainTicker.C:
			if !b.alive || b.game.winner != nil {
				visitedTicker.Stop()
				return
			}
			curX, _ := b.currentPosition.Load("x")
			curY, _ := b.currentPosition.Load("y")
			curRotation, _ := b.currentPosition.Load("rotation")
			curRotationDir, _ := b.currentPosition.Load("rotationDir")

			angle := b.findAngleToFarthestIntersection(curX.(int), curY.(int))
			diff := angle - curRotation.(int)
			if diff > 0 {
				direction := directionRight
				if diff > 180 {
					direction = directionLeft
				}
				if curRotationDir != direction {
					b.StopRotation()
					b.StartRotation(direction)
				}
			} else if diff < 0 {
				direction := directionLeft
				if diff < -180 {
					direction = directionRight
				}
				if curRotationDir != direction {
					b.StopRotation()
					b.StartRotation(direction)
				}
			} else {
				b.StopRotation()
			}

			rotationRad := float64(curRotation.(int)) * math.Pi / 180
			moveBresenham(b, curX.(int), curY.(int), rotationRad)
		case <-visitedTicker.C:
			trace, _ := b.currentPosition.Load("trace")
			b.currentPosition.Store("trace", !trace.(bool))
			visitedTicker = createRandomIntervalTicker(1000, 2000)
		}
	}
}

func (b *Bot) findAngleToFarthestIntersection(x0, y0 int) int {
	channel := make(chan *intersection)
	for i := 0; i < 36; i++ {
		go b.getDistanceToWall(channel, x0, y0, i*10)
	}
	var farthestIntersection *intersection
	finished := 0
	for {
		select {
		case intersection := <-channel:
			if farthestIntersection == nil || farthestIntersection.distance < intersection.distance {
				farthestIntersection = intersection
			}
			finished++
			if finished < 36 {
				continue
			}
			return farthestIntersection.angle
		}
	}
}

func (b *Bot) getDistanceToWall(channel chan *intersection, x0, y0, rotationDeg int) {
	distance := 0
	rotationRad := float64(rotationDeg) * math.Pi / 180
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
	for {
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
		if !b.game.board.isValidMove(fromX, fromY, x0, y0) {
			channel <- &intersection{distance, rotationDeg}
			return
		}
		distance++
	}
}
