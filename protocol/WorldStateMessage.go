package protocol

import (
	"encoding/json"
	"time"

	sf "bitbucket.org/krepa098/gosfml2"
)

// A WorldStateMessage is used to send a list of entities to each client after every server tick.
// It is the structure used to convey to the clients what the state of the world according to the
// server is.
type WorldStateMessage struct {
	MessageType MessageType
	Timestamp   int64
	Entities    []MessageEntity
}

// Convert the message into JSON representation and return the raw bytes
func (this *WorldStateMessage) Encode() []byte {
	bytes, err := json.Marshal(this)
	if err != nil {
		panic(err.Error())
	}

	return AddNewlineToByteSlice(bytes)
}

// Satisfy the Message interface
func (this *WorldStateMessage) GetTimestamp() int64 {
	return this.Timestamp
}

// Satisfy the Message interface
func (this *WorldStateMessage) GetMessageType() MessageType {
	return this.MessageType
}

// Constructor function to create a new WorldStateMessage and return a pointer to it
func CreateWorldStateMessage(entities []MessageEntity) *WorldStateMessage {
	t := time.Now()
	msg := new(WorldStateMessage)
	msg.Timestamp = t.UnixNano()
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
	Position sf.Vector2f
	LastSeq  int64
}

// Create a new MessageEntity.  Don't bother making a pointer to it, it's a very small struct.  If
// ever we need to send hundreds of these at once we might consider making it a pointer for memory
// efficiency.
func CreateMessageEntity(id int64, pos sf.Vector2f, seq int64) MessageEntity {
	return MessageEntity{Id: id, Position: pos, LastSeq: seq}
}
