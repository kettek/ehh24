package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kettek/ehh24/pkg/game/ables"
	"github.com/kettek/ehh24/pkg/game/context"
	"github.com/kettek/ehh24/pkg/outro"
	"github.com/kettek/ehh24/pkg/statemachine"
	input "github.com/quasilyte/ebitengine-input"
)

// State is our absolutely amazing game with so many features and fun.
type State struct {
	insys  input.System
	geom   ebiten.GeoM
	midlay *ebiten.Image

	debugUI *TargetOverlay

	gctx ContextGame
	dctx context.Draw
}

// NewState does exactly what you should think.
func NewState() *State {
	g := &State{}
	g.insys.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})
	g.gctx.Places = make(map[string]*Place)
	// Setup input system
	// Make our lil cursor?
	c := NewThinger("cursor")
	c.controller = NewCursorController()
	c.originX = -0.5
	c.originY = -0.5
	c.SetPriority(ables.PriorityBeyond)
	c.SetTag("cursor")

	t := NewThinger("test")
	t.controller = NewPlayerController(&g.insys)
	t.originX = -0.5
	t.originY = -1
	t.SetPriority(ables.PriorityMiddle)
	t.SetTag("qi")

	geom := ebiten.GeoM{}
	geom.Scale(3, 3)

	g.geom = geom

	g.gctx.Zoom = g.geom.Element(0, 0)

	vis := NewVisibilityOverlay(320, 240)
	vis.SetPriority(ables.PriorityOverlay + 1000)
	vis.SetTag("visibility")

	sno := NewSnoverlay(320, 240)
	sno.SetPriority(ables.PriorityOverlay)
	sno.SetTag("snow")

	fadein := NewFadeInOverlay(320, 240, 100)
	fadein.SetPriority(ables.PriorityOverlay + 100)

	g.debugUI = NewTargetOverlay(320, 240)

	g.gctx.Referables = Referables{t /*vis, sno,*/, fadein, c}

	// Some boids of testing.
	/*roboid := NewThinger("boid")
	roboid.SetX(200)
	roboid.SetY(200)
	rc := NewBoidController(1)
	rc.settles = true
	rc.targetID = t.ID()
	roboid.controller = rc
	roboid.centerX = 0.5
	roboid.centerY = 0.5
	roboid.SetPriority(ables.PriorityMiddle)
	roboid.SetTag("boid")
	roboid.Stack("roboid")
	g.gctx.Referables = append(g.gctx.Referables, roboid)
	for i := 0; i < 20; i++ {
		b := NewThinger("boid")
		b.SetX(400)
		b.SetY(200)
		bc := NewBoidController(1)
		bc.settles = true
		bc.targetID = roboid.ID()
		b.controller = bc
		b.centerX = 0.5
		b.centerY = 0.5
		b.SetPriority(ables.PriorityMiddle)
		b.SetTag("boid")
		if i%2 == 0 {
			b.Stack("boid2")
		}
		g.gctx.Referables = append(g.gctx.Referables, b)
	}*/

	// Some more boids of testing.
	/*for i := 0; i < 20; i++ {
		b := NewThinger("boid")
		b.Stack("boid2")
		b.controller = NewBoidController(2)
		b.centerX = 0.5
		b.centerY = 0.5
		b.SetPriority(ables.PriorityMiddle)
		b.SetTag("boid")
		g.gctx.Referables = append(g.gctx.Referables, b)
	}*/

	g.midlay = ebiten.NewImage(320, 240)

	inventory := NewInventory("qi")
	inventory.SetPriority(ables.PriorityUI)
	inventory.SetTag("inventory")
	g.gctx.Referables = append(g.gctx.Referables, inventory)

	// for now, just try to load in test place.
	g.gctx.Place = NewPlace("cells")
	g.gctx.Places["cells"] = g.gctx.Place
	if start := g.gctx.Place.GetAreaByFirstTag("spawn"); start != nil {
		cx, cy := start.Center()
		t.SetX(cx)
		t.SetY(cy)
	}

	return g
}

// Init initializes the game.
func (g *State) Init() {
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
}

// Update updates the game.
func (g *State) Update() statemachine.State {
	g.insys.Update()

	startProfile("update")
	updateables := g.gctx.Referables.Updateables()
	var changes []Change
	for _, t := range updateables {
		changes = append(changes, t.Update(&g.gctx)...)
	}
	// Also do place.
	changes = append(changes, g.gctx.Place.Update(&g.gctx)...)
	endProfile("update")

	startProfile("changes")
	for _, c := range changes {
		c.Apply(&g.gctx)
		// I'm sorry for this...
		if c, ok := c.(*ChangeState); ok {
			if c.State == "end" {
				return outro.NewState()
			}
		}
	}
	endProfile("changes")

	startProfile("sort drawables")
	// Probably shouldn't do this, but...
	for _, t := range g.gctx.Referables.Drawables() {
		t.SetOffset(int(t.Y()))
	}
	for _, t := range g.gctx.Place.referables.Drawables() {
		t.SetOffset(int(t.Y()))
	}
	endProfile("sort drawables")

	return nil
}

// Draw draws the game.
func (g *State) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM = g.geom

	g.dctx.Target = g.midlay
	g.dctx.Op = op

	g.midlay.Clear()
	g.midlay.Fill(color.NRGBA{10, 10, 10, 255})

	g.debugUI.Draw(&g.dctx)

	// A bit terrible to merge like this, but oh wel..
	startProfile("draw drawables")
	referables := append(g.gctx.Place.referables, g.gctx.Referables...)
	for _, t := range referables.SortedDrawables() {
		t.Draw(&g.dctx)
	}
	endProfile("draw drawables")

	op = &ebiten.DrawImageOptions{}
	screen.DrawImage(g.midlay, op)

	op.Blend = ebiten.BlendDestinationAtop

	startProfile("draw overlays")
	for _, t := range referables.Overlays() {
		t.DrawTo(screen)
	}
	endProfile("draw overlays")

	// Print our debuggies
	if debug {
		for _, t := range referables.Debugables() {
			ebitenutil.DebugPrintAt(g.debugUI.img, t.String(), int(t.X()*g.gctx.Zoom), int(t.Y()*g.gctx.Zoom))
		}
		for i, p := range profiles {
			ebitenutil.DebugPrintAt(g.debugUI.img, fmt.Sprintf("%03d %s", p.duration.Milliseconds(), p.name), 0, 30+i*10)
		}
		g.debugUI.DrawTo(screen)
	}
}

// Layout is a thing, yo.
func (g *State) Layout(ow, oh int) (int, int) {
	ow = 1280
	oh = 720
	if g.dctx.Width != float64(ow) || g.dctx.Height != float64(oh) {
		g.dctx.Width = float64(ow)
		g.dctx.Height = float64(oh)
		g.gctx.Width = float64(ow)
		g.gctx.Height = float64(oh)
		g.midlay = ebiten.NewImage(ow, oh)
		for _, t := range g.gctx.Referables.Resizables() {
			t.Resize(ow, oh)
		}
		for _, t := range g.gctx.Place.referables.Resizables() {
			t.Resize(ow, oh)
		}
		g.debugUI.Resize(ow, oh)
	}
	return ow, oh
}
