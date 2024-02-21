package logging

type LogName int64

const (
	LogIota string = "LogIota"

	APIPush_MessagePushedToQueueService LogName = iota
	APIPop_MessagePulledFromQueueService
	APIPop_QueueIsEmptyNoMessagePulled

	// Batch endpoints
	APIPushBatch_BatchSize
)

// Map is more verbose than using an array/slice
var logNameToString = map[LogName]string{
	APIPush_MessagePushedToQueueService:  "APIPush_MessagePushedToQueueService",
	APIPop_MessagePulledFromQueueService: "APIPop_MessagePulledFromQueueService",
	APIPop_QueueIsEmptyNoMessagePulled:   "APIPop_QueueIsEmptyNoMessagePulled",

	// Batch endpoints
	APIPushBatch_BatchSize: "APIPushBatch_BatchSize",
}

var logNameToMessageString = map[LogName]string{
	APIPush_MessagePushedToQueueService:  "API Push: Message pushed to queue service",
	APIPop_MessagePulledFromQueueService: "API Pop: Message pulled from queue service",
	APIPop_QueueIsEmptyNoMessagePulled:   "API Pop: Queue is empty",

	// Batch endpoints
	APIPushBatch_BatchSize: "API Push Batch: batch size - %d",
}

func (logName LogName) String() string {
	return logNameToString[logName]
}

func (logName LogName) Message() string {
	return logNameToMessageString[logName]
}
