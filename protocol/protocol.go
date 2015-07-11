package protocol

import (
	"encoding/json"
	"errors"
	"time"
)

const (
	PLAYER_UUID_MESSAGE MessageType = iota + 1
	SEND_INPUT_MESSAGE
	WORLD_STATE_MESSAGE
)

// Enum to keep track of message types
type MessageType int

// Interface for generic network messages which can be serialized to JSON
type Message interface {
	GetSentTime() time.Time
	GetRcvdTime() time.Time
	SetRcvdTime(t time.Time)
	GetMessageType() MessageType
	Encode() []byte
}

// Figure out what a message is from its JSON representation and return the specific instance of
// it.
func DecodeMessage(raw []byte) (Message, error) {
	// First we're going to marshall the raw JSON into a blank interface.  This will let
	// us get a peek at the message type field.
	unknown := make(map[string]interface{})
	json.Unmarshal(raw, &unknown)

	// Assuming the JSON is a map, we should have a MessageType key
	mUntyped, ok := unknown["MessageType"]
	if !ok {
		return nil, errors.New("Invalid message recieved (type key not present)")
	}

	// All JSON numbers get stored as float64, so assert it is present
	mFloat, ok := mUntyped.(float64)
	if !ok {
		return nil, errors.New("Invalid message recieved (type nan)")
	}

	// Finally if we've got the float we can cast it to MessageType (int)
	mType := MessageType(mFloat)

	switch mType {
	case PLAYER_UUID_MESSAGE:
		return DecodePlayerUUIDMessage(raw), nil
	case SEND_INPUT_MESSAGE:
		return DecodeSendInputMessage(raw), nil
	case WORLD_STATE_MESSAGE:
		return DecodeWorldStateMessage(raw), nil
	}

	return nil, errors.New("The message type matched nothing")
}

// Newlines are easier delimiters
func AddNewlineToByteSlice(raw []byte) []byte {
	str := string(raw)
	str += "\n"
	return []byte(str)
}
