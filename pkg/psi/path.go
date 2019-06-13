package psi

import (
	parserlib "github.com/vilterp/go-parserlib/pkg"
)

type Path struct {
	Nodes    []Node
	NodeName string
	AttrName string
	AttrText *parserlib.TextNode
}

// top to bottom
func GetPath(node Node, pos parserlib.Position) *Path {
	return getPathRecurse(node, nil, pos)
}

func getPathRecurse(node Node, nodes []Node, pos parserlib.Position) *Path {
	if !node.SourceSpan().Contains(pos) {
		return nil
	}

	attrs := node.AttrNodes()
	for name, attr := range attrs {
		if attr.Span.Contains(pos) {
			return &Path{
				NodeName: node.TypeName(),
				AttrName: name,
				AttrText: attr,
				Nodes:    nodes,
			}
		}
	}

	children := node.Children()
	for _, child := range children {
		if child.SourceSpan().Contains(pos) {
			nodes := append(nodes, node)
			res := getPathRecurse(child, nodes, pos)
			if res != nil {
				return res
			}
			nodes = nodes[:len(nodes)-1]
		}
	}

	return nil
}
