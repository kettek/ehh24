package game

import (
	"github.com/kettek/ehh24/pkg/res"
	"github.com/kettek/ehh24/pkg/stax"
)

type Staxer struct {
	stax       res.StaxImage // hmm
	stack      *stax.Stack
	animation  *stax.Animation
	frame      *stax.Frame
	frameIndex int
	frameTimer int
}

func NewStaxer(name string) Staxer {
	st, err := res.GetStax(name)
	if err != nil {
		panic(err)
	}
	stack := st.Stax.Stack("base")
	if stack == nil {
		panic("stack not found")
	}
	animation := stack.Animation("idle")
	if animation == nil {
		panic("animation not found")
	}
	frame := animation.Frame(0)
	if frame == nil {
		panic("frame not found")
	}
	return Staxer{
		stax:      st,
		stack:     stack,
		animation: animation,
		frame:     frame,
	}
}

func (s *Staxer) Update() {
	s.frameTimer++
	if s.frameTimer > s.animation.FrameTime {
		s.frameTimer = 0
		s.frameIndex++
		if s.frameIndex >= len(s.animation.Frames) {
			s.frameIndex = 0
		}
		s.frame = s.animation.Frame(s.frameIndex)
	}
}
