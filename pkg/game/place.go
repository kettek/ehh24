package game

import (
	"github.com/kettek/ehh24/pkg/game/ables"
	"github.com/kettek/ehh24/pkg/res"
	"github.com/solarlune/resolv"
	"github.com/traefik/yaegi/interp"
)

// Place is where things do be happen, tho.
type Place struct {
	Name       string
	referables Referables // Da referables in da place.
	areas      []*Area    // Da collision areas.
	space      *resolv.Space
	entered    bool
	// interpreter stuff
	interp  *interp.Interpreter
	OnEnter func(p *Place)
	OnLeave func(p *Place)
	OnTick  func(p *Place)
}

// NewPlace does a thingie.
func NewPlace(name string) *Place {
	p := &Place{}

	// Load from res.
	rp, ok := res.Places["places/"+name]
	if !ok {
		panic("place not found: " + name)
	}
	p.Name = rp.Name

	// Setup interpreter stuff
	if script, ok := res.Scripts["places/"+name]; ok {
		p.interp = interp.New(interp.Options{})
		setupInterp(p.interp, script)
		if fn, err := p.interp.Eval("Enter"); err == nil {
			p.OnEnter = fn.Interface().(func(p *Place))
		}
		if fn, err := p.interp.Eval("Leave"); err == nil {
			p.OnLeave = fn.Interface().(func(p *Place))
		}
		if fn, err := p.interp.Eval("Tick"); err == nil {
			p.OnTick = fn.Interface().(func(p *Place))
		}
	}

	// Load in the floors.
	for _, floor := range rp.Floor {
		fl := NewFloor(floor.Name)
		fl.SetX(float64(floor.Point.X))
		fl.SetY(float64(floor.Point.Y))
		fl.SetPriority(ables.PriorityBack)
		p.referables = append(p.referables, fl)
	}

	// Load in the staticers.
	for _, static := range rp.Statics {
		st := NewStaticer(static.Name)
		st.SetX(float64(static.Point.X))
		st.SetY(float64(static.Point.Y))
		st.SetPriority(ables.PriorityMiddle)
		p.referables = append(p.referables, st)
	}

	// Load in things.
	// TODO!

	// Load in collision areas.
	p.space = resolv.NewSpace(1280, 1280, 8, 8)
	for _, poly := range rp.Polygons {
		area := &Area{
			original: poly,
		}
		// I guess we have to transform this crap for resolv...
		var points []float64
		cx, cy := 0.0, 0.0
		for _, point := range poly.Points {
			points = append(points, float64(point.X), float64(point.Y))
			cx += float64(point.X)
			cy += float64(point.Y)
		}
		cx /= float64(len(poly.Points))
		cy /= float64(len(poly.Points))
		for i := 0; i < len(points); i += 2 {
			points[i] -= cx
			points[i+1] -= cy
		}
		area.shape = resolv.NewConvexPolygon(cx, cy, points)
		area.shape.SetData(area)

		p.space.Add(area.shape)

		p.areas = append(p.areas, area)
	}

	return p
}

func (p *Place) Update(ctx *ContextGame) []Change {
	changes := []Change{}

	if p.OnTick != nil {
		p.OnTick(p)
	}

	for _, t := range p.referables.Updateables() {
		changes = append(changes, t.Update(ctx)...)
	}

	return changes
}

func (p *Place) GetAreaByFirstTag(tag string) *Area {
	for _, area := range p.areas {
		if area.original.Tag == tag {
			return area
		}
	}
	return nil
}

func (p *Place) RemoveAreaByFirstTag(tag string) {
	for i, area := range p.areas {
		if area.original.Tag == tag {
			p.areas = append(p.areas[:i], p.areas[i+1:]...)
			return
		}
	}
	return
}
