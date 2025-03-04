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

// CellState represents the different states a cell can be in
type CellState int

// Constants for the different cell states
const (
	CellStateEmpty CellState = iota
	CellStateFilled
	CellStateTetromino
	CellStateShadow
)

func (c CellState) String() string {
	switch c {
	case CellStateEmpty:
		return " "
	case CellStateFilled:
		return "X"
	case CellStateTetromino:
		return "T"
	case CellStateShadow:
		return "S"
	}

	return "?"
}

// Status represents the current state of the game
type Status string

// Constants for the different game states
const (
	StatusInGame   Status = "in game"
	StatusPaused   Status = "paused"
	StatusGameOver Status = "game over"
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
	message      []string
	messageTimer *time.Timer
	rng          *rand.Rand
	status       Status
	config       *Config
	invalidate   chan struct{}
	left, top    int
}

func NewGame() *Game {
	return NewGameWithScreen(nil, nil)
}

func NewGameWithScreen(screen tcell.Screen, config *Config) *Game {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	if config == nil {
		config = NewConfig()
	}

	w, h := 80, 25
	if screen != nil {
		w, h = screen.Size()
	}

	left := (w - gameWidth) / 2
	top := (h - gameHeight) / 2

	game := &Game{
		screen:       screen,
		grid:         NewGrid(),
		quit:         make(chan struct{}),
		rng:          rng,
		score:        0,
		linesCleared: 0,
		level:        1,
		status:       StatusInGame, // Initially the game is active
		config:       config,
		invalidate:   make(chan struct{}, 1),
		left:         left,
		top:          top,
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
	defer g.Stop()

	for {
		ev := g.screen.PollEvent()
		switch e := ev.(type) {
		case *tcell.EventKey:
			if e.Key() == tcell.KeyEscape {
				return
			}

			stateChanged := false

			switch g.status {
			case StatusInGame:
				if eqKey(e, g.config.KeyBindings.Left) {
					stateChanged = g.grid.MoveTetromino(g.tetromino, "left")
				} else if eqKey(e, g.config.KeyBindings.Right) {
					stateChanged = g.grid.MoveTetromino(g.tetromino, "right")
				} else if eqKey(e, g.config.KeyBindings.Down) {
					g.softDrop = true
					stateChanged = true
				} else if eqKey(e, g.config.KeyBindings.Rotate) {
					stateChanged = g.RotateTetromino()
				} else if eqKey(e, g.config.KeyBindings.HardDrop) {
					g.hardDrop()
					stateChanged = true
				} else if eqKey(e, g.config.KeyBindings.Pause) {
					g.status = StatusPaused
					stateChanged = true
				} else if eqKey(e, g.config.KeyBindings.ToggleShadow) {
					g.config.ShowShadow = !g.config.ShowShadow
					stateChanged = true
				}

			case StatusPaused:
				if eqKey(e, g.config.KeyBindings.Pause) {
					g.status = StatusInGame
					stateChanged = true
				}

			case StatusGameOver:
				if e.Key() == tcell.KeyEnter || eqKey(e, "Y") {
					g.resetGame()
					stateChanged = true
				} else if e.Key() == tcell.KeyEscape || eqKey(e, "N") {
					return
				}
			}

			if stateChanged {
				g.Invalidate()
			}
		}
	}
}

func (g *Game) Invalidate() {
	select {
	case g.invalidate <- struct{}{}:
	default: // Non-blocking send; ignore if already full
	}
}

// Run executes the main game loop, handling gravity and rendering
func (g *Game) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(g.currentFallSpeed())
	defer ticker.Stop()

	g.spawnTetromino()
	g.Invalidate()

	for {
		select {
		case <-ticker.C:
			if g.status != StatusInGame {
				// If not in the "in game" state, just wait and rerender
				g.Invalidate()
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

			if g.score > g.config.HighScore {
				g.config.HighScore = g.score
			}

			g.softDrop = false
			g.Invalidate()

		case <-g.invalidate:
			g.Render()

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

	g.drawGameBorder()
	g.drawGrid()
	g.drawShadowTetromino()
	g.drawActiveTetromino()
	g.drawGameInfo()
	g.drawGameStatusMessages()
	g.drawTemporaryMessage()

	g.screen.Show()
}

func (g *Game) drawGrid() {
	var char []string
	for y, row := range g.grid.cells {
		for x, cell := range row {
			char = CellGridEmpty
			color := tcell.ColorGray
			if cell != CellStateEmpty {
				char = CellGridFilled
				color = tcell.ColorWhite
			}
			g.drawCell(x, y, char, color, tcell.ColorDefault)
		}
	}
}

// drawCell draws a single cell on the screen with the specified colors and piece
func (g *Game) drawCell(cx int, cy int, piece []string, fgColor, bgColor tcell.Color) {
	// Calculate actual screen position based on glass and brick dimensions
	screenX := g.left + (cx * brickWidth)
	screenY := g.top + (cy * brickHeight)

	// Set up the cell style with optional foreground and background colors
	style := tcell.StyleDefault.Foreground(fgColor).Background(bgColor)

	// Iterate through the piece to draw it on the screen
	for dy := range brickHeight {
		sx := screenX
		for range brickWidth {
			line := piece[dy]
			for _, r := range line {
				if sx < screenX+brickWidth {
					g.screen.SetContent(sx, screenY+dy, r, nil, style)
					sx++
				}
			}
		}
	}
}

// drawTetromino renders tetromino pieces, both active and shadow, based on given parameters
func (g *Game) drawTetromino(position [2]int, piece []string, color tcell.Color) {
	for i, row := range g.tetromino.GetCurrentShape() {
		for j, cell := range row {
			if cell == 1 {
				x := (position[1] + j)
				y := position[0] + i
				g.drawCell(x, y, piece, color, tcell.ColorDefault)
			}
		}
	}
}

// drawActiveTetromino renders the active tetromino on the grid
func (g *Game) drawActiveTetromino() {
	if g.status == StatusInGame || g.status == StatusPaused {
		g.drawTetromino(g.tetromino.position, CellTetromino, tcell.ColorYellow)
	}
}

// drawShadowTetromino calculates and renders the shadow piece on the grid
func (g *Game) drawShadowTetromino() {
	if !g.config.ShowShadow {
		return
	}

	if g.status == StatusInGame || g.status == StatusPaused {
		shadowPosition := g.CalculateShadowPosition()
		g.drawTetromino(shadowPosition, CellShadow, tcell.ColorGray)
	}
}
func (g *Game) drawGameInfo() {
	highScoreText := fmt.Sprintf("High Score: %d", g.config.HighScore)

	scoreText := fmt.Sprintf("Score: %d", g.score)
	levelText := fmt.Sprintf("Level: %d", g.level)
	linesText := fmt.Sprintf("Lines: %d", g.linesCleared)
	shadowText := fmt.Sprintf("[S]hadow: %s", boolToStr(g.config.ShowShadow))

	g.drawText(2+(glassCols*brickWidth), 0, highScoreText, tcell.ColorDarkGreen)
	g.drawText(2+(glassCols*brickWidth), 2, scoreText, tcell.ColorGreen)
	g.drawText(2+(glassCols*brickWidth), 3, linesText, tcell.ColorGreen)
	g.drawText(2+(glassCols*brickWidth), 4, levelText, tcell.ColorGreen)
	g.drawText(2+(glassCols*brickWidth), 6, shadowText, tcell.ColorGreen)
}

func boolToStr(b bool) string {
	if b {
		return "on"
	}
	return "off"
}

func (g *Game) drawGameStatusMessages() {
	var msgs []string
	switch g.status {
	case StatusGameOver:
		msgs = []string{"Game Over", "Play again? Y/n"}
	case StatusPaused:
		msgs = []string{"Paused", "Resume press 'P'"}
	}
	g.drawCenteredMessages(msgs, tcell.ColorRed)
}

func (g *Game) drawCenteredMessages(msgs []string, color tcell.Color) {
	for i, msg := range msgs {
		g.drawText(((glassCols*2)-len(msg))/2, (glassRows/2)-(len(msgs)/2)+i, msg, color)
	}
}

func (g *Game) drawTemporaryMessage() {
	if len(g.message) > 0 {
		for i, msg := range g.message {
			g.drawText(2+(glassCols*brickWidth), 2+infoHeight+i, msg, tcell.ColorRed)
		}
	}
}

func (g *Game) drawGameBorder() {
	style := tcell.StyleDefault.Foreground(tcell.Color100)
	for y := -2; y < gameHeight; y++ {
		for x := -2; x < gameWidth+1; x++ {
			ch := ' '
			if x == -2 && y == -2 {
				ch = '┌'
			} else if x == -2 && y+1 == gameHeight {
				ch = '└'
			} else if x == gameWidth && y == -2 {
				ch = '┐'
			} else if x == gameWidth && y+1 == gameHeight {
				ch = '┘'
			} else if x == -2 || x == gameWidth {
				ch = '│'
			} else if y == -2 || y+1 == gameHeight {
				ch = '─'
			}

			g.screen.SetContent(g.left+x, g.top+y, ch, nil, style)
		}
	}
}

// drawText draws text at a specific location using the specified color
func (g *Game) drawText(x, y int, text string, color tcell.Color) {
	style := tcell.StyleDefault.Foreground(color)
	pos := 0
	for _, c := range text {
		g.screen.SetContent(g.left+x+pos, g.top+y, c, nil, style)
		pos++
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
	g.message = nil
	g.softDrop = false

	g.status = StatusInGame // Reset status to in game

	// Restart game loop, spawning a new tetromino first
	g.spawnTetromino()
}

// Calculate the shadow position by simulating a hard drop
func (g *Game) CalculateShadowPosition() [2]int {
	shadowPosition := g.tetromino.position
	for g.grid.CanMove(g.tetromino, [2]int{shadowPosition[0] + 1, shadowPosition[1]}) {
		shadowPosition[0]++
	}
	return shadowPosition
}

// Spawn the next tetromino and check for game over
func (g *Game) spawnTetromino() {
	tetrominoType := string(tetrominoTypes[g.rng.Intn(len(tetrominoTypes))])
	tetromino, err := NewTetromino(tetrominoType)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	tetromino.position = [2]int{-1, glassCols/2 - len(tetromino.GetCurrentShape()[0])/2}

	// Check if spawn position is occupied
	if g.grid.CanMove(tetromino, tetromino.position) {
		g.tetromino = tetromino
	} else {
		g.status = StatusGameOver // Update status to game over
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
				y, x := pos[0]+i, pos[1]+j
				if y >= 0 {
					g.grid.cells[y][x] = CellStateFilled
				}
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
}

// Update level and manage message display
func (g *Game) updateLevel(linesCleared int) {
	g.linesCleared += linesCleared
	if g.linesCleared >= 10 {
		g.level++
		g.linesCleared -= 10
		g.message = []string{"  Level Up!", fmt.Sprintf("  New Level: %d", g.level)}
		g.setMessageTimer(3 * time.Second) // Display message for 3 seconds
	}
}

// Set a message to display temporarily
func (g *Game) setMessageTimer(duration time.Duration) {
	if g.messageTimer != nil {
		g.messageTimer.Stop()
	}
	g.messageTimer = time.AfterFunc(duration, func() {
		g.message = nil
		g.Render()
	})
}

// hardDrop instantly places the tetromino to the lowest possible position
func (g *Game) hardDrop() {
	for g.grid.MoveTetromino(g.tetromino, "down") {
	}
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
