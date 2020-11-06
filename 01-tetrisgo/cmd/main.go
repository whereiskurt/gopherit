package main

import (
	"fmt"
	"os"
	"strconv"
	"tetrigo/pkg"
	"time"
)

func main() {
	seed := time.Now().UnixNano()

	if len(os.Args) > 1 {
		argseed, _ := strconv.Atoi(os.Args[1])
		seed = int64(argseed) // Seed will be zero if fails to parse strconv
	}

	game := pkg.NewClassicTetris(8, 5, seed)

	//fmt.Printf("Tetris Game Summary:\n %s\n------------------\n", spew.Sdump(game))
	fmt.Printf("%s------------------\n", game.Board)
}
