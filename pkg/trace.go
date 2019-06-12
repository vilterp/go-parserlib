package parserlib

import (
	"fmt"

	pp "github.com/vilterp/go-pretty-print"
)

type TraceTree struct {
	grammar *Grammar

	RuleID    RuleID
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
	// If it's a keyword
	KeywordMatch string
	// If it's a ref
	RefTrace *TraceTree `json:",omitempty"`
	// If it's a success
	Success bool
}

func (tt *TraceTree) Format() pp.Doc {
	rule := tt.grammar.ruleForID[tt.RuleID]

	switch tRule := rule.(type) {
	case *choice:
		return pp.Seq([]pp.Doc{
			pp.Textf("CHOICE(%d, ", tt.ChoiceIdx),
			pp.Newline,
			pp.Indent(2, tt.ChoiceTrace.Format()),
			pp.Newline,
			pp.Text(")"),
		})
	case *sequence:
		seqDocs := make([]pp.Doc, len(tt.ItemTraces))
		for idx, item := range tt.ItemTraces {
			seqDocs[idx] = item.Format()
		}
		return pp.Seq([]pp.Doc{
			pp.Text("SEQUENCE("),
			pp.Newline,
			pp.Indent(2, pp.Join(seqDocs, pp.CommaNewline)),
			pp.Newline,
			pp.Text(")"),
		})
	case *regex:
		return pp.Textf("REGEX(%#v)", tt.RegexMatch)
	case *succeed:
		return pp.Text("SUCCESS")
	case *ref:
		return pp.Seq([]pp.Doc{
			pp.Textf("REF(%s,", tRule.name),
			pp.Newline,
			pp.Indent(2, tt.RefTrace.Format()),
			pp.Newline,
			pp.Text(")"),
		})
	case *keyword:
		return pp.Textf("%#v", tRule.value)
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

const (
	PosToLeft = iota
	PosLeftEdge
	PosWithin
	PosRightEdge
	PosToRight
)
