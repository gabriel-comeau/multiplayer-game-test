package protocol

import (
	"encoding/json"
	"time"

	"github.com/gabriel-comeau/multiplayer-game-test/shared"
)

// Message sent to server by client indicating their current input state this tick
type SendInputMessage struct {
	MessageType MessageType
	Timestamp   int64
	Input       *shared.InputState
	Dt          time.Duration
	Seq         int64
	PlayerId    int64
}

// Encode the message to JSON format and get the raw bytes
func (this *SendInputMessage) Encode() []byte {
	bytes, err := json.Marshal(this)
	if err != nil {
		panic(err.Error())
	}

	return AddNewlineToByteSlice(bytes)
}

// Satisfy the Message interface
func (this *SendInputMessage) GetTimestamp() int64 {
	return this.Timestamp
}

// Satisfy the Message interface
func (this *SendInputMessage) GetMessageType() MessageType {
	return this.MessageType
}

// Constructor, returns a pointer to a SendInputMessage
func CreateSendInputMessage(inputState *shared.InputState, seq int64, dt time.Duration, playerId int64) *SendInputMessage {
	t := time.Now()
	msg := new(SendInputMessage)
	msg.Timestamp = t.UnixNano()
	msg.MessageType = SEND_INPUT_MESSAGE
	msg.Input = inputState
	msg.Seq = seq
	msg.Dt = dt
	msg.PlayerId = playerId
	return msg
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
