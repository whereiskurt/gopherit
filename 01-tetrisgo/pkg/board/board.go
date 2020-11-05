package board

import (
	"fmt"
)

// Board holds coords for blocks
type Board struct {
	Width  int
	Height int
	Bits   [][]BitState // Each bit of the board, occupied by a potential block.
	Stats
}

// Stats help determine the value/cost of the board
type Stats struct {
	Marked  int //How many are occupied?
	Trapped int //The amount of 'unreachable' blocks
}

// BitState maintains occupied state and associated block.
// If the board is 8 wide by 10 tall there will be 80 BlockStates stored
// in the Board coords
type BitState struct {
	occupied bool
	block    TetrisBlock
}

// NewBoard constructs a Tetris board state
func NewBoard(width int, height int) (b *Board) {
	b = new(Board)
	b.Height = height
	b.Width = width
	b.makeEmpty()
	return
}

// String outputs basic ASCII board
func (b *Board) String() (s string) {
	for h := 0; h < b.Height; h++ {
		for w := 0; w < b.Width; w++ {
			occupied := b.Bits[w][h].occupied
			if occupied == false {
				s += fmt.Sprintf("0")
			} else {
				s += fmt.Sprintf("%v", b.Bits[w][h].block.Label)
			}
		}
		s += fmt.Sprintf("\n")
	}

	return s
}

func (b *Board) makeEmpty() {
	b.Bits = make([][]BitState, b.Width)
	for i := range b.Bits {
		b.Bits[i] = make([]BitState, b.Height)
	}
	return
}

func (b *Board) copy() *Board {
	copy := NewBoard(b.Width, b.Height)
	for h := 0; h < b.Height; h++ {
		for w := 0; w < b.Width; w++ {
			occupied := b.Bits[w][h].occupied
			block := b.Bits[w][h].block
			if occupied == true {
				copy.set(w, h, occupied, block)
			}
		}
	}
	return copy
}

func (b *Board) set(w int, h int, isOccupied bool, block TetrisBlock) {
	b.Bits[w][h].occupied = isOccupied
	b.Bits[w][h].block = block
	if isOccupied {
		b.Marked++
	}
}

func (b *Board) unset(w int, h int, isOccupied bool) {
	b.Bits[w][h].occupied = false
	b.Bits[w][h].block = TetrisBlock{}
	if isOccupied {
		b.Marked--
	}
}
