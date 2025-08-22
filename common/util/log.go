package util

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

func InitLog(name string, level slog.Level) (*slog.Logger, error) {
	file, err := os.OpenFile("./log/test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	l := slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}))
	return l.With("ServiceName", name), nil
}

func SetTrace(ctx context.Context, logger *slog.Logger) *slog.Logger {
	span := trace.SpanFromContext(ctx)
	return logger.With("TraceId", span.SpanContext().TraceID().String()).
		With("SpanId", span.SpanContext().SpanID().String()).
		WithGroup("detail")
}
