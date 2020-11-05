package board

import (
	"fmt"
	"testing"
)

func TestBuildBoard(t *testing.T) {
	b := NewBoard(8, 5)
	if b.String() != "00000000\n00000000\n00000000\n00000000\n00000000\n" {
		t.Fatalf("Board not full of zeros when.")
	}
}

func TestPlaceElleShape(t *testing.T) {
	b := NewBoard(8, 5)

	b.PlaceBlock(0, 0, b.ElleShape(Green))
	if b.String() != "LLL00000\n00L00000\n00000000\n00000000\n00000000\n" {
		t.Fatalf("Failed to place ElleShape at spot 0,0")
	}

	b.RotatePiece(0, 0)
	if b.String() != "0L000000\n0L000000\nLL000000\n00000000\n00000000\n" {
		t.Fatalf("Failed to Up->Right rotate ElleShape at spot 0,0")
	}

	b.RotatePiece(0, 0)
	if b.String() != "L0000000\nLLL00000\n00000000\n00000000\n00000000\n" {
		t.Fatalf("Failed to Right->Down rotate ElleShape at spot 0,0\n%s", b)
	}

	b.RotatePiece(0, 0)
	if b.String() != "LL000000\nL0000000\nL0000000\n00000000\n00000000\n" {
		t.Fatalf("Failed to Down->Left rotate ElleShape at spot 0,0\n%s", b)
	}

	b.RotatePiece(0, 0)
	if b.String() != "LLL00000\n00L00000\n00000000\n00000000\n00000000\n" {
		t.Fatalf("Failed to Left->Up rotate ElleShape at spot 0,0\n%s", b)
	}

	b.DropToBottom(0, 0)
	if b.String() != "00000000\n00000000\n00000000\nLLL00000\n00L00000\n" {
		t.Fatalf("Failed to drop ElleShape at spot 0,0 to bottom\n%s", b)
	}
	fmt.Printf("%s\n", b)
}

func TestPlacePipe(t *testing.T) {
	b := NewBoard(8, 5)

	b.PlaceBlock(0, 0, b.PipeShape(Green))
	if b.String() != "P0000000\nP0000000\nP0000000\nP0000000\n00000000\n" {
		t.Fatalf("Failed to place PipeShape at spot 0,0:\n%s", b)
	}

	b.DropToBottom(0, 0)
	if b.String() != "00000000\nP0000000\nP0000000\nP0000000\nP0000000\n" {
		t.Fatalf("Failed to place PipeShape at spot 0,0:\n%s", b)
	}

	b.PlaceBlock(7, 0, b.PipeShape(Green))
	if b.String() != "0000000P\nP000000P\nP000000P\nP000000P\nP0000000\n" {
		t.Fatalf("Failed to place PipeShape at spot 7,0:\n%s", b)
	}

	b.DropToBottom(7, 0)
	if b.String() != "00000000\nP000000P\nP000000P\nP000000P\nP000000P\n" {
		t.Fatalf("Failed to drop PipeShape at spot 7,0:\n%s", b)
	}

	b.PlaceBlock(6, 0, b.PipeShape(Green))
	if b.String() != "000000P0\nP00000PP\nP00000PP\nP00000PP\nP000000P\n" {
		t.Fatalf("Failed to place PipeShape at spot 6,0:\n%s", b)
	}

	b.DropToBottom(6, 0)
	if b.String() != "00000000\nP00000PP\nP00000PP\nP00000PP\nP00000PP\n" {
		t.Fatalf("Failed to drop PipeShape at spot 6,0:\n%s", b)
	}

	b.PlaceBlock(5, 0, b.PipeShape(Green))
	b.DropToBottom(5, 0)
	if b.String() != "00000000\nP0000PPP\nP0000PPP\nP0000PPP\nP0000PPP\n" {
		t.Fatalf("Failed to drop PipeShape at spot 5,0:\n%s", b)
	}

	b.PlaceBlock(1, 0, b.PipeShape(Green))
	b.RotatePiece(1, 0)
	b.DropToBottom(1, 0)
	if b.String() != "00000000\nP0000PPP\nP0000PPP\nP0000PPP\nPPPPPPPP\n" {
		t.Fatalf("Failed to drop PipeShape at spot 1,0:\n%s", b)
	}

	b.TetrisMatch(b.TetrisReduce)
	if b.String() != "00000000\n00000000\nP0000PPP\nP0000PPP\nP0000PPP\n" {
		t.Fatalf("Failed to Tetris match and reduce on board:\n%s", b)
	}

	b1 := b.copy()
	if b1.String() != "00000000\n00000000\nP0000PPP\nP0000PPP\nP0000PPP\n" {
		t.Fatalf("Failed to Tetris copy board:\n%s", b)
	}

	fmt.Printf("%s\n", b)
}

func TestMoveSquare(t *testing.T) {
	b := NewBoard(8, 5)

	b.PlaceBlock(0, 0, b.SquareShape(Red))
	if b.String() != "BB000000\nBB000000\n00000000\n00000000\n00000000\n" {
		t.Fatalf("Failed to place PipeShape at spot 0,0:\n%s", b)
	}

	b.MoveBlock(0, 0, 1, 0)
	if b.String() != "0BB00000\n0BB00000\n00000000\n00000000\n00000000\n" {
		t.Fatalf("Failed to place PipeShape at spot 0,0:\n%s", b)
	}

	wasMoved := b.MoveBlock(1, 0, 7, 0)
	if wasMoved {
		t.Fatalf("Failed placed square out of bounds:\n%s", b)
	}
	wasMoved = b.MoveBlock(1, 0, 8, 0)
	if wasMoved {
		t.Fatalf("Failed placed square out of bounds:\n%s", b)
	}

	wasMoved = b.MoveBlock(1, 0, 6, 0)
	if !wasMoved || b.String() != "000000BB\n000000BB\n00000000\n00000000\n00000000\n" {
		t.Fatalf("Failed placed square out of bounds:\n%s", b)
	}

	b.DropToBottom(6, 0)
	if b.String() != "00000000\n00000000\n00000000\n000000BB\n000000BB\n" {
		t.Fatalf("Failed placed square out of bounds:\n%s", b)
	}

	wasMoved = b.DropToBottom(6, 0)
	if wasMoved {
		t.Fatalf("Dropped a piece that wasn't there:\n%s", b)
	}

	fmt.Printf("%s\n", b)
}
