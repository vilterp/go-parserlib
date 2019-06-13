package treesql

import (
	parserlib "github.com/vilterp/go-parserlib/pkg"
	"github.com/vilterp/go-parserlib/pkg/psi"
)

func Complete(t *parserlib.Node, cursorPos parserlib.Position) []*psi.Completion {
	if !t.Span.Contains(cursorPos) {
		return nil
	}
	return nil
}
