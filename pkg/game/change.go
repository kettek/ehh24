package game

// Change is a requested change to the game state originating from an action.
type Change interface {
	Apply(g *ContextGame)
}

// ChangeVisibilityOverlay changes the visibility overlay.
type ChangeVisibilityOverlay struct {
	X, Y  float64
	Angle float64
}

// Apply applies changes to the visibility overlay.
func (c *ChangeVisibilityOverlay) Apply(ctx *ContextGame) {
	if v, ok := ctx.Referables.ByFirstTag("visibility").(*VisibilityOverlay); ok {
		v.SetX(c.X)
		v.SetY(c.Y)
		v.TargetAngle = c.Angle
	}
}

type ChangeTravel struct {
	Place string
}

func (c *ChangeTravel) Apply(ctx *ContextGame) {
	if _, ok := ctx.Places[c.Place]; ok {
		ctx.Place = ctx.Places[c.Place]
		return
	}
	place := NewPlace(c.Place)
	ctx.Places[c.Place] = place
	ctx.Place = place
}
