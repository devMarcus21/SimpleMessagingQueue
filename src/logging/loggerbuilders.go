package logging

import (
	"log/slog"
	"os"
)

type LoggerBuilder func(LogEventType) *slog.Logger

func addLogAttributes(logger *slog.Logger, logEventType LogEventType) *slog.Logger {
	return logger.With(
		LogEventTypeName, logEventType)
}

func BuildTextLogger(logEventType LogEventType) *slog.Logger {
	return addLogAttributes(slog.New(slog.NewTextHandler(os.Stdout, nil)), logEventType)
}

func BuildJsonLogger(logEventType LogEventType) *slog.Logger {
	return addLogAttributes(slog.New(slog.NewJSONHandler(os.Stdout, nil)), logEventType)
}
