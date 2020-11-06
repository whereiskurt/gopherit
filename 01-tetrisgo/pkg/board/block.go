package board

// TetrisBlock holds bit pattern and orientation. The pattern and orientation are related.
type TetrisBlock struct {
	Label       string
	Type        BlockType
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

// Elle is Tetris 2x2 block square
const Elle BlockType = "Elle"

// IElle is Tetris block shape
const IElle BlockType = "InvertedElle"

// Rotate takes the block right Up->Right->Down->Left->Up...
func (blk *TetrisBlock) Rotate() {
	pattern := blk.Pattern

	// 3. Rotate the pattern for the shape and set orientation
	patwidth := len(pattern)
	patheight := len(pattern[0])
	rpat := make([][]bool, patheight)
	for h := 0; h < patheight; h++ {
		rpat[h] = make([]bool, patwidth)
	}

	for w := 0; w < patwidth; w++ {
		for h := 0; h < patheight; h++ {
			rpat[patheight-h-1][w] = pattern[w][h]
		}
	}
	blk.Pattern = rpat

	switch blk.Orientation {
	case Up:
		blk.Orientation = Right
	case Right:
		blk.Orientation = Down
	case Down:
		blk.Orientation = Left
	case Left:
		blk.Orientation = Up
	}

}

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

	return pattern
}

// BlockOrientation describes how the block is applied to the board
type BlockOrientation string

// Up is the default block orientation
const Up BlockOrientation = "Up"

// Down is a block orientation
const Down BlockOrientation = "Right"

// Left is a block orientation
const Left BlockOrientation = "Down"

// Right is a block orientation
const Right BlockOrientation = "Left"

// BlockColour describes what colour the block has on the board
type BlockColour string

// Red is a block colour
const Red BlockColour = "Red"

// Blue is a block colour
const Blue BlockColour = "Blue"

// Green is a block colour
const Green BlockColour = "Green"

// Orange is a block colour
const Orange BlockColour = "Orange"

// Purple is a block colour
const Purple BlockColour = "Purple"

// MakeTetrisBlock creates a TetrisBlock
func (b *Board) MakeTetrisBlock(label string, t BlockType, c BlockColour, o BlockOrientation) TetrisBlock {
	return TetrisBlock{
		Label:       label,
		Type:        t,
		Orientation: o,
		Colour:      c,
		Pattern:     makeBlockPattern(t, Up)}
}

// TeeShape builds a default T Tetris block
func (b *Board) TeeShape(c BlockColour) (blk TetrisBlock) {
	blk = b.MakeTetrisBlock("T", Tee, c, Up)
	return blk
}

// ElleShape builds a default L Tetris block
func (b *Board) ElleShape(c BlockColour) (blk TetrisBlock) {
	blk = b.MakeTetrisBlock("L", Elle, c, Up)
	return blk
}

// IElleShape builds a default inverted-L Tetris block
func (b *Board) IElleShape(c BlockColour) (blk TetrisBlock) {
	blk = b.MakeTetrisBlock("I", IElle, c, Up)
	return blk
}

// PipeShape builds a default | Tetris block
func (b *Board) PipeShape(c BlockColour) (blk TetrisBlock) {
	blk = b.MakeTetrisBlock("P", Pipe, c, Up)
	return blk
}

// SquareShape builds a default [] Tetris block
func (b *Board) SquareShape(c BlockColour) (blk TetrisBlock) {
	blk = b.MakeTetrisBlock("B", Square, c, Up)
	return blk
}
