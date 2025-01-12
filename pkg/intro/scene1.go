package intro

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ehh24/pkg/game/context"
	"github.com/kettek/ehh24/pkg/res"
)

type Scene1 struct {
	dropVelocity float64
	dropY        float64
	targetY      float64
	textTimer    int
	textFadeIn   int
	textFadeOut  int
}

func (s *Scene1) Init() {
	s.dropY = 0
	s.dropVelocity = 0
	s.targetY = 600
}

func (s *Scene1) Update() Scene {
	if s.dropY < s.targetY {
		s.dropVelocity += 0.05
		s.dropY += s.dropVelocity
	} else {
		sc := &Scene2{}
		sc.Init()
		return sc
	}
	return s
}

func (s *Scene1) Draw(ctx context.Draw) {
	op := &ebiten.DrawImageOptions{}
	if s.dropY < s.targetY {
		op.ColorScale.ScaleAlpha(float32(s.dropY+50) / float32(s.targetY))
		drop := res.Images["drop"]
		op.GeoM.Translate(220, s.dropY)
		op.GeoM.Translate(-float64(drop.Bounds().Dx())/2, -float64(drop.Bounds().Dy())/2)
		op.GeoM.Concat(ctx.Op.GeoM)
		ctx.Target.DrawImage(drop, op)
	}
}
