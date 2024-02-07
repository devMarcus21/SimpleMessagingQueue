package logging

import (
	"github.com/google/uuid"
)

type LogName int64

const (
	MessagePushedToQueueService LogName = iota
	MessagePulledFromQueueService
)

type Logger interface {
	// Each log context will be mapped to a unqiue message id
	// <UUID> Id of the log context <LogName> Name of the log row <any> value of the log row
	// All log rows will be stored under the same log context id along with the timestamp
	Log(uuid.UUID, LogName, any)
}
