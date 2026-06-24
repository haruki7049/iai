package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// MenuScene is the in-game menu screen, reached from the title screen.
type MenuScene struct {
	gameButton *Button
	exitButton *Button
}

// NewMenuScene creates a new MenuScene.
func NewMenuScene() *MenuScene {
	return &MenuScene{
		gameButton: &Button{
			Bounds: image.Rect(20, 60, 140, 92),
			Label:  "Game",
		},
		exitButton: &Button{
			Bounds: image.Rect(20, 100, 140, 132),
			Label:  "Exit",
		},
	}
}

func (s *MenuScene) Update() (Scene, error) {
	if s.gameButton.Clicked() {
		return NewGameScene(), nil
	}
	if s.exitButton.Clicked() {
		// ebiten.Termination signals a normal, intentional exit; the
		// caller (main) treats it differently from a real error.
		return nil, ebiten.Termination
	}
	return nil, nil
}

func (s *MenuScene) Draw(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(20, 20)
	op.LineSpacing = lineSpacing()
	text.Draw(screen, "Menu", fontFace, op)

	s.gameButton.Draw(screen, fontFace)
	s.exitButton.Draw(screen, fontFace)
}
