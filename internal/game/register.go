package game

import (
	"bytes"
	"encoding/json"
	"image"
	"io"
	"log"
	"net/http"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// My actual Cloud Run URL
const apiURL = "https://iai-server-63587581912.europe-west1.run.app/scores"

type ScoreRecord struct {
	Nickname string `json:"nickname"`
	Score    int    `json:"score"`
}

type RegisterScene struct {
	score     int
	nickname  []rune
	isAsking  bool
	isSaving  bool
	isDone    bool
	yesButton *Button
	noButton  *Button
}

func NewRegisterScene(score int) *RegisterScene {
	return &RegisterScene{
		score:    score,
		isAsking: true,
		yesButton: &Button{
			Bounds: image.Rect(20, 60, 100, 92),
			Label:  "Yes",
		},
		noButton: &Button{
			Bounds: image.Rect(120, 60, 200, 92),
			Label:  "No",
		},
	}
}

func (s *RegisterScene) Update() (Scene, error) {
	if s.isDone {
		return NewMenuScene(), nil
	}

	if s.isAsking {
		if s.yesButton.Clicked() {
			s.isAsking = false
		}
		if s.noButton.Clicked() {
			return NewMenuScene(), nil
		}
		return nil, nil
	}

	if s.isSaving {
		return nil, nil
	}

	s.nickname = ebiten.AppendInputChars(s.nickname)
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(s.nickname) > 0 {
		s.nickname = s.nickname[:len(s.nickname)-1]
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && len(s.nickname) > 0 {
		s.isSaving = true
		go s.saveScore()
	}

	return nil, nil
}

func (s *RegisterScene) saveScore() {
	record := ScoreRecord{
		Nickname: string(s.nickname),
		Score:    s.score,
	}

	data, err := json.Marshal(record)
	if err == nil {
		resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(data))
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusCreated {
				// Read and log the error message from the server
				body, _ := io.ReadAll(resp.Body)
				log.Printf("Server error: %d %s\nBody: %s", resp.StatusCode, resp.Status, string(body))
			}
		} else {
			// Log network errors
			log.Printf("Network error: %v", err)
		}
	} else {
		log.Printf("JSON encode error: %v", err)
	}
	s.isDone = true
}

func (s *RegisterScene) Draw(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(20, 20)
	op.LineSpacing = lineSpacing()

	if s.isAsking {
		text.Draw(screen, "Register your score?", fontFace, op)
		s.yesButton.Draw(screen, fontFace)
		s.noButton.Draw(screen, fontFace)
		return
	}

	if s.isSaving {
		text.Draw(screen, "Saving...", fontFace, op)
		return
	}

	text.Draw(screen, "Enter nickname and press Enter:\n"+string(s.nickname), fontFace, op)
}
