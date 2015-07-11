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
func (eh *EntityHolder) AddEntity(entity *PlayerEntity) {
	eh.lock.Lock()
	defer eh.lock.Unlock()
	eh.entities[entity.entityId] = entity
}

// Remove a entity from the map
func (eh *EntityHolder) RemoveEntity(id int64) {
	eh.lock.Lock()
	defer eh.lock.Unlock()
	delete(eh.entities, id)
}

// Gets a specific entity, by ID, out of the map.  Returns nil if entity is not available.
//
// TODO: to make this more of an idiomatic Go call, change this to return entity, err like regular
// maps do.
func (eh *EntityHolder) GetEntity(id int64) *PlayerEntity {
	eh.lock.RLock()
	defer eh.lock.RUnlock()
	e, ok := eh.entities[id]
	if ok {
		return e
	} else {
		return nil
	}
}

// Get all of the entities as a slice.
func (eh *EntityHolder) GetEntities() []*PlayerEntity {
	eh.lock.RLock()
	defer eh.lock.RUnlock()
	eSlice := make([]*PlayerEntity, 0)
	for _, entity := range eh.entities {
		eSlice = append(eSlice, entity)
	}

	return eSlice
}

// Constructor to init the holder
func CreateEntityHolder() *EntityHolder {
	return &EntityHolder{
		lock:     new(sync.RWMutex),
		entities: make(map[int64]*PlayerEntity),
	}
}
