package game

import (
	"fmt"
	"image/color"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type gameState int

const (
	gameStateWaiting gameState = iota
	gameStateMistake
	gameStateSlash
	gameStateClear
	gameStateDead
)

const (
	slashWaitMin               = 5 * time.Second
	slashWaitMax               = 15 * time.Second
	initialReactionLimit       = 500 * time.Millisecond
	reactionLimitFloor         = 10 * time.Millisecond
	reactionLimitStepLarge     = 50 * time.Millisecond
	reactionLimitStepThreshold = 50 * time.Millisecond
	reactionLimitStepSmall     = 10 * time.Millisecond
	mistakeDisplayDuration     = 1500 * time.Millisecond
	clearDisplayDuration       = 5 * time.Second
	deadFadeDuration           = 3 * time.Second
)

type GameScene struct {
	state         gameState
	stateTime     time.Time
	slashWait     time.Duration
	reactionLimit time.Duration
	score         int // Track the number of clears
}

func NewGameScene() *GameScene {
	return &GameScene{
		state:         gameStateWaiting,
		stateTime:     time.Now(),
		slashWait:     randomSlashWait(),
		reactionLimit: initialReactionLimit,
		score:         0, // Initialize score
	}
}

func randomSlashWait() time.Duration {
	span := slashWaitMax - slashWaitMin
	return slashWaitMin + time.Duration(rand.Int64N(int64(span)+1))
}

func nextReactionLimit(current time.Duration) time.Duration {
	var next time.Duration
	if current > reactionLimitStepThreshold {
		next = max(current-reactionLimitStepLarge, reactionLimitStepThreshold)
	} else {
		next = current - reactionLimitStepSmall
	}
	if next < reactionLimitFloor {
		next = reactionLimitFloor
	}
	return next
}

func (s *GameScene) Update() (Scene, error) {
	switch s.state {
	case gameStateWaiting:
		switch {
		case inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft):
			playCutSound()
			s.enterState(gameStateMistake)
		case time.Since(s.stateTime) >= s.slashWait:
			playCutSound()
			s.enterState(gameStateSlash)
		}

	case gameStateMistake:
		if time.Since(s.stateTime) >= mistakeDisplayDuration {
			playCutSound()
			s.enterState(gameStateDead)
		}

	case gameStateSlash:
		clicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
		elapsed := time.Since(s.stateTime)

		switch {
		case clicked && elapsed < s.reactionLimit:
			playCutSound()
			s.enterState(gameStateClear)
		case clicked || elapsed >= s.reactionLimit:
			playCutSound()
			s.enterState(gameStateDead)
		}

	case gameStateClear:
		if time.Since(s.stateTime) >= clearDisplayDuration {
			s.score++ // Increment score on clear
			s.reactionLimit = nextReactionLimit(s.reactionLimit)
			s.slashWait = randomSlashWait()

			s.enterState(gameStateWaiting)
		}

	case gameStateDead:
		if time.Since(s.stateTime) >= deadFadeDuration {
			return NewRegisterScene(s.score), nil // Transition to register scene
		}
	}

	return nil, nil
}

func (s *GameScene) enterState(next gameState) {
	s.state = next
	s.stateTime = time.Now()

	if s.state == gameStateWaiting {
		playDecisionSound()
	}
}

func (s *GameScene) Draw(screen *ebiten.Image) {
	screen.Fill(s.backgroundColor())

	if msg := s.message(); msg != "" {
		op := &text.DrawOptions{}
		op.GeoM.Translate(20, 20)
		op.LineSpacing = lineSpacing()
		text.Draw(screen, msg, fontFace, op)
	}

	if s.state == gameStateSlash {
		s.drawRemainingMillis(screen)
	}
}

func (s *GameScene) drawRemainingMillis(screen *ebiten.Image) {
	remaining := max(s.reactionLimit-time.Since(s.stateTime), 0)
	msg := fmt.Sprintf("%d ms", remaining.Milliseconds())

	w, h := text.Measure(msg, fontFace, lineSpacing())
	bounds := screen.Bounds()

	op := &text.DrawOptions{}
	op.GeoM.Translate(
		float64(bounds.Dx())/2-w/2,
		float64(bounds.Dy())/2-h/2,
	)
	op.LineSpacing = lineSpacing()
	text.Draw(screen, msg, fontFace, op)
}

func (s *GameScene) message() string {
	switch s.state {
	case gameStateMistake:
		return "Mistake!!"
	case gameStateSlash:
		return "Slash!!"
	case gameStateClear:
		return "Clear!! Next game will start right away..."
	case gameStateDead:
		return "You are dead."
	default:
		return ""
	}
}

func (s *GameScene) backgroundColor() color.Color {
	if s.state != gameStateDead {
		return color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
	}

	progress := float64(time.Since(s.stateTime)) / float64(deadFadeDuration)
	if progress > 1 {
		progress = 1
	}

	return color.RGBA{
		R: uint8(0xff * progress),
		G: 0x00,
		B: 0x00,
		A: 0xff,
	}
}
