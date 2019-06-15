package parserlib

import (
	"fmt"
	"strings"
)

type Position struct {
	Line   int
	Col    int
	Offset int
}

var StartPos = Position{
	Line:   1,
	Col:    1,
	Offset: 0,
}

func (pos *Position) String() string {
	return fmt.Sprintf("line %d, col %d", pos.Line, pos.Col)
}

func (pos *Position) CompactString() string {
	return fmt.Sprintf("%d:%d", pos.Line, pos.Col)
}

func (pos *Position) MoreOnLine(n int) Position {
	return Position{
		Col:    pos.Col + n,
		Line:   pos.Line,
		Offset: pos.Offset + n,
	}
}

func (pos *Position) Newline() Position {
	return Position{
		Col:    1,
		Line:   pos.Line + 1,
		Offset: pos.Offset + 1,
	}
}

func (pos *Position) Lt(other Position) bool {
	return pos.Offset < other.Offset
}

func (pos *Position) ShowInContext(input string) string {
	lines := strings.Split(input, "\n")
	inputLine := lines[pos.Line-1]
	return fmt.Sprintf(
		"%s\n%s",
		inputLine,
		strings.Repeat(" ", pos.Col-1)+"^",
	)
}

type SourceSpan struct {
	From Position
	To   Position
}

func (ss SourceSpan) Length() int {
	return ss.To.Offset - ss.From.Offset
}

func (ss SourceSpan) String() string {
	return fmt.Sprintf("[%s - %s]", ss.From.CompactString(), ss.To.CompactString())
}

func (ss SourceSpan) Contains(p Position) bool {
	return (ss.From == p || ss.From.Lt(p)) && (p.Lt(ss.To) || p == ss.To)
}

func (ss SourceSpan) GetText(input string) string {
	lines := strings.Split(input, "\n")
	if ss.From.Line == ss.To.Line {
		line := lines[ss.From.Line-1]
		return line[ss.From.Col-1 : ss.To.Col-1]
	}
	panic("TODO: GetText across multiple lines")
}

// PositionFromOffset gives a position from a 0-indexed offset
func PositionFromOffset(doc string, offset int) Position {
	pos := StartPos
	for i := 0; i < offset; i++ {
		char := doc[i]
		if char == '\n' {
			pos = pos.Newline()
		} else {
			pos = pos.MoreOnLine(1)
		}
	}
	return pos
}

func (p Position) EnsurePosition(origInput string) Position {
	offset := 0
	lines := strings.Split(origInput, "\n")
	for i := 0; i < p.Line-1; i++ {
		offset += len(lines[i])
		offset++
	}
	offset += p.Col - 1
	return Position{
		Line:   p.Line,
		Col:    p.Col,
		Offset: offset,
	}
}
