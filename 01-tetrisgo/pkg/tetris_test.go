package pkg

import (
	"os"
	"strconv"
	"testing"
	"time"
)

func TestClassicTetris(t *testing.T) {
	seed := time.Now().UnixNano()

	if len(os.Args) > 1 {
		argseed, _ := strconv.Atoi(os.Args[1])
		seed = int64(argseed) // Seed will be zero if fails to parse strconv
	}

	game := NewClassicTetris(8, 5, seed)

	t.Logf("Board:\n%s", game.Board)

	return
}
