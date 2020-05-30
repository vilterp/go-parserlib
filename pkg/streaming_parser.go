package parserlib

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"
)

type StreamingParser struct {
	grammar  *Grammar
	reader   *bufio.Reader
	stackTop *StackFrame
	pos      Position
}

type StackFrame struct {
	rule     Rule
	parent   *StackFrame
	startPos Position

	seqItem      int
	choicePushed bool
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
			rule:     rule,
			parent:   nil,
			startPos: startPos,
		},
	}, nil
}

type EvtType = int

const (
	PushRule EvtType = iota
	PopRule
)

var nameForType = map[EvtType]string{
	PushRule: "Push",
	PopRule:  "Pop",
}

type Event struct {
	Type EvtType
	Rule Rule
	Pos  Position
	Text string
}

func (e *Event) String() string {
	return fmt.Sprintf(
		"Evt{type: %s, rule: %s, text: %s, pos: %s}",
		nameForType[e.Type], e.Rule.String(), strconv.Quote(e.Text), e.Pos.CompactString(),
	)
}

func (sp *StreamingParser) NextEvent() (*Event, error) {
	switch tRule := sp.stackTop.rule.(type) {
	case *ChoiceRule:
		// just popped from choice; we're done
		if sp.stackTop.choicePushed {
			sp.popStack()
			return &Event{
				Type: PopRule,
				Rule: tRule,
				Pos:  sp.pos,
			}, nil
		}
		// have to make choice immediately, right?
		r, err := sp.peekRune()
		if err != nil {
			return nil, err
		}
		for _, choice := range tRule.choices {
			if sp.matches(choice, r) {
				sp.stackTop.choicePushed = true
				sp.pushStack(choice)
				return &Event{
					Type: PushRule,
					Rule: choice,
					Pos:  sp.pos,
				}, nil
			}
		}
		return nil, makeParseError(fmt.Sprintf("no choice matched %s", strconv.QuoteRune(r)), sp.pos, sp.stackTop)
	case *SeqRule:
		if sp.stackTop.seqItem == len(tRule.items) {
			return &Event{
				Type: PopRule,
				Rule: tRule,
				Pos:  sp.pos,
			}, nil
		}
		item := tRule.items[sp.stackTop.seqItem]
		sp.stackTop.seqItem++
		sp.pushStack(item)
		return &Event{
			Type: PushRule,
			Rule: item,
			Pos:  sp.pos,
		}, nil
	case *TextRule:
		for _, expRune := range tRule.value {
			actualRune, err := sp.nextRune()
			if err != nil {
				return nil, err
			}
			if actualRune != expRune {
				return nil, makeParseError(fmt.Sprintf("expected %s; got %s", strconv.QuoteRune(expRune), strconv.QuoteRune(actualRune)), sp.pos, sp.stackTop)
			}
		}
		sp.popStack()
		return &Event{
			Type: PopRule,
			Rule: tRule,
			Pos:  sp.pos,
			Text: tRule.value,
		}, nil
	case *RefRule:
		rule, ok := sp.grammar.rules[tRule.Name]
		if !ok {
			panic(fmt.Sprintf("unknown rule: %v", tRule.Name))
		}
		ret := &Event{
			Type: PushRule,
			Rule: rule,
			Pos:  sp.pos,
		}
		sp.pushStack(rule)
		return ret, nil
	case *SucceedRule:
		sp.popStack()
		return &Event{
			Type: PopRule,
			Rule: tRule,
			Pos:  sp.pos,
		}, nil
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

func (sp *StreamingParser) pushStack(rule Rule) *StackFrame {
	newTop := &StackFrame{
		rule:     rule,
		startPos: sp.pos,
		parent:   sp.stackTop,
	}
	sp.stackTop = newTop
	return newTop
}

// returns old top
func (sp *StreamingParser) popStack() *StackFrame {
	oldTop := sp.stackTop
	sp.stackTop = oldTop.parent
	return oldTop
}

func (sp *StreamingParser) nextRune() (rune, error) {
	r, _, err := sp.reader.ReadRune()
	if err != nil {
		return 0, err
	}
	if r == '\n' {
		sp.pos.Line++
		sp.pos.Col = 1
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

// errors

type parseError struct {
	msg   string
	pos   Position
	stack *StackFrame
}

func makeParseError(msg string, pos Position, stack *StackFrame) *parseError {
	return &parseError{
		msg:   msg,
		pos:   pos,
		stack: stack,
	}
}

func (e *parseError) Error() string {
	return fmt.Sprintf("parse error at %v: %v.\nStack:\n%v", e.pos.CompactString(), e.msg, formatStackTrace(e.stack))
}

func formatStackTrace(frame *StackFrame) string {
	if frame == nil {
		return ""
	}
	frameFmt := frame.rule.String()
	return "  " + frameFmt + "\n" + formatStackTrace(frame.parent)
}
