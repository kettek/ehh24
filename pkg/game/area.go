package game

import "github.com/solarlune/resolv"

// Area represents an invisible polygonal shape for blocking movement or causing triggers.
type Area struct {
	// ??
	shape *resolv.ConvexPolygon
}
