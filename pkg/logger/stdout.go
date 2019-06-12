package logger

import (
	"fmt"
	"strings"
)

type StdoutLogger struct {
	indentation int
}

var _ Logger = &StdoutLogger{}

func NewStdoutLogger() *StdoutLogger {
	return &StdoutLogger{
		indentation: 0,
	}
}

func (l *StdoutLogger) Log(line ...interface{}) {
	fmt.Print(strings.Repeat("  ", l.indentation))
	fmt.Println(line...)
}

func (l *StdoutLogger) Logf(format string, things ...interface{}) {
	fmt.Print(strings.Repeat("  ", l.indentation))
	fmt.Printf(format, things...)
	fmt.Println()
}

func (l *StdoutLogger) Indent() {
	l.indentation += 1
}

func (l *StdoutLogger) Outdent() {
	if l.indentation == 0 {
		panic("indent already 0")
	}
	l.indentation -= 1
}
