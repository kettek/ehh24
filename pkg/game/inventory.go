package game

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kettek/ehh24/pkg/game/ables"
	"github.com/kettek/ehh24/pkg/game/context"
)

type Inventory struct {
	ables.IDable
	ables.Priorityable
	ables.Tagable
	ables.Positionable
	tw        float64
	th        float64
	width     float64
	height    float64
	targetTag string // eh, whatever
	items     []InvItem
}

type InvItem struct {
	item     ables.StorageItem
	staticer *Staticer
}

func NewInventory(tag string) *Inventory {
	inv := &Inventory{
		targetTag: tag,
	}
	inv.SetPriority(ables.PriorityUI)
	inv.SetTag("inventory")

	inv.width = 11 * 2
	inv.height = 11 * 4

	return inv
}

func (inv *Inventory) ItemBounds(index int) (x, y, x2, y2 float64) {
	x = inv.X()
	y = inv.Y()

	x += float64(index%2) * 11.0
	y += math.Floor(float64(index)/2.0) * 11 // ??

	return x, y, x + 9, y + 9
}

func (inv *Inventory) Update(ctx *ContextGame) []Change {
	inv.SetX(ctx.Width/ctx.Zoom - inv.width - 8)
	inv.SetY(ctx.Height/ctx.Zoom - inv.height - 8)

	inv.items = nil
	x, y := ctx.MousePosition()
	if t, ok := ctx.Referables.ByFirstTag(inv.targetTag).(*Thinger); ok {
		for index, item := range t.Storagable {
			ix, iy, ix2, iy2 := inv.ItemBounds(index)
			ix /= ctx.Zoom
			iy /= ctx.Zoom
			ix2 /= ctx.Zoom
			iy2 /= ctx.Zoom
			if x >= ix && x <= ix2 && y >= iy && y <= iy2 {
				fmt.Println("hit item", item)
			}
			inv.items = append(inv.items, InvItem{
				item:     item,
				staticer: NewStaticer(item.Tag),
			})
		}
	}
	return nil
}

func (inv *Inventory) Draw(ctx *context.Draw) {
	zoom := ctx.Op.GeoM.Element(0, 0)
	x := inv.X()
	y := inv.Y()
	vector.DrawFilledRect(ctx.Target, float32(x*zoom), float32(y*zoom), float32(inv.width*zoom), float32(inv.height*zoom), color.NRGBA{255, 0, 0, 255}, true)
	for index, _ := range inv.items {
		x1, y1, x2, y2 := inv.ItemBounds(index)
		w := x2 - x1
		h := y2 - y1
		//x1 += w / 2
		//y1 += h / 2
		vector.DrawFilledRect(ctx.Target, float32(x1*zoom), float32(y1*zoom), float32(w*zoom), float32(h*zoom), color.White, true)
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
