package main

import (
	"log"
	"sync"

	"github.com/gdamore/tcell"
)

func main() {
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Println("Could not load config, using defaults:", err)
		config = NewConfig()
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := screen.Init(); err != nil {
		panic(err)
	}
	defer screen.Fini()

	var wg sync.WaitGroup
	wg.Add(2)

	// Create a new random number generator with a seed based on the current time
	game := NewGameWithScreen(screen, config)
	go game.Run(&wg)
	go game.ListenForInput(&wg)

	wg.Wait()

	if err := config.SaveConfig("config.json"); err != nil {
		log.Println("Error saving config:", err)
	}
}
