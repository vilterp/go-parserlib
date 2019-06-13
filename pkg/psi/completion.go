package psi

import (
	"fmt"
	"strings"
)

type Completion struct {
	Kind    string
	Content string
}

func (c *Completion) String() string {
	return fmt.Sprintf("%s: %s", c.Kind, c.Content)
}

type Completions []*Completion

func (c Completions) String() string {
	var lines []string
	for _, line := range c {
		lines = append(lines, line.String())
	}
	return strings.Join(lines, "\n")
}
