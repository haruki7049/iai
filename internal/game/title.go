package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const TITLE = "iai"

// TitleScene is the first screen shown when the game starts. It shows the
// game title and a button that leads to the in-game menu.
type TitleScene struct {
	menuButton *Button
}

// NewTitleScene creates a new TitleScene.
func NewTitleScene() *TitleScene {
	return &TitleScene{
		menuButton: &Button{
			Bounds: image.Rect(20, 60, 140, 92),
			Label:  "Menu",
		},
	}
}

func (s *TitleScene) Update() (Scene, error) {
	if s.menuButton.Clicked() {
		playDecisionSound()
		return NewMenuScene(), nil
	}
	return nil, nil
}

func (s *TitleScene) Draw(screen *ebiten.Image) {
	titleOp := &text.DrawOptions{}
	titleOp.GeoM.Translate(20, 20)
	titleOp.LineSpacing = lineSpacing()
	text.Draw(screen, TITLE, fontFace, titleOp)

	s.menuButton.Draw(screen, fontFace)
}
