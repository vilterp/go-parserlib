package psi

import (
	"fmt"

	parserlib "github.com/vilterp/go-parserlib/pkg"
)

type ErrorAnnotation struct {
	Span    parserlib.SourceSpan
	Message string
}

func (e *ErrorAnnotation) String() string {
	return fmt.Sprintf("%v: %v", e.Span.String(), e.Message)
}
