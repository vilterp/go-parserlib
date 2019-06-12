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

func (tt *TraceTree) ToTree() *Node {
	rule := tt.grammar.ruleForID[tt.RuleID]
	name := ""
	switch tRule := rule.(type) {
	case *ref:
		name = tRule.name
		return &Node{
			Name:     name,
			StartPos: tt.StartPos,
			EndPos:   tt.EndPos,
			Children: tt.RefTrace.getChildren(),
		}
	default:
		// TODO(vilterp): this is janky.
		//   needs some kind of refactoring; I'm not sure what
		name, ok := tt.grammar.nameForID[tt.RuleID]
		if !ok {
			panic(fmt.Sprintf("name not found for rule %s", rule.String()))
		}
		return &Node{
			Name:     name,
			StartPos: tt.StartPos,
			EndPos:   tt.EndPos,
			Children: tt.getChildren(),
		}
	}
}

func (tt *TraceTree) getChildren() []*Node {
	if len(tt.ItemTraces) > 0 {
		var out []*Node
		for _, itemTrace := range tt.ItemTraces {
			if itemTrace == nil {
				continue
			}
			out = append(out, itemTrace.getChildren()...)
		}
		return out
	} else if tt.ChoiceTrace != nil {
		return tt.ChoiceTrace.getChildren()
	} else if tt.KeywordMatch != "" {
		return nil
	} else if tt.RegexMatch != "" {
		return nil
	} else if tt.Success {
		return nil
	} else if tt.RefTrace != nil {
		return []*Node{tt.RefTrace.ToTree()}
	} else {
		return nil
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
			pp.Newline,
			pp.Indent(2, pp.Join(children, pp.Newline)),
		)
	}
	return pp.Seq(docs)
}

func (n *Node) Text() string {
	panic("implement me")
}
