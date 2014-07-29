package protocol

import (
	"encoding/json"
	"time"
)

// A message sent to a client upon their initial connection in order to let them know
// what their unique ID is.
type PlayerUUIDMessage struct {
	MessageType MessageType
	Timestamp   int64
	UUID        int64
}

// Encode the message to JSON format and get the raw bytes
func (this *PlayerUUIDMessage) Encode() []byte {
	bytes, err := json.Marshal(this)
	if err != nil {
		panic(err.Error())
	}

	return AddNewlineToByteSlice(bytes)
}

// Satisfy the Message interface
func (this *PlayerUUIDMessage) GetTimestamp() int64 {
	return this.Timestamp
}

// Satisfy the Message interface
func (this *PlayerUUIDMessage) GetMessageType() MessageType {
	return this.MessageType
}

// Constructor for PlayerUUIDMessage, returns pointer to one
func CreatePlayerUUIDMessage(uuid int64) *PlayerUUIDMessage {
	t := time.Now()
	msg := new(PlayerUUIDMessage)
	msg.Timestamp = t.UnixNano()
	msg.MessageType = PLAYER_UUID_MESSAGE
	msg.UUID = uuid
	return msg
}

// Decode a PlayerUUIDMessage from raw bytes of JSON data and return a pointer to it
func DecodePlayerUUIDMessage(raw []byte) *PlayerUUIDMessage {
	msg := new(PlayerUUIDMessage)
	err := json.Unmarshal(raw, msg)
	if err != nil {
		panic(err.Error())
	}

	return msg
}
