package ables

// Tagable b tag
type Tagable struct {
	tag string
}

// Tag returns the tag of the referable.
func (t Tagable) Tag() string {
	return t.tag
}

// SetTag sets the tag of the referable.
func (t *Tagable) SetTag(tag string) {
	t.tag = tag
}

// MakeTagable makes a new Tagable.
func MakeTagable(tag string) Tagable {
	return Tagable{tag: tag}
}
