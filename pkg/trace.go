package parserlib

import (
	"fmt"

	pp "github.com/vilterp/go-pretty-print"
)

type TraceTree struct {
	grammar   *Grammar
	origInput string

	Rule      Rule
	StartPos  Position
	CursorPos int // cursor offset relative to StartPos
	EndPos    Position

	// If it's a choice node.
	ChoiceIdx   int
	ChoiceTrace *TraceTree `json:",omitempty"`
	// If it's a sequence
	AtItemIdx  int
	ItemTraces []*TraceTree `json:",omitempty"`
	// If it's a regex
	RegexMatch string
	// If it's a text
	TextMatch string
	// If it's a ref
	RefTrace *TraceTree `json:",omitempty"`
	// If it's a success
	Success bool
	// If it's a named rule
	InnerTrace *TraceTree `json:",omitempty"`
}

func (tt *TraceTree) Format() pp.Doc {
	rule := tt.Rule

	switch tRule := rule.(type) {
	case *ChoiceRule:
		return pp.SeqV(
			pp.Textf("CHOICE(%d, ", tt.ChoiceIdx),
			pp.Newline,
			pp.Indent(2, tt.ChoiceTrace.Format()),
			pp.Newline,
			pp.Text(")"),
		)
	case *SeqRule:
		seqDocs := make([]pp.Doc, len(tt.ItemTraces))
		for idx, item := range tt.ItemTraces {
			seqDocs[idx] = item.Format()
		}
		return pp.SeqV(
			pp.Text("SEQUENCE("),
			pp.Newline,
			pp.Indent(2, pp.Join(seqDocs, pp.CommaNewline)),
			pp.Newline,
			pp.Text(")"),
		)
	case *SucceedRule:
		return pp.Text("SUCCESS")
	case *RefRule:
		return pp.SeqV(
			pp.Textf("REF(%s,", tRule.Name),
			pp.Newline,
			pp.Indent(2, tt.RefTrace.Format()),
			pp.Newline,
			pp.Text(")"),
		)
	case *TextRule:
		return pp.Textf("%#v", tRule.value)
	case *NamedRule:
		return pp.SeqV(
			pp.Textf("NAMED(%v", tRule.name),
			pp.CommaNewline,
			pp.Indent(2, tt.InnerTrace.Format()),
			pp.Newline,
			pp.Text(")"),
		)
	default:
		panic(fmt.Sprintf("don't know how to format a %T trace", rule))
	}
}

func (tt *TraceTree) OptWhitespaceSurroundRes() *TraceTree {
	whitespaceSeq := tt
	return whitespaceSeq.ItemTraces[1]
}

func (tt *TraceTree) GetSpan() SourceSpan {
	return SourceSpan{
		From: tt.StartPos,
		To:   tt.EndPos,
	}
}
