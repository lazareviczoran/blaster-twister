package main

import (
	"fmt"
	"math/rand"
	"time"
)

const height int = 400
const width int = 300

// Field represents a pixel of the arena
type Field struct {
	isUsed bool
}

// Board is a model of the arena
type Board struct {
	fields [][]Field
}

func initBoard(height, width int) *Board {
	var fields = make([][]Field, height)
	for i := 0; i < height; i++ {
		fields[i] = make([]Field, width)
		for j := 0; j < width; j++ {
			fields[i][j] = Field{isUsed: false}
		}
	}
	return &Board{fields: fields}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 2; i++ {
		fmt.Printf("Player%d's staring rotation is:%d\n", i, getStartRotation())
	}
	var board = initBoard(height, width)
	board.fields[30][50].isUsed = true

	fmt.Printf("This is the value of the 30th row and 50th column: %+v\n", &board.fields[30][50])
}

func getStartRotation() int {
	return rand.Intn(360)
}
