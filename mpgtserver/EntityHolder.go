package main

import (
	"sync"
)

// Basically a wrapper around a map of PlayerEntity structs to make it thread safe.
type EntityHolder struct {
	lock     *sync.RWMutex
	entities map[int64]*PlayerEntity
}

// Add a new entity to the map
func (this *EntityHolder) AddEntity(entity *PlayerEntity) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.entities[entity.entityId] = entity
}

// Remove a entity from the map
func (this *EntityHolder) RemoveEntity(id int64) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.entities, id)
}

// Gets a specific entity, by ID, out of the map.  Returns nil if entity is not available.
//
// TODO: to make this more of an idiomatic Go call, change this to return entity, err like regular
// maps do.
func (this *EntityHolder) GetEntity(id int64) *PlayerEntity {
	this.lock.RLock()
	defer this.lock.RUnlock()
	e, ok := this.entities[id]
	if ok {
		return e
	} else {
		return nil
	}
}

// Get all of the entities as a slice.
func (this *EntityHolder) GetEntities() []*PlayerEntity {
	this.lock.RLock()
	defer this.lock.RUnlock()
	eSlice := make([]*PlayerEntity, 0)
	for _, entity := range this.entities {
		eSlice = append(eSlice, entity)
	}

	return eSlice
}

// Constructor to init the holder
func CreateEntityHolder() *EntityHolder {
	lck := new(sync.RWMutex)
	holder := new(EntityHolder)
	holder.lock = lck
	holder.entities = make(map[int64]*PlayerEntity)

	return holder
}
