package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"tic-tac-toe/ui"
)

func main() {
	// Create a new game UI (no params now!)
	gameUI := ui.NewGameUI()

	// Configure the Ebiten window.
	ebiten.SetWindowSize(600, 700) // adjust height for the new menu/buttons space
	ebiten.SetWindowTitle("Tic Tac Toe")

	// Start the game loop.
	if err := ebiten.RunGame(gameUI); err != nil {
		log.Fatal(err)
	}
}
