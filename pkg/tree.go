package parserlib

import (
	"fmt"

	pp "github.com/vilterp/go-pretty-print"
)

type Node struct {
	Name     string
	Span     SourceSpan
	Children []*Node
}

type TextNode struct {
	Span SourceSpan
	Text string
}

func (tn *TextNode) String() string {
	return fmt.Sprintf(`"%s"@%v`, tn.Text, tn.Span)
}

func (tt *TraceTree) ToTree() *Node {
	rule := tt.Rule
	name := ""
	switch tRule := rule.(type) {
	case *ref:
		name = tRule.name
		return &Node{
			Name: name,
			// TODO(vilterp): use SourceSpan in TraceTree too?
			Span: SourceSpan{
				From: tt.StartPos,
				To:   tt.EndPos,
			},
			Children: tt.RefTrace.getChildren(),
		}
	default:
		panic(fmt.Sprintf("only should on Ref, not %T", rule))
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
		return []*Node{tt.ToTree()}
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
		pp.Textf(" %s", n.Span.String()),
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

// TODO(vilterp): point to origInput somewhere from Node??
func (n *Node) Text(origInput string) *TextNode {
	return &TextNode{
		Span: n.Span,
		Text: n.Span.GetText(origInput),
	}
}
