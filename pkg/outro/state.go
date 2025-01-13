package outro

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ehh24/pkg/game/context"
	"github.com/kettek/ehh24/pkg/statemachine"
)

type Scene interface {
	Init()
	Update() Scene
	Draw(ctx context.Draw)
}

// State be our intro.
type State struct {
	scene       Scene
	fadeIn      int
	fadeOut     int
	width       int
	height      int
	drawContext context.Draw
}

// NewState creates a new state.
func NewState() *State {
	return &State{
		fadeIn:  100,
		fadeOut: 100,
	}
}

// Init is called when the state is to be first entered.
func (s *State) Init() {
	s.scene = &Scene1{}
	s.scene.Init()
}

// Update updates the state.
func (s *State) Update() statemachine.State {
	s.scene = s.scene.Update()
	return nil
}

// Draw draws the state.
func (s *State) Draw(screen *ebiten.Image) {
	s.drawContext.Target = screen
	s.drawContext.Op = &ebiten.DrawImageOptions{}
	s.drawContext.Op.GeoM.Scale(3, 3)
	s.scene.Draw(s.drawContext)
}

// Layout does a layout.
func (s *State) Layout(ow, oh int) (int, int) {
	s.width, s.height = ow, oh
	return ow, oh
}
