package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"tic-tac-toe/ui"
)

//go:embed assets/icon.png
var iconBytes []byte

func main() {
	// Create a new game UI
	gameUI := ui.NewGameUI()

	// Try to load and set the icon
	icon, _, err := image.Decode(bytes.NewReader(iconBytes))
	if err != nil {
		log.Printf("Warning: Could not load icon: %v", err)
		// Continue without an icon
	} else {
		ebiten.SetWindowIcon([]image.Image{icon})
	}

	// Set a reasonable default window size with a 1:1 aspect ratio
	ebiten.SetWindowSize(600, 600)
	ebiten.SetFullscreen(false)
	ebiten.SetWindowTitle("Tic Tac Toe")

	// Start the game loop
	if err := ebiten.RunGame(gameUI); err != nil {
		log.Fatal(err)
	}
}
