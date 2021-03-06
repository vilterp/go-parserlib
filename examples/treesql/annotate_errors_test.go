package treesql_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vilterp/go-parserlib/examples/treesql"
)

type annotateTestCase struct {
	input  string
	errors []string
}

func TestAnnotate(t *testing.T) {
	cases := []annotateTestCase{
		{
			`MANY boop {}`,
			[]string{
				"[1:6 - 1:10]: no table named `boop`",
			},
		},
		{
			`MANY posts {
  id, body, foo,
  comments: MANY comments {
    id, bar
  }
}`,
			[]string{
				"[2:13 - 2:16]: no column `foo` in table `posts`",
				"[4:9 - 4:12]: no column `bar` in table `comments`",
			},
		},
		{
			`MANY foo { id, c: MANY comments { doop } }`,
			[]string{
				"[1:6 - 1:9]: no table named `foo`",
				"[1:35 - 1:39]: no column `doop` in table `comments`",
			},
		},
	}

	for idx, testCase := range cases {
		t.Run(fmt.Sprintf("case %d", idx), func(t *testing.T) {
			traceTree, err := treesql.Grammar.Parse("select", testCase.input, 0, nil)
			if err != nil {
				t.Fatal(err)
			}

			tree := traceTree.ToRuleTree()
			selectPsi := treesql.ToSelect(tree)

			errors := treesql.Annotate(treesql.BlogSchema, selectPsi)
			var errorStrings []string
			for _, err := range errors {
				errorStrings = append(errorStrings, err.String())
			}
			actualErrorLines := strings.Join(errorStrings, "\n")
			expectedErrorLines := strings.Join(testCase.errors, "\n")
			if actualErrorLines != expectedErrorLines {
				t.Fatalf("EXPECTED\n\n%v\n\nGOT\n\n%v", expectedErrorLines, actualErrorLines)
			}
		})
	}
}
