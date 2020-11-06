package board

//CanPlace determines if the spot on the board can have a piece
func (b *Board) CanPlace(x int, y int, block TetrisBlock) (canPlace bool) {
	pattern := block.Pattern
	//Given this block's pattern can we place it on the board?
	for w := 0; w < len(pattern); w++ {
		for h := 0; h < len(pattern[w]); h++ {
			// Out of bounds checking - shape overlaying on the board
			if x+w >= len(b.Bits) {
				return false
			}
			if y+h >= len(b.Bits[x+w]) {
				return false
			}
			// If the space is already occupied we cannot place there
			if b.Bits[x+w][y+h].occupied == true && pattern[w][h] == true {
				return false
			}
		}
	}
	return true
}

// PlaceBlock anchors a piece to the board a spot x,y given the blocks pattern
func (b *Board) PlaceBlock(x int, y int, block TetrisBlock) (canFit bool) {
	if !b.CanPlace(x, y, block) {
		return false
	}

	pattern := block.Pattern
	// If we are with-in boundaries, apply the OccupyPatten to the Board
	for w := 0; w < len(pattern); w++ {
		for h := 0; h < len(pattern[w]); h++ {
			// We set the Board bit to the block, and occupied pattern
			b.set(x+w, y+h, pattern[w][h], block)
		}
	}

	return true
}

// MoveBlock takes block at x,y and tries to move it to x1,y1
func (b *Board) MoveBlock(x, y, x1, y1 int) (wasMoved bool) {

	// 1. Boudary checking for x,y
	// TODO: Add more boundary checks for x,y >=0 and x1+patwidth<b.width...
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
	if patwidth == 0 {
		return false
	}
	if x+patwidth >= b.Width {
		return false
	}

	patheight := len(pattern[0])
	if patheight == 0 {
		return false
	}
	if x+patheight >= b.Height {
		return false
	}

	// 2. Remove piece and try to place at x1,y1 - otherwise put back at x,y
	b.RemovePiece(x, y)
	wasMoved = b.PlaceBlock(x1, y1, block)
	if !wasMoved {
		b.PlaceBlock(x, y, block)
	}

	return wasMoved
}

// DropToBottom looksup the block at x,y and tries to 'move it' to the bottom of the board.
// In Tetris this is when you push the 'down' arrow on the currently moving block.
func (b *Board) DropToBottom(x, y int) (wasMoved bool) {
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
				return false
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

	return true
}

// RemovePiece looksup the piece at x,y and unsets each bit of the pattern
func (b *Board) RemovePiece(x, y int) (delete bool) {
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
			b.unset(x+w, y+h, pattern[w][h])
		}
	}
	return true
}

// RotatePiece takes the active piece at x,y and rotates to the right Up->Right->Down->Left->Up->Right...
func (b *Board) RotatePiece(x, y int) (rotated bool) {
	//1. Boundary checks and block/pattern lookup
	if y >= b.Height || x >= b.Width {
		return false
	}
	bit := b.Bits[x][y]

	block := bit.block
	pattern := block.Pattern
	if pattern == nil {
		return false
	}

	// 2. Remove the piece from the board we are rotating
	b.RemovePiece(x, y)

	// 3. Rotate the block pattern
	block.Rotate()

	//4. Place the rotated block back on the board.
	return b.PlaceBlock(x, y, block)
}

// TetrisMatch replaces rows of all occupied with empty blocks and calls TetrisReduce
func (b *Board) TetrisMatch(onTetris func(row int)) {
next_row:
	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			//1. If any of the bits on the row aren't 1, then skip row.
			if b.Bits[x][y].occupied == false {
				continue next_row
			}
		}

		//2. All set, so clear Clear the row of all occupied
		for x := 0; x < b.Width; x++ {
			b.unset(x, y, true)
		}

		// All the bits are occupied, so Tetris!
		if onTetris != nil {
			onTetris(y)
		}

	}
}

// TetrisReduce implements classic Tetris rule of 'all zero rows replace with the row above'
func (b *Board) TetrisReduce(startrow int) {
skiprow:
	for h := startrow; h < b.Height; h++ {
		//1. Check each bit in the row
		for w := 0; w < b.Width; w++ {
			// If any bit is occupied the row cannot reduce
			if b.Bits[w][h].occupied == true {
				continue skiprow
			}
		}

		// The h row is all unoccupied and we work to row 0 zero copying rows
		for s := h; s > 0; s-- {
			for w := 0; w < b.Width; w++ {
				//1. Assign current row the bit from the row 'above' (closer to zero)
				b.Bits[w][s] = b.Bits[w][s-1]
				//2. Unoccupy the bit and remove block reference
				b.Bits[w][s-1].occupied = false
				b.Bits[w][s-1].block = TetrisBlock{}
			}

		}
	}
	return
}
