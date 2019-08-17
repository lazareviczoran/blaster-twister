package main

// Field represents a pixel of the arena
type Field struct {
	player Player
}

// Board is a model of the arena
type Board struct {
	fields [][]Field
}

func initBoard(height, width int) *Board {
	fields := make([][]Field, width)
	for i := 0; i < width; i++ {
		fields[i] = make([]Field, height)
		for j := 0; j < height; j++ {
			fields[i][j] = Field{player: nil}
		}
	}
	return &Board{fields: fields}
}

func (b *Board) isValidMove(fromX, fromY, toX, toY int) bool {
	if toX < 0 || toX >= width || toY < 0 || toY >= height {
		return false
	}

	if toX != fromX &&
		toY != fromY &&
		b.fields[toX][fromY].player != nil &&
		b.fields[fromX][toY].player != nil {
		return false
	}

	return b.fields[toX][toY].player == nil
}

func (f *Field) setUsed(p Player) {
	f.player = p
}
