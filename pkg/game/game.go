package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	thingers []*Thinger
	geom     ebiten.GeoM
	width    float64
	height   float64
}

func NewGame() *Game {
	// Make our lil cursor?
	c := NewThinger("cursor")
	c.controller = NewCursorController()
	c.originX = -0.5
	c.originY = -0.5
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	t := NewThinger("test")
	t.controller = NewPlayerController()
	t.originX = -0.5
	t.originY = -1
	geom := ebiten.GeoM{}
	geom.Scale(3, 3)
	return &Game{
		thingers: []*Thinger{t, c},
		geom:     geom,
	}
}

func (g *Game) Update() error {
	for _, t := range g.thingers {
		t.Update(&DrawContext{
			Target: nil,
			GeoM:   g.geom,
			Width:  g.width,
			Height: g.height,
		})
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	dctx := &DrawContext{
		Target: screen,
		GeoM:   g.geom,
		Width:  g.width,
		Height: g.height,
	}
	for _, t := range g.thingers {
		t.Draw(dctx)
	}
}

func (g *Game) Layout(ow, oh int) (int, int) {
	g.width = float64(ow)
	g.height = float64(oh)
	return ow, oh
}
