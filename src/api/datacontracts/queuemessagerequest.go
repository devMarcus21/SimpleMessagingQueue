package datacontracts

import (
	queueUtils "github.com/devMarcus21/SimpleMessagingQueue/src/datastructures/queue"

	"github.com/google/uuid"
)

type QueueMessageRequest struct {
	Data []byte `json:"data"`
}

func ConvertQueueMessageRequestToQueueMessage(message QueueMessageRequest, id uuid.UUID, createdTimestamp int64) queueUtils.QueueMessage {
	return queueUtils.QueueMessage{
		MessageId: id,
		Timestamp: createdTimestamp,
		Data:      message.Data,
	}
}