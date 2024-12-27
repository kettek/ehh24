package game

// Change is a requested change to the game state originating from an action.
type Change interface {
	Apply(g *Game)
}

// ChangeVisibilityOverlay changes the visibility overlay.
type ChangeVisibilityOverlay struct {
	X, Y  float64
	Angle float64
}

// Apply applies changes to the visibility overlay.
func (c *ChangeVisibilityOverlay) Apply(g *Game) {
	if v, ok := g.referables.ByFirstTag("visibility").(*VisibilityOverlay); ok {
		v.X = c.X
		v.Y = c.Y
		v.TargetAngle = c.Angle
	}
}
