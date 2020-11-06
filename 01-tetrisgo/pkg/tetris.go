package pkg

import "tetrigo/pkg/board"

// TetrisGame using pkg/board/ primatives to play tetris
type TetrisGame struct {
	Board     *board.Board
	Seed      int64 // Random seed for piece generation
	Block     board.TetrisBlock
	BlockX    int
	BlockY    int
	NextBlock board.TetrisBlock
	State     GameState
	Score     int
}

// GameState tracks if
type GameState string

const (
	// Init is the when the board is init'd but the timer hasn't started
	Init GameState = "Init"
	// Start happens at the begging of the game before
	Start GameState = "Start"
	// TickBegin is the top of the pendulum tick
	TickBegin GameState = "TickBegin"
	// Tick is the middle of the tick
	Tick GameState = "Tick"
	// TickEnd happens at the end eg. looking tetris/reduce
	TickEnd GameState = "TickEnd"
	// End is called when no more pieces can be added
	End GameState = "End"
)

// NewClassicTetris provides Tetris board/blocks/gameplay/timer etc. based on classic arcade Tetris.
func NewClassicTetris(width int, height int, seed int64) TetrisGame {
	tg := TetrisGame{}
	tg.State = Init

	tg.Board = board.NewBoard(width, height, seed)
	tg.Seed = seed
	tg.NextBlock = tg.Board.RandomBlock()
	tg.AddNextBlock()

	return tg
}

// AddNextBlock advances NextBlock to Block
func (tg *TetrisGame) AddNextBlock() (wasAdded bool) {

	// 1. Calc middle of the row given block pattern size
	x := ((tg.Board.Width/2 - 1) - (len(tg.Block.Pattern)-1)/2)
	y := 0

	// 2. Check if we can place it on the board
	if !tg.Board.CanPlace(x, y, tg.NextBlock) {
		return false
	}

	// 3. Make the 'next block' the current block and generate next random block
	tg.Block = tg.NextBlock
	tg.NextBlock = tg.Board.RandomBlock()
	tg.BlockX = x
	tg.BlockY = y

	wasAdded = tg.Board.PlaceBlock(x, y, tg.Block)
	return wasAdded
}
