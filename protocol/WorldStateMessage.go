package protocol

import (
	"encoding/json"
	"time"

	"github.com/gabriel-comeau/multiplayer-game-test/shared"
)

// A WorldStateMessage is used to send a list of entities to each client after every server tick.
// It is the structure used to convey to the clients what the state of the world according to the
// server is.
type WorldStateMessage struct {
	MessageType MessageType
	SentTime    time.Time
	RcvdTime    time.Time
	Entities    []MessageEntity
}

// Convert the message into JSON representation and return the raw bytes
func (m *WorldStateMessage) Encode() []byte {
	bytes, err := json.Marshal(m)
	if err != nil {
		panic(err.Error())
	}

	return AddNewlineToByteSlice(bytes)
}

// Satisfy the Message interface
func (m *WorldStateMessage) GetSentTime() time.Time {
	return m.SentTime
}

// Satisfy the Message interface
func (m *WorldStateMessage) GetRcvdTime() time.Time {
	return m.RcvdTime
}

// Satisfy the Message interface
func (m *WorldStateMessage) GetMessageType() MessageType {
	return m.MessageType
}

// Constructor function to create a new WorldStateMessage and return a pointer to it
func CreateWorldStateMessage(entities []MessageEntity) *WorldStateMessage {
	msg := new(WorldStateMessage)
	msg.SentTime = time.Now()
	msg.MessageType = WORLD_STATE_MESSAGE
	msg.Entities = entities
	return msg
}

// Take in a raw byte slice and convert it back to a WorldStateMessage pointer
func DecodeWorldStateMessage(raw []byte) *WorldStateMessage {
	msg := new(WorldStateMessage)
	err := json.Unmarshal(raw, msg)
	if err != nil {
		panic(err.Error())
	}

	return msg
}

// A MessageEntity represents the state of an entity on the server as it is conveyed to the client.
// There should be one of these for each player currently in the server's world state.
type MessageEntity struct {
	Id       int64
	Position shared.FloatVector
	LastSeq  int64
}

// Create a new MessageEntity.  Don't bother making a pointer to it, it's a very small struct.  If
// ever we need to send hundreds of these at once we might consider making it a pointer for memory
// efficiency.
func CreateMessageEntity(id int64, pos shared.FloatVector, seq int64) MessageEntity {
	return MessageEntity{Id: id, Position: pos, LastSeq: seq}
}
