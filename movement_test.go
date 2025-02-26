package main

import "testing"

func TestLeftRightMovementBlockedAtWalls(t *testing.T) {
	game := NewGame()
	tetromino, _ := NewTetromino("I")
	game.tetromino = tetromino

	// Attempt to place tetromino at the left-most column and try to move left
	game.tetromino.position = [2]int{0, 0}
	canMoveLeft := game.grid.MoveTetromino(game.tetromino, "left")
	if canMoveLeft {
		t.Error("Expected tetromino to be blocked moving left at the wall")
	}

	// Attempt to place tetromino at the right-most column and try to move right
	rightMostPosition := cols - len(tetromino.GetCurrentShape()[0])
	game.tetromino.position = [2]int{0, rightMostPosition}
	canMoveRight := game.grid.MoveTetromino(game.tetromino, "right")
	if canMoveRight {
		t.Error("Expected tetromino to be blocked moving right at the wall")
	}
}

func TestLeftRightMovementBlockedByBlocks(t *testing.T) {
	game := NewGame()
	tetromino, _ := NewTetromino("I")
	game.tetromino = tetromino

	// Place tetromino and block to the left, then try to move left
	game.tetromino.position = [2]int{1, 1}
	game.grid.cells[1][0] = "X" // Blocking left movement
	canMoveLeft := game.grid.MoveTetromino(game.tetromino, "left")
	if canMoveLeft {
		t.Error("Expected tetromino to be blocked moving left by another block")
	}

	// Place tetromino and block to the right, then try to move right
	game.tetromino.position = [2]int{1, 0}
	game.grid.cells[1][4] = "X" // Blocking right movement
	canMoveRight := game.grid.MoveTetromino(game.tetromino, "right")
	if canMoveRight {
		t.Error("Expected tetromino to be blocked moving right by another block")
	}
}

// TestSoftDropSpeedsUpGravity checks that the soft drop speeds up tetromino descent
func TestSoftDropSpeedsUpGravity(t *testing.T) {
	game := NewGame()
	tetromino, _ := NewTetromino("I")
	game.tetromino = tetromino

	initialPosition := game.tetromino.position

	// Simulate holding the soft drop for faster descent
	game.softDrop = true
	for i := 0; i < 5; i++ { // let it fall for 5 ticks
		game.grid.MoveTetromino(game.tetromino, "down")
	}

	if game.tetromino.position[0] <= initialPosition[0] {
		t.Errorf("Expected tetromino to have moved downward with soft drop active, but got position %v", game.tetromino.position)
	}
}

// TestHardDropPlacesPieceInstantly checks that the hard drop places a tetromino instantly at the bottom
func TestHardDropPlacesPieceInstantly(t *testing.T) {
	game := NewGame()
	tetromino, _ := NewTetromino("I")
	game.tetromino = tetromino

	game.hardDrop()

	// Tetromino should be locked immediately at the bottom
	expectedRow := rows - len(tetromino.GetCurrentShape()) // Lowest possible position
	if tetromino.position[0] != expectedRow {
		t.Errorf("Expected tetromino to hard-drop to bottom (row %d), but was %d", expectedRow, game.tetromino.position[0])
	}

	// Verify the grid has the tetromino locked as expected
	for i, row := range tetromino.GetCurrentShape() {
		for j, cell := range row {
			if cell == 1 && game.grid.cells[expectedRow+i][tetromino.position[1]+j] != "X" {
				t.Errorf("Expected cell %v to be locked but was not", [2]int{expectedRow + i, tetromino.position[1] + j})
			}
		}
	}
}
