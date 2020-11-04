package board

// PlaceBlock anchors a piece to the board a spot x,y given the blocks pattern
func (b *Board) PlaceBlock(x int, y int, block TetrisBlock) (canFit bool) {
	pattern := block.Pattern
	for w := 0; w < len(pattern); w++ {
		for h := 0; h < len(pattern[w]); h++ {
			// Out of bounds checking - shape overlaying on the board
			if x+w >= len(b.Bits) {
				return false
			}
			if y+h >= len(b.Bits[x+w]) {
				return false
			}
			// If the space is already occupied we cannot Place here
			if b.Bits[x+w][y+h].occupied == true && pattern[w][h] == true {
				return false
			}
		}
	}

	// If we are with-in boundaries, apply the OccupyPatten to the Board
	for w := 0; w < len(pattern); w++ {
		for h := 0; h < len(pattern[w]); h++ {
			// We set the Board bit to the block, and occupied pattern
			b.set(x+w, y+h, pattern[w][h], block)
		}
	}

	return true
}

// TetrisReduce implements classic Tetris rule of 'all zero rows replace with the row above'
func (b *Board) TetrisReduce(startrow int) {
next_row:
	for h := startrow; h < b.Height; h++ {
		for w := 0; w < b.Width; w++ {
			occupied := b.Bits[w][h].occupied
			// Classic Tetris if any bit is occupied the row cannot reduce
			if occupied == true {
				continue next_row
			}
		}

		// The h row is all unoccupied and we work to row 0 zero copying rows
		for s := h; s > 0; s-- {
			for w := 0; w < b.Width; w++ {
				b.Bits[w][s] = b.Bits[w][s-1]
				b.Bits[w][s-1].occupied = false
				b.Bits[w][s-1].block = TetrisBlock{}
			}

		}
	}
	return
}

// DropToBottom looksup the block at x,y and tries to 'move it' to the bottom of the board.
// In Tetris this is when you push the 'down' arrow on the currently moving block.
func (b *Board) DropToBottom(x, y int) {
	if y >= b.Height || x >= b.Width {
		return
	}
	bit := b.Bits[x][y]

	block := bit.block
	pattern := block.Pattern
	if pattern == nil {
		return
	}

	patwidth := len(pattern)
	patheight := len(pattern[0])

	for w := 0; w < patwidth; w++ {
		for h := 0; h < patheight; h++ {
			if pattern[w][h] == true {
				//A bit of this block is occupy the space below, so we could move down.
				if h+1 < patheight && pattern[w][h+1] == true {
					continue
				}

				//The bit below is not occupied
				if h+1+y < b.Height && b.Bits[x+w][h+1+y].occupied == false {
					continue
				}

				//The bit below is occupied and our pattern has bit that needs the spot
				return
			}
		}
	}

	//Remove the block from the board
	for w := 0; w < patwidth; w++ {
		for h := 0; h < patheight; h++ {
			b.unset(x+w, y+h, pattern[w][h])
		}
	}

	//Move the block down one row, which we know is inbounds from above
	b.PlaceBlock(x, y+1, block)

	// RECURSE!! :-)
	b.DropToBottom(x, y+1)

	return
}

// RotatePiece takes the active piece at x,y and transposes/mirrors as appropriate
func (b *Board) RotatePiece(x, y int) (rotated bool) {
	if y >= b.Height || x >= b.Width {
		return false
	}
	bit := b.Bits[x][y]

	block := bit.block
	pattern := block.Pattern
	if pattern == nil {
		return false
	}

	patwidth := len(pattern)
	patheight := len(pattern[0])

	// Remove from board
	for w := 0; w < patwidth; w++ {
		for h := 0; h < patheight; h++ {
			//if pattern[w][h] == true {
			b.unset(x+w, y+h, pattern[w][h])
			//}
		}
	}

	rpat := make([][]bool, patheight)
	for h := 0; h < patheight; h++ {
		rpat[h] = make([]bool, patwidth)
	}

	for w := 0; w < patwidth; w++ {
		for h := 0; h < patheight; h++ {
			rpat[patheight-h-1][w] = pattern[w][h]
		}
	}
	block.Pattern = rpat

	switch block.Orientation {
	case Up:
		block.Orientation = Right
	case Right:
		block.Orientation = Down
	case Down:
		block.Orientation = Left
	case Left:
		block.Orientation = Up
	}

	return b.PlaceBlock(x, y, block)
}

// TetrisMatch replaces rows of all occupied with empty blocks and calls TetrisReduce
func (b *Board) TetrisMatch(onTetris func(row int)) {
next_row:
	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			//If any of the bits on the row aren't set, skip row.
			if b.Bits[x][y].occupied == false {
				continue next_row
			}
		}
		// All the bits are occupied, so Tetris!
		if onTetris != nil {
			onTetris(y)
		}

		// Clear the row of all occupied
		for x := 0; x < b.Width; x++ {
			b.unset(x, y, true)
		}
		b.TetrisReduce(y)

	}
}
