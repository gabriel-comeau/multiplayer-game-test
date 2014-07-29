package main

import (
	"sync"
)

// A simple thread safe structure for generating unique IDs for players as they join.
type IdGenerator struct {
	lock   *sync.RWMutex
	nextId int64
}

// The main method for this struct - gets the next available ID.  Locks the mutex while
// this is happening so anyone else calling at the same time will have to wait.  No one will
// end up getting the same ID though.
func (this *IdGenerator) GetNextId() int64 {
	this.lock.Lock()
	defer this.lock.Unlock()
	ret := this.nextId
	this.nextId++
	return ret
}

// Constructor which initializes the ID generator and returns a pointer to it.
func CreateIdGenerator() *IdGenerator {
	gen := new(IdGenerator)
	lck := new(sync.RWMutex)
	gen.lock = lck
	gen.nextId = 1
	return gen
}
