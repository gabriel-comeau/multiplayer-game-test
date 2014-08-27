package main

import (
	sf "bitbucket.org/krepa098/gosfml2"

	"github.com/gabriel-comeau/multiplayer-game-test/shared"
)

// Converts a shared.FloatVector to an SFML Vector2f
func ConvertToSFMLVector(vec shared.FloatVector) sf.Vector2f {
	return sf.Vector2f{X: vec.X, Y: vec.Y}
}

// Converts an SFML Vector2f to a shared.FloatVector
func ConvertToFloatVector(vec sf.Vector2f) shared.FloatVector {
	return shared.FloatVector{X: vec.X, Y: vec.Y}
}
