package game

import (
	"fmt"
	"math"
)

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
func (a *ActionPosition) Apply(t *Thinger) (c []Change) {
	c = append(c, &ChangeThingerPosition{
		Force:   true,
		Thinger: t,
		X:       a.X,
		Y:       a.Y,
	})

	return c
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

	c = append(c, &ChangeThingerPosition{
		Thinger: t,
		X:       t.X() + dx/dist*a.Speed,
		Y:       t.Y() + dy/dist*a.Speed*0.6,
	})

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

// ActionPickup moves towards a position and attempts to pick up something with a given tag.
type ActionPickup struct {
	ActionMoveTo
	Target string
}

// Apply does the thing.
func (a *ActionPickup) Apply(t *Thinger) []Change {
	dx := a.X - t.X()
	dy := a.Y - t.Y()
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist < 20 { // Eh... 10 seems good enough
		a.done = true
		return []Change{
			&ChangeAcquireItem{
				Tag: a.Target,
			},
		}
	}
	// Otherwise...
	return a.ActionMoveTo.Apply(t)
}

// Done also does a thing.
func (a *ActionPickup) Done() bool {
	return a.done
}

// ActionUse moves towards a position and attempts to use something with a given tag.
type ActionUse struct {
	ActionMoveTo
	Target string
	Item   string
}

// Apply does the thing.
func (a *ActionUse) Apply(t *Thinger) (c []Change) {
	dx := a.X - t.X()
	dy := a.Y - t.Y()
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist < 10 { // Eh... 10 seems good enough
		a.done = true
		c = append(c, &ChangeUse{
			Tag: a.Target,
		})
		if a.Item != "" {
			fmt.Println("remove", a.Item)
			c = append(c, &ChangeLoseItem{
				Tag: a.Item,
			})
		}
		return c
	}
	// Otherwise...
	return a.ActionMoveTo.Apply(t)
}

// Done also does a thing.
func (a *ActionUse) Done() bool {
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
	Place string
}

func (a *ActionTravel) Apply(t *Thinger) []Change {
	return []Change{&ChangeTravel{Place: a.Place}}
}

func (a *ActionTravel) Done() bool {
	return true
}

type ActionMonologue struct {
	Text  string
	Timer int
}

func (a *ActionMonologue) Apply(t *Thinger) []Change {
	t.monologue = a.Text
	a.Timer--
	if a.Timer <= 0 {
		t.monologue = ""
	}
	return nil
}

func (a *ActionMonologue) Done() bool {
	return a.Timer <= 0
}

type ActionState struct {
	State string
}

func (a *ActionState) Apply(t *Thinger) []Change {
	return []Change{&ChangeState{State: a.State}}
}

func (a *ActionState) Done() bool {
	return true
}
