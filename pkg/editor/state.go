package editor

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ehh24/pkg/statemachine"
)

// State is the editor state.
type State struct {
}

// NewState creates a new editor state.
func NewState() *State {
	return &State{}
}

// Update updates the editor state.
func (s *State) Update() statemachine.State {
	return nil
}

// Draw draws the editor state.
func (s *State) Draw(screen *ebiten.Image) {
}

func (s *State) Layout(ow, oh int) (int, int) {
	return ow, oh
}
