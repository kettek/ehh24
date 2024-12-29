package game

import (
	"math"
	"math/rand"
)

// BoidController is boid logic attached to a thinger.
type BoidController struct {
	flockID     int
	dx          float64
	dy          float64
	visualRange float64
	speedLimit  float64
}

// NewBoidController makes a new boid controller.
func NewBoidController(flockID int) *BoidController {
	return &BoidController{
		flockID:     flockID,
		visualRange: 100,
		speedLimit:  3,
		dx:          rand.Float64()*10 - 5,
		dy:          rand.Float64()*10 - 5,
	}
}

// Update updates.
func (b *BoidController) Update(ctx *ContextGame, t *Thinger) (a []Action) {
	boidRefs := ctx.Referables.ByTag("boid")

	boids := make([]*Thinger, 0, len(boidRefs))
	for _, boidRef := range boidRefs {
		boid, ok := boidRef.(*Thinger)
		if !ok {
			continue
		}
		bc, ok := boid.controller.(*BoidController)
		if !ok {
			continue
		}
		if bc.flockID != b.flockID {
			continue
		}
		if boid.ID() == t.ID() {
			continue
		}
		boids = append(boids, boid)
	}

	b.flyTowardsCenter(t, boids)
	b.avoidOthers(t, boids)
	b.matchVelocity(t, boids)
	b.limitSpeed()
	b.keepInBounds(ctx, t)

	a = append(a, &ActionPosition{
		X: t.X() + b.dx,
		Y: t.Y() + b.dy,
	})

	a = append(a, &ActionFace{
		Radians: math.Atan2(b.dy, b.dx) + math.Pi/2,
	})

	return a
}

func (b *BoidController) flyTowardsCenter(self *Thinger, boids []*Thinger) {
	const centerFactor = 0.002

	cx := 0.0
	cy := 0.0
	neighborCount := 0.0

	for _, boid := range boids {
		if b.distance(self, boid) < b.visualRange {
			cx += boid.X()
			cy += boid.Y()
			neighborCount++
		}
	}

	if neighborCount > 0.0 {
		cx = cx / neighborCount
		cy = cy / neighborCount

		b.dx += (cx - self.X()) * centerFactor
		b.dy += (cy - self.Y()) * centerFactor
	}
}

func (b *BoidController) avoidOthers(self *Thinger, boids []*Thinger) {
	const avoidFactor = 0.02
	const minDistance = 15.0

	moveX := 0.0
	moveY := 0.0

	for _, boid := range boids {
		if b.distance(self, boid) < minDistance {
			moveX += self.X() - boid.X()
			moveY += self.Y() - boid.Y()
		}
	}
	b.dx += moveX * avoidFactor
	b.dy += moveY * avoidFactor
}

func (b *BoidController) matchVelocity(self *Thinger, boids []*Thinger) {
	const matchFactor = 0.05

	avgDX := 0.0
	avgDY := 0.0
	neighborCount := 0.0

	for _, boid := range boids {
		if b.distance(self, boid) < b.visualRange {
			otherController := boid.controller.(*BoidController)
			avgDX += otherController.dx
			avgDY += otherController.dy
			neighborCount++
		}
	}

	if neighborCount > 0.0 {
		avgDX = avgDX / neighborCount
		avgDY = avgDY / neighborCount

		b.dx += (avgDX - b.dx) * matchFactor
		b.dy += (avgDY - b.dy) * matchFactor
	}
}

func (b *BoidController) limitSpeed() {
	speed := math.Sqrt(b.dx*b.dx + b.dy*b.dy)
	if speed > b.speedLimit {
		b.dx = (b.dx / speed) * b.speedLimit
		b.dy = (b.dy / speed) * b.speedLimit
	}
}

func (b *BoidController) keepInBounds(ctx *ContextGame, self *Thinger) {
	const margin = -10.0
	const turnFactor = 0.5

	w, h := ctx.Size()

	if self.X() < margin {
		b.dx += turnFactor
	}
	if self.Y() < margin {
		b.dy += turnFactor
	}
	if self.X() > w-margin {
		b.dx -= turnFactor
	}
	if self.Y() > h-margin {
		b.dy -= turnFactor
	}
}

func (b *BoidController) distance(self, other *Thinger) float64 {
	return math.Sqrt(
		(self.X()-other.X())*(self.X()-other.X()) +
			(self.Y()-other.Y())*(self.Y()-other.Y()),
	)
}
