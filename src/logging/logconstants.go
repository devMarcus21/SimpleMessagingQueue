package logging

type LogName int64

const (
	MessagePushedToQueueService LogName = iota
	MessagePulledFromQueueService
)

// Map is more verbose than using an array/slice
var logNameToString = map[LogName]string{
	MessagePushedToQueueService:   "MessagePushedToQueueService",
	MessagePulledFromQueueService: "MessagePulledFromQueueService",
}

var logNameToMessageString = map[LogName]string{
	MessagePushedToQueueService:   "Message pushed to queue service",
	MessagePulledFromQueueService: "Message pulled from queue service",
}

func (logName LogName) String() string {
	return logNameToString[logName]
}

func (logName LogName) Message() string {
	return logNameToMessageString[logName]
}
