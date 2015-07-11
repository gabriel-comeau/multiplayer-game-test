package main

import (
	"sync"
)

// A simple thread safe structure for generating unique IDs for players as they join.
type IdGenerator struct {
	lock   *sync.Mutex
	nextId int64
}

// The main method for this struct - gets the next available ID.  Locks the mutex while
// this is happening so anyone else calling at the same time will have to wait.  No one will
// end up getting the same ID though.
func (idg *IdGenerator) GetNextId() int64 {
	idg.lock.Lock()
	defer idg.lock.Unlock()
	ret := idg.nextId
	idg.nextId++
	return ret
}

// Constructor which initializes the ID generator and returns a pointer to it.
func CreateIdGenerator() *IdGenerator {
	return &IdGenerator{
		lock:   new(sync.Mutex),
		nextId: 1,
	}
}
