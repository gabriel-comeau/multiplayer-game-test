package main

import (
	"github.com/gabriel-comeau/multiplayer-game-test/shared"
)

// A simpler, server-side version of the client-side Unit structure.  This one doesn't worry
// about textures but does keep track of the owning player's UUID, position and the last acknowledged
// sequence number.
type PlayerEntity struct {
	entityId int64
	position shared.FloatVector
	lastSeq  int64
}

// Move the entity by a given offset.
func (this *PlayerEntity) Move(offset shared.FloatVector) {
	this.position.X += offset.X
	this.position.Y += offset.Y
}

// Create a new entity and return a pointer to it.
func CreatePlayerEntity(id int64, initialPos shared.FloatVector) *PlayerEntity {
	ent := new(PlayerEntity)
	ent.position = initialPos
	ent.entityId = id
	ent.lastSeq = 0
	return ent
}
