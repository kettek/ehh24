package game

// Change is a requested change to the game state originating from an action.
type Change interface {
	Apply(g *State)
}

// ChangeVisibilityOverlay changes the visibility overlay.
type ChangeVisibilityOverlay struct {
	X, Y  float64
	Angle float64
}

// Apply applies changes to the visibility overlay.
func (c *ChangeVisibilityOverlay) Apply(g *State) {
	if v, ok := g.referables.ByFirstTag("visibility").(*VisibilityOverlay); ok {
		v.SetX(c.X)
		v.SetY(c.Y)
		v.TargetAngle = c.Angle
	}
}
