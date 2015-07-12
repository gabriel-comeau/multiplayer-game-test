package shared

const (

	// *******************************************
	//                                           *
	// Network configuration                     *
	//                                           *
	// *******************************************

	// Host/interface
	HOST string = "localhost"

	// Port
	PORT string = "1339"

	// *******************************************
	//                                           *
	// Paths to folders and files needed         *
	//                                           *
	// *******************************************

	// Path to the root of the textures folder
	TEXTURE_ROOT = "./data/images/"

	// *******************************************
	//                                           *
	// Velocities and speeds                     *
	//                                           *
	// *******************************************

	// How fast (pixels per second)
	SPEED float32 = 300

	// Max allowable difference between the reported frame delta in a message and the
	// time between when that message was receieved and the previous one was.
	MAX_DT_DIFF_MILLIS = 8

	// Divide a number of nanoseconds by this number to get the value in millis
	NANO_TO_MILLI = 1000000
)
