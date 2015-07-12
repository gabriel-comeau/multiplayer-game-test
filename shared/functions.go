package shared

// Taking into account the max speed constant (pixels per second an object can move),
// use the current input and the frame time to create a vector representing the offset for how
// far a unit moved and in which direction.
//
// Both the client and the server use this calculation so it belongs to the shared package.
func GetVectorFromInputAndDt(inputState *InputState, dt MDuration) FloatVector {
	dtFloatSeconds := float32(dt.Seconds())
	velocity := FloatVector{X: 0, Y: 0}

	if inputState.KeyDownDown && !inputState.KeyUpDown {
		velocity.Y = (SPEED * dtFloatSeconds)
	}

	if inputState.KeyUpDown && !inputState.KeyDownDown {
		velocity.Y = (SPEED * dtFloatSeconds) * -1
	}

	if inputState.KeyLeftDown && !inputState.KeyRightDown {
		velocity.X = (SPEED * dtFloatSeconds) * -1
	}

	if inputState.KeyRightDown && !inputState.KeyLeftDown {
		velocity.X = (SPEED * dtFloatSeconds)
	}

	return velocity
}
