package game

import (
	"github.com/kettek/ehh24/pkg/res"
	"github.com/solarlune/resolv"
)

// Area represents an invisible polygonal shape for blocking movement or causing triggers.
type Area struct {
	// ??
	shape *resolv.ConvexPolygon
	// Just store a ref to original polygon, I guess?
	original *res.Polygon
}

func (a *Area) ContainsPoint(x, y float64) bool {
	isInside := false
	// Calculate using original.
	for i, j := 0, len(a.original.Points)-1; i < len(a.original.Points); j, i = i, i+1 {
		px := float64(a.original.Points[i].X)
		py := float64(a.original.Points[i].Y)
		qx := float64(a.original.Points[j].X)
		qy := float64(a.original.Points[j].Y)
		if ((py > y) != (qy > y)) && (x < (qx-px)*(y-py)/(qy-py)+px) {
			isInside = !isInside
		}
	}
	return isInside
}
