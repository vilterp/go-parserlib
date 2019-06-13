package psi

import (
	"sort"

	parserlib "github.com/vilterp/go-parserlib/pkg"
	pp "github.com/vilterp/go-pretty-print"
)

type Node interface {
	TypeName() string
	Children() []Node
	// general attributes
	Attributes() map[string]string
	// attributes whose values are text nodes
	AttrNodes() map[string]*parserlib.TextNode
}

func Format(n Node) pp.Doc {
	// TODO(vilterp): ugh, so much code. Go, sorted map plz.

	var attrDocs []pp.Doc

	// sort attr names
	attrs := n.Attributes()
	var attrNames []string
	for name := range attrs {
		attrNames = append(attrNames, name)
	}
	sort.Strings(attrNames)

	// sort attr node names
	attrNodes := n.AttrNodes()
	var attrNodeNames []string
	for name := range attrNodes {
		attrNodeNames = append(attrNodeNames, name)
	}
	sort.Strings(attrNodeNames)

	// format attrs
	for _, attrName := range attrNames {
		attrDocs = append(attrDocs, pp.Textf("%v: %v", attrName, attrs[attrName]))
	}

	// format attr nodes
	for _, attrName := range attrNodeNames {
		attrDocs = append(attrDocs, pp.Textf("%v: %v", attrName, attrNodes[attrName].String()))
	}

	var childDocs []pp.Doc
	for _, child := range n.Children() {
		childDocs = append(childDocs, Format(child))
	}

	header := pp.Seq([]pp.Doc{
		pp.Text(n.TypeName()),
		pp.Text(" <"),
		pp.Join(attrDocs, pp.CommaSpace),
		pp.Text(">"),
	})

	if len(childDocs) == 0 {
		return header
	}
	return pp.Seq([]pp.Doc{
		header,
		pp.Newline,
		pp.Indent(2, pp.Join(childDocs, pp.Newline)),
	})
}
