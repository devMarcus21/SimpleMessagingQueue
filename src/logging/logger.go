package logging

import (
	"github.com/google/uuid"
)

type LogName int64

const (
	MessagePushedToQueueService LogName = iota
	MessagePulledFromQueueService
)

type LogMessage struct {
	LogContextId uuid.UUID
	Timestamp    int64
	Log          LogName
	LogValue     any
}

type Logger interface {
	// Each log context will be mapped to a unqiue message id
	// <UUID> Id of the log context <LogName> Name of the log row <any> value of the log row
	// All log rows will be stored under the same log context id along with the timestamp
	Log(uuid.UUID, LogName, any)

	// Gets log data for a specific message id
	// Returns a list of all logs for the message id
	GetLogs(uuid.UUID) []LogMessage
}
