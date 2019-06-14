package psi

import parserlib "github.com/vilterp/go-parserlib/pkg"

type Language struct {
	Grammar        *parserlib.Grammar
	Extract        func(n *parserlib.Node) Node
	AnnotateErrors func(n Node) []*ErrorAnnotation
	Complete       func(n Node, p parserlib.Position) []*Completion
}

func Parse(l *Language, input string) (Node, error) {
	traceTree, err := l.Grammar.Parse(l.Grammar.StartRule, input, 0, nil)
	if err != nil {
		return nil, err
	}
	ruleTree := traceTree.ToRuleTree()
	return l.Extract(ruleTree), nil
}
