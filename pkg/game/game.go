package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ehh24/pkg/game/ables"
	"github.com/kettek/ehh24/pkg/game/context"
	input "github.com/quasilyte/ebitengine-input"
)

// Game is our absolutely amazing game with so many features and fun.
type Game struct {
	insys  input.System
	geom   ebiten.GeoM
	midlay *ebiten.Image

	referables Referables

	gctx context.Game
	dctx context.Draw
}

// NewGame does exactly what you should think.
func NewGame() *Game {
	g := &Game{}
	g.insys.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})
	// Setup input system
	// Make our lil cursor?
	c := NewThinger("cursor")
	c.controller = NewCursorController()
	c.originX = -0.5
	c.originY = -0.5
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	c.SetPriority(ables.PriorityBeyond)
	c.SetTag("cursor")

	t := NewThinger("test")
	t.controller = NewPlayerController(&g.insys)
	t.originX = -0.5
	t.originY = -1
	t.SetPriority(ables.PriorityMiddle)

	// Make some test stuf.
	t1 := NewStaticer("term-small")
	t1.originX = -0.5
	t1.originY = -1
	t1.SetPriority(ables.PriorityMiddle)

	t2 := NewStaticer("term-large")
	t2.originX = -0.5
	t2.originY = -1
	t2.SetPriority(ables.PriorityMiddle)

	geom := ebiten.GeoM{}
	geom.Scale(3, 3)

	g.geom = geom

	g.gctx.Zoom = g.geom.Element(0, 0)

	vis := NewVisibilityOverlay(320, 240)
	vis.SetPriority(ables.PriorityOverlay)
	vis.SetTag("visibility")

	g.referables = Referables{t, t1, t2, vis, c}

	g.midlay = ebiten.NewImage(320, 240)

	return g
}

// Update updates the game.
func (g *Game) Update() error {
	g.insys.Update()

	var changes []Change
	for _, t := range g.referables.Updateables() {
		changes = append(changes, t.Update(&g.gctx)...)
	}

	for _, c := range changes {
		c.Apply(g)
	}

	return nil
}

// Draw draws the game.
func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM = g.geom

	g.dctx.Target = g.midlay
	g.dctx.Op = op

	g.midlay.Clear()
	g.midlay.Fill(color.NRGBA{20, 20, 20, 255})

	for _, t := range g.referables.SortedDrawables() {
		t.Draw(&g.dctx)
	}
	op = &ebiten.DrawImageOptions{}
	screen.DrawImage(g.midlay, op)

	op.Blend = ebiten.BlendDestinationAtop

	for _, t := range g.referables.Overlays() {
		t.DrawTo(screen)
	}

	g.dctx.Target = screen
}

// Layout is a thing, yo.
func (g *Game) Layout(ow, oh int) (int, int) {
	if g.dctx.Width != float64(ow) || g.dctx.Height != float64(oh) {
		g.dctx.Width = float64(ow)
		g.dctx.Height = float64(oh)
		g.gctx.Width = float64(ow)
		g.gctx.Height = float64(oh)
		g.midlay = ebiten.NewImage(ow, oh)
		for _, t := range g.referables.Overlays() {
			t.Resize(ow, oh)
		}
	}
	return ow, oh
}
