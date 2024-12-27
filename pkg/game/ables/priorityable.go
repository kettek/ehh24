package ables

// Priorityable is a terrible thing.
type Priorityable struct {
	priority int
	offset   int // Offset within.. a priority range.
}

// Priority returns the priority + offset of the referable.
func (p Priorityable) Priority() int {
	return p.priority + p.offset
}

// SetPriority sets the priority of the referable.
func (p *Priorityable) SetPriority(priority int) {
	p.priority = priority
}

// Offset returns the offset of the referable.
func (p Priorityable) Offset() int {
	return p.offset
}

// SetOffset sets the offset of the referable.
func (p *Priorityable) SetOffset(offset int) {
	p.offset = offset
}

// MakePriorityable makes a new Priorityable.
func MakePriorityable(priority int) Priorityable {
	return Priorityable{priority: priority}
}

// Default priority ranges.
const (
	PriorityBack    = 0
	PriorityMiddle  = 100000
	PriorityFront   = 200000
	PriorityOverlay = 300000
	PriorityUI      = 400000
	PriorityBeyond  = 500000
)
