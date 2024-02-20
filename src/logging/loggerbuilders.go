package logging

import (
	"log/slog"
	"os"
)

type LoggerBuilder func() *slog.Logger

func BuildTextLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func BuildJsonLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
