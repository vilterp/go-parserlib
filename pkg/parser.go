package parserlib

import (
	"fmt"
	"strings"

	"github.com/vilterp/go-parserlib/pkg/logger"
)

// TODO: structured parse errors
// each one has a position
// print out with position
// maybe store whole trace

type ParserState struct {
	grammar *Grammar
	input   string
	stack   []*ParserStackFrame
	trace   *TraceTree

	logger logger.Logger
}

type ParserStackFrame struct {
	input string
	// position we're at, exclusive
	// TODO: record start pos
	pos Position

	rule Rule
}

func (g *Grammar) Parse(startRuleName string, input string, cursor int, log logger.Logger) (*TraceTree, error) {
	if log == nil {
		log = logger.NewNoopLogger()
	}
	ps := ParserState{
		grammar: g,
		input:   input,
		logger:  log,
	}
	initPos := Position{Line: 1, Col: 1, Offset: 0}
	startRule := &RefRule{
		Name: startRuleName,
	}
	traceTree, err := ps.callRule(startRule, initPos, cursor)
	if err != nil {
		return traceTree, err
	}
	if traceTree.EndPos.Offset != len(input) {
		return traceTree, fmt.Errorf("%d extra chars at end of input", len(input)-traceTree.EndPos.Offset)
	}
	return traceTree, nil
}

func (ps *ParserState) callRule(rule Rule, pos Position, cursor int) (*TraceTree, *ParseError) {
	// Create and push stack frame.
	stackFrame := &ParserStackFrame{
		input: ps.input,
		rule:  rule,
		pos:   pos,
	}
	ps.stack = append(ps.stack, stackFrame)
	ps.logger.Indent()
	// Run the rule.
	traceTree, err := ps.runRule(cursor)
	// Pop the stack frame.
	ps.stack = ps.stack[:len(ps.stack)-1]
	ps.logger.Outdent()
	if traceTree == nil {
		panic(fmt.Sprintf("nil trace tree returned for rule %v", rule))
	}
	// Return.
	if err != nil {
		return traceTree, err
	}
	return traceTree, nil
}

func (sf *ParserStackFrame) Errorf(
	innerErr *ParseError, fmtString string, params ...interface{},
) *ParseError {
	return &ParseError{
		input:    sf.input,
		innerErr: innerErr,
		msg:      fmt.Sprintf(fmtString, params...),
		pos:      sf.pos,
	}
}

func (ps *ParserState) runRule(cursor int) (*TraceTree, *ParseError) {
	frame := ps.stack[len(ps.stack)-1]
	rule := frame.rule
	startPos := frame.pos
	minimalTrace := &TraceTree{
		origInput: ps.input,
		grammar:   ps.grammar,
		Rule:      rule,
		StartPos:  startPos,
		CursorPos: 0,
		EndPos:    startPos,
	}
	switch tRule := rule.(type) {
	case *ChoiceRule:
		trace := &TraceTree{
			origInput: ps.input,
			grammar:   ps.grammar,
			Rule:      rule,
			StartPos:  startPos,
			CursorPos: cursor,
		}
		maxAdvancement := 0
		maxAdvancementTraceIndex := 0
		var maxAdvancementTrace *TraceTree
		ps.logger.Log("CHOICE:", rule.String())
		ps.logger.Indent()
		defer ps.logger.Outdent()
		for choiceIdx, choice := range tRule.choices {
			ps.logger.Logf("trying choice %d", choiceIdx)
			choiceTrace, err := ps.callRule(choice, frame.pos, cursor)
			advancement := choiceTrace.EndPos.Offset - choiceTrace.StartPos.Offset
			ps.logger.Logf("choice %d advanced %d", choiceIdx, advancement)
			if advancement >= maxAdvancement {
				maxAdvancement = advancement
				maxAdvancementTrace = choiceTrace
				maxAdvancementTraceIndex = choiceIdx
				ps.logger.Logf("max advancement now %d", maxAdvancement)
				if err == nil {
					ps.logger.Logf("choice %d is a match: %s", choiceIdx, tRule.choices[choiceIdx].String())
					// We found a match!
					trace.EndPos = choiceTrace.EndPos
					trace.ChoiceIdx = choiceIdx
					trace.ChoiceTrace = choiceTrace
					return trace, nil
				}
			}
		}
		trace.EndPos = maxAdvancementTrace.EndPos
		trace.ChoiceIdx = maxAdvancementTraceIndex
		trace.ChoiceTrace = maxAdvancementTrace
		ps.logger.Logf(
			"err; going with max advancement trace %d: %s",
			maxAdvancementTraceIndex, tRule.choices[maxAdvancementTraceIndex].String(),
		)
		return trace, frame.Errorf(nil, "no match for rule `%s`", rule.String())
	case *SeqRule:
		trace := &TraceTree{
			origInput:  ps.input,
			grammar:    ps.grammar,
			Rule:       rule,
			StartPos:   startPos,
			CursorPos:  cursor,
			ItemTraces: make([]*TraceTree, len(tRule.items)),
		}
		advancement := 0
		// bla [foo, bar, baz]
		// ---  ---  ---  ---
		//   3
		//            ^
		ps.logger.Log("SEQ:", rule.String())
		ps.logger.Indent()
		defer ps.logger.Outdent()
		for itemIdx, item := range tRule.items {
			ps.logger.Logf("item %d", itemIdx)
			trace.AtItemIdx = itemIdx
			itemTrace, err := ps.callRule(item, frame.pos, cursor-advancement)
			advancement += itemTrace.GetSpan().Length()
			trace.EndPos = itemTrace.EndPos
			trace.ItemTraces[itemIdx] = itemTrace
			ps.logger.Logf("got to %d %+v advanced %d", itemIdx, trace.EndPos, advancement)
			if err != nil {
				return trace, frame.Errorf(err, "no match for sequence item %d", itemIdx)
			}
			frame.pos = itemTrace.EndPos
		}
		trace.EndPos = frame.pos
		return trace, nil
	case *KeywordRule:
		ps.logger.Log("KEYWORD:", tRule.value)
		minimalTrace.KeywordMatch = tRule.value
		remainingInput := ps.input[frame.pos.Offset:]
		if len(tRule.value) > len(remainingInput) {
			trimmed := strings.TrimPrefix(tRule.value, remainingInput)
			if len(trimmed) < len(tRule.value) {
				minimalTrace.EndPos = minimalTrace.StartPos.MoreOnLine(len(trimmed))
			}
			return minimalTrace, frame.Errorf(
				nil, `expected "%s"; got "%s"<EOF>`, tRule.value, remainingInput,
			)
		}
		trimmed := strings.TrimPrefix(remainingInput, tRule.value)
		advancement := len(remainingInput) - len(trimmed)
		minimalTrace.EndPos = minimalTrace.StartPos.MoreOnLine(advancement)
		if advancement == len(tRule.value) {
			minimalTrace.CursorPos = cursor
			return minimalTrace, nil
		}
		return minimalTrace, frame.Errorf(nil, `expected "%s"; got "%s"`, tRule.value, remainingInput)
	case *RefRule:
		ps.logger.Log("REF:", tRule.Name)
		refRule, ok := ps.grammar.rules[tRule.Name]
		if !ok {
			panic(fmt.Sprintf("nonexistent rule slipped through validation: %s", tRule.Name))
		}
		refTrace, err := ps.callRule(refRule, frame.pos, cursor)
		minimalTrace.RefTrace = refTrace
		minimalTrace.EndPos = refTrace.EndPos
		if err != nil {
			return minimalTrace, frame.Errorf(err, `no match for rule "%s"`, tRule.Name)
		}
		return &TraceTree{
			origInput: ps.input,
			grammar:   ps.grammar,
			Rule:      rule,
			StartPos:  startPos,
			CursorPos: cursor,
			EndPos:    refTrace.EndPos,
			RefTrace:  refTrace,
		}, nil
	case *NamedRule:
		ps.logger.Log("NAMED:", tRule.Name)
		innerTrace, err := ps.callRule(tRule.Inner, frame.pos, cursor)
		if err != nil {
			return minimalTrace, err
		}
		return &TraceTree{
			origInput:  ps.input,
			grammar:    ps.grammar,
			Rule:       rule,
			StartPos:   startPos,
			EndPos:     innerTrace.EndPos,
			CursorPos:  cursor,
			InnerTrace: innerTrace,
		}, nil
	case *RegexRule:
		ps.logger.Log("REGEX:", tRule.regex)
		loc := tRule.regex.FindStringIndex(ps.input[frame.pos.Offset:])
		if loc == nil || loc[0] != 0 {
			return minimalTrace, frame.Errorf(nil, "no match found for regex %s", tRule.regex)
		}
		matchText := ps.input[frame.pos.Offset : frame.pos.Offset+loc[1]]
		endPos := frame.pos
		for _, char := range matchText {
			if char == '\n' {
				endPos = endPos.Newline()
			} else {
				endPos = endPos.MoreOnLine(1)
			}
		}
		return &TraceTree{
			origInput:  ps.input,
			grammar:    ps.grammar,
			Rule:       rule,
			StartPos:   startPos,
			CursorPos:  cursor,
			EndPos:     endPos,
			RegexMatch: matchText,
		}, nil
	case *SucceedRule:
		minimalTrace.Success = true
		return minimalTrace, nil
	default:
		panic(fmt.Sprintf("not implemented: %T", rule))
	}
}
