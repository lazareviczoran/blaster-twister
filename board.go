package main

// Field represents a pixel of the arena
type Field struct {
	player *Player
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

func (b *Board) isValidMove(x, y int) bool {
	if x < 0 || x >= width || y < 0 || y >= height {
		return false
	}
	return b.fields[x][y].player == nil
}

func (f *Field) setUsed(p *Player) {
	f.player = p
}
