package game

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"net/http"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type ScoreboardScene struct {
	scores    []string
	isLoading bool
	backBtn   *Button
}

func NewScoreboardScene() *ScoreboardScene {
	s := &ScoreboardScene{
		isLoading: true,
		backBtn: &Button{
			Bounds: image.Rect(20, 20, 100, 52),
			Label:  "Back",
		},
	}
	go s.fetchScores()
	return s
}

func (s *ScoreboardScene) fetchScores() {
	defer func() { s.isLoading = false }()

	resp, err := http.Get(apiURL)
	if err != nil {
		s.scores = []string{"Failed to load scores."}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.scores = []string{"Failed to read response."}
		return
	}

	var records []ScoreRecord
	if err := json.Unmarshal(body, &records); err != nil {
		s.scores = []string{"Failed to parse response."}
		return
	}

	for i, r := range records {
		s.scores = append(s.scores, fmt.Sprintf("%d. %s: %d", i+1, r.Nickname, r.Score))
	}
}

func (s *ScoreboardScene) Update() (Scene, error) {
	if s.backBtn.Clicked() {
		playCancelSound()
		return NewMenuScene(), nil
	}
	return nil, nil
}

func (s *ScoreboardScene) Draw(screen *ebiten.Image) {
	s.backBtn.Draw(screen, fontFace)

	op := &text.DrawOptions{}
	op.GeoM.Translate(20, 80)
	op.LineSpacing = lineSpacing()

	if s.isLoading {
		text.Draw(screen, "Loading...", fontFace, op)
		return
	}

	if len(s.scores) == 0 {
		text.Draw(screen, "No scores yet.", fontFace, op)
		return
	}

	for i, sc := range s.scores {
		op.GeoM.Reset()
		op.GeoM.Translate(20, float64(80+i*24))
		text.Draw(screen, sc, fontFace, op)
	}
}
