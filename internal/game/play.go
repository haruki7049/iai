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

// gameState represents the current phase of the iai-giri mini-game.
type gameState int

const (
	gameStateWaiting gameState = iota // waiting for the random cue delay
	gameStateMistake                  // player clicked too early (false start)
	gameStateSlash                    // "Slash!!" is shown, waiting for player reaction
	gameStateClear                    // player reacted in time
	gameStateDead                     // player reacted too late, or made a mistake
)

const (
	// slashWaitMin and slashWaitMax bound the random delay before
	// "Slash!!" is shown. This stays the same across rounds.
	slashWaitMin = 5 * time.Second
	slashWaitMax = 15 * time.Second

	// initialReactionLimit is the reaction window for the very first
	// round. It shrinks after every successful clear.
	initialReactionLimit = 500 * time.Millisecond

	// reactionLimitFloor is the lowest the reaction window can ever
	// shrink to.
	reactionLimitFloor = 10 * time.Millisecond

	// reactionLimitStepLarge is the per-clear decrement while the
	// reaction window is still above reactionLimitStepThreshold.
	reactionLimitStepLarge = 50 * time.Millisecond

	// reactionLimitStepThreshold is the reaction window value at which
	// the decrement step switches from reactionLimitStepLarge to
	// reactionLimitStepSmall.
	reactionLimitStepThreshold = 50 * time.Millisecond

	// reactionLimitStepSmall is the per-clear decrement once the
	// reaction window has reached reactionLimitStepThreshold.
	reactionLimitStepSmall = 10 * time.Millisecond

	// mistakeDisplayDuration is how long "Mistake!!" stays on screen
	// before the false start counts as a death.
	mistakeDisplayDuration = 1500 * time.Millisecond

	// clearDisplayDuration is how long "Clear!!" stays on screen
	// before the next round begins.
	clearDisplayDuration = 5 * time.Second

	// deadFadeDuration is how long the background takes to fade from
	// black to red after the player fails, before returning to the menu.
	deadFadeDuration = 3 * time.Second
)

// GameScene is the iai-giri (quick-draw slash) mini-game screen.
type GameScene struct {
	state         gameState
	stateTime     time.Time     // when the current state started
	slashWait     time.Duration // randomly chosen delay before the slash cue
	reactionLimit time.Duration // current round's reaction window
}

// NewGameScene creates a new GameScene and rolls the random delay before
// the slash cue appears.
func NewGameScene() *GameScene {
	return &GameScene{
		state:         gameStateWaiting,
		stateTime:     time.Now(),
		slashWait:     randomSlashWait(),
		reactionLimit: initialReactionLimit,
	}
}

// randomSlashWait returns a random duration in [slashWaitMin, slashWaitMax].
func randomSlashWait() time.Duration {
	span := slashWaitMax - slashWaitMin
	return slashWaitMin + time.Duration(rand.Int64N(int64(span)+1))
}

// nextReactionLimit returns the reaction window for the round after a
// successful clear. It shrinks by reactionLimitStepLarge per clear until
// it reaches reactionLimitStepThreshold, then shrinks by
// reactionLimitStepSmall, never going below reactionLimitFloor.
func nextReactionLimit(current time.Duration) time.Duration {
	var next time.Duration
	if current > reactionLimitStepThreshold {
		next = current - reactionLimitStepLarge
		if next < reactionLimitStepThreshold {
			next = reactionLimitStepThreshold
		}
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
			// Clicked before "Slash!!" appeared: a false start.
			s.enterState(gameStateMistake)
		case time.Since(s.stateTime) >= s.slashWait:
			s.enterState(gameStateSlash)
		}

	case gameStateMistake:
		if time.Since(s.stateTime) >= mistakeDisplayDuration {
			s.enterState(gameStateDead)
		}

	case gameStateSlash:
		clicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
		elapsed := time.Since(s.stateTime)

		switch {
		case clicked && elapsed < s.reactionLimit:
			s.enterState(gameStateClear)
		case clicked || elapsed >= s.reactionLimit:
			s.enterState(gameStateDead)
		}

	case gameStateClear:
		if time.Since(s.stateTime) >= clearDisplayDuration {
			// Start the next round: shrink the reaction window, roll
			// a fresh random wait, and go back to waiting.
			s.reactionLimit = nextReactionLimit(s.reactionLimit)
			s.slashWait = randomSlashWait()
			s.enterState(gameStateWaiting)
		}

	case gameStateDead:
		if time.Since(s.stateTime) >= deadFadeDuration {
			return NewMenuScene(), nil
		}
	}

	return nil, nil
}

// enterState switches to the given state and resets the state timer.
func (s *GameScene) enterState(next gameState) {
	s.state = next
	s.stateTime = time.Now()
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

// drawRemainingMillis renders the time left, in milliseconds, before the
// reaction window for the current round runs out. It is only shown while
// the player's input is being accepted.
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

// message returns the text to show for the current state.
func (s *GameScene) message() string {
	switch s.state {
	case gameStateMistake:
		return "Mistake!!"
	case gameStateSlash:
		return "Slash!!"
	case gameStateClear:
		return "Clear!!"
	case gameStateDead:
		return "You are dead."
	default:
		return ""
	}
}

// backgroundColor returns the background fill color for the current
// state. Only the dead state animates, fading from black to red over
// deadFadeDuration.
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
