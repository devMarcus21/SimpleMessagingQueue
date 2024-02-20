package queue

import (
	"github.com/google/uuid"
)

type QueueMessage struct {
	MessageId          uuid.UUID
	ProducerIdentifier string
	Timestamp          int64
	Data               map[string]any
}

type Queue interface {
	// Adds a message to the back of the queue data structure
	Offer(QueueMessage)

	// Removes and returns the message at the front of the queue
	// Returns true if a message was returned and false if not
	Poll() (QueueMessage, bool)

	// Returns the number of messages currently in the queue
	Size() int

	// Returns true if the queue has no messages in it and false if not
	IsEmpty() bool
}
