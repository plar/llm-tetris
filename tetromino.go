// file: tetromino.go
package main

import "fmt"

const tetrominoTypes = "IOTJLSZ"

// Tetromino shape definitions
var tetrominoShapes = map[string][][][]int{
	"I": {
		{{1, 1, 1, 1}},       // 0 degrees
		{{1}, {1}, {1}, {1}}, // 90 degrees
		{{1, 1, 1, 1}},       // 180 degrees
		{{1}, {1}, {1}, {1}}, // 270 degrees
	},
	"O": {
		{{1, 1}, {1, 1}}, // 0 degrees
		{{1, 1}, {1, 1}}, // 90 degrees
		{{1, 1}, {1, 1}}, // 180 degrees
		{{1, 1}, {1, 1}}, // 270 degrees
	},
	"T": {
		{{0, 1, 0}, {1, 1, 1}},   // 0 degrees
		{{1, 0}, {1, 1}, {1, 0}}, // 90 degrees
		{{1, 1, 1}, {0, 1, 0}},   // 180 degrees
		{{0, 1}, {1, 1}, {0, 1}}, // 270 degrees
	},
	"S": {
		{{0, 1, 1}, {1, 1, 0}},   // 0 degrees
		{{1, 0}, {1, 1}, {0, 1}}, // 90 degrees
		{{0, 1, 1}, {1, 1, 0}},   // 180 degrees
		{{1, 0}, {1, 1}, {0, 1}}, // 270 degrees
	},
	"Z": {
		{{1, 1, 0}, {0, 1, 1}},   // 0 degrees
		{{0, 1}, {1, 1}, {1, 0}}, // 90 degrees
		{{1, 1, 0}, {0, 1, 1}},   // 180 degrees
		{{0, 1}, {1, 1}, {1, 0}}, // 270 degrees
	},
	"J": {
		{{1, 0, 0}, {1, 1, 1}},   // 0 degrees
		{{1, 1}, {1, 0}, {1, 0}}, // 90 degrees
		{{1, 1, 1}, {0, 0, 1}},   // 180 degrees
		{{0, 1}, {0, 1}, {1, 1}}, // 270 degrees
	},
	"L": {
		{{0, 0, 1}, {1, 1, 1}},   // 0 degrees
		{{1, 0}, {1, 0}, {1, 1}}, // 90 degrees
		{{1, 1, 1}, {1, 0, 0}},   // 180 degrees
		{{1, 1}, {0, 1}, {0, 1}}, // 270 degrees
	},
}

// Tetromino represents a tetromino piece
type Tetromino struct {
	shape         [][][]int
	position      [2]int
	rotationState int
}

// NewTetromino creates a new Tetromino of given type
func NewTetromino(tType string) (*Tetromino, error) {
	shape, exists := tetrominoShapes[tType]
	if !exists {
		return nil, fmt.Errorf("invalid tetromino type: %s", tType)
	}
	return &Tetromino{
		shape:         shape,
		position:      [2]int{0, 0}, // default position at top-left corner
		rotationState: 0,            // default to 0 degrees
	}, nil
}

// Rotate changes the tetromino to its next rotation state
func (t *Tetromino) Rotate() {
	t.rotationState = (t.rotationState + 1) % len(t.shape)
}

// GetCurrentShape returns the current shape of the tetromino based on its rotation
func (t *Tetromino) GetCurrentShape() [][]int {
	return t.shape[t.rotationState]
}
