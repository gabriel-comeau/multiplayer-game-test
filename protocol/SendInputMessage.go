package protocol

import (
	"encoding/json"
	"time"

	"github.com/gabriel-comeau/multiplayer-game-test/shared"
)

// Message sent to server by client indicating their current input state this tick
type SendInputMessage struct {
	MessageType MessageType
	SendTime    time.Time
	RcvdTime    time.Time
	Input       *shared.InputState
	Dt          shared.MDuration
	Seq         int64
	PlayerId    int64
}

// Encode the message to JSON format and get the raw bytes
func (m *SendInputMessage) Encode() []byte {
	bytes, err := json.Marshal(m)
	if err != nil {
		panic(err.Error())
	}

	return AddNewlineToByteSlice(bytes)
}

// Message interface
func (m *SendInputMessage) GetSentTime() time.Time {
	return m.SendTime
}

// Message interface
func (m *SendInputMessage) GetRcvdTime() time.Time {
	return m.RcvdTime
}

// Message interface
func (m *SendInputMessage) SetRcvdTime(t time.Time) {
	m.RcvdTime = t
}

// Message interface
func (m *SendInputMessage) GetMessageType() MessageType {
	return m.MessageType
}

// Constructor, returns a pointer to a SendInputMessage
func CreateSendInputMessage(inputState *shared.InputState, seq int64, dt time.Duration, playerId int64) *SendInputMessage {
	return &SendInputMessage{
		SendTime:    time.Now(),
		MessageType: SEND_INPUT_MESSAGE,
		Input:       inputState,
		Seq:         seq,
		Dt:          shared.MDuration{dt},
		PlayerId:    playerId,
	}
}

// Decode a SendInputMessage from raw bytes of JSON data and return a pointer to it
func DecodeSendInputMessage(raw []byte) *SendInputMessage {
	msg := new(SendInputMessage)
	err := json.Unmarshal(raw, msg)
	if err != nil {
		panic(err.Error())
	}

	return msg
}
