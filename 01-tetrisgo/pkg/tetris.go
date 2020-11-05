package pkg

// TetrisGame using pkg/board/ primatives to play tetris
type TetrisGame struct {
	Seed int // Random seed for piece generation
}

// NewClassicTetris provides Tetris board/blocks/gameplay/timer etc. based on classic arcade Tetris.
func NewClassicTetris(seed int) TetrisGame {
	tg := TetrisGame{}
	tg.Seed = seed
	return tg
}
