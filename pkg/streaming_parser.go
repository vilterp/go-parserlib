package parserlib

import (
	"bufio"
	"fmt"
	"io"
	"unicode/utf8"
)

type StreamingParser struct {
	grammar  *Grammar
	reader   *bufio.Reader
	stackTop *StackFrame
	pos      Position
}

type StackFrame struct {
	Rule     Rule
	parent   *StackFrame
	startPos Position

	seqItem   int
	choiceIdx int // shit, I guess we can't do backtracking... yolo
}

func NewStreamingParser(in io.Reader, g *Grammar, startRule string) (*StreamingParser, error) {
	// TODO: validate that there are no regex rules...?
	rule, ok := g.rules[startRule]
	if !ok {
		return nil, fmt.Errorf("no such rule: %v", startRule)
	}
	startPos := Position{
		Line:   1,
		Col:    1,
		Offset: 0,
	}
	return &StreamingParser{
		reader:  bufio.NewReader(in),
		grammar: g,
		pos:     startPos,
		stackTop: &StackFrame{
			Rule:      rule,
			parent:    nil,
			startPos:  startPos,
			seqItem:   0,
			choiceIdx: 0,
		},
	}, nil
}

type EvtType = int

const (
	PushRule EvtType = iota
	PopRule
)

type Event struct {
	Type EvtType
	Rule Rule
	Pos  Position
	Text string
}

func (sp *StreamingParser) NextEvent() (*Event, error) {

}

func (sp *StreamingParser) runRule() error {
	switch tRule := sp.stackTop.Rule.(type) {
	case *ChoiceRule:
		// have to make choice immediately, right?
		XXX
	case *SeqRule:
		XXX
	case *TextRule:
		for _, expRune := range tRule.value {
			actualRune, err := sp.nextRune()
			if err != nil {
				return err
			}
			if actualRune != expRune {
				return fmt.Errorf("expected %v; got %v", expRune, actualRune)
			}
		}
	case *RefRule:
		XXX
	default:
		panic(fmt.Sprintf("unhandled rule type: %T", tRule))
	}
}

func (sp *StreamingParser) matches(rule Rule, r rune) bool {
	switch tRule := rule.(type) {
	case *TextRule:
		expRune, _ := utf8.DecodeRuneInString(tRule.value)
		return expRune == r
	case *SeqRule:
		return sp.matches(tRule.items[0], r)
	case *ChoiceRule:
		for _, choice := range tRule.choices {
			if sp.matches(choice, r) {
				return true
			}
		}
		return false
	case *RefRule:
		rule, ok := sp.grammar.rules[tRule.Name]
		if !ok {
			panic(fmt.Sprintf("unknown rule: %v", tRule.Name))
		}
		return sp.matches(rule, r)
	case *NamedRule:
		return sp.matches(tRule.Inner, r)
	case *SucceedRule:
		return true
	default:
		panic(fmt.Sprintf("unhandled rule type: %T", tRule))
	}
}

func (sp *StreamingParser) nextRune() (rune, error) {
	r, _, err := sp.reader.ReadRune()
	if err != nil {
		return 0, err
	}
	if r == '\n' {
		sp.pos.Line++
		sp.pos.Col = 0
	} else {
		sp.pos.Col++
	}
	sp.pos.Offset++
	return r, nil
}

func (sp *StreamingParser) peekRune() (rune, error) {
	r, _, err := sp.reader.ReadRune()
	if err != nil {
		return 0, err
	}
	if err := sp.reader.UnreadRune(); err != nil {
		return 0, err
	}
	return r, nil
}
