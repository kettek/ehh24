package game

import (
	"github.com/kettek/ehh24/pkg/res"
)

// Area represents an invisible polygonal shape for blocking movement or causing triggers.
type Area struct {
	// Just store a ref to original polygon, I guess?
	original *res.Polygon
}

func (a *Area) ContainsPoint(x, y float64) bool {
	return a.original.ContainsPoint(x, y)
}

func (a *Area) Center() (float64, float64) {
	cx := 0.0
	cy := 0.0
	for _, point := range a.original.Points {
		cx += float64(point.X)
		cy += float64(point.Y)
	}
	cx /= float64(len(a.original.Points))
	cy /= float64(len(a.original.Points))
	return cx, cy
}

func (a *Area) Bounds() (float64, float64, float64, float64) {
	minX := 999999.0
	minY := 999999.0
	maxX := -999999.0
	maxY := -999999.0
	for _, point := range a.original.Points {
		x := float64(point.X)
		y := float64(point.Y)
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}
	return minX, minY, maxX, maxY
}
