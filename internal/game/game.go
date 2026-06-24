package game

import (
	"github.com/hajimehoshi/bitmapfont/v4"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var fontFace = text.NewGoXFace(bitmapfont.Face)

const WINDOW_WIDTH = 1280
const WINDOW_HEIGHT = 720

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	textOp := &text.DrawOptions{}
	textOp.LineSpacing = fontFace.Metrics().HLineGap + fontFace.Metrics().HAscent + fontFace.Metrics().HDescent
	text.Draw(screen, "Hello!! It's \"iai\" game!!", fontFace, textOp)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth / 4, outsideHeight / 4
}
