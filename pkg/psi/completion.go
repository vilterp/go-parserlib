package psi

import "strings"

//type Completion struct {
//	// TODO(vilterp): how does this represent completion partway into a token?
//	Pos  parserlib.Position
//	Text string
//}

type Completion string

type Completions []Completion

func (c Completions) String() string {
	var lines []string
	for _, line := range c {
		lines = append(lines, string(line))
	}
	return strings.Join(lines, "\n")
}
