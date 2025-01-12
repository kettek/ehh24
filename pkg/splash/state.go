package splash

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ehh24/pkg/intro"
	"github.com/kettek/ehh24/pkg/res"
	"github.com/kettek/ehh24/pkg/statemachine"
)

// State be our intro.
type State struct {
	fadeIn  int
	fadeOut int
	width   int
	height  int
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
}

// Update updates the state.
func (s *State) Update() statemachine.State {
	if s.fadeIn > 0 {
		s.fadeIn--
	} else if s.fadeOut > 0 {
		s.fadeOut--
	} else {
		return intro.NewState()
	}
	return nil
}

// Draw draws the state.
func (s *State) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	if s.fadeIn > 0 {
		op.ColorScale.ScaleAlpha(1.0 - float32(s.fadeIn)/100)
	} else if s.fadeOut > 0 {
		op.ColorScale.ScaleAlpha(float32(s.fadeOut) / 100)
	} else {
		op.ColorScale.ScaleAlpha(0)
	}
	{
		ebi := res.Images["ebiten"]
		op.GeoM.Translate(float64(s.width)/2-float64(ebi.Bounds().Dy()/2), float64(s.height)/2-float64(ebi.Bounds().Dy())/2)
		screen.DrawImage(ebi, op)
	}
}

// Layout does a layout.
func (s *State) Layout(ow, oh int) (int, int) {
	s.width, s.height = ow, oh
	return ow, oh
}
