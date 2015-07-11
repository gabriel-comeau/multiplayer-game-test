package main

import (
	"time"

	"github.com/gabriel-comeau/multiplayer-game-test/shared"
)

// A simpler, server-side version of the client-side Unit structure.  This one doesn't worry
// about textures but does keep track of the owning player's UUID, position and the last acknowledged
// sequence number.
type PlayerEntity struct {
	entityId    int64
	position    shared.FloatVector
	lastSeq     int64
	lastSeqTime time.Time
}

// Move the entity by a given offset.
func (p *PlayerEntity) Move(offset shared.FloatVector) {
	p.position.X += offset.X
	p.position.Y += offset.Y
}

func (p *PlayerEntity) GetTimeOffset(current time.Time) time.Duration {
	return p.lastSeqTime.Sub(current)
}

// Create a new entity and return a pointer to it.
func CreatePlayerEntity(id int64, initialPos shared.FloatVector) *PlayerEntity {
	return &PlayerEntity{
		position: initialPos,
		entityId: id,
		lastSeq:  0,
	}
}
