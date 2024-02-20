package logging

import (
	"context"
	"log/slog"
)

type EmptyLoggerHandler struct{}

func NewEmptyLoggerHandler() slog.Handler {
	return &EmptyLoggerHandler{}
}

func (*EmptyLoggerHandler) Enabled(_ context.Context, level slog.Level) bool {
	return true
}

func (*EmptyLoggerHandler) Handle(ctx context.Context, r slog.Record) error {
	return nil
}

func (*EmptyLoggerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return nil
}

func (*EmptyLoggerHandler) WithGroup(name string) slog.Handler {
	return nil
}

func (*EmptyLoggerHandler) Handler() slog.Handler {
	return nil
}
