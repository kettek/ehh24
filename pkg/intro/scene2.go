package intro

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ehh24/pkg/game/context"
)

type Scene2 struct {
	fadeIn   int
	fadeWait int
	fadeOut  int
}

func (s *Scene2) Init() {
	s.fadeWait = 100
	s.fadeIn = 100
	s.fadeOut = 100
}

func (s *Scene2) Update() Scene {
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

func (s *Scene2) Draw(ctx context.Draw) {
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

	ctx.Text("ミル・・・ネ", op.GeoM, clr)
}
