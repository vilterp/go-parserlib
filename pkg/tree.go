package parserlib

import (
	"fmt"

	pp "github.com/vilterp/go-pretty-print"
)

type Node struct {
	Name     string
	StartPos Position
	EndPos   Position
	Children []*Node
}

func (tt *TraceTree) ToTree(name string) *Node {
	rule := tt.grammar.ruleForID[tt.RuleID]
	n := &Node{
		Name:     name,
		StartPos: tt.StartPos,
		EndPos:   tt.EndPos,
	}
	return n
	switch tRule := rule.(type) {
	case *ref:
		name = tRule.name
	}
	if len(tt.ItemTraces) > 0 {
		for _, itemTrace := range tt.ItemTraces {
			tree := itemTrace.ToTree(name)
			if tree == nil {
				continue
			}
			n.Children = append(n.Children, tree)
		}
		return n
	} else if tt.ChoiceTrace != nil {
		return tt.ChoiceTrace.ToTree(name)
	} else if tt.KeywordMatch != "" {
		return nil
	} else if tt.RegexMatch != "" {
		return nil
	} else if tt.Success {
		return nil
	} else if tt.RefTrace != nil {
		return tt.RefTrace.ToTree(name)
	} else {
		panic(fmt.Sprintf("don't know how to en-tree a %+v", tt))
	}
}

func (n *Node) Format() pp.Doc {
	var children []pp.Doc
	for _, child := range n.Children {
		children = append(children, child.Format())
	}
	var docs = []pp.Doc{
		pp.Text(n.Name),
		pp.Textf(" [%v - %v]", n.StartPos.CompactString(), n.EndPos.CompactString()),
	}
	if len(n.Children) > 0 {
		docs = append(
			docs,
			pp.Text(" {"),
			pp.Newline,
			pp.Indent(2, pp.Join(children, pp.Newline)),
			pp.Newline,
			pp.Text("}"),
		)
	}
	return pp.Seq(docs)
}

func (n *Node) Text() string {
	panic("implement me")
}
