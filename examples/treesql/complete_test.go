package treesql_test

import (
	"fmt"
	"testing"

	"github.com/vilterp/go-parserlib/examples/treesql"
	parserlib "github.com/vilterp/go-parserlib/pkg"
	"github.com/vilterp/go-parserlib/pkg/psi"
)

type completionTestCase struct {
	input       string
	cursorPos   parserlib.Position
	completions psi.Completions
}

func TestComplete(t *testing.T) {
	cases := []completionTestCase{
		{
			`MANY po {}`,
			parserlib.Position{Line: 1, Col: 7, Offset: 6},
			psi.Completions{"posts"},
		},
		{
			`MANY posts {
  id,
  comments: MANY co {
  }
}`,
			parserlib.Position{Line: 3, Col: 20, Offset: 38},
			psi.Completions{"comments"},
		},
	}

	for idx, testCase := range cases {
		t.Run(fmt.Sprintf("case %d", idx), func(t *testing.T) {
			traceTree, _ := treesql.Grammar.Parse("select", testCase.input, 0, nil)
			tree := traceTree.ToTree()
			sel := treesql.ToSelect(tree)

			completions := treesql.Complete(sel, blogSchema, testCase.cursorPos)
			if completions.String() != testCase.completions.String() {
				t.Fatalf("EXPECTED\n\n%v\n\nGOT\n\n%v", testCase.completions.String(), completions.String())
			}
		})
	}
}
