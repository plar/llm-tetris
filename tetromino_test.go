// file: tetromino_test.go
package main

import (
	"testing"
)

// TestTetrominoCreation verifies tetrominos are created correctly with defined shapes
func TestTetrominoCreation(t *testing.T) {
	for tetrominoType, expectedShapes := range tetrominoShapes {
		tetromino, err := NewTetromino(tetrominoType)
		if err != nil {
			t.Fatalf("failed to create tetromino %s: %v", tetrominoType, err)
		}

		// Check each rotation state shape
		for state, expectedShape := range expectedShapes {
			tetromino.rotationState = state
			if got := tetromino.GetCurrentShape(); !equalShapes(got, expectedShape) {
				t.Errorf("expected shape for tetromino %s at rotation %d to be %v, got %v",
					tetrominoType, state, expectedShape, got)
			}
		}
	}
}

// TestTetrominoRotation checks if the rotation state updates correctly
func TestTetrominoRotation(t *testing.T) {
	tetromino, _ := NewTetromino("T")

	if tetromino.rotationState != 0 {
		t.Errorf("expected initial rotation state to be 0, got %d", tetromino.rotationState)
	}

	// Rotate 4 times and ensure it cycles back to initial state
	for i := 1; i <= 3; i++ {
		tetromino.Rotate()
		if tetromino.rotationState != i {
			t.Errorf("expected rotation state after %d rotates to be %d, got %d", i, i, tetromino.rotationState)
		}
	}
	tetromino.Rotate() // This should bring it back to state 0
	if tetromino.rotationState != 0 {
		t.Errorf("expected rotation state after back to 0, got %d", tetromino.rotationState)
	}
}

// Helper method to compare if two shapes are equal
func equalShapes(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func TestRotationChangesShape(t *testing.T) {
	tetromino, _ := NewTetromino("T")
	originalShape := tetromino.GetCurrentShape()

	tetromino.Rotate()
	newShape := tetromino.GetCurrentShape()

	if &originalShape == &newShape {
		t.Error("Expected tetromino shape to change upon rotation")
	}
}

func TestRotationBlockedByWall(t *testing.T) {
	game := NewGame()
	tetromino, _ := NewTetromino("L")
	game.tetromino = tetromino

	// Place tetromino near the right wall
	tetromino.position = [2]int{0, glassCols - 1}
	canRotate := game.grid.CanRotate(tetromino)

	if canRotate {
		t.Error("Expected rotation to be blocked by wall, but it was not")
	}
}

func TestRotationBlockedByBlocks(t *testing.T) {
	game := NewGame()
	tetromino, _ := NewTetromino("L")
	game.tetromino = tetromino

	// Place a blocking block on rotation path
	tetromino.position = [2]int{1, 1}
	game.grid.cells[2][1] = CellStateFilled
	canRotate := game.grid.CanRotate(tetromino)

	if canRotate {
		t.Error("Expected rotation to be blocked by another block")
	}
}
