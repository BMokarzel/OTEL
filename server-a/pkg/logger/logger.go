package logger

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

type Logger struct {
	logger *zerolog.Logger
}

func New(build, applicationName string) *Logger {
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("build", build).
		Str("application-name", applicationName).
		Logger()
	return &Logger{
		logger: &logger,
	}
}

func (l *Logger) log(ctx context.Context, logFunc *zerolog.Event) *zerolog.Event {
	span := trace.SpanFromContext(ctx)
	traceId := span.SpanContext().TraceID().String()
	logFunc = logFunc.Str("trace-id", traceId)
	return logFunc
}

type Event struct {
	event *zerolog.Event
}

func (e *Event) Msg(message string, args ...any) {
	e.event.Msgf(message, args...)
}

func (l *Logger) Debug(ctx context.Context) *Event {
	return &Event{l.log(ctx, l.logger.Debug())}
}

func (l *Logger) Info(ctx context.Context) *Event {
	return &Event{l.log(ctx, l.logger.Info())}
}

func (l *Logger) Warn(ctx context.Context) *Event {
	return &Event{l.log(ctx, l.logger.Warn())}
}

func (l *Logger) Error(ctx context.Context) *Event {
	return &Event{l.log(ctx, l.logger.Error())}
}
