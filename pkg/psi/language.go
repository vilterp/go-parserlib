package psi

import parserlib "github.com/vilterp/go-parserlib/pkg"

type Language struct {
	Grammar        *parserlib.Grammar
	AnnotateErrors func(n Node) []*ErrorAnnotation
	Complete       func(n Node, p parserlib.Position) []*Completion
}
