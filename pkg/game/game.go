package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
)

type Game struct {
	insys           input.System
	updateDrawers   []UpdateDrawer
	geom            ebiten.GeoM
	midlay          *ebiten.Image
	cursor          UpdateDrawer
	darknessOverlay *DarknessOverlay

	gctx GameContext
	dctx DrawContext
}

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
	g.cursor = c

	t := NewThinger("test")
	t.controller = NewPlayerController(&g.insys)
	t.originX = -0.5
	t.originY = -1

	// Make some test stuf.
	t1 := NewStaticer("term-small")
	t1.originX = -0.5
	t1.originY = -1

	t2 := NewStaticer("term-large")
	t2.originX = -0.5
	t2.originY = -1

	geom := ebiten.GeoM{}
	geom.Scale(3, 3)

	g.geom = geom
	g.updateDrawers = []UpdateDrawer{t, t1, t2}

	g.gctx.Zoom = g.geom.Element(0, 0)

	g.darknessOverlay = NewDarknessOverlay(320, 240)
	g.midlay = ebiten.NewImage(320, 240)

	return g
}

func (g *Game) Update() error {
	g.insys.Update()

	var changes []Change
	for _, t := range g.updateDrawers {
		changes = append(changes, t.Update(&g.gctx)...)
	}

	for _, c := range changes {
		c.Apply(g)
	}

	g.cursor.Update(&g.gctx)
	g.darknessOverlay.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM = g.geom

	g.dctx.Target = g.midlay
	g.dctx.Op = op

	g.midlay.Clear()
	g.midlay.Fill(color.NRGBA{20, 20, 20, 255})
	g.darknessOverlay.Draw(&g.dctx)

	for _, t := range g.updateDrawers {
		t.Draw(&g.dctx)
	}
	op = &ebiten.DrawImageOptions{}
	screen.DrawImage(g.midlay, op)

	op.Blend = ebiten.BlendDestinationAtop

	screen.DrawImage(g.darknessOverlay.Image, op)

	g.dctx.Target = screen
	g.cursor.Draw(&g.dctx)
}

func (g *Game) Layout(ow, oh int) (int, int) {
	if g.dctx.Width != float64(ow) || g.dctx.Height != float64(oh) {
		g.dctx.Width = float64(ow)
		g.dctx.Height = float64(oh)
		g.gctx.Width = float64(ow)
		g.gctx.Height = float64(oh)
		g.darknessOverlay.Resize(ow, oh)
		g.midlay = ebiten.NewImage(ow, oh)
	}
	return ow, oh
}
