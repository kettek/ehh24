package game

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ehh24/pkg/game/ables"
	"github.com/kettek/ehh24/pkg/game/context"
	"github.com/kettek/ehh24/pkg/res"
)

type Inventory struct {
	ables.IDable
	ables.Priorityable
	ables.Tagable
	ables.Positionable
	tw            float64
	th            float64
	width         float64
	height        float64
	targetTag     string // eh, whatever
	items         []InvItem
	hoveredName   string
	fade          int
	heldItemIndex int
}

const fadeMax = 40
const fadeMin = 20

type InvItem struct {
	item   ables.StorageItem
	staxer Staxer
}

func NewInventory(tag string) *Inventory {
	inv := &Inventory{
		targetTag: tag,
		fade:      fadeMin,
	}
	inv.SetPriority(ables.PriorityUI)
	inv.SetTag("inventory")

	inv.width = 11 * 4
	inv.height = 11 * 4

	return inv
}

func (inv *Inventory) ItemBounds(index int) (x, y, x2, y2 float64) {
	cx := inv.X() + inv.width/2
	cy := inv.Y() + inv.height/2

	ratio := float64(index) / 6
	angle := ratio*math.Pi*2 + math.Pi
	x = cx + math.Cos(angle)*16
	y = cy + math.Sin(angle)*16

	x -= 4
	y -= 4

	return x, y, x + 9, y + 9
}

func (inv *Inventory) SyncTo(storage ables.Storagable) {
	// Grow/shrink.
	if len(inv.items) < len(storage) {
		for i := len(inv.items); i < len(storage); i++ {
			inv.items = append(inv.items, InvItem{})
		}
	} else if len(inv.items) > len(storage) {
		inv.items = inv.items[:len(storage)]
	}
	// Sync items.
	for i, item := range storage {
		if item.Tag != inv.items[i].item.Tag {
			inv.items[i].item = item
			inv.items[i].staxer = NewStaxer(item.Tag)
		}
	}
}

func (inv *Inventory) Update(ctx *ContextGame) []Change {
	w, h := ctx.Size()
	inv.SetX(w/2 - inv.width/2)
	inv.SetY(h - inv.height - 8)

	// I guess just get our target and sync to it.
	if t, ok := ctx.Referables.ByFirstTag(inv.targetTag).(*Thinger); ok {
		inv.SyncTo(t.Storagable)
	}

	x, y := ctx.MousePosition()
	inv.hoveredName = ""
	if t, ok := ctx.Referables.ByFirstTag(inv.targetTag).(*Thinger); ok {
		if t.controller != nil {
			t.controller.Unblock()
		}
		pc := t.controller.(*PlayerController) // hackiness abounds.
		inv.heldItemIndex = -1
		for index, item := range t.Storagable {
			ix, iy, ix2, iy2 := inv.ItemBounds(index)
			if x >= ix && x <= ix2 && y >= iy && y <= iy2 {
				inv.hoveredName = item.Name
				// Get our cursor and show pickup.
				if c, ok := ctx.Referables.ByFirstTag("cursor").(*Thinger); ok {
					c.Animation("grab")
				}
				if t.controller != nil {
					t.controller.Block()
				}
				if pc != nil {
					if pc.input.ActionIsPressed(InputMoveTo) {
						pc.heldItem = &inv.items[index] // uh-oh!!!
					}
				}
			}
			if pc != nil && pc.heldItem != nil {
				if pc.heldItem == &inv.items[index] {
					inv.heldItemIndex = index
				}
			}
		}
	}

	// Fade in the orb while hovered.
	if x >= inv.X() && x <= inv.X()+inv.width && y >= inv.Y() && y <= inv.Y()+inv.height {
		if inv.fade < fadeMax {
			inv.fade++
		}
	} else {
		if inv.fade > fadeMin {
			inv.fade--
		}
	}

	return nil
}

func (inv *Inventory) Draw(ctx *context.Draw) {
	{
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(inv.X(), inv.Y())
		op.GeoM.Concat(ctx.Op.GeoM)
		op.ColorScale.ScaleAlpha(float32(inv.fade) / fadeMax)
		ctx.Target.DrawImage(res.Images["orb"], op)
	}

	for index, item := range inv.items {
		op := &ebiten.DrawImageOptions{}
		ix, iy, _, _ := inv.ItemBounds(index)
		op.GeoM.Translate(ix, iy)
		op.GeoM.Concat(ctx.Op.GeoM)
		inv.DrawItem(ctx, op, item, true)
	}

	if inv.hoveredName != "" {
		geom := ebiten.GeoM{}
		geom.Translate(ctx.MousePosition())
		geom.Translate(0, -16)
		geom.Concat(ctx.Op.GeoM)
		alpha := float32(inv.fade) / fadeMax
		ctx.Text(inv.hoveredName, geom, color.NRGBA{139, 98, 16, uint8(alpha * 255)})
	}

	if inv.heldItemIndex != -1 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(ctx.MousePosition())
		op.GeoM.Translate(0, -8)
		op.GeoM.Translate(-float64(inv.items[inv.heldItemIndex].staxer.stax.Stax.SliceWidth/2), -float64(inv.items[inv.heldItemIndex].staxer.stax.Stax.SliceHeight/2))
		op.GeoM.Concat(ctx.Op.GeoM)
		inv.DrawItem(ctx, op, inv.items[inv.heldItemIndex], false)
	}
}

func (inv *Inventory) DrawItem(ctx *context.Draw, op *ebiten.DrawImageOptions, item InvItem, fade bool) {
	opts := &ebiten.DrawImageOptions{}

	const sliceDistance = 1
	sliceDistanceEnd := math.Max(1, sliceDistance*op.GeoM.Element(0, 0))

	for i, slice := range item.staxer.frame.Slices {
		for j := 0; j < int(sliceDistanceEnd); j++ {
			opts.GeoM.Reset()

			opts.GeoM.Translate(0, -sliceDistance*float64(i))

			opts.GeoM.Concat(op.GeoM)

			opts.GeoM.Translate(0, float64(j))
			if fade {
				opts.ColorScale.ScaleAlpha(float32(inv.fade) / fadeMax)
			}

			opts.Blend = ctx.Op.Blend
			sub := item.staxer.stax.EbiImage.SubImage(image.Rect(slice.X, slice.Y, slice.X+item.staxer.stax.Stax.SliceWidth, slice.Y+item.staxer.stax.Stax.SliceHeight)).(*ebiten.Image)
			ctx.Target.DrawImage(sub, opts)
		}
	}
}

/*func (inv *Inventory) Resize(width, height int) {
	if inv.tw != float64(width) || inv.th != float64(height) {
		inv.tw = float64(width)
		inv.th = float64(height)
		inv.width = 9.0*2 + 16
		inv.height = 9.0*4 + 16
		inv.SetX(inv.tw - inv.width - 16)
		inv.SetY(inv.th - inv.height - 16)
	}
}*/
