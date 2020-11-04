package main

import (
	"fmt"

	"tetrigo/pkg/board"
)

func main() {
	b := board.NewBoard(8, 5)
	b.PlaceBlock(0, 0, b.MakeElle(board.Green, board.Up))
	fmt.Printf("%s(Marked:%d/%d)\n", b, b.Marked, b.Height*b.Width)
	b.RotatePiece(0, 0)
	// fmt.Printf("%s(Marked:%d/%d)\n", b, b.Marked, b.Height*b.Width)
	// b.RotatePiece(0, 0)
	// fmt.Printf("%s(Marked:%d/%d)\n", b, b.Marked, b.Height*b.Width)
	// b.RotatePiece(0, 0)
	// fmt.Printf("%s(Marked:%d/%d)\n", b, b.Marked, b.Height*b.Width)
	// b.RotatePiece(0, 0)
	// fmt.Printf("%s(Marked:%d/%d)\n", b, b.Marked, b.Height*b.Width)

	//b.PlaceBlock(4, 0, b.MakeElle(board.Green, board.Up))
	b.DropToBottom(0, 0)
	b.PlaceBlock(0, 0, b.MakePipe(board.Green, board.Up))
	b.DropToBottom(0, 0)
	fmt.Printf("%s(Marked:%d/%d)\n", b, b.Marked, b.Height*b.Width)

	b.PlaceBlock(3, 0, b.MakePipe(board.Green, board.Up))
	b.RotatePiece(3, 0)
	b.RotatePiece(3, 0)
	b.RotatePiece(3, 0)

	fmt.Printf("%s(Marked:%d/%d)\n", b, b.Marked, b.Height*b.Width)

	//b.RotatePiece(1, 0)
	//b.RotatePiece(1, 2)
	//b.RotatePiece(1, 2)

}
func game1() {

	b := board.NewBoard(8, 12)
	b.PlaceBlock(0, 0, b.MakeIElle(board.Blue, board.Up))
	b.PlaceBlock(1, 2, b.MakePipe(board.Green, board.Up))
	b.PlaceBlock(2, 2, b.MakePipe(board.Green, board.Up))
	b.PlaceBlock(3, 2, b.MakePipe(board.Green, board.Up))
	b.PlaceBlock(4, 2, b.MakePipe(board.Green, board.Up))
	b.PlaceBlock(5, 2, b.MakePipe(board.Green, board.Up))
	b.PlaceBlock(6, 2, b.MakePipe(board.Green, board.Up))
	b.PlaceBlock(7, 2, b.MakePipe(board.Green, board.Up))

	b.DropToBottom(0, 0)
	b.DropToBottom(1, 2)
	b.DropToBottom(2, 2)
	b.DropToBottom(3, 2)
	b.DropToBottom(4, 2)
	b.DropToBottom(5, 2)
	b.DropToBottom(6, 2)
	b.DropToBottom(7, 2)

	b.DropToBottom(0, 1)

	// b.TetrisReduce(0)

	// b.Place(3, 0, b.NewIElle(Purple, Up))
	// b.Place(0, 4, b.NewElle(Orange, Up))
	// b.Place(4, 1, b.NewTee(Red, Up))

	// fmt.Printf("%s(Marked:%d/%d)\n", b, b.marked, b.height*b.width)
	// b.TetrisReduce(0)

	b.TetrisMatch(nil)

	fmt.Printf("%s(Marked:%d/%d)\n", b, b.Marked, b.Height*b.Width)
	return
}
