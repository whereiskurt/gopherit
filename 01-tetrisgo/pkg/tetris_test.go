package pkg

import "testing"

func TestClassicTetris(t *testing.T) {
	tetris := NewClassicTetris(0)

	if tetris.Seed == 0 {
		t.Logf("Tetris seed not set - default zero random() game.")
	}

	return
}
