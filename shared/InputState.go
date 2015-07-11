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
func (i *InputState) HasInput() bool {

	if i.KeyLeftDown {
		return true
	}

	if i.KeyRightDown {
		return true
	}

	if i.KeyDownDown {
		return true
	}

	if i.KeyUpDown {
		return true
	}

	return false
}
