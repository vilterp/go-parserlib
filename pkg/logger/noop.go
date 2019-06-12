package logger

type NoopLogger struct{}

var _ Logger = &NoopLogger{}

func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

func (*NoopLogger) Log(line ...interface{}) {
}

func (*NoopLogger) Logf(format string, things ...interface{}) {
}

func (*NoopLogger) Indent() {
}

func (*NoopLogger) Outdent() {
}
