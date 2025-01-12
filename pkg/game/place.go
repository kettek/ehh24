package game

import (
	"github.com/kettek/ehh24/pkg/game/ables"
	"github.com/kettek/ehh24/pkg/res"
	"github.com/traefik/yaegi/interp"
)

// Place is where things do be happen, tho.
type Place struct {
	Name       string
	referables Referables // Da referables in da place.
	areas      []*Area    // Da collision areas.
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
		st.SetTag(static.Tag)
		p.referables = append(p.referables, st)
	}

	// Load in things.
	// TODO!

	// Load in collision areas.
	for _, poly := range rp.Polygons {
		area := &Area{
			original: poly,
		}

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
