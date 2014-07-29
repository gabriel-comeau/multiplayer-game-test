package main

import (
	sf "bitbucket.org/krepa098/gosfml2"
)

// A simpler, server-side version of the client-side Unit structure.  This one doesn't worry
// about textures but does keep track of the owning player's UUID, position and the last acknowledged
// sequence number.
type PlayerEntity struct {
	entityId int64
	position sf.Vector2f
	lastSeq  int64
}

// Move the entity by a given offset.
func (this *PlayerEntity) Move(offset sf.Vector2f) {
	this.position.X += offset.X
	this.position.Y += offset.Y
}

// Create a new entity and return a pointer to it.
func CreatePlayerEntity(id int64, initialPos sf.Vector2f) *PlayerEntity {
	ent := new(PlayerEntity)
	ent.position = initialPos
	ent.entityId = id
	ent.lastSeq = 0
	return ent
}
