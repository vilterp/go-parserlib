package logger

type Logger interface {
	Log(line ...interface{})
	Logf(format string, things ...interface{})
	Indent()
	Outdent()
}
