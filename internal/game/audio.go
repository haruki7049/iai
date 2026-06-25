package game

import (
	"io"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"

	"github.com/haruki7049/iai/assets"
)

// sampleRate must match the sample rate the sound effects were generated with.
const sampleRate = 44100

var audioContext = audio.NewContext(sampleRate)

// soundPCM holds the raw, decoded PCM bytes for each loaded sound effect.
// audio.Context.NewPlayerFromBytes() expects exactly this format, so each WAV
// file is decoded once at startup instead of on every playback.
var soundPCM = map[string][]byte{}

func init() {
	loadSound("accept", "se/accept.wav")
	loadSound("deny", "se/deny.wav")
	loadSound("cut", "se/cut.wav")
}

// loadSound decodes the WAV file at path and stores its PCM bytes under name.
// A failure is logged rather than treated as fatal, so a missing sound asset
// does not crash the whole game; it simply stays silent.
func loadSound(name, path string) {
	f, err := assets.Assets.Open(path)
	if err != nil {
		log.Printf("sound %q: failed to open %q: %v", name, path, err)
		return
	}
	defer f.Close()

	stream, err := wav.DecodeWithSampleRate(sampleRate, f)
	if err != nil {
		log.Printf("sound %q: failed to decode %q: %v", name, path, err)
		return
	}

	pcm, err := io.ReadAll(stream)
	if err != nil {
		log.Printf("sound %q: failed to read %q: %v", name, path, err)
		return
	}

	soundPCM[name] = pcm
}

// playSound starts a fresh playback of the named sound effect. Each call
// creates a new player, so the same sound effect can overlap with itself
// (e.g. rapid menu navigation) without cutting the previous playback short.
func playSound(name string) {
	pcm, ok := soundPCM[name]
	if !ok {
		return
	}
	audioContext.NewPlayerFromBytes(pcm).Play()
}

// playDecisionSound plays the confirmation sound. Use it whenever the player
// confirms something: pressing a "proceed"-style button, submitting the
// nickname, or starting a new game.
func playDecisionSound() {
	playSound("accept")
}

// playCancelSound plays the cancel sound. Use it whenever the player goes
// back to a previous screen.
func playCancelSound() {
	playSound("deny")
}

// playCutSound plays the blade-cutting sound. Use it both at the moment the
// player must slash, and at the moment the player dies.
func playCutSound() {
	playSound("cut")
}
