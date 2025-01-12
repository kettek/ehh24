package game

import (
	"fmt"
	"strings"

	"github.com/kettek/ehh24/pkg/res"
)

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
	parts := strings.SplitN(c.Place, ":", 2)
	var placeName string
	var enter string

	if len(parts) > 0 {
		placeName = parts[0]
	}
	if len(parts) > 1 {
		enter = parts[1]
	}

	if placeName == "" {
		return
	}

	if _, ok := ctx.Places[placeName]; ok {
		ctx.Place = ctx.Places[placeName]
	} else {
		place := NewPlace(placeName)
		ctx.Places[placeName] = place
		ctx.Place = place
	}
	ctx.Place.referables = append(ctx.Place.referables, NewFadeInOverlay(int(ctx.Width), int(ctx.Height), 50))
	// Move player into position.
	if enter != "" {
		if area := ctx.Place.GetAreaByFirstTag(enter); area != nil {
			if pl, ok := ctx.Referables.ByFirstTag("qi").(*Thinger); ok {
				x, y := area.Center()
				pl.SetX(x)
				pl.SetY(y)
			}
		}
	}
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

type ChangeUse struct {
	Tag string
}

func (c *ChangeUse) Apply(ctx *ContextGame) {
	fmt.Println("... using ", c.Tag)
	// Alright, let's see what the given area does.
	if area := ctx.Place.GetAreaByFirstTag(c.Tag); area != nil {
		targets := strings.Split(area.original.TargetTag, ";")
		actions := strings.Split(area.original.TargetAction, ";")
		for i, target := range targets {
			var act string
			if i < len(actions) {
				act = actions[i]
			} else if len(actions) > 0 {
				act = actions[0]
			}
			if area2 := ctx.Place.GetAreaByFirstTag(target); area2 != nil {
				if act == "del" {
					ctx.Place.RemoveAreaByFirstTag(target)
				}
			}
			// Check referables too, I guess.
			if ref := ctx.Place.referables.ByFirstTag(target); ref != nil {
				if act == "del" {
					ctx.Place.referables.RemoveByFirstTag(target)
				}
			}
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

type ChangeRemoveReferable struct {
	ID int
}

func (c *ChangeRemoveReferable) Apply(ctx *ContextGame) {
	ctx.Referables.RemoveByID(c.ID)
}
