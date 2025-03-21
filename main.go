package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"tic-tac-toe/ui"
)

func main() {
	// Create a new game UI
	gameUI := ui.NewGameUI()

	// Set a reasonable default window size with a 1:1 aspect ratio
	ebiten.SetWindowSize(600, 600)
	ebiten.SetFullscreen(false)
	ebiten.SetWindowTitle("Tic Tac Toe")

	// Start the game loop
	if err := ebiten.RunGame(gameUI); err != nil {
		log.Fatal(err)
	}
}
