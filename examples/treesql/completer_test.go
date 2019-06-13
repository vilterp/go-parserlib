package treesql_test

import (
	"fmt"
	"testing"

	"github.com/vilterp/go-parserlib/examples/treesql"
)

type selectTestCase struct {
	input  string
	output string
}

func TestToSelect(t *testing.T) {
	cases := []selectTestCase{
		{
			`MANY posts {
	 id,
   title,
   body,
	 comments: MANY comments {
	   id,
     body
	 }
	}`,
			`MANY "posts"@[1:6 - 1:11] {
  "id"@[2:3 - 2:5]
  "title"@[3:4 - 3:9]
  "body"@[4:4 - 4:8]
  "comments"@[5:3 - 5:11]: MANY "comments"@[5:18 - 5:26] {
    "id"@[6:5 - 6:7]
    "body"@[7:6 - 7:10]
  }
}`,
		},
	}

	for idx, testCase := range cases {
		t.Run(fmt.Sprintf("case %d", idx), func(t *testing.T) {
			traceTree, err := treesql.Grammar.Parse("select", testCase.input, 0, nil)
			if err != nil {
				t.Fatal(err)
			}
			tree := traceTree.ToTree()
			sel := treesql.ToSelect(tree)
			actual := sel.Format().String()
			if actual != testCase.output {
				t.Fatalf("EXPECTED\n\n%v\n\nGOT\n\n%v", testCase.output, actual)
			}
		})
	}
}
