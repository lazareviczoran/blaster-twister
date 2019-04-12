package main

// Field represents a pixel of the arena
type Field struct {
	isUsed bool
}

// Board is a model of the arena
type Board struct {
	fields [][]Field
}

func initBoard(height, width int) *Board {
	fields := make([][]Field, height)
	for i := 0; i < height; i++ {
		fields[i] = make([]Field, width)
		for j := 0; j < width; j++ {
			fields[i][j] = Field{isUsed: false}
		}
	}
	return &Board{fields: fields}
}
