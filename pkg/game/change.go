package game

type Change interface {
	Apply(g *Game)
}

type ChangeDarknessOverlay struct {
	X, Y  float64
	Angle float64
}

func (c *ChangeDarknessOverlay) Apply(g *Game) {
	g.darknessOverlay.X = c.X
	g.darknessOverlay.Y = c.Y
	g.darknessOverlay.TargetAngle = c.Angle
}
