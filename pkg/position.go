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
