package main

const (
	glassCols = 10
	glassRows = 20

	brickWidth  = 2
	brickHeight = 1

	infoWidth  = 20
	infoHeight = 8

	gameWidth  = 2 + (glassCols * brickWidth) + infoWidth
	gameHeight = 1 + (glassRows * brickHeight)
)

var (
	CellGridEmpty  = []string{"˓˒"}
	CellGridFilled = []string{"░░"}
	CellTetromino  = []string{"[]"}
	CellShadow     = []string{"[]"}
)
