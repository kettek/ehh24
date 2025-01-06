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
