package game

import (
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ehh24/pkg/game/context"
)

// Referable refers to anything in za warudo that can be referred to. This will generally be types that implement Updateable and Drawable.
type Referable interface {
	Tag() string
	ID() int
}

// Referables is a collection of Referable objects.
type Referables []Referable

// ByTag returns a new Referables object that contains only the Referables with the given tag.
func (r Referables) ByTag(tag string) Referables {
	var res Referables
	for _, t := range r {
		if t.Tag() == tag {
			res = append(res, t)
		}
	}
	return res
}

// ByFirstTag just returns the first Referable found by tag.
func (r Referables) ByFirstTag(tag string) Referable {
	for _, t := range r {
		if t.Tag() == tag {
			return t
		}
	}
	return nil
}

// ByID returns the Referable with the given ID.
func (r Referables) ByID(id int) Referable {
	for _, t := range r {
		if t.ID() == id {
			return t
		}
	}
	return nil
}

// Updateable refers to anything in za warudo that can be updated.
type Updateable interface {
	Update(ctx *ContextGame) []Change
}

// Updateables returns a list of all the Updateable objects in the Referables.
func (r Referables) Updateables() []Updateable {
	var res []Updateable
	for _, t := range r {
		if u, ok := t.(Updateable); ok {
			res = append(res, u)
		}
	}
	return res
}

// Drawable refers to anything in za warudo that can be drawn.
type Drawable interface {
	Draw(ctx *context.Draw)
	Tag() string
	ID() int
	X() float64
	Y() float64
	Priority() int
	SetOffset(int)
}

// Drawables returns a list of all the Drawable objects in the Referables.
func (r Referables) Drawables() []Drawable {
	var res []Drawable
	for _, t := range r {
		if d, ok := t.(Drawable); ok {
			res = append(res, d)
		}
	}
	return res
}

// SortedDrawables returns a list of all the Drawables in the Referables, sorted by Priority.
func (r Referables) SortedDrawables() []Drawable {
	drawables := r.Drawables()
	slices.SortFunc(drawables, func(a, b Drawable) int {
		return a.Priority() - b.Priority()
	})
	return drawables
}

// Overlayable represents a drawable that has an ebiten.Image
type Overlayable interface {
	Drawable
	Updateable
	DrawTo(*ebiten.Image)
	Resize(int, int)
}

// Overlays returns a list of all the Overlayable objects in the Referables.
func (r Referables) Overlays() []Overlayable {
	var res []Overlayable
	for _, t := range r {
		if d, ok := t.(Overlayable); ok {
			res = append(res, d)
		}
	}
	return res
}

// Debugable is a referable that can be debuggied.
type Debugable interface {
	String() string
	X() float64
	Y() float64
}

// Debugables returns a list of all the Debugable objects in the Referables.
func (r Referables) Debugables() []Debugable {
	var res []Debugable
	for _, t := range r {
		if d, ok := t.(Debugable); ok {
			res = append(res, d)
		}
	}
	return res
}
