package logging

type LogEventType string
type LogEvent int64
type LogProperty string

const (
	LogEventIota string = "LogEventIota"

	LogEventTypeName string       = "LogEventTypeName"
	Request          LogEventType = "Request"
	Service          LogEventType = "Service"

	HandlerActionName LogEvent = iota
	JsonDecodeError

	// Push Events
	APIPush
	MessagePushedToQueueService

	// Pop Events
	APIPop
	MessagePulledFromQueueService
	QueueIsEmptyNoMessagePulled

	// Batching
	BatchMessageProperties

	// Batch Push Events
	APIPushBatch
	PushBatchBatchSize

	// Log Properties
	NewMessageId    LogProperty = "NewMessageId"
	PulledMessageId LogProperty = "PulledMessageId"

	IsBatched    LogProperty = "IsBatched"
	BatchedIndex LogProperty = "BatchedIndex"
)

// Map is more verbose than using an array/slice
var logNameToString = map[LogEvent]string{
	HandlerActionName: "HandlerActionName",
	JsonDecodeError:   "JsonDecodeError",

	// Push Events
	APIPush:                     "APIPush",
	MessagePushedToQueueService: "MessagePushedToQueueService",

	// Pop Events
	APIPop:                        "APIPop",
	MessagePulledFromQueueService: "MessagePulledFromQueueService",
	QueueIsEmptyNoMessagePulled:   "QueueIsEmptyNoMessagePulled",

	// Batching
	BatchMessageProperties: "BatchMessageProperties",

	// Batch Push Events
	APIPushBatch:       "APIPushBatch",
	PushBatchBatchSize: "PushBatchBatchSize",
}

var logNameToMessageString = map[LogEvent]string{
	JsonDecodeError: "JsonDecodeError: %s",

	// Push Events
	MessagePushedToQueueService: "Message pushed to queue service",

	// Pop Events
	MessagePulledFromQueueService: "Message pulled from queue service",
	QueueIsEmptyNoMessagePulled:   "Queue is empty",

	// Batching
	BatchMessageProperties: "Batched message property data",

	// Batch Push Events
	PushBatchBatchSize: "Push Batch: batch size - %d",
}

func (logName LogEvent) String() string {
	return logNameToString[logName]
}

func (logName LogEvent) Message() string {
	return logNameToMessageString[logName]
}
