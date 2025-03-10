package main

import (
	"strings"
	"testing"
)

// TestNewGrid checks if the grid initializes correctly with empty cells.
func TestNewGrid(t *testing.T) {
	grid := NewGrid()

	// Verify the grid is of the correct size
	if len(grid.cells) != glassRows {
		t.Errorf("Expected grid to have %d rows, got %d", glassRows, len(grid.cells))
	}
	if len(grid.cells[0]) != glassCols {
		t.Errorf("Expected grid to have %d cols, got %d", glassCols, len(grid.cells[0]))
	}

	// Verify all cells are initialized to a space
	for i := range grid.cells {
		for j := range grid.cells[i] {
			if grid.cells[i][j] != CellStateEmpty {
				t.Errorf("Expected cell [%d][%d] to be empty, got '%v'", i, j, grid.cells[i][j])
			}
		}
	}
}

// TestPrintGrid checks the output format of the Print method.
func TestPrintGrid(t *testing.T) {
	grid := NewGrid()
	expectedLine := strings.Repeat("| ", glassCols) + "|"

	for i := 0; i < glassRows; i++ {
		rowOutput := ""
		for j := 0; j < glassCols; j++ {
			rowOutput += "|" + grid.cells[i][j].String()
		}
		rowOutput += "|"

		if rowOutput != expectedLine {
			t.Errorf("Grid line %d did not match expected pattern. Got: %v; Want: %v", i, rowOutput, expectedLine)
		}
	}
}

// TestMoveTetromino tests movement constraints against grid edges and obstacles
func TestMoveTetromino(t *testing.T) {
	grid := NewGrid()
	tetromino, _ := NewTetromino("I")

	// Move down initially to test edge constraints
	success := grid.MoveTetromino(tetromino, "down")
	if !success {
		t.Fatal("expected to move down successfully")
	}

	// Try to move out of left bounds
	tetromino.position = [2]int{0, 0} // Reset position
	success = grid.MoveTetromino(tetromino, "left")
	if success {
		t.Fatal("expected failure moving left out of bounds")
	}

	// Try to move out of right bounds
	tetromino.position = [2]int{0, glassCols - len(tetromino.GetCurrentShape()[0])}
	success = grid.MoveTetromino(tetromino, "right")
	if success {
		t.Fatal("expected failure moving right out of bounds")
	}

	// Try to move out of bottom bounds
	tetromino.position = [2]int{glassRows - len(tetromino.GetCurrentShape()), 0}
	success = grid.MoveTetromino(tetromino, "down")
	if success {
		t.Fatal("expected failure moving down out of bounds")
	}

	// Place an obstacle and try to move into it
	grid.cells[3][0] = CellStateFilled // Making a block at (3, 0)
	tetromino.position = [2]int{2, 0}
	success = grid.MoveTetromino(tetromino, "down")
	if success {
		t.Fatal("expected failure moving into an occupied cell")
	}
}
