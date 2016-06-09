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
func (ch *ClientHolder) AddClient(client *Client) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	ch.clients[client.clientId] = client
}

// Remove a client from the map
func (ch *ClientHolder) RemoveClient(client *Client) {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	delete(ch.clients, client.clientId)
}

// Gets a specific client, by ID, out of the map.  Returns nil if client is not available.
//
// TODO: to make this more of an idiomatic Go call, change this to return client, err like regular
// maps do.
func (ch *ClientHolder) GetClient(id int64) *Client {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	c, ok := ch.clients[id]
	if ok {
		return c
	} else {
		return nil
	}
}

// Get all of the clients as a slice.
func (ch *ClientHolder) GetClients() []*Client {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	cSlice := make([]*Client, 0)
	for _, client := range ch.clients {
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
