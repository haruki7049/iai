package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// MenuScene is the in-game menu screen, reached from the title screen.
type MenuScene struct {
	gameButton *Button
}

// NewMenuScene creates a new MenuScene.
func NewMenuScene() *MenuScene {
	return &MenuScene{
		gameButton: &Button{
			Bounds: image.Rect(20, 60, 140, 92),
			Label:  "Game",
		},
	}
}

func (s *MenuScene) Update() (Scene, error) {
	if s.gameButton.Clicked() {
		return NewGameScene(), nil
	}
	return nil, nil
}

func (s *MenuScene) Draw(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(20, 20)
	op.LineSpacing = fontFace.Metrics().HLineGap + fontFace.Metrics().HAscent + fontFace.Metrics().HDescent
	text.Draw(screen, "Menu", fontFace, op)

	s.gameButton.Draw(screen, fontFace)
}
