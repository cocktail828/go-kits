package slogx

import (
	"golang.org/x/exp/slog"
)

// slogWrapper implements Logger interface
type slogWrapper struct {
	l *slog.Logger
}

// NewLoggerWithSlog creates a new logger which wraps
// the given logrus.Logger
func NewLoggerWithSlog(logger *slog.Logger) Logger {
	return &slogWrapper{
		l: logger,
	}
}

func (c slogWrapper) Debug(msg string, args ...any) {
	c.l.Debug(msg, args...)
}

func (c slogWrapper) Info(msg string, args ...any) {
	c.l.Info(msg, args...)
}

func (c slogWrapper) Warn(msg string, args ...any) {
	c.l.Warn(msg, args...)
}

func (c slogWrapper) Error(msg string, args ...any) {
	c.l.Error(msg, args...)
}

func (c slogWrapper) With(args ...any) Logger {
	return NewLoggerWithSlog(c.l.With(args...))
}

func (c slogWrapper) WithGroup(name string) Logger {
	return NewLoggerWithSlog(c.l.WithGroup(name))
}
