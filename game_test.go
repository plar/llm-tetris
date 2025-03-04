// file: game_test.go
package main

import (
	"os"
	"testing"
)

func TestAutoFall(t *testing.T) {
	game := NewGame()
	tetromino, _ := NewTetromino("I")
	game.grid = NewGrid()
	game.tetromino = tetromino
	tetromino.position = [2]int{0, 0} // Start from top left

	// Simulate the game loop manually without goroutines for test purposes
	canMove := true
	for canMove {
		canMove = game.grid.MoveTetromino(tetromino, "down")
	}

	finalPos := tetromino.position

	// Since the grid height is 20 and tetromino "I" in horizontal position counts 1,
	// the final position should be the row 19th
	if finalPos[0] != glassRows-1 {
		t.Errorf("Expected final position on row %d, but got %d", glassRows-1, finalPos[0])
	}

	// Check that moving further down stops after landing
	canMove = game.grid.MoveTetromino(tetromino, "down")
	if canMove {
		t.Error("Expected no further downward movement at the grid bottom")
	}
}

func TestAutoFallWithCollision(t *testing.T) {
	game := NewGame()
	tetromino, _ := NewTetromino("O")
	game.grid = NewGrid()
	game.tetromino = tetromino
	tetromino.position = [2]int{17, 0}       // Position it close to the bottom
	game.grid.cells[19][0] = CellStateFilled // Create a block to simulate collision
	game.grid.cells[19][1] = CellStateFilled

	// Simulate gravity falling
	canMove := true
	for canMove {
		canMove = game.grid.MoveTetromino(tetromino, "down")
	}

	finalPos := tetromino.position

	// Since there's a block on row 19, tetromino "O" should stop on row 17
	if finalPos[0] != 17 {
		t.Errorf("Expected final position on row 17, but got %d", finalPos[0])
	}

	// Ensure it doesn't move further due to collision
	canMove = game.grid.MoveTetromino(tetromino, "down")
	if canMove {
		t.Error("Expected tetromino to stop moving due to collision with settled blocks")
	}
}

// Test that a full row is properly detected and cleared
func TestClearFullRow(t *testing.T) {
	grid := NewGrid()
	// Fill the entire row
	for j := 0; j < glassCols; j++ {
		grid.cells[10][j] = CellStateFilled
	}

	// Clear full rows
	grid.clearFullRows()

	// Check that row 10 is cleared
	for j := 0; j < glassCols; j++ {
		if grid.cells[10][j] != CellStateEmpty {
			t.Errorf("expected row 10 to be empty, but found %s at col %d", grid.cells[10][j], j)
		}
	}

	// Ensure rows above are shifted correctly
	expectedEmptyRow := [glassCols]CellState{}
	for i := 0; i < 10; i++ {
		if grid.cells[i] != expectedEmptyRow {
			t.Errorf("expected row %d to be shifted to empty, but got %v", i, grid.cells[i])
		}
	}
}

// Test cascading clear of multiple rows
func TestClearMultipleFullRows(t *testing.T) {
	grid := NewGrid()

	// Fill several rows
	for j := 0; j < glassCols; j++ {
		grid.cells[9][j] = CellStateFilled
		grid.cells[10][j] = CellStateFilled
		grid.cells[11][j] = CellStateFilled
	}

	grid.clearFullRows()

	// Check rows that should be cleared
	for i := 9; i <= 11; i++ {
		for j := 0; j < glassCols; j++ {
			if grid.cells[i][j] != CellStateEmpty {
				t.Errorf("expected row %d to be empty, but found %s at col %d", i, grid.cells[i][j], j)
			}
		}
	}

	// Ensure rows above (0-8) are shifted to empty rows
	for i := 0; i < 9; i++ {
		for j := 0; j < glassCols; j++ {
			if grid.cells[i][j] != CellStateEmpty {
				t.Errorf("expected row %d to be empty after shifting, but got %v", i, grid.cells[i][j])
			}
		}
	}
}

func TestTetrominoLocking(t *testing.T) {
	game := NewGame()
	tetromino, _ := NewTetromino("I")
	game.grid = NewGrid()
	game.tetromino = tetromino
	tetromino.position = [2]int{glassRows - 1, 0}

	game.placeTetromino()

	totalCellFilled := 0
	for _, c := range game.grid.cells[glassRows-1][:] {
		if c == CellStateFilled {
			totalCellFilled++
		}
	}

	if totalCellFilled != 4 {
		t.Error("Expected last row to be filled with 'XXXX'")
	}

	game.grid.clearFullRows() // Should not clear since it's only filled partway
	if game.grid.cells[glassRows-1][0] != CellStateFilled {
		t.Error("Expected 'X' still at bottom after locking")
	}
}

func TestPieceSpawnAfterLock(t *testing.T) {
	game := NewGame()

	// Simulate the tetromino moving down to the bottom
	game.spawnTetromino() // Initial spawn

	// Manually move tetromino to the bottom
	for game.grid.MoveTetromino(game.tetromino, "down") {
		// Continue moving down until it can no longer move
	}

	// Now lock the tetromino into the grid
	game.placeTetromino()
	firstPosition := game.tetromino.position

	// Spawn a new tetromino, ensure it's in a different position
	game.spawnTetromino()
	secondPosition := game.tetromino.position

	if firstPosition == secondPosition {
		t.Error("Expected new tetromino to spawn; positions should be different after locking")
	}

	if game.gameOver {
		t.Error("Did not expect game over immediately after spawning a new tetromino")
	}
}

func TestGameOverCondition(t *testing.T) {
	game := NewGame()

	// Simulate the spawn area as filled
	for j := 0; j < glassCols; j++ {
		game.grid.cells[0][j] = CellStateFilled
	}

	// Try to spawn a new tetromino and verify game over
	game.spawnTetromino()

	if !game.gameOver {
		t.Error("Expected game over when spawning into a filled area")
	}
}

func TestScoreOnLineClear(t *testing.T) {
	game := NewGame()

	// Manually fill rows to simulate line clear
	for i := 0; i < 4; i++ {
		for j := 0; j < glassCols; j++ {
			game.grid.cells[i][j] = CellStateFilled
		}
	}
	linesCleared := game.grid.clearFullRows()
	game.updateScore(linesCleared)

	expectedScore := 800 // Four lines clear

	if game.score != expectedScore {
		t.Errorf("Expected score %d, got %d", expectedScore, game.score)
	}
}

func TestMultipleLineClears(t *testing.T) {
	game := NewGame()

	// Fill two rows
	for i := 0; i < 2; i++ {
		for j := 0; j < glassCols; j++ {
			game.grid.cells[i][j] = CellStateFilled
		}
	}

	linesCleared := game.grid.clearFullRows()
	game.updateScore(linesCleared)

	expectedScore := 300 // Two lines clear

	if game.score != expectedScore {
		t.Errorf("Expected score %d, got %d", expectedScore, game.score)
	}

	// Clear two more lines
	for i := 0; i < 2; i++ {
		for j := 0; j < glassCols; j++ {
			game.grid.cells[i][j] = CellStateFilled
		}
	}

	linesCleared = game.grid.clearFullRows()
	game.updateScore(linesCleared)

	expectedScore += 300 // Accumulate new score

	if game.score != expectedScore {
		t.Errorf("Expected cumulative score %d, got %d", expectedScore, game.score)
	}
}

// TestLevelProgression checks if level increases after 10 lines are cleared
func TestLevelProgression(t *testing.T) {
	game := NewGame()

	// Simulate clearing 10 lines
	for i := 0; i < 10; i++ {
		for j := 0; j < glassCols; j++ {
			game.grid.cells[i][j] = CellStateFilled
		}
	}
	linesCleared := game.grid.clearFullRows()
	game.updateLevel(linesCleared)

	if game.level != 2 {
		t.Errorf("Expected level 2 after clearing 10 lines, but got level %d", game.level)
	}
}

// TestIncreasedGravityAtHigherLevels verifies that the fall speed decreases as levels increase
func TestIncreasedGravityAtHigherLevels(t *testing.T) {
	game := NewGame()
	game.level = 5
	initialSpeed := game.currentFallSpeed()

	game.level = 6
	newSpeed := game.currentFallSpeed()

	if newSpeed >= initialSpeed {
		t.Error("Expected fall speed to increase (duration decrease) as level increases")
	}
}

// TestGameOverAndRestart verifies that the game correctly detects game over and resets the game state
func TestGameOverAndRestart(t *testing.T) {
	game := NewGame()

	// Simulate a filled grid
	for j := 0; j < glassCols; j++ {
		game.grid.cells[0][j] = CellStateFilled
	}

	// Attempt to spawn a new Tetromino to trigger game over
	game.spawnTetromino()
	if !game.gameOver {
		t.Error("Expected game over when spawning into a filled area")
	}

	// Simulate restart by calling reset function
	game.resetGame()
	if game.gameOver {
		t.Error("Expected game to reset, but gameOver state remains true")
	}
	if game.score != 0 || game.linesCleared != 0 || game.level != 1 {
		t.Error("Expected game state to reset, but some values did not reset")
	}
}

func TestHighScorePersistence(t *testing.T) {
	game := NewGame()
	game.score = 800 // Set a test score
	game.UpdateHighScore()
	game.config.SaveConfig("test_config.json")
	defer os.Remove("test_config.json")

	if game.config.HighScore != 800 {
		t.Errorf("Expected high score 800, but got %d", game.config.HighScore)
	}

	config, _ := LoadConfig("test_config.json")
	if config.HighScore != 800 {
		t.Errorf("Expected high score 800, but got %d", game.config.HighScore)
	}

	// Clean up the test artifact
}
