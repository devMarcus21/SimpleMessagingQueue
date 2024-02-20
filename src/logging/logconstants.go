package logging

type LogName int64

const (
	MessagePushedToQueueService LogName = iota
	MessagePulledFromQueueService
)

var logNameToString = [...]string{
	"MessagePushedToQueueService",
	"MessagePulledFromQueueService",
}

var logNameToMessageString = [...]string{
	"Message pushed to queue service",
	"Message pulled from queue service",
}

func (logName LogName) String() string {
	return logNameToString[logName]
}

func (logName LogName) Message() string {
	return logNameToMessageString[logName]
}
