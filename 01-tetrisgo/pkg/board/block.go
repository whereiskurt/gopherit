package board

// TetrisBlock holds bit pattern and orientation. The pattern and orientation are related.
type TetrisBlock struct {
	Label       string
	Orientation BlockOrientation
	Colour      BlockColour
	Pattern     [][]bool
}

// BlockType is one of the 5x possible
type BlockType string

// Square is 4x4 Tetris block shape
const Square BlockType = "Square"

// Pipe is 4x1 Tetris block shape
const Pipe BlockType = "Pipe"

// Tee is Tetris block shape
const Tee BlockType = "Tee"

// Elle is Tetris block shape
const Elle BlockType = "Elle"

// IElle is Tetris block shape
const IElle BlockType = "InvertedElle"

func makeBlockPattern(t BlockType, o BlockOrientation) (pattern [][]bool) {
	//Create a default "UP" and then rotate if not o==Up
	switch t {
	case Square:
		pattern = make([][]bool, 2)
		for i := range pattern {
			pattern[i] = make([]bool, 2)
		}
		pattern[0][0] = true
		pattern[0][1] = true
		pattern[1][0] = true
		pattern[1][1] = true
	case Pipe:
		pattern = make([][]bool, 1)
		for i := range pattern {
			pattern[i] = make([]bool, 4)
		}
		pattern[0][0] = true
		pattern[0][1] = true
		pattern[0][2] = true
		pattern[0][3] = true
	case Tee:
		pattern = make([][]bool, 3)
		for i := range pattern {
			pattern[i] = make([]bool, 2)
		}
		pattern[0][0] = true
		pattern[1][0] = true
		pattern[1][1] = true
		pattern[2][0] = true
	case Elle:
		pattern = make([][]bool, 3)
		for i := range pattern {
			pattern[i] = make([]bool, 2)
		}
		pattern[0][0] = true
		pattern[1][0] = true
		pattern[2][0] = true
		pattern[2][1] = true
	case IElle:
		pattern = make([][]bool, 3)
		for i := range pattern {
			pattern[i] = make([]bool, 2)
		}
		pattern[0][0] = true
		pattern[0][1] = true
		pattern[1][0] = true
		pattern[2][0] = true
	}

	switch o {
	case Up:
		break
	}

	return pattern
}

// BlockOrientation describes how the block is applied to the board
type BlockOrientation int

// Up is the default block orientation
const Up BlockOrientation = 6

// Down is a block orientation
const Down BlockOrientation = 1

// Left is a block orientation
const Left BlockOrientation = 2

// Right is a block orientation
const Right BlockOrientation = 3

// BlockColour describes what colour the block has on the board
type BlockColour int

// Red is a block colour
const Red BlockColour = 1

// Blue is a block colour
const Blue BlockColour = 2

// Green is a block colour
const Green BlockColour = 3

// Orange is a block colour
const Orange BlockColour = 4

// Purple is a block colour
const Purple BlockColour = 5

// makeTetrisBlock creates a TetrisBlock
func (b *Board) makeTetrisBlock(label string, t BlockType, colour BlockColour, o BlockOrientation) TetrisBlock {
	return TetrisBlock{
		Label:       label,
		Orientation: o,
		Colour:      colour,
		Pattern:     makeBlockPattern(t, Up)}
}

// MakeTee builds a default T Tetris block
func (b *Board) MakeTee(colour BlockColour, o BlockOrientation) TetrisBlock {
	block := b.makeTetrisBlock("T", Tee, colour, o)
	return block
}

// MakeElle builds a default L Tetris block
func (b *Board) MakeElle(colour BlockColour, o BlockOrientation) TetrisBlock {
	block := b.makeTetrisBlock("L", Elle, colour, o)
	return block
}

// MakeIElle builds a default inverted-L Tetris block
func (b *Board) MakeIElle(colour BlockColour, o BlockOrientation) TetrisBlock {
	block := b.makeTetrisBlock("I", IElle, colour, o)
	return block
}

// MakePipe builds a default | Tetris block
func (b *Board) MakePipe(colour BlockColour, o BlockOrientation) TetrisBlock {
	block := b.makeTetrisBlock("P", Pipe, colour, o)
	return block
}

// MakeSquare builds a default [] Tetris block
func (b *Board) MakeSquare(colour BlockColour, o BlockOrientation) TetrisBlock {
	block := b.makeTetrisBlock("B", Square, colour, o)
	return block
}
