package pkg

// TetrisGame using pkg/board/ primatives to play tetris
type TetrisGame struct {
	Seed int // Random seed for piece generation
}

// NewTetrisGame returns a TetrisGame to play
func NewTetrisGame() TetrisGame {
	tg := TetrisGame{}
	return tg
}
