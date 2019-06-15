package treesql

import (
	parserlib "github.com/vilterp/go-parserlib/pkg"
	"github.com/vilterp/go-parserlib/pkg/psi"
)

func MakeLanguage(schema *SchemaDesc) *psi.Language {
	return &psi.Language{
		Complete: func(n psi.Node, p parserlib.Position) []*psi.Completion {
			return Complete(n.(*Select), schema, p)
		},
		AnnotateErrors: func(n psi.Node) []*psi.ErrorAnnotation {
			return Annotate(schema, n.(*Select))
		},
		Extract: func(n *parserlib.Node) psi.Node {
			return ToSelect(n)
		},
		GetHighlightedElement: func(n psi.Node, pos parserlib.Position) *psi.HighlightedElement {
			return GetHighlightedElement(n, pos)
		},
		Grammar: Grammar,
	}
}
