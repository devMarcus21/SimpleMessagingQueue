package queue

import (
	"github.com/google/uuid"
)

type QueueMessage struct {
	MessageId uuid.UUID
	Timestamp int64
	Data      []byte
}

type Queue interface {
	// Adds a message to the back of the queue data structure
	Offer(QueueMessage)

	// Removes and returns the message at the front of the queue
	Poll() (QueueMessage, bool)

	// Returns the number of messages currently in the queue
	Size() int

	// Returns true if the queue has no messages in it and false if not
	IsEmpty() bool
}
