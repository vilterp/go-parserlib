package parserlib

import (
	"fmt"

	pp "github.com/vilterp/go-pretty-print"
)

type RuleNode struct {
	OrigInput string `json:"-"`
	Name      string
	Span      SourceSpan
	Children  []*RuleNode
}

type TextNode struct {
	Span SourceSpan
	Text string
}

func (tn *TextNode) String() string {
	return fmt.Sprintf(`"%s"@%v`, tn.Text, tn.Span)
}

func (tt *TraceTree) ToRuleTree() *RuleNode {
	rule := tt.Rule
	name := ""
	switch tRule := rule.(type) {
	case *ref:
		name = tRule.name
		return &RuleNode{
			OrigInput: tt.origInput,
			Name:      name,
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

func (tt *TraceTree) getChildren() []*RuleNode {
	if len(tt.ItemTraces) > 0 {
		var out []*RuleNode
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
		return []*RuleNode{tt.ToRuleTree()}
	} else {
		return nil
	}
}

func (n *RuleNode) Format() pp.Doc {
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

func (n *RuleNode) Text() *TextNode {
	return &TextNode{
		Span: n.Span,
		Text: n.Span.GetText(n.OrigInput),
	}
}

func (n *RuleNode) GetChildrenWithName(name string) []*RuleNode {
	var out []*RuleNode
	for _, child := range n.Children {
		if child.Name == name {
			out = append(out, child)
		}
	}
	return out
}

func (n *RuleNode) GetChildWithName(name string) *RuleNode {
	children := n.GetChildrenWithName(name)
	if len(children) != 1 {
		return nil
	}
	return children[0]
}
