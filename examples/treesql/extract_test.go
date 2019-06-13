package treesql_test

import (
	"fmt"
	"testing"

	"github.com/vilterp/go-parserlib/examples/treesql"
	"github.com/vilterp/go-parserlib/pkg/psi"
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
			`Select <many: true, table_name: "posts"@[1:6 - 1:11]>
  Selection <name: "id"@[2:3 - 2:5]>
  Selection <name: "title"@[3:4 - 3:9]>
  Selection <name: "body"@[4:4 - 4:8]>
  Selection <name: "comments"@[5:3 - 5:11]>
    Select <many: true, table_name: "comments"@[5:18 - 5:26]>
      Selection <name: "id"@[6:5 - 6:7]>
      Selection <name: "body"@[7:6 - 7:10]>`,
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
			actual := psi.Format(sel).String()
			if actual != testCase.output {
				t.Fatalf("EXPECTED\n\n%v\n\nGOT\n\n%v", testCase.output, actual)
			}
		})
	}
}
