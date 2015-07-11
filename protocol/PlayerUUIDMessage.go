package protocol

import (
	"encoding/json"
	"time"
)

// A message sent to a client upon their initial connection in order to let them know
// what their unique ID is.
type PlayerUUIDMessage struct {
	MessageType MessageType
	SentTime    time.Time
	RcvdTime    time.Time
	UUID        int64
}

// Encode the message to JSON format and get the raw bytes
func (m *PlayerUUIDMessage) Encode() []byte {
	bytes, err := json.Marshal(m)
	if err != nil {
		panic(err.Error())
	}

	return AddNewlineToByteSlice(bytes)
}

// Message interface
func (m *PlayerUUIDMessage) GetSentTime() time.Time {
	return m.SentTime
}

// Message interface
func (m *PlayerUUIDMessage) GetRcvdTime() time.Time {
	return m.RcvdTime
}

// Message interface
func (m *PlayerUUIDMessage) SetRcvdTime(t time.Time) {
	m.RcvdTime = t
}

// Message interface
func (m *PlayerUUIDMessage) GetMessageType() MessageType {
	return m.MessageType
}

// Constructor for PlayerUUIDMessage, returns pointer to one
func CreatePlayerUUIDMessage(uuid int64) *PlayerUUIDMessage {
	return &PlayerUUIDMessage{
		SentTime:    time.Now(),
		MessageType: PLAYER_UUID_MESSAGE,
		UUID:        uuid,
	}
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
