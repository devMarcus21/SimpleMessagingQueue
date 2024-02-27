package queue

import (
	"github.com/google/uuid"
)

type BatchProperties struct {
	isBatchedMessage bool
	batchIndex       int // index within the batch
}

type QueueMessage struct {
	MessageId          uuid.UUID
	ProducerIdentifier string
	Timestamp          int64
	Data               map[string]any
	batchProperties    BatchProperties
}

func (message QueueMessage) IsBatchedMessage() (bool, int) {
	return message.batchProperties.isBatchedMessage, message.batchProperties.batchIndex
}

// Converts the given message into a message from a batch by adding batch properties to it
func (message *QueueMessage) MakeBatchedMessage(batchIndex int) bool {
	if message.batchProperties.isBatchedMessage {
		return false
	}
	message.batchProperties.isBatchedMessage = true
	message.batchProperties.batchIndex = batchIndex

	return true
}

type Queue interface {
	// Adds a message to the back of the queue data structure
	Offer(QueueMessage)

	OfferAll([]QueueMessage)

	AddFirst(QueueMessage)

	// Removes and returns the message at the front of the queue
	// Returns true if a message was returned and false if not
	Poll() (QueueMessage, bool)

	// Returns the number of messages currently in the queue
	Size() int

	// Returns true if the queue has no messages in it and false if not
	IsEmpty() bool
}
