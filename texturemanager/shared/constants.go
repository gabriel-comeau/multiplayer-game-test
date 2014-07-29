package shared

import (
	"time"
)

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

	// Maximum allowable DT per move
	MAX_DT time.Duration = time.Second / 20
)
