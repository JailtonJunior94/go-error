package o11y

import (
	"context"
	"log/slog"
)

type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
}

type logger struct {
	log *slog.Logger
}

func NewLogger(log *slog.Logger) Logger {
	return &logger{
		log: log,
	}
}

func (l *logger) Debug(ctx context.Context, msg string, args ...any) {
	l.log.DebugContext(ctx, msg, args...)
}

func (l *logger) Info(ctx context.Context, msg string, args ...any) {
	l.log.InfoContext(ctx, msg, args...)
}

func (l *logger) Warn(ctx context.Context, msg string, args ...any) {
	l.log.WarnContext(ctx, msg, args...)
}

func (l *logger) Error(ctx context.Context, msg string, args ...any) {
	l.log.ErrorContext(ctx, msg, args...)
}
