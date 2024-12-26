package game

import (
	"github.com/kettek/ehh24/pkg/res"
	"github.com/kettek/ehh24/pkg/stax"
)

type StackLayer int

const (
	Base  StackLayer = 0
	Heart            = 1
	Head             = 2
	Eyes             = 3
)

type Staxer struct {
	stax       res.StaxImage // hmm
	stack      *stax.Stack
	lastAnim   string
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
	if len(st.Stax.Stacks) == 0 {
		panic("no stacks found")
	}
	stack := &st.Stax.Stacks[0]
	if len(stack.Animations) == 0 {
		panic("no animations found")
	}
	animation := &stack.Animations[0]
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

func (s *Staxer) Animation(name string) {
	if s.lastAnim == name {
		return
	}
	s.lastAnim = name
	animation := s.stack.Animation(name)
	if animation == nil {
		panic("animation not found")
	}
	s.animation = animation
	s.frame = animation.Frame(0)
	s.frameIndex = 0
	s.frameTimer = 0
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
