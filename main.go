package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/haruki7049/iai/internal/game"
)

func main() {
	ebiten.SetWindowSize(game.WINDOW_WIDTH, game.WINDOW_HEIGHT)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	if err := ebiten.RunGame(&game.Game{}); err != nil {
		log.Fatal(err)
	}
}
