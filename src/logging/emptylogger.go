package logging

import (
	"github.com/google/uuid"
)

// Temporary empty logger until a real solution is found
// For now this is scaffolding
type EmptyLogger struct{}

func NewEmptyLogger() Logger {
	return &EmptyLogger{}
}

func (EmptyLogger) Log(uuid.UUID, LogName, any) {}
