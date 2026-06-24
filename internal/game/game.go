package game

import (
	"github.com/hajimehoshi/bitmapfont/v4"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// fontFace is the shared bitmap font face used by every scene.
var fontFace = text.NewGoXFace(bitmapfont.Face)

const WINDOW_WIDTH = 1280
const WINDOW_HEIGHT = 720

// Game is the root ebiten.Game implementation. It only owns the currently
// active Scene and delegates Update/Draw to it.
type Game struct {
	scene Scene
}

func (g *Game) Update() error {
	if g.scene == nil {
		g.scene = NewTitleScene()
	}

	next, err := g.scene.Update()
	if err != nil {
		return err
	}
	if next != nil {
		g.scene = next
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.scene == nil {
		return
	}
	g.scene.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth / 4, outsideHeight / 4
}
