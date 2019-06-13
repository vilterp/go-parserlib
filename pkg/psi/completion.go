package psi

import parserlib "github.com/vilterp/go-parserlib/pkg"

type Completion struct {
	// TODO(vilterp): how does this represent completion partway into a token?
	Pos  parserlib.Position
	Text string
}

// top to bottom
func GetPath(node Node, pos parserlib.Position) []Node {
	return nil
}
