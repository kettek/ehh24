package game

import (
	"fmt"

	"github.com/kettek/ehh24/pkg/game/ables"
	"github.com/kettek/ehh24/pkg/res"
)

// Place is where things do be happen, tho.
type Place struct {
	referables Referables // Da referables in da place.
}

// NewPlace does a thingie.
func NewPlace(name string) *Place {
	p := &Place{}

	// Load from res.
	rp, ok := res.Places["places/"+name]
	if !ok {
		panic("place not found: " + name)
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
	for _, poly := range rp.Polygons {
		// TODO!
		fmt.Println(poly)
	}

	return p
}
