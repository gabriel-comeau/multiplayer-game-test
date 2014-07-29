package main

import (
	"sync"
)

// Basically a wrapper around a map of Client structs to make it thread safe.
type ClientHolder struct {
	lock    *sync.RWMutex
	clients map[int64]*Client
}

// Add a new client to the map
func (this *ClientHolder) AddClient(client *Client) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.clients[client.clientId] = client
}

// Remove a client from the map
func (this *ClientHolder) RemoveClient(client *Client) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.clients, client.clientId)
}

// Returns how many elements are in the map (at call time - thread safety!)
func (this *ClientHolder) Count() int {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return len(this.clients)
}

// Gets a specific client, by ID, out of the map.  Returns nil if client is not available.
//
// TODO: to make this more of an idiomatic Go call, change this to return client, err like regular
// maps do.
func (this *ClientHolder) GetClient(id int64) *Client {
	this.lock.RLock()
	defer this.lock.RUnlock()
	c, ok := this.clients[id]
	if ok {
		return c
	} else {
		return nil
	}
}

// Get all of the clients as a slice.
func (this *ClientHolder) GetClients() []*Client {
	this.lock.RLock()
	defer this.lock.RUnlock()
	cSlice := make([]*Client, 0)
	for _, client := range this.clients {
		cSlice = append(cSlice, client)
	}

	return cSlice
}

// Constructor to init the holder
func CreateClientHolder() *ClientHolder {
	lck := new(sync.RWMutex)
	holder := new(ClientHolder)
	holder.lock = lck
	holder.clients = make(map[int64]*Client)

	return holder
}
