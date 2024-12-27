package ables

// IDable be id. It is unique per referable.
type IDable struct {
	id int
}

// ID returns the ID of the referable.
func (i IDable) ID() int {
	return i.id
}

// SetID sets the ID of the referable.
func (i *IDable) SetID(id int) {
	i.id = id
}

// MakeIDable makes a new IDable.
func MakeIDable(id int) IDable {
	return IDable{id: id}
}

// NextIDable returns a new IDable with the next ID.
func NextIDable() IDable {
	return MakeIDable(NextID())
}

// idCounter is used for global id counting, yo.
var idCounter int

// NextID returns the next ID.
func NextID() int {
	idCounter++
	return idCounter
}
