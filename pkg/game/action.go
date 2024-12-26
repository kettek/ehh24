package game

import "math"

type Action interface {
	Apply(t *Thinger)
	Done() bool
}

type ActionLook struct {
	LookX, LookY float64
	ShouldFace   bool
}

func (a *ActionLook) Done() bool {
	return true
}

func (a *ActionLook) Apply(t *Thinger) {
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
}

type ActionPosition struct {
	X, Y float64
}

func (a *ActionPosition) Done() bool {
	return true
}

func (a *ActionPosition) Apply(t *Thinger) {
	t.X = a.X
	t.Y = a.Y
}

type ActionMoveTo struct {
	X, Y, Speed float64
	done        bool
}

func (a *ActionMoveTo) Apply(t *Thinger) {
	dx := a.X - t.X
	dy := a.Y - t.Y
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist < a.Speed {
		t.X = a.X
		t.Y = a.Y
		a.done = true
		t.Animation("center")
		t.walking = false
		t.walkTicker = 0
		t.faceUp = false
		return
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
	t.X += dx / dist * a.Speed
	t.Y += dy / dist * a.Speed * 0.6
}

func (a *ActionMoveTo) Done() bool {
	return a.done
}
