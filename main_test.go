package main

import (
	"testing"
)

func TestLifeGameInit(t *testing.T) {
	game := new(LifeGame)
	result := game.Init(5, 5)

	if result != game {
		t.Error("Init should return the same instance")
	}

	if game.row != 5 || game.col != 5 {
		t.Errorf("Expected dimensions 5x5, got %dx%d", game.row, game.col)
	}

	if game.time != 0 {
		t.Errorf("Expected time to be 0, got %d", game.time)
	}

	if len(game.board) != 5 {
		t.Errorf("Expected board to have 5 rows, got %d", len(game.board))
	}

	for i, row := range game.board {
		if len(row) != 5 {
			t.Errorf("Expected row %d to have 5 columns, got %d", i, len(row))
		}
	}
}

func TestLifeGameIsEmpty(t *testing.T) {
	game := new(LifeGame).Init(3, 3)

	for r := 0; r < game.row; r++ {
		for c := 0; c < game.col; c++ {
			game.board[r][c] = false
		}
	}

	if !game.isEmpty() {
		t.Error("Expected empty board to return true")
	}

	game.board[1][1] = true

	if game.isEmpty() {
		t.Error("Expected non-empty board to return false")
	}
}

func TestLifeGameEquals(t *testing.T) {
	game1 := new(LifeGame).Init(3, 3)
	game2 := new(LifeGame).Init(3, 3)

	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			game1.board[r][c] = false
			game2.board[r][c] = false
		}
	}

	if !game1.equals(game2) {
		t.Error("Expected identical empty boards to be equal")
	}

	game1.board[1][1] = true
	game2.board[1][1] = true

	if !game1.equals(game2) {
		t.Error("Expected identical boards with same pattern to be equal")
	}

	game2.board[0][0] = true

	if game1.equals(game2) {
		t.Error("Expected different boards to not be equal")
	}

	game3 := new(LifeGame).Init(2, 2)
	if game1.equals(game3) {
		t.Error("Expected boards with different dimensions to not be equal")
	}
}

func TestLifeGameCountNowAlive(t *testing.T) {
	game := new(LifeGame).Init(3, 3)

	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			game.board[r][c] = false
		}
	}

	if game.count_now_alive(1, 1) != 0 {
		t.Error("Expected count of dead cell to be 0")
	}

	game.board[1][1] = true
	if game.count_now_alive(1, 1) != 1 {
		t.Error("Expected count of live cell to be 1")
	}

	game.board[2][2] = true
	if game.count_now_alive(-1, -1) != 1 {
		t.Error("Expected wrapping to work correctly")
	}

	game.board[0][0] = true
	if game.count_now_alive(3, 3) != 1 {
		t.Error("Expected wrapping to work correctly for positive overflow")
	}
}

func TestLifeGameIsDeadOrAlive(t *testing.T) {
	game := new(LifeGame).Init(5, 5)

	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			game.board[r][c] = false
		}
	}

	game.board[1][1] = true
	game.board[1][2] = true
	game.board[2][1] = true

	if !game.is_dead_or_alive(2, 2) {
		t.Error("Dead cell with 3 neighbors should become alive")
	}

	game.board[2][2] = true // make center cell alive
	// Remove one neighbor to have exactly 2 neighbors
	game.board[2][1] = false

	if !game.is_dead_or_alive(2, 2) {
		t.Error("Live cell with 2 neighbors should stay alive")
	}

	game.board[2][1] = true // add back the neighbor

	if !game.is_dead_or_alive(2, 2) {
		t.Error("Live cell with 3 neighbors should stay alive")
	}

	game.board[1][1] = false
	game.board[1][2] = false

	if game.is_dead_or_alive(2, 2) {
		t.Error("Live cell with < 2 neighbors should die")
	}

	game.board[1][1] = true
	game.board[1][2] = true
	game.board[1][3] = true
	game.board[2][3] = true
	game.board[3][3] = true

	if game.is_dead_or_alive(2, 2) {
		t.Error("Live cell with > 3 neighbors should die")
	}
}

func TestConwaysPatterns(t *testing.T) {
	game := new(LifeGame).Init(4, 4)
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			game.board[r][c] = false
		}
	}

	game.board[1][1] = true
	game.board[1][2] = true
	game.board[2][1] = true
	game.board[2][2] = true

	for r := 1; r <= 2; r++ {
		for c := 1; c <= 2; c++ {
			if !game.is_dead_or_alive(r, c) {
				t.Errorf("Block pattern should be stable at (%d,%d)", r, c)
			}
		}
	}

	if game.is_dead_or_alive(0, 0) || game.is_dead_or_alive(0, 1) || game.is_dead_or_alive(0, 2) {
		t.Error("Cells around block should remain dead")
	}
}

func TestBlinkerPattern(t *testing.T) {
	game := new(LifeGame).Init(5, 5)

	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			game.board[r][c] = false
		}
	}

	game.board[1][2] = true
	game.board[2][2] = true
	game.board[3][2] = true

	if !game.is_dead_or_alive(2, 1) || !game.is_dead_or_alive(2, 2) || !game.is_dead_or_alive(2, 3) {
		t.Error("Vertical blinker should become horizontal")
	}

	if game.is_dead_or_alive(1, 2) || game.is_dead_or_alive(3, 2) {
		t.Error("Top and bottom of vertical blinker should die")
	}
}

func TestStabilizationEmptyBoard(t *testing.T) {
	game := new(LifeGame).Init(3, 3)

	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			game.board[r][c] = false
		}
	}

	if !game.isEmpty() {
		t.Error("Empty board should be detected as empty")
	}

	game.board[1][1] = true
	if game.isEmpty() {
		t.Error("Non-empty board should not be detected as empty")
	}
}

func TestStabilizationStaticPattern(t *testing.T) {
	game := new(LifeGame).Init(4, 4)

	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			game.board[r][c] = false
		}
	}

	game.board[1][1] = true
	game.board[1][2] = true
	game.board[2][1] = true
	game.board[2][2] = true

	next := new(LifeGame).Init(4, 4)
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			next.board[r][c] = game.is_dead_or_alive(r, c)
		}
	}

	if !next.equals(game) {
		t.Error("Block pattern should remain stable")
	}
}

func TestOscillationDetection(t *testing.T) {
	game1 := new(LifeGame).Init(5, 5)
	game2 := new(LifeGame).Init(5, 5)

	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			game1.board[r][c] = false
			game2.board[r][c] = false
		}
	}

	game1.board[1][2] = true
	game1.board[2][2] = true
	game1.board[3][2] = true

	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			game2.board[r][c] = game1.is_dead_or_alive(r, c)
		}
	}

	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			game3.board[r][c] = game2.is_dead_or_alive(r, c)
		}
	}

	if !game3.equals(game1) {
		t.Error("Blinker should oscillate back to original state")
	}

	if game2.equals(game1) {
		t.Error("Blinker intermediate state should be different")
	}
}

func TestPatternDiesOut(t *testing.T) {
	game := new(LifeGame).Init(5, 5)

	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			game.board[r][c] = false
		}
	}

	game.board[2][2] = true

	next := new(LifeGame).Init(5, 5)
	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			next.board[r][c] = game.is_dead_or_alive(r, c)
		}
	}

	if !next.isEmpty() {
		t.Error("Isolated single cell should die")
	}
}

func TestEdgeCases(t *testing.T) {
	game := new(LifeGame).Init(1, 1)

	if game.row != 1 || game.col != 1 {
		t.Error("1x1 board should be created correctly")
	}

	bigGame := new(LifeGame).Init(100, 100)

	if bigGame.row != 100 || bigGame.col != 100 {
		t.Error("Large board should be created correctly")
	}

	if len(bigGame.board) != 100 {
		t.Error("Large board should have correct number of rows")
	}
}

// Benchmark test for performance
func BenchmarkLifeGameIsDeadOrAlive(b *testing.B) {
	game := new(LifeGame).Init(50, 50)

	for r := 0; r < 50; r++ {
		for c := 0; c < 50; c++ {
			game.board[r][c] = (r+c)%3 == 0
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		game.is_dead_or_alive(25, 25)
	}
}

func BenchmarkGeneration(b *testing.B) {
	game := new(LifeGame).Init(50, 50)

	for r := 0; r < 50; r++ {
		for c := 0; c < 50; c++ {
			game.board[r][c] = (r+c)%3 == 0
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		next := new(LifeGame).Init(50, 50)
		for r := 0; r < 50; r++ {
			for c := 0; c < 50; c++ {
				next.board[r][c] = game.is_dead_or_alive(r, c)
			}
		}
	}
}
