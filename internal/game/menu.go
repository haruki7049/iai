package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type MenuScene struct {
	gameButton       *Button
	exitButton       *Button
	scoreboardButton *Button
}

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
		scoreboardButton: &Button{
			Bounds: image.Rect(20, 140, 180, 172),
			Label:  "Score Board",
		},
	}
}

func (s *MenuScene) Update() (Scene, error) {
	if s.gameButton.Clicked() {
		playDecisionSound()
		return NewGameScene(), nil
	}
	if s.exitButton.Clicked() {
		return nil, ebiten.Termination
	}
	if s.scoreboardButton.Clicked() {
		playDecisionSound()
		return NewScoreboardScene(), nil
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
	s.scoreboardButton.Draw(screen, fontFace)
}
