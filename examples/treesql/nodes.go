package treesql

import (
	"fmt"

	parserlib "github.com/vilterp/go-parserlib/pkg"
	"github.com/vilterp/go-parserlib/pkg/psi"
)

// can do completion on this representation
// and annotate it with errors
type Select struct {
	Many       bool
	TableName  *parserlib.TextNode
	Selections []*Selection
}

var _ psi.Node = &Select{}

func (*Select) TypeName() string {
	return "Select"
}

func (s *Select) Children() []psi.Node {
	// ugh wai go
	var out []psi.Node
	for _, sel := range s.Selections {
		out = append(out, sel)
	}
	return out
}

func (s *Select) Attributes() map[string]string {
	return map[string]string{
		"many": fmt.Sprintf("%v", s.Many),
	}
}

func (s *Select) AttrNodes() map[string]*parserlib.TextNode {
	return map[string]*parserlib.TextNode{
		"table_name": s.TableName,
	}
}

type Selection struct {
	Name      *parserlib.TextNode
	SubSelect *Select
}

var _ psi.Node = &Selection{}

func (s *Selection) TypeName() string {
	return "Selection"
}

func (s *Selection) Children() []psi.Node {
	if s.SubSelect == nil {
		return nil
	}
	return []psi.Node{s.SubSelect}
}

func (s *Selection) Attributes() map[string]string {
	return nil
}

func (s *Selection) AttrNodes() map[string]*parserlib.TextNode {
	return map[string]*parserlib.TextNode{
		"name": s.Name,
	}
}
