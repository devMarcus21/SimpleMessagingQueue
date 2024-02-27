package logging

type LogName int64

const (
	LogIota string = "LogIota"

	HandlerActionName LogName = iota
	// API Push endpoints
	APIPush
	APIPush_MessagePushedToQueueService
	// API Pop endpoints
	APIPop
	APIPop_MessagePulledFromQueueService
	APIPop_QueueIsEmptyNoMessagePulled
	// Batching
	BatchMessageProperties
	// API Batch Push endpoints
	APIPPushBatch
	APIPushBatch_BatchSize
)

// Map is more verbose than using an array/slice
var logNameToString = map[LogName]string{
	HandlerActionName: "HandlerActionName",
	// API Push endpoints
	APIPush:                             "APIPush",
	APIPush_MessagePushedToQueueService: "APIPush_MessagePushedToQueueService",
	// API Pop endpoints
	APIPop:                               "APIPop",
	APIPop_MessagePulledFromQueueService: "APIPop_MessagePulledFromQueueService",
	APIPop_QueueIsEmptyNoMessagePulled:   "APIPop_QueueIsEmptyNoMessagePulled",
	// Batching
	BatchMessageProperties: "BatchMessageProperties",
	// API Batch Push endpoints
	APIPPushBatch:          "APIPPushBatch",
	APIPushBatch_BatchSize: "APIPushBatch_BatchSize",
}

var logNameToMessageString = map[LogName]string{
	// API Push endpoints
	APIPush_MessagePushedToQueueService: "API Push: Message pushed to queue service",
	// API Pop endpoints
	APIPop_MessagePulledFromQueueService: "API Pop: Message pulled from queue service",
	APIPop_QueueIsEmptyNoMessagePulled:   "API Pop: Queue is empty",
	// Batching
	BatchMessageProperties: "Batched message property data",
	// API Batch Push endpoints
	APIPushBatch_BatchSize: "API Push Batch: batch size - %d",
}

func (logName LogName) String() string {
	return logNameToString[logName]
}

func (logName LogName) Message() string {
	return logNameToMessageString[logName]
}
