package treesql

import parserlib "github.com/vilterp/go-parserlib/pkg"

type Completion struct {
	// TODO(vilterp): how does this represent completion partway into a token?
	Pos  parserlib.Position
	Text string
}

func Complete(t *parserlib.Node, cursorPos int) []*Completion {
	return nil
}
