package test_harness

import (
	"fmt"
	"testing"

	"github.com/vilterp/go-parserlib/examples/treesql"
)

func TestTestHarness(t *testing.T) {
	// mainly testing that none of these panic the server
	inputs := []string{
		``,
		`M`,
		`MANY `,
		`MANY foo`,
		`MANY foo {`,
		`MANY foo {}`,
		`MANY foo { id`,
		`MANY foo { id }`,
		`MANY foo { id, }`,
		`MANY foo { id, c }`,
		`MANY foo { id, c: }`,
		`MANY foo { id, c: M }`,
		`MANY foo { id, c: MANY  }`,
		`MANY foo { id, c: MANY b }`,
	}

	s := NewServer(treesql.MakeLanguage(treesql.BlogSchema))

	for idx, input := range inputs {
		t.Run(fmt.Sprintf("case %d", idx), func(t *testing.T) {
			_, err := s.getResp(&completionsRequest{
				CursorPos: 0,
				Input:     input,
			})
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
