package game

import (
	"github.com/kettek/ehh24/pkg/res"
	"github.com/kettek/ehh24/pkg/stax"
)

// Our bits and bobs of most thingers.
const (
	BackLeg = iota
	BackArm
	Heart
	FrontLeg
	FrontArm
	Head
	Eyes
)

// Staxer is a contained state structure for rendering staxie files.
type Staxer struct {
	stax       res.StaxImage // hmm
	stack      *stax.Stack
	lastAnim   string
	animation  *stax.Animation
	frame      *stax.Frame
	frameIndex int
	frameTimer int
}

// NewStaxer does exactly what u thinkie.
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

// Stack sets the staxer's stack to the given string. Panics if no string exists.
func (s *Staxer) Stack(name string) {
	stack := s.stax.Stax.Stack(name)
	if stack == nil {
		panic("stack not found")
	}
	s.stack = stack
	s.lastAnim = "ieee"
	s.Animation(stack.Animations[0].Name)
}

// Animation sets the staxer's animation to the given string. Panics if no string exists.
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

// Update updates the animation timer.
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
