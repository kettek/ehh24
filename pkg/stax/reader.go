package stax

import (
	"errors"
)

func ReadStaxFromPNG(data []byte) (*Stax, error) {
	if len(data) < 8 {
		return nil, ErrDataTooShort
	}
	if data[0] != 137 || data[1] != 80 || data[2] != 78 || data[3] != 71 || data[4] != 13 || data[5] != 10 || data[6] != 26 || data[7] != 10 {
		return nil, ErrInvalidPNG
	}
	offset := 8
	// Read out chunks.
	for {
		if len(data) < offset+12 {
			break
		}
		chunkLength := int(data[offset])<<24 | int(data[offset+1])<<16 | int(data[offset+2])<<8 | int(data[offset+3])
		chunkType := string(data[offset+4 : offset+8])
		if chunkType == "stAx" {
			st := &Stax{}
			if err := st.UnmarshalBinary(data[offset+8 : offset+8+int(chunkLength)]); err != nil {
				return nil, err
			}
			return st, nil
		}
		offset += 12 + int(chunkLength)
	}
	return nil, ErrNoStaxChunk
}

type Stax struct {
	SliceWidth  int
	SliceHeight int
	Stacks      []Stack
}

func (s *Stax) Stack(name string) *Stack {
	for i := range s.Stacks {
		if s.Stacks[i].Name == name {
			return &s.Stacks[i]
		}
	}
	return nil
}

type decodeContext struct {
	SliceWidth  int
	SliceHeight int
	SliceCount  int
	x           int
	y           int
}

func (s *Stax) UnmarshalBinary(data []byte) error {
	if len(data) < 7 {
		return ErrStaxDataTooShort
	}
	offset := 0
	version := data[offset]
	if version != 0 {
		return ErrInvalidVersion
	}
	offset++

	s.SliceWidth = int(data[offset])<<8 | int(data[offset+1])
	s.SliceHeight = int(data[offset+2])<<8 | int(data[offset+3])
	offset += 4

	count := int(data[offset])<<8 | int(data[offset+1])
	offset += 2

	for i := 0; i < count; i++ {
		var stack Stack
		n, err := stack.UnmarshalBinary(data[offset:], &decodeContext{
			SliceWidth:  s.SliceWidth,
			SliceHeight: s.SliceHeight,
		})
		if err != nil {
			return err
		}
		offset += n
		s.Stacks = append(s.Stacks, stack)
	}

	return nil
}

type Stack struct {
	Name       string
	Animations []Animation
}

func (s *Stack) Animation(name string) *Animation {
	for i := range s.Animations {
		if s.Animations[i].Name == name {
			return &s.Animations[i]
		}
	}
	return nil
}

func (s *Stack) UnmarshalBinary(data []byte, ctx *decodeContext) (int, error) {
	if len(data) < 1 {
		return 0, ErrStackDataTooShort
	}
	offset := 0
	nameLength := int(data[offset])
	offset++
	if len(data) < offset+nameLength {
		return 0, ErrStackDataTooShort
	}
	s.Name = string(data[offset : offset+nameLength])
	offset += nameLength

	if len(data) < offset+2 {
		return 0, ErrStackDataTooShort
	}
	sliceCount := int(data[offset])<<8 | int(data[offset+1])
	offset += 2

	if len(data) < offset+2 {
		return 0, ErrStackDataTooShort
	}
	animationCount := int(data[offset])<<8 | int(data[offset+1])
	offset += 2

	ctx.SliceCount = sliceCount

	for i := 0; i < animationCount; i++ {
		var anim Animation
		n, err := anim.UnmarshalBinary(data[offset:], ctx)
		if err != nil {
			return 0, err
		}
		offset += n
		s.Animations = append(s.Animations, anim)
	}

	return offset, nil
}

type Animation struct {
	Name      string
	Frames    []Frame
	FrameTime int
}

func (a *Animation) Frame(index int) *Frame {
	if index < 0 || index >= len(a.Frames) {
		return nil
	}
	return &a.Frames[index]
}

func (a *Animation) UnmarshalBinary(data []byte, ctx *decodeContext) (int, error) {
	if len(data) < 1 {
		return 0, ErrAnimationDataTooShort
	}
	offset := 0
	nameLength := int(data[offset])
	offset++
	if len(data) < offset+nameLength {
		return 0, ErrAnimationDataTooShort
	}
	a.Name = string(data[offset : offset+nameLength])
	offset += nameLength

	if len(data) < offset+4 {
		return 0, ErrAnimationDataTooShort
	}
	// Frametime in little endian
	a.FrameTime = int(data[offset])<<24 | int(data[offset+1])<<16 | int(data[offset+2])<<8 | int(data[offset+3])
	offset += 4

	// frame count uint16
	if len(data) < offset+2 {
		return 0, ErrAnimationDataTooShort
	}
	frameCount := int(data[offset])<<8 | int(data[offset+1])
	offset += 2

	for i := 0; i < frameCount; i++ {
		var frame Frame
		ctx.x = 0
		if err := frame.UnmarshalBinary(data[offset:], ctx); err != nil {
			return 0, err
		}
		ctx.y += ctx.SliceHeight
		a.Frames = append(a.Frames, frame)
		offset += ctx.SliceCount
	}
	return offset, nil
}

type Frame struct {
	Slices []Slice
}

func (f *Frame) Slice(index int) *Slice {
	if index < 0 || index >= len(f.Slices) {
		return nil
	}
	return &f.Slices[index]
}

func (f *Frame) UnmarshalBinary(data []byte, ctx *decodeContext) error {
	for i := 0; i < ctx.SliceCount; i++ {
		slice := Slice{
			X:       ctx.x,
			Y:       ctx.y,
			Shading: data[0],
		}
		ctx.x += ctx.SliceWidth
		f.Slices = append(f.Slices, slice)
	}
	return nil
}

type Slice struct {
	X       int
	Y       int
	Shading uint8
}

var ErrInvalidPNG = errors.New("invalid png")
var ErrNoStaxChunk = errors.New("no stAx chunk")
var ErrInvalidVersion = errors.New("invalid version")
var ErrDataTooShort = errors.New("data too short")
var ErrStaxDataTooShort = errors.New("stax data too short")
var ErrStackDataTooShort = errors.New("stack data too short")
var ErrAnimationDataTooShort = errors.New("animation data too short")
