package game

import (
	"math"
	"math/rand"

	"github.com/kettek/ehh24/pkg/game/ables"
)

// BoidController is boid logic attached to a thinger.
type BoidController struct {
	flockID      int
	dx           float64
	dy           float64
	visualRange  float64
	speedLimit   float64
	targetID     int
	settled      bool
	shouldSettle bool
	settles      bool
	meander      bool
	block        bool
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
	if b.block {
		return
	}
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

	var target *Thinger
	if b.targetID == 0 {
		b.flyTowardsCenter(t, boids)
	} else if tt := ctx.Referables.ByID(b.targetID); tt != nil {
		var ok bool
		target, ok = tt.(*Thinger)
		if ok {
			b.flyTowardsPosition(t, target.X(), target.Y())
		}
	}
	b.avoidOthers(t, boids)
	b.matchVelocity(t, boids)
	b.limitSpeed()
	b.keepInBounds(ctx, t)
	b.doSettle(t, boids)
	b.doMeander(t, target)

	a = append(a, &ActionPosition{
		X: t.X() + b.dx,
		Y: t.Y() + b.dy,
	})

	a = append(a, &ActionFace{
		Radians: math.Atan2(b.dy, b.dx) + math.Pi/2,
	})

	return a
}

func (b *BoidController) flyTowardsPosition(self *Thinger, x, y float64) {
	const centerFactor = 0.002

	b.dx += (x - self.X()) * centerFactor
	b.dy += (y - self.Y()) * centerFactor

	if b.settles {
		dist := b.distance(self, &Thinger{
			Positionable: ables.MakePositionable(x, y),
		})
		if dist < 60 {
			b.shouldSettle = true
		}
	}
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
	const minDistance = 8.0

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

func (b *BoidController) doSettle(self *Thinger, boids []*Thinger) {
	if !b.shouldSettle || b.settled {
		return
	}

	b.dx *= 0.9
	b.dy *= 0.9
	if math.Abs(b.dx) < 0.5 && math.Abs(b.dy) < 0.5 {
		b.setSettle(self, true)
	} else {
		// Check against our other boids settle count.
		settleCount := 0
		for _, boid := range boids {
			bc := boid.controller.(*BoidController)
			if bc.settled {
				settleCount++
			}
		}
		if settleCount >= len(boids)-len(boids)/4 {
			b.setSettle(self, true)
		}
	}
}

func (b *BoidController) setSettle(self *Thinger, settle bool) {
	if settle {
		b.settled = true
		b.meander = true
		self.Animation("settled")
		b.speedLimit = 0.5
	} else {
		b.settled = false
		b.meander = false
		self.Animation("fly")
		b.speedLimit = 3
	}
}

func (b *BoidController) doMeander(self *Thinger, target *Thinger) {
	if !b.meander {
		return
	}

	if target != nil {
		if b.distance(self, target) > 60 {
			b.setSettle(self, false)
		}
	}

	b.dx += rand.Float64()*0.2 - 0.1
	b.dy += rand.Float64()*0.2 - 0.1
}

func (b *BoidController) distance(self, other *Thinger) float64 {
	return math.Sqrt(
		(self.X()-other.X())*(self.X()-other.X()) +
			(self.Y()-other.Y())*(self.Y()-other.Y()),
	)
}

func (b *BoidController) Block() {
	b.block = true
}

func (b *BoidController) Unblock() {
	b.block = false
}
