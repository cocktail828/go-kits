package slogx

// DefaultNopLogger returns a nop logger.
func DefaultNopLogger() Logger {
	return nopLogger{}
}

type nopLogger struct{}

func (c nopLogger) Debug(msg string, args ...any) {}
func (c nopLogger) Info(msg string, args ...any)  {}
func (c nopLogger) Warn(msg string, args ...any)  {}
func (c nopLogger) Error(msg string, args ...any) {}
func (c nopLogger) With(args ...any) Logger       { return c }
func (c nopLogger) WithGroup(name string) Logger  { return c }
