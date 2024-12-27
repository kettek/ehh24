package ables

// Priorityable is a terrible thing.
type Priorityable struct {
	priority int
}

// Priority returns the priority of the referable.
func (p Priorityable) Priority() int {
	return p.priority
}

// SetPriority sets the priority of the referable.
func (p *Priorityable) SetPriority(priority int) {
	p.priority = priority
}

// MakePriorityable makes a new Priorityable.
func MakePriorityable(priority int) Priorityable {
	return Priorityable{priority: priority}
}

// Default priority ranges.
const (
	PriorityBack    = 0
	PriorityMiddle  = 10000
	PriorityFront   = 20000
	PriorityOverlay = 30000
	PriorityUI      = 40000
	PriorityBeyond  = 50000
)
