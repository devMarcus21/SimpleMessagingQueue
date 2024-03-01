package logging

type LogEventType string
type LogEvent int64
type LogProperty string
type HandlerAction string

const (
	LogEventIota string = "LogEventIota"

	LogEventTypeName string       = "LogEventTypeName"
	Request          LogEventType = "Request"
	Service          LogEventType = "Service"
)

// Log Events
const (
	JsonDecodeError LogEvent = iota

	// Push Events
	MessagePushedToQueueService

	// Pop Events
	MessagePulledFromQueueService
	QueueIsEmptyNoMessagePulled

	// Batching
	BatchMessageProperties

	// Batch Push Events
	PushBatchBatchSize
)

// Log Properties
const (
	HandlerActionName LogProperty = "HandlerActionName"
	NewMessageId      LogProperty = "NewMessageId"
	PulledMessageId   LogProperty = "PulledMessageId"
	RequestId         LogProperty = "RequestId"

	IsBatched    LogProperty = "IsBatched"
	BatchedIndex LogProperty = "BatchedIndex"
)

// Handler Actions
const (
	APIPush      HandlerAction = "APIPush"
	APIPop       HandlerAction = "APIPop"
	APIPushBatch HandlerAction = "APIPushBatch"
)

// Map is more verbose than using an array/slice
var logNameToString = map[LogEvent]string{
	JsonDecodeError: "JsonDecodeError",

	// Push Events
	MessagePushedToQueueService: "MessagePushedToQueueService",

	// Pop Events
	MessagePulledFromQueueService: "MessagePulledFromQueueService",
	QueueIsEmptyNoMessagePulled:   "QueueIsEmptyNoMessagePulled",

	// Batching
	BatchMessageProperties: "BatchMessageProperties",

	// Batch Push Events
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
