package game

import "math"

// Action represents an action that can be applied to a thinger, optionally looping with Done() checks.
type Action interface {
	Apply(t *Thinger) []Change
	Done() bool
}

// ActionLook changes the look direction of a thinger.
type ActionLook struct {
	LookX, LookY float64
	ShouldFace   bool
}

// Done is true.
func (a *ActionLook) Done() bool {
	return true
}

// Apply changes the look direction of a thinger.
func (a *ActionLook) Apply(t *Thinger) (c []Change) {
	t.lookX = a.LookX
	t.lookY = a.LookY

	if a.ShouldFace {
		if t.lookX < -0.4 {
			t.Animation("left")
		} else if t.lookX > 0.5 {
			t.Animation("right")
		} else {
			t.Animation("center")
		}
		if t.lookY < -0.4 {
			t.faceUp = true
			t.lookUp = true
		} else {
			t.faceUp = false
			t.lookUp = false
		}
	}

	c = append(c, &ChangeVisibilityOverlay{
		X:     t.X(),
		Y:     t.Y() - float64(t.stax.Stax.SliceHeight)/2,
		Angle: math.Atan2(t.lookY+0.4, t.lookX), // Fix this hardcoded 0.4... it's the offset we need for eye position
	})

	return c
}

// ActionPosition changes the position of a thinger.
type ActionPosition struct {
	X, Y float64
}

// Done is true.
func (a *ActionPosition) Done() bool {
	return true
}

// Apply changes the position of a thinger.
func (a *ActionPosition) Apply(t *Thinger) []Change {
	t.SetX(a.X)
	t.SetY(a.Y)
	return nil
}

// ActionMoveTo moves a thinger to a position. This occurs over time.
type ActionMoveTo struct {
	X, Y, Speed float64
	done        bool
}

// Apply moves a thinger to a position over time.
func (a *ActionMoveTo) Apply(t *Thinger) (c []Change) {
	dx := a.X - t.X()
	dy := a.Y - t.Y()
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist < a.Speed {
		t.SetX(a.X)
		t.SetY(a.Y)
		a.done = true
		t.Animation("center")
		t.walking = false
		t.walkTicker = 0
		t.faceUp = false
		return nil
	}
	t.walking = true
	t.walkTicker++
	if dx < 0 {
		t.Animation("left")
	} else if dx > 0 {
		t.Animation("right")
	}
	if dy < 0 {
		t.faceUp = true
	} else {
		t.faceUp = false
	}
	t.SetX(t.X() + dx/dist*a.Speed)
	t.SetY(t.Y() + dy/dist*a.Speed*0.6)

	c = append(c, &ChangeVisibilityOverlay{
		X:     t.X(),
		Y:     t.Y() - float64(t.stax.Stax.SliceHeight)/2,
		Angle: math.Atan2(t.lookY+0.4, t.lookX), // Fix this hardcoded 0.4... it's the offset we need for eye position
	})

	return c
}

// Done is true when the thinger has moved to the target position.
func (a *ActionMoveTo) Done() bool {
	return a.done
}

type ActionFace struct {
	Radians float64
}

func (a *ActionFace) Apply(t *Thinger) []Change {
	diff := a.Radians - t.rotation
	if diff > math.Pi {
		diff = diff - math.Pi*2
	} else if diff < -math.Pi {
		diff = diff + math.Pi*2
	}
	if diff > 0.05 {
		t.rotation += diff / 10
	} else if diff < -0.05 {
		t.rotation += diff / 10
	}
	if t.rotation > math.Pi {
		t.rotation -= math.Pi * 2
	} else if t.rotation < -math.Pi {
		t.rotation += math.Pi * 2
	}
	return nil
}

func (a *ActionFace) Done() bool {
	return true
}

type ActionTravel struct {
	place string
}

func (a *ActionTravel) Apply(t *Thinger) []Change {
	return []Change{&ChangeTravel{Place: a.place}}
}

func (a *ActionTravel) Done() bool {
	return true
}
