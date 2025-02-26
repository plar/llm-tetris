package main

import "fmt"

const (
	rows = 20
	cols = 10
)

// Cell states in the grid
const emptyCell = " "

type Grid struct {
	cells [rows][cols]string
}

// NewGrid initializes a 20x10 grid with each cell set to empty.
func NewGrid() *Grid {
	grid := &Grid{}
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			grid.cells[i][j] = emptyCell // Using a space to denote an empty cell
		}
	}
	return grid
}

// Print displays the grid to the console in a structured format.
func (g *Grid) Print() {
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			fmt.Printf("|%s", g.cells[i][j])
		}
		fmt.Println("|") // Close the row with a "|"
	}
}

// Check if a tetromino can move to a specified position
func (g *Grid) CanMove(tetromino *Tetromino, newPos [2]int) bool {
	newShape := tetromino.GetCurrentShape()
	for i, row := range newShape {
		for j, cell := range row {
			if cell == 1 {
				x := newPos[0] + i
				y := newPos[1] + j

				if x < 0 || x >= rows || y < 0 || y >= cols {
					return false // Out of bounds
				}
				if g.cells[x][y] != emptyCell {
					return false // Collision with settled block
				}
			}
		}
	}
	return true
}

// Move the tetromino within the grid if possible
func (g *Grid) MoveTetromino(tetromino *Tetromino, direction string) bool {
	var newPosition [2]int

	switch direction {
	case "left":
		newPosition = [2]int{tetromino.position[0], tetromino.position[1] - 1}
	case "right":
		newPosition = [2]int{tetromino.position[0], tetromino.position[1] + 1}
	case "down":
		newPosition = [2]int{tetromino.position[0] + 1, tetromino.position[1]}
	default:
		return false
	}

	if g.CanMove(tetromino, newPosition) {
		tetromino.position = newPosition
		return true
	}
	return false
}

// Check if a row is fully filled
func (g *Grid) isRowFull(row int) bool {
	for _, cell := range g.cells[row] {
		if cell == emptyCell {
			return false
		}
	}
	return true
}

// Remove a row and shift all above rows down
func (g *Grid) clearRow(row int) {
	// Shift all rows above the specified row down by one
	for i := row; i > 0; i-- {
		g.cells[i] = g.cells[i-1]
	}
	// Clear the top row after shifting
	for j := 0; j < cols; j++ {
		g.cells[0][j] = emptyCell
	}
}

// clearFullRows removes full rows and returns the number of rows cleared
func (g *Grid) clearFullRows() int {
	linesCleared := 0
	for i := 0; i < rows; i++ {
		if g.isRowFull(i) {
			g.clearRow(i)
			linesCleared++
		}
	}
	return linesCleared
}

func (g *Grid) CanRotate(tetromino *Tetromino) bool {
	nextRotation := (tetromino.rotationState + 1) % len(tetromino.shape)
	shape := tetromino.shape[nextRotation]
	for i, row := range shape {
		for j, cell := range row {
			if cell == 1 {
				x := tetromino.position[0] + i
				y := tetromino.position[1] + j
				if x < 0 || x >= rows || y < 0 || y >= cols || g.cells[x][y] != emptyCell {
					return false // Out of bounds or collision
				}
			}
		}
	}
	return true
}

// Rotate attempts to rotate the active tetromino clockwise
func (g *Game) RotateTetromino() {
	if g.grid.CanRotate(g.tetromino) {
		g.tetromino.Rotate()
		g.Render()
	}
}
