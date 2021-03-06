package main

import (
	"bufio"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/gabriel-comeau/multiplayer-game-test/protocol"
	"github.com/gabriel-comeau/multiplayer-game-test/shared"
)

type TestPlayer struct {
	conn     net.Conn
	lastSeq  int64
	playerId int64
}

const (
	NUM_CLIENTS               = 1
	SLEEP_TIME  time.Duration = 33 * time.Millisecond
	COUNTER_MAX int           = 5
)

func main() {
	for i := 0; i < NUM_CLIENTS; i++ {
		log.Print("Launching client:", i)
		launchClient()
	}

	for {
		time.Sleep(SLEEP_TIME)
	}
}

func launchClient() {
	testPlayer := new(TestPlayer)
	testPlayer.conn, testPlayer.playerId = connectToServer()
	log.Printf("Got a uuid of: %v\n", testPlayer.playerId)
	go listenForMessages(testPlayer)
	go runTestPlayer(testPlayer)

}

// This is the "main" game loop for each test player
func runTestPlayer(testPlayer *TestPlayer) {

	// Preset up the timestep stuff so there's a value for the first iteration of the loop
	lastTick := time.Now()
	var dt time.Duration = 0
	var inputState *shared.InputState
	var counter int = 0

	for {
		// Generate a new random input state if the counter is zero.  We don't want to send
		// a new input state each tick or the movement is too wacky.
		if counter == 0 {
			inputState = generateRandomInputState()
		}

		msg := protocol.CreateSendInputMessage(inputState, testPlayer.lastSeq, dt, testPlayer.playerId)
		testPlayer.conn.Write(msg.Encode())
		testPlayer.lastSeq++

		log.Printf("Sending message from client %v: %+v\n", testPlayer.playerId, msg)

		counter++
		if counter > COUNTER_MAX {
			counter = 0
		}

		// Get how long it took to do all of this
		now := time.Now()
		dt = now.Sub(lastTick)
		lastTick = now

		if dt < SLEEP_TIME {
			time.Sleep(SLEEP_TIME - dt)
		}

		log.Printf("Client %v made it out of sleep loop\n", testPlayer.playerId)
	}
}

// Get a random input state
func generateRandomInputState() *shared.InputState {
	is := new(shared.InputState)
	is.KeyUpDown = coinToss()
	is.KeyDownDown = coinToss()
	is.KeyLeftDown = coinToss()
	is.KeyRightDown = coinToss()

	return is
}

// Do a random "coin toss", returning either true or false randomly
func coinToss() bool {
	rand.Seed(time.Now().UnixNano())
	result := rand.Intn(2)

	if result == 0 {
		return false
	}

	return true
}

// Establish a connection to the game server, return the network connection and the uuid
func connectToServer() (net.Conn, int64) {
	var playerId int64
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
			log.Print("Error while trying to accept player id")
			os.Exit(1)
			break
		}

		if string(line) == "" || string(line) == "\n" {
			continue
		}

		message, err := protocol.DecodeMessage(line)
		if err != nil {
			log.Print("ERROR during decode: " + err.Error())
			continue
		}

		if message.GetMessageType() == protocol.PLAYER_UUID_MESSAGE {
			typed, ok := message.(*protocol.PlayerUUIDMessage)
			if !ok {
				log.Print("Message couldn't be asserted into PlayerUUIDMessage though that was message id")
				conn.Close()
				os.Exit(1)
			}

			playerId = typed.UUID
			break
		} else {
			log.Print("Got the wrong type of message - expected PLAYER_UUID_MESSAGE")
			conn.Close()
			os.Exit(1)
		}
	}

	return conn, playerId
}

// Listens for incoming messages from the server, decodes the serialized versions into
// message objects and then pushes them into the message queue.
//
// This is a concurrent function - it runs simultaneously to the main game loop as a goroutine
func listenForMessages(testPlayer *TestPlayer) {
	b := bufio.NewReader(testPlayer.conn)
	for {
		line, err := b.ReadBytes('\n')

		if err != nil {
			testPlayer.conn.Close()
			log.Print("ERROR, CLOSING CONN: " + err.Error())
			break
		}

		if string(line) == "" || string(line) == "\n" {
			continue
		}

		// Deal with incoming messages from the server
		message, err := protocol.DecodeMessage(line)
		if err != nil {
			log.Print("Error decoding message: " + err.Error())
			continue
		}

		// We don't really care about the messages right now, just print it out
		log.Printf("Client: %v recieved world state message: %v\n", testPlayer.playerId, message)
	}
}
