package shared

// Basically a rewrite of GoSFML2's Vector2f type.  It's the only thing from the entire
// package being used by the server side of the application which means that the server is
// basically dependant on having the SFML2 libs installed (which have a long dependancy chain of
// things which don't apply to servers like opengl libs).  Instead we just use this type on the
// server and convert back and forth on the client.
type FloatVector struct {
	X float32
	Y float32
}

// Add this vector to another
func (v FloatVector) Plus(other FloatVector) FloatVector {
	return FloatVector{X: v.X + other.X, Y: v.Y + other.Y}
}

// Subtract a vector from this vector
func (v FloatVector) Minus(other FloatVector) FloatVector {
	return FloatVector{X: v.X - other.X, Y: v.Y - other.Y}
}
