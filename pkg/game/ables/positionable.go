package ables

// Positionable b x and y
type Positionable struct {
	x, y float64
}

// X returns the x of the position.
func (p Positionable) X() float64 {
	return p.x
}

// Y returns the y of the position.
func (p Positionable) Y() float64 {
	return p.y
}

// SetX sets the x of the position.
func (p *Positionable) SetX(x float64) {
	p.x = x
}

// SetY sets the y of the position.
func (p *Positionable) SetY(y float64) {
	p.y = y
}

// MakePositionable makes a new Positionable.
func MakePositionable(x, y float64) Positionable {
	return Positionable{x: x, y: y}
}
