package statemachine

import "github.com/hajimehoshi/ebiten/v2"

// State is the interface for a game state.
type State interface {
	Init()
	Update() State
	Draw(screen *ebiten.Image)
	Layout(ow, oh int) (int, int)
}
