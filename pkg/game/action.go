package game

type Action interface {
	Apply(t *Thinger)
}

type ActionLook struct {
	LookX, LookY float64
	ShouldFace   bool
}

func (a *ActionLook) Apply(t *Thinger) {
	t.lookX = a.LookX
	t.lookY = a.LookY

	if a.ShouldFace {
		if t.lookX < -0.1 {
			t.Animation("left")
		} else if t.lookX > 0.2 {
			t.Animation("right")
		} else {
			t.Animation("center")
		}
	}
}

type ActionPosition struct {
	X, Y float64
}

func (a *ActionPosition) Apply(t *Thinger) {
	t.X = a.X
	t.Y = a.Y
}
