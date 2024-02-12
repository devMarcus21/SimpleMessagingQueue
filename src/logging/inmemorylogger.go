package logging

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type InMemoryLogger struct {
	logContexts map[uuid.UUID][]LogMessage
	muLock      *sync.Mutex
}

func NewInMemoryLogger() Logger {
	return &InMemoryLogger{
		logContexts: map[uuid.UUID][]LogMessage{},
		muLock:      new(sync.Mutex),
	}
}

func (logger *InMemoryLogger) Log(logContextId uuid.UUID, log LogName, logValue any) {
	logger.muLock.Lock()
	defer logger.muLock.Unlock()

	now := time.Now().Unix()
	message := LogMessage{
		LogContextId: logContextId,
		Timestamp:    now,
		Log:          log,
		LogValue:     logValue,
	}

	_, idExists := logger.logContexts[logContextId]
	if !idExists {
		logger.logContexts[logContextId] = []LogMessage{}
	}

	logger.logContexts[logContextId] = append(logger.logContexts[logContextId], message)
}

func (logger *InMemoryLogger) GetLogs(logContextId uuid.UUID) []LogMessage {
	logger.muLock.Lock()
	defer logger.muLock.Unlock()

	return logger.logContexts[logContextId]
}
