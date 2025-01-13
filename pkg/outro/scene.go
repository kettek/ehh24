package outro

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ehh24/pkg/game/context"
	"github.com/kettek/ehh24/pkg/res"
)

type Scene1 struct {
	fadeIn   int
	fadeWait int
	fadeOut  int
}

func (s *Scene1) Init() {
	s.fadeWait = 1000
	s.fadeIn = 100
	s.fadeOut = 100
}

func (s *Scene1) Update() Scene {
	if s.fadeIn > 0 {
		s.fadeIn--
	} else if s.fadeWait > 0 {
		s.fadeWait--
	} else if s.fadeOut > 0 {
		s.fadeOut--
	} else {
		return nil
	}
	return s
}

func (s *Scene1) Draw(ctx context.Draw) {
	op := &ebiten.DrawImageOptions{}

	alpha := 1.0

	if s.fadeIn > 0 {
		alpha = 1.0 - float64(s.fadeIn)/100
	} else if s.fadeOut > 0 {
		alpha = float64(s.fadeOut) / 100
	} else {
		alpha = 0
	}

	a := uint8(alpha * 255)

	clr := color.RGBA{a, a, a, a}

	op.GeoM.Translate(220, 100)
	op.GeoM.Scale(3, 3)

	ctx.Text("アリガトウゴザイマス！", op.GeoM, clr)

	op.GeoM.Translate(40, 20)
	ctx.Target.DrawImage(res.Images["thanks"], op)
}
