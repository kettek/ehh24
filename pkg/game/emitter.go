package game

import (
	"math/rand/v2"

	"github.com/kettek/ehh24/pkg/game/ables"
	"github.com/kettek/ehh24/pkg/game/context"
)

// Emitter sprays stuff.
type Emitter struct {
	ables.Positionable
	ables.IDable
	ables.Tagable
	ables.Priorityable
	particles []Particle
}

type Particle struct {
	ables.Positionable
	lifetime int
}

// NewEmitter makes a staticer, wow.
func NewEmitter(name string) *Emitter {
	return &Emitter{
		Positionable: ables.MakePositionable(32+rand.Float64()*256, 32+rand.Float64()*256),
	}
}

// Update updates the particles
func (t *Emitter) Update(ctx *ContextGame) []Change {
	return nil
}

// Draw draws the emitter bits.
func (t *Emitter) Draw(ctx *context.Draw) {
	// TODODO
}
