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
	Offer(QueueMessage)
	Poll() (QueueMessage, bool)
	Size() int
	IsEmpty() bool
}
