package datacontracts

import (
	queueUtils "github.com/devMarcus21/SimpleMessagingQueue/src/datastructures/queue"

	"github.com/google/uuid"
)

type QueueMessageRequest struct {
	ProducerIdentifier string         `json:"producerIdentifier"`
	Data               map[string]any `json:"data"`
}

func ConvertQueueMessageRequestToQueueMessage(message QueueMessageRequest, id uuid.UUID, createdTimestamp int64) queueUtils.QueueMessage {
	return queueUtils.QueueMessage{
		MessageId:          id,
		ProducerIdentifier: message.ProducerIdentifier,
		Timestamp:          createdTimestamp,
		Data:               message.Data,
	}
}

type BatchQueueMessageRequest struct {
	Messages []QueueMessageRequest `json:"messages"`
}
