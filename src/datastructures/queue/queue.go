package queue

import (
	"github.com/google/uuid"
)

type BatchProperties struct {
	isBatchedMessage bool
	batchIndex       int // index within the batch
	batchSize        int
}

type RetryProperties struct {
	retries               int
	visibleAfterTimestamp int64 // time at which this message should be available for retry
}

type QueueMessage struct {
	MessageId          uuid.UUID
	ProducerIdentifier string
	Timestamp          int64
	Data               map[string]any
	Acknowledged       bool
	batchProperties    BatchProperties
	retryProperties    RetryProperties
}

func (message QueueMessage) IsBatchedMessage() (bool, int) {
	return message.batchProperties.isBatchedMessage, message.batchProperties.batchIndex
}

// Converts the given message into a message from a batch by adding batch properties to it
func (message *QueueMessage) MakeBatchedMessage(batchIndex int, batchSize int) bool {
	if message.batchProperties.isBatchedMessage {
		return false
	}
	message.batchProperties.isBatchedMessage = true
	message.batchProperties.batchIndex = batchIndex
	message.batchProperties.batchSize = batchSize

	return true
}

// Get number of previous retries of the message
func (message QueueMessage) RetryCount() int {
	return message.retryProperties.retries
}

func (message QueueMessage) NextRetryTimestamp() int64 {
	return message.retryProperties.visibleAfterTimestamp
}

// Indicates another retry attempt retryAvailableTimestamp will be the timestamp at which
// this message is next available for reconsumption again
func (message *QueueMessage) AddRetry(retryAvailableTimestamp int64) {
	message.retryProperties.visibleAfterTimestamp = retryAvailableTimestamp
	message.retryProperties.retries = message.retryProperties.retries + 1
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
