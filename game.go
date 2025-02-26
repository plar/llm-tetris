// file: game.go
package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell"
)

type Game struct {
	screen       tcell.Screen
	grid         *Grid
	tetromino    *Tetromino
	quit         chan struct{}
	gameOver     bool
	softDrop     bool
	score        int
	linesCleared int
	level        int
	message      string
	messageTimer *time.Timer
	rng          *rand.Rand
	status       string
	config       *Config
}

func NewGame() *Game {
	return NewGameWithScreen(nil, nil)
}

func NewGameWithScreen(screen tcell.Screen, config *Config) *Game {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	if config == nil {
		config = NewConfig()
	}

	game := &Game{
		screen:       screen,
		grid:         NewGrid(),
		quit:         make(chan struct{}),
		rng:          rng,
		score:        0,
		linesCleared: 0,
		level:        1,
		status:       "in game", // Initially the game is active
		config:       config,
	}
	return game
}

var keyNamesToKey = func() map[string]tcell.Key {
	kv := make(map[string]tcell.Key)
	for k, v := range tcell.KeyNames {
		kv[v] = k
	}
	return kv
}()

func eqKey(k *tcell.EventKey, bindingKey string) bool {
	key, ok := keyNamesToKey[bindingKey]
	if ok && k.Key() == key {
		return true
	}

	if len(bindingKey) > 0 {
		lwr := strings.ToLower(bindingKey)
		upr := strings.ToUpper(bindingKey)
		return k.Key() == tcell.KeyRune && (k.Rune() == rune(lwr[0]) || k.Rune() == rune(upr[0]))
	}

	return false
}

// Listen for player input to move tetromino
func (g *Game) ListenForInput(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		ev := g.screen.PollEvent()
		switch e := ev.(type) {
		case *tcell.EventKey:
			if e.Key() == tcell.KeyEscape {
				g.Stop()
				return
			}

			switch g.status {
			case "in game":
				if eqKey(e, g.config.KeyBindings.Left) {
					g.grid.MoveTetromino(g.tetromino, "left")
				} else if eqKey(e, g.config.KeyBindings.Right) {
					g.grid.MoveTetromino(g.tetromino, "right")
				} else if eqKey(e, g.config.KeyBindings.Down) {
					g.softDrop = true
				} else if eqKey(e, g.config.KeyBindings.Rotate) {
					g.RotateTetromino()
				} else if eqKey(e, g.config.KeyBindings.HardDrop) {
					g.hardDrop()
				} else if eqKey(e, g.config.KeyBindings.Pause) {
					g.status = "paused"
				}

			case "paused":
				if eqKey(e, g.config.KeyBindings.Pause) {
					g.status = "in game"
				}

			case "game over":
				if e.Rune() == 'y' || e.Rune() == 'Y' {
					g.resetGame()
				} else {
					g.Stop()
					return
				}
			}
		}
	}
}

// Run executes the main game loop, handling gravity and rendering
func (g *Game) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(g.currentFallSpeed())
	defer ticker.Stop()

	g.spawnTetromino()
	g.Render()

	for {
		select {
		case <-ticker.C:
			if g.status != "in game" {
				// If not in the "in game" state, just wait and rerender
				g.Render()
				continue
			}

			fallSpeed := g.currentFallSpeed()
			if g.softDrop {
				fallSpeed = 50 * time.Millisecond
			}
			ticker.Reset(fallSpeed)

			// Process the game state if it is active
			if !g.grid.MoveTetromino(g.tetromino, "down") {
				g.placeTetromino()
				g.spawnTetromino()
			}
			g.Render()
			g.softDrop = false

		case <-g.quit:
			return
		}
	}
}

// currentFallSpeed returns the fall speed based on the current level
func (g *Game) currentFallSpeed() time.Duration {
	var baseDuration = 500 * time.Millisecond

	speedIncrease := 50 * time.Millisecond * time.Duration(g.level-1)
	if speedIncrease >= baseDuration {
		return 50 * time.Millisecond // Cap the speed to a reasonable minimum
	}
	return baseDuration - speedIncrease
}

// Render draws the current game state to the screen
func (g *Game) Render() {
	if g.screen == nil {
		return
	}

	g.screen.Clear()

	// Draw the grid with double-width characters
	for y, row := range g.grid.cells {
		for x, cell := range row {
			char := ".."             // Two dots for empty cells
			color := tcell.ColorGray // Default color for empty cells
			if cell != emptyCell {
				char = "[]"              // Two block characters for occupied cells
				color = tcell.ColorWhite // Color for occupied cells
			}
			drawText(g.screen, x*2, y, char, color) // Multiply column by 2 for alignment with double characters
		}
	}

	// If game is in progress, draw the active tetromino using double-width characters
	if g.status == "in game" || g.status == "paused" {
		for i, row := range g.tetromino.GetCurrentShape() {
			for j, cell := range row {
				if cell == 1 {
					x := (g.tetromino.position[1] + j)
					y := g.tetromino.position[0] + i
					drawText(g.screen, x*2, y, "[]", tcell.ColorYellow)
				}
			}
		}
	}

	// Display score, level, and lines cleared
	scoreText := fmt.Sprintf("Score: %d (%d)", g.score, g.config.HighScore)
	levelText := fmt.Sprintf("Level: %d", g.level)
	linesText := fmt.Sprintf("Lines: %d", g.linesCleared)
	drawText(g.screen, 2+(cols*2), 0, scoreText, tcell.ColorGreen)
	drawText(g.screen, 2+(cols*2), 1, linesText, tcell.ColorGreen)
	drawText(g.screen, 2+(cols*2), 2, levelText, tcell.ColorGreen)

	// Handle different messages based on game status
	switch g.status {
	case "game over":
		msgs := []string{"Game Over", "Play again? Y/n"}
		for i, msg := range msgs {
			drawText(g.screen, ((cols*2)-len(msg))/2, (rows/2)+i, msg, tcell.ColorRed)
		}
	case "paused":
		msgs := []string{"Paused", "Resume press 'P'"}
		for i, msg := range msgs {
			drawText(g.screen, ((cols*2)-len(msg))/2, (rows/2)+i, msg, tcell.ColorRed)
		}
	}

	// Display a temporary message if any
	if g.message != "" {
		drawText(g.screen, 0, rows+3, g.message, tcell.ColorRed)
	}

	g.screen.Show()
}

// drawText draws text at a specific location using the specified color
func drawText(screen tcell.Screen, x, y int, text string, color tcell.Color) {
	style := tcell.StyleDefault.Foreground(color)
	for i, c := range text {
		screen.SetContent(x+i, y, c, nil, style)
	}
}

// resetGame resets the game state
func (g *Game) resetGame() {
	g.grid = NewGrid()
	g.tetromino = nil
	g.gameOver = false
	g.score = 0
	g.linesCleared = 0
	g.level = 1
	g.message = ""
	g.softDrop = false

	g.status = "in game" // Reset status to in game

	// Restart game loop, spawning a new tetromino first
	g.spawnTetromino()
}

// Spawn the next tetromino and check for game over
func (g *Game) spawnTetromino() {
	tetrominoType := string(tetrominoTypes[g.rng.Intn(len(tetrominoTypes))])
	tetromino, err := NewTetromino(tetrominoType)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	tetromino.position = [2]int{0, cols/2 - len(tetromino.GetCurrentShape()[0])/2}

	// Check if spawn position is occupied
	if g.grid.CanMove(tetromino, tetromino.position) {
		g.tetromino = tetromino
	} else {
		g.status = "game over" // Update status to game over
		g.gameOver = true
		g.UpdateHighScore()
	}
}

// placeTetromino locks the tetromino in place and clears lines if possible
func (g *Game) placeTetromino() {
	shape := g.tetromino.GetCurrentShape()
	pos := g.tetromino.position
	for i, row := range shape {
		for j, cell := range row {
			if cell == 1 {
				g.grid.cells[pos[0]+i][pos[1]+j] = "X"
			}
		}
	}

	linesCleared := g.grid.clearFullRows()
	g.updateScore(linesCleared)
	g.updateLevel(linesCleared)
}

// updateScore updates the score based on the number of lines cleared
func (g *Game) updateScore(linesCleared int) {
	switch linesCleared {
	case 1:
		g.score += 100
	case 2:
		g.score += 300
	case 3:
		g.score += 500
	case 4:
		g.score += 800
	default:
		// No lines cleared, no score change
	}
	// fmt.Printf("Score: %d\n", g.score) // Debug message; use tcell for real display
}

// Update level and manage message display
func (g *Game) updateLevel(linesCleared int) {
	g.linesCleared += linesCleared
	if g.linesCleared >= 10 {
		g.level++
		g.linesCleared -= 10
		g.message = fmt.Sprintf("Level Up! New Level: %d", g.level)
		g.setMessageTimer(2 * time.Second) // Display message for 2 seconds
	}
}

// Set a message to display temporarily
func (g *Game) setMessageTimer(duration time.Duration) {
	if g.messageTimer != nil {
		g.messageTimer.Stop()
	}
	g.messageTimer = time.AfterFunc(duration, func() {
		g.message = ""
		g.Render()
	})
}

// hardDrop instantly places the tetromino to the lowest possible position
func (g *Game) hardDrop() {
	for g.grid.MoveTetromino(g.tetromino, "down") {
	}
	g.placeTetromino()
	if !g.gameOver {
		g.spawnTetromino()
	}
	g.Render()
}

func (g *Game) Stop() {
	g.UpdateHighScore()
	g.quit <- struct{}{}
	close(g.quit)
}

func (g *Game) UpdateHighScore() {
	if g.score < g.config.HighScore {
		return
	}
	g.config.HighScore = g.score
}
