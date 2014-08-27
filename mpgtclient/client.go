package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"runtime"
	"time"

	sf "bitbucket.org/krepa098/gosfml2"

	"github.com/gabriel-comeau/multiplayer-game-test/protocol"
	"github.com/gabriel-comeau/multiplayer-game-test/shared"
)

var (
	// This block of variables is shared global state throughout the client.  Obviously not great
	// but since our client program is pretty simple, this is quick and effective.

	// Our current movement velocity
	velocity sf.Vector2f

	// Current state of input - which buttons are being pressed
	inputState *shared.InputState

	// Socket connection to the server
	conn net.Conn

	// Unique identifying ID sent over by the server on connection
	myPlayerId int64

	// Keep track of the entities we need to draw.  The key is their UUID.  Our player entity
	// is just another in this list.
	entities map[int64]*Unit

	// Hold messages from the server in a queue
	messageQueue *protocol.MessageQueue

	// Channel to send outgoing messages on
	outgoing chan protocol.Message

	// Keep a list of inputs that have been processed locally through client-side prediction
	// but haven't yet been acknowledged on the server
	unacked []*protocol.SendInputMessage

	// Current sequence number for input messages we'll be sending
	currentSeq int64
)

func init() {
	runtime.LockOSThread()
	inputState = new(shared.InputState)
	entities = make(map[int64]*Unit)
	messageQueue = protocol.CreateMessageQueue()
	outgoing = make(chan protocol.Message)
	currentSeq = 0
	unacked = make([]*protocol.SendInputMessage, 0)
}

func main() {

	// Open the game window.
	renderWindow := sf.NewRenderWindow(sf.VideoMode{1024, 768, 32}, "Wow! Much client-side-interpretation", sf.StyleDefault, sf.DefaultContextSettings())

	// Because we send a message to the server every frame where input is present, we'll limit the frames so that
	// machines capable of rendering hundreds of frames per second don't try to send hundreds of network message
	// per second.
	//
	// This could be handled better of course, having artificial limiting on only the networking
	// portion, but for simplicity's sake, this will do the job.
	renderWindow.SetFramerateLimit(60)

	// establish connection to server
	conn = connectToServer()

	// Preset up the timestep stuff so there's a value for the first rendered frame
	lastTick := time.Now()
	var dt time.Duration = 0

	// This is the start of our main game loop.  As long as the window remains open, this will continue.
	for renderWindow.IsOpen() {

		// process user input, changing the value of the inputstate struct
		inputState = handleUserInput(renderWindow, inputState)
		if inputState.HasInput() {
			velocity = ConvertToSFMLVector(shared.GetVectorFromInputAndDt(inputState, dt))

			// client side prediction
			player, ok := entities[myPlayerId]
			if ok {
				player.Move(velocity)
			}

			// We need to send this input to the server, so build a message object
			// stick it in the unacked map and then transmit
			inputMsg := protocol.CreateSendInputMessage(inputState, currentSeq, dt, myPlayerId)
			unacked = append(unacked, inputMsg)
			currentSeq++
			outgoing <- inputMsg
		}

		// process incoming messages
		incoming := messageQueue.PopAll()
		for _, message := range incoming {
			switch message.GetMessageType() {

			case protocol.WORLD_STATE_MESSAGE:
				typed, ok := message.(*protocol.WorldStateMessage)
				if !ok {
					fmt.Println("Got a message with WORLD_STATE_MESSAGE id but couldn't be cast")
					continue
				}

				// Now we're going to iterate through the entities contained in the message.  On the
				// first pass we're going to do a couple of different things.  If the entity doesn't
				// exist, we're going to add it.  If the entity is ours, we'll apply interpolation based
				// on our past inputs.  If the entity is someone else's and exists we'll move it.

				for _, msgEnt := range typed.Entities {
					// first check if this thing exists at all
					existingEnt, ok := entities[msgEnt.Id]
					if !ok {
						// not in the map, let's create it - we'll bail after this because even if
						// this is our own entity we'll start to worry about interpolation on the next
						// pass only
						addEntityToGameWorld(msgEnt.Id, ConvertToSFMLVector(msgEnt.Position))
						continue
					}

					// Is this ours or someone else's?
					if msgEnt.Id == myPlayerId {

						// First, set the position to wherever the server thinks it was
						existingEnt.SetPosition(ConvertToSFMLVector(msgEnt.Position))

						// Next, let's go through our pending inputs list and get rid of everything older
						// than this seq number
						newUnacked := make([]*protocol.SendInputMessage, 0)
						for _, oldMsg := range unacked {
							if oldMsg.Seq > msgEnt.LastSeq {
								// not processed yet, so reapply and keep it in the list
								newUnacked = append(newUnacked, oldMsg)
								existingEnt.Move(ConvertToSFMLVector(shared.GetVectorFromInputAndDt(oldMsg.Input, oldMsg.Dt)))
							}
						}
						unacked = newUnacked

					} else {
						// This is someone else's entity, so just move it
						existingEnt.SetPosition(ConvertToSFMLVector(msgEnt.Position))
					}
				}

				// Now, we've also got to check to see if we still have any remaining entities
				// that belong to players who've left.  This means we have to do another ugly iteration
				// but that's life.
				removeDisconnectedPlayers(typed.Entities)
			}
		}

		now := time.Now()
		dt = now.Sub(lastTick)
		lastTick = now

		renderWindow.Clear(sf.Color{0, 0, 0, 0})

		// Draw all the units but draw the player last so it's always on top
		var playerUnit *Unit
		for unitId, unit := range entities {
			if unitId == myPlayerId {
				playerUnit = unit
			} else {
				unit.Draw(renderWindow, sf.DefaultRenderStates())
			}
		}

		if playerUnit != nil {
			playerUnit.Draw(renderWindow, sf.DefaultRenderStates())
		}

		renderWindow.Display()
	}
}

// Look over the events coming in, check them against the current keystates, and then update
// the keystates to match.  This is one of those bad functions which mutates the package-wide
// keystate struct but really this is the only function which writes to it so why bother copying
// things?
func handleUserInput(renderWindow *sf.RenderWindow, inputState *shared.InputState) *shared.InputState {

	// Handle user input
	for event := renderWindow.PollEvent(); event != nil; event = renderWindow.PollEvent() {
		switch ev := event.(type) {
		case sf.EventKeyReleased:
			switch ev.Code {

			case sf.KeyEscape:
				renderWindow.Close()

			case sf.KeyLeft:
				if inputState.KeyLeftDown {
					inputState.KeyLeftDown = false
				}

			case sf.KeyRight:
				if inputState.KeyRightDown {
					inputState.KeyRightDown = false
				}

			case sf.KeyDown:
				if inputState.KeyDownDown {
					inputState.KeyDownDown = false
				}

			case sf.KeyUp:
				if inputState.KeyUpDown {
					inputState.KeyUpDown = false
				}

			}

		case sf.EventKeyPressed:
			switch ev.Code {

			case sf.KeyLeft:
				if !inputState.KeyLeftDown {
					inputState.KeyLeftDown = true
				}

			case sf.KeyRight:
				if !inputState.KeyRightDown {
					inputState.KeyRightDown = true
				}

			case sf.KeyDown:
				if !inputState.KeyDownDown {
					inputState.KeyDownDown = true
				}

			case sf.KeyUp:
				if !inputState.KeyUpDown {
					inputState.KeyUpDown = true
				}
			}
		}
	}

	return inputState
}

// Add a new entity to the game world.  It will use the Player constructor (for the player texture)
// if the entity ID matches the player ID and the Other constructor otherwise.
func addEntityToGameWorld(id int64, pos sf.Vector2f) {
	if id == myPlayerId {
		_, ok := entities[myPlayerId]
		if ok {
			fmt.Println("ERROR - tried to add a new player with the same ID")
			return
		}

		player := NewPlayer(pos)
		entities[myPlayerId] = player
	} else {
		_, ok := entities[id]
		if ok {
			fmt.Println("ERROR - tried to add a new other entity with an ID that was already in the system")
			return
		}

		other := NewOther(pos)
		entities[id] = other
	}
}

// Go through the entities which the server is tracking and compare to the entities the client is
// tracking.  Remove any client-side entities that aren't in the server's list.
func removeDisconnectedPlayers(serverEntities []protocol.MessageEntity) {
	// We'll make a map of the IDs so we can do more convenient lookups
	// I used bool because it's small.
	serverEnts := make(map[int64]bool)
	for _, ent := range serverEntities {
		serverEnts[ent.Id] = true
	}

	// Now iterate through the map of stored client-side entities.  Remove anything
	// that doesn't appear in the server entities map
	for id, _ := range entities {
		_, ok := serverEnts[id]
		if !ok {
			if id == myPlayerId {
				panic("Whoops trying to remove myself - this is an error condition so we out")
			}

			delete(entities, id)
		}
	}
}

// Establish a connection to the game server and start the two goroutines which require it
func connectToServer() net.Conn {
	conn, err := net.Dial("tcp", shared.HOST+":"+shared.PORT)
	if err != nil {
		panic(err.Error())
	}

	// Now we're going to wait for the server to give us an entity ID
	b := bufio.NewReader(conn)
	for {
		line, err := b.ReadBytes('\n')

		if err != nil {
			conn.Close()
			fmt.Println("Error while trying to accept player id")
			os.Exit(1)
			break
		}

		if string(line) == "" || string(line) == "\n" {
			continue
		}

		message, err := protocol.DecodeMessage(line)
		if err != nil {
			fmt.Println("ERROR during decode: " + err.Error())
			continue
		}

		if message.GetMessageType() == protocol.PLAYER_UUID_MESSAGE {
			typed, ok := message.(*protocol.PlayerUUIDMessage)
			if !ok {
				fmt.Println("Message couldn't be asserted into PlayerUUIDMessage though that was message id")
				conn.Close()
				os.Exit(1)
			}
			myPlayerId = typed.UUID
			break
		} else {
			fmt.Println("Got the wrong type of message - expected PLAYER_UUID_MESSAGE")
			conn.Close()
			os.Exit(1)
		}
	}

	go listenForMessages(conn)
	go writeMessages(conn, outgoing)

	return conn
}

// Listens for incoming messages from the server, decodes the serialized versions into
// message objects and then pushes them into the message queue.
//
// This is a concurrent function - it runs simultaneously to the main game loop as a goroutine
func listenForMessages(conn net.Conn) {
	b := bufio.NewReader(conn)
	for {
		line, err := b.ReadBytes('\n')

		if err != nil {
			conn.Close()
			fmt.Println("ERROR, CLOSING CONN: " + err.Error())
			break
		}

		if string(line) == "" || string(line) == "\n" {
			continue
		}

		// Deal with incoming messages from the server
		message, err := protocol.DecodeMessage(line)
		if err != nil {
			fmt.Println("Error decoding message: " + err.Error())
			continue
		}
		messageQueue.PushMessage(message)
	}
}

// This function writes outgoing messages to the connection.
//
// This is a concurrent function - it runs simultaneously to the main game loop as a goroutine
func writeMessages(conn net.Conn, msgChan chan protocol.Message) {
	for {
		msg := <-msgChan
		conn.Write(msg.Encode())
	}
}
