package game

import "github.com/kettek/ehh24/pkg/res"

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

// ChangeAcquireItem finds the given area with the tag, deletes the area, and adds the item as an inventory to the player using Tag for its identifier and the area's Message for its name.
type ChangeAcquireItem struct {
	Tag string
}

// Apply applies the change to the game.
func (c *ChangeAcquireItem) Apply(ctx *ContextGame) {
	// Just get our player.
	if pl, ok := ctx.Referables.ByFirstTag("qi").(*Thinger); ok {
		if area := ctx.Place.GetAreaByFirstTag(c.Tag); area != nil {
			pl.AddItem(area.original.Text, area.original.Tag)
			// Delete the area...
			ctx.Place.RemoveAreaByFirstTag(c.Tag)
			// Also delete any associated referable in the map.
			ctx.Place.referables.RemoveByFirstTag(c.Tag)
		}
	}
}

type ChangeThingerPosition struct {
	Force   bool
	Thinger *Thinger
	X, Y    float64
}

func (c *ChangeThingerPosition) Apply(ctx *ContextGame) {
	if c.Force {
		c.Thinger.SetX(c.X)
		c.Thinger.SetY(c.Y)
		return
	}
	for _, area := range ctx.Place.areas {
		if area.ContainsPoint(c.X, c.Y) && area.original.Kind == res.PolygonKindBlock {
			return
		}
	}
	c.Thinger.SetX(c.X)
	c.Thinger.SetY(c.Y)
}
