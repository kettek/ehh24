package statemachine

import "github.com/hajimehoshi/ebiten/v2"

// Machine is the state machine for the game, wow.
type Machine struct {
	state  State
	checks []func()
	w, h   int
}

// NewMachine does a thing.
func NewMachine(s State) *Machine {
	m := &Machine{
		state: s,
	}
	m.state.Init()
	return m
}

// Update updates the game state.
func (g *Machine) Update() error {
	for _, f := range g.checks {
		f()
	}
	next := g.state.Update()
	if next != nil {
		g.SetState(next)
		g.state = next
	}
	return nil
}

// Draw draws the game state.
func (g *Machine) Draw(screen *ebiten.Image) {
	g.state.Draw(screen)
}

// Layout does a layout.
func (g *Machine) Layout(ow, oh int) (int, int) {
	g.w, g.h = ow, oh
	return g.state.Layout(ow, oh)
}

// SetState sets the state.
func (g *Machine) SetState(s State) {
	g.state = s
	g.state.Init()
	g.state.Layout(g.w, g.h)
}

// AddCheck adds a check at the beginning of Update.
func (g *Machine) AddCheck(f func()) {
	g.checks = append(g.checks, f)
}

// ClearChecks clears all checks.
func (g *Machine) ClearChecks() {
	g.checks = []func(){}
}
