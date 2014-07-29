package shared

// A simple structure to keep track of the given input from the client.  Instead of just sending
// keypresses directly, we can use this so multiple (or no) keypresses can be sent in one go.
type InputState struct {
	KeyLeftDown  bool
	KeyRightDown bool
	KeyDownDown  bool
	KeyUpDown    bool
}

// Check if any of the key in the state are actually being pressed.
func (this *InputState) HasInput() bool {
	hasInputs := false

	if this.KeyLeftDown {
		hasInputs = true
	}

	if this.KeyRightDown {
		hasInputs = true
	}

	if this.KeyDownDown {
		hasInputs = true
	}

	if this.KeyUpDown {
		hasInputs = true
	}

	return hasInputs
}
