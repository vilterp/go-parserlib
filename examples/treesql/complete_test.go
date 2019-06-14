package treesql_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vilterp/go-parserlib/examples/treesql"
	parserlib "github.com/vilterp/go-parserlib/pkg"
)

type completionTestCase struct {
	input       string
	cursorPos   parserlib.Position
	completions []string
}

func TestComplete(t *testing.T) {
	cases := []completionTestCase{
		// table names
		{
			`MANY `,
			parserlib.Position{Line: 1, Col: 6, Offset: 5},
			[]string{"table: posts", "table: comments"},
		},
		{
			`MANY po {}`,
			parserlib.Position{Line: 1, Col: 7, Offset: 6},
			[]string{"table: posts"},
		},
		{
			`MANY posts {
  id,
  comments: MANY co {
  }
}`,
			parserlib.Position{Line: 3, Col: 20, Offset: 38},
			[]string{"table: comments"},
		},
		// column names
		{
			`MANY posts { p }`,
			parserlib.Position{Line: 1, Col: 15, Offset: 14},
			[]string{"column: pics"},
		},
		{
			`MANY comments { p }`,
			parserlib.Position{Line: 1, Col: 18, Offset: 17},
			[]string{"column: post_id"},
		},
	}

	for idx, testCase := range cases {
		t.Run(fmt.Sprintf("case %d", idx), func(t *testing.T) {
			traceTree, _ := treesql.Grammar.Parse("select", testCase.input, 0, nil)
			tree := traceTree.ToTree()
			sel := treesql.ToSelect(tree)

			actual := treesql.Complete(sel, treesql.BlogSchema, testCase.cursorPos)
			expected := strings.Join(testCase.completions, "\n")
			if actual.String() != expected {
				t.Fatalf("EXPECTED\n\n%v\n\nGOT\n\n%v", expected, actual.String())
			}
		})
	}
}
