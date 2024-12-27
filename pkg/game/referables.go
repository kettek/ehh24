package game

import "github.com/kettek/ehh24/pkg/game/context"

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
	Update(ctx *context.Game) []Change
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
