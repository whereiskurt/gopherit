package main

import (
	"fmt"
	"tetrigo/pkg"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	t := pkg.TetrisGame{}

	fmt.Printf("Tetris Game Summary: %s\n------------------\n", spew.Sdump(t))
}
