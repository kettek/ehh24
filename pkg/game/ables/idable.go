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
