package main

import (
	"bufio"
	"log"
	"net"
	"time"

	"github.com/gabriel-comeau/multiplayer-game-test/protocol"
	"github.com/gabriel-comeau/multiplayer-game-test/shared"
)

// Link a client-id to a network connection
type Client struct {
	clientId int64
	conn     net.Conn
}

const (
	// How long we want to have between iterations of the main server loop.  The loop will sleep
	// for this time minus however long it took (assuming that difference is positive of course)
	SLEEP_TIME time.Duration = 33 * time.Millisecond
)

var (
	// The unique ID generator - see IdGenerator.go
	idGen *IdGenerator

	// The thread-safe client holder.  See ClientHolder.go
	clientHolder *ClientHolder

	// The thread-safe map of entities present in the game
	entityHolder *EntityHolder

	// The queue of messages to process each iteration of the main loop
	messageQueue *protocol.MessageQueue
)

func init() {
	idGen = CreateIdGenerator()
	clientHolder = CreateClientHolder()
	entityHolder = CreateEntityHolder()
	messageQueue = protocol.CreateMessageQueue()
}

func main() {

	// Start listening on the socket for incoming connections
	go listenForConns()

	// Preset up the timestep stuff so there's a value for the first iteration of the loop
	lastTick := time.Now()
	var dt time.Duration

	for {
		messages := messageQueue.PopAll()
		for _, message := range messages {

			// The only message expected FROM the client is the move message
			// so lets look for that one
			if message.GetMessageType() != protocol.SEND_INPUT_MESSAGE {
				log.Print("Got an invalid message type from client:", string(message.GetMessageType()))
				break
			}

			typed, ok := message.(*protocol.SendInputMessage)
			if !ok {
				log.Print("Message couldn't be asserted into SendInputMessage")
				break
			}

			// If the message is valid.  It checks that the DT is not too big basically -
			// this is a lame validation but for now it prevents any obviously bad manipulation.
			//
			// In a better system we could check the times of when we received this message and the
			// ones before it, to make sure the time delta fits within that frame.  For now we'll do
			// a simple clamp type approach.
			if !validateMessageDt(typed.Dt) {
				log.Print("Dropping invalid send input message - DT sucks: ", typed.Dt)
				break

			}
			// Get the vector for the move
			moveVec := shared.GetVectorFromInputAndDt(typed.Input, typed.Dt)

			// Get the seq
			seq := typed.Seq

			// Move the unit -- it could have been removed so make sure it's still there
			ent := entityHolder.GetEntity(typed.PlayerId)
			if ent == nil {
				continue
			}

			ent.Move(moveVec)

			// Apply the new last sequence number
			if ent.lastSeq < seq {
				ent.lastSeq = seq
			}

			//if ent.lastSeqTime.Before(u)
		}

		// OK, all messages processed for this tick, send out an entity message
		// We'll take stock of where all the entities are and send out an updated world state.
		msgEnts := make([]protocol.MessageEntity, 0)
		for _, ent := range entityHolder.GetEntities() {
			msgEnts = append(msgEnts, protocol.MessageEntity{Id: ent.entityId, Position: ent.position, LastSeq: ent.lastSeq})
		}

		worldStateMessage := protocol.CreateWorldStateMessage(msgEnts)
		broadcastMessage(worldStateMessage)

		// Get how long it took to do all of this
		now := time.Now()
		dt = now.Sub(lastTick)
		lastTick = now

		// If it took less long than SLEEP_TIME, sleep for the difference, otherwise this ends
		// up sending A LOT of messages with no changes to the clients.
		//
		// Alternatively, we could check to see if the world state has changed and only send it out
		// when something new has happened.
		if dt < SLEEP_TIME {
			time.Sleep(SLEEP_TIME - dt)
		}
	}
}

// Concurrent function which spins in a loop, listening for new connections on the socket.  When it
// gets one, it generates an ID for the new user, creates a client object which gets put into the
// ClientHolder and then sends that client their ID.
func listenForConns() {
	server, err := net.Listen("tcp", ":"+shared.PORT)
	if server == nil || err != nil {
		panic("couldn't start listening: " + err.Error())
	}

	log.Print("SERVER LISTENING")

	for {
		newConn, err := server.Accept()

		if err != nil {
			log.Printf("ERROR DURING ACCEPT: %v", err)
		}

		if newConn != nil {
			playerId := idGen.GetNextId()
			log.Printf("ACCEPTED: %v <-> %v\n", newConn.LocalAddr(), newConn.RemoteAddr())
			log.Printf("Player # is: %v\n", playerId)

			player := CreatePlayerEntity(playerId, shared.FloatVector{X: 30, Y: 30})
			entityHolder.AddEntity(player)

			client := new(Client)
			client.conn = newConn
			client.clientId = playerId
			clientHolder.AddClient(client)

			sendUUIDToPlayer(playerId, client)

			go handleClient(client)
		}
	}
}

// Handle an individual client connection.  Runs concurrently in a goroutine.  As it recieves new
// input messages, it puts them in the global MessageQueue so they'll be processed by the main server
// loop.  Also responsible for handling client disconnection.
func handleClient(client *Client) {
	log.Println("Handlin' client")
	b := bufio.NewReader(client.conn)
	for {
		line, err := b.ReadBytes('\n')
		if err != nil {
			break
		}

		if string(line) == "" || string(line) == "\n" {
			continue
		}

		// Dispatch client messages
		message, err := protocol.DecodeMessage(line)
		if err != nil {
			log.Println("Error when reading message:", err.Error())
			continue
		}

		if validateMessageClientId(message, client.clientId) {
			message.SetRcvdTime(time.Now())
			messageQueue.PushMessage(message)
		}
	}

	// EOF happened - this client has disconnected
	log.Printf("Player: %v left\n", client.clientId)
	clientHolder.RemoveClient(client)

	// remove the entity from the holder
	entityHolder.RemoveEntity(client.clientId)
}

// See note in handle client for how this could be improved.
func validateMessageDt(dt time.Duration) bool {
	if dt > shared.MAX_DT {
		return false
	}
	return true
}

// Ensure that the message is coming from the right client so no one tries any funny
// business.
func validateMessageClientId(message protocol.Message, clientId int64) bool {
	// check if this is a SendInputMessage
	if message.GetMessageType() == protocol.SEND_INPUT_MESSAGE {
		typed, ok := message.(*protocol.SendInputMessage)
		if !ok {
			log.Println("Message couldn't be asserted into SendInputMessage")
			return false
		}

		if typed.PlayerId == clientId {
			return true
		} else {
			return false
		}
	}

	// The other messages don't come from players so this doesn't make any sense.
	log.Println("Someone sent a bad message to the server - only expecting SEND_INPUT_MESSAGE, got: ", message.GetMessageType())
	return false
}

// Sends a UUID message to a player.
func sendUUIDToPlayer(id int64, client *Client) {
	msg := protocol.CreatePlayerUUIDMessage(id)
	sendMessageToClient(msg, id)
}

// Send a message to all players
func broadcastMessage(msg protocol.Message) {
	encoded := msg.Encode()
	for _, c := range clientHolder.GetClients() {
		c.conn.Write(encoded)
	}
}

// Send a message to a specific player
func sendMessageToClient(msg protocol.Message, cid int64) {
	c := clientHolder.GetClient(cid)
	if c != nil {
		c.conn.Write(msg.Encode())
	}
}
