package protocol

import (
	"sync"
)

// A thread-safe structure to hold messages.  It uses a mutex to ensure that if the
// messages are being read no new ones get added until the read has ended.  Additionally,
// it is a queue so a read operation is also a write operation as it removes the element.
type MessageQueue struct {
	lock     *sync.RWMutex
	messages []Message
}

// Pushes a message into the back of the queue
func (mq *MessageQueue) PushMessage(msg Message) {
	mq.lock.Lock()
	defer mq.lock.Unlock()
	mq.messages = append(mq.messages, msg)
}

// Remove the oldest message from the queue.  Returns nil on
// an empty queue
func (mq *MessageQueue) PopMessage() Message {
	mq.lock.Lock()
	defer mq.lock.Unlock()
	if len(mq.messages) > 0 {
		ev := mq.messages[0]
		mq.messages = mq.messages[1:]
		return ev
	}
	return nil
}

// Since the messages will be read in a loop, we don't want new messages being pushed in
// while reading the old ones (effectively making the loop continue to grow) so instead we
// pull all the messages out in a single operation behind the locked mutex and then clear
// the queue.
func (mq *MessageQueue) PopAll() []Message {
	mq.lock.Lock()
	defer mq.lock.Unlock()

	// copy the existing messages to a new slice
	msgs := make([]Message, len(mq.messages))
	for i, m := range mq.messages {
		msgs[i] = m
	}

	// replace the existing slice with an empty one
	mq.messages = make([]Message, 0)
	return msgs
}

// Creates a new MessageQueue object, inits it and returns it
func CreateMessageQueue() *MessageQueue {
	mq := new(MessageQueue)
	mq.lock = new(sync.RWMutex)
	mq.messages = make([]Message, 0)

	return mq
}
