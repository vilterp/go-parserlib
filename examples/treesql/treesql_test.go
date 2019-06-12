package treesql_test

import (
	"fmt"
	"testing"

	"github.com/vilterp/go-parserlib/examples/treesql"
	"github.com/vilterp/go-parserlib/pkg/logger"
)

type testCase struct {
	input  string
	output string
	err    string
}

func TestTreeSQL(t *testing.T) {
	testCases := []testCase{
		{
			`MANY posts {
  id
}`,
			`select [1:1 - 1:18]
         table_name`,
			``,
		},
		{
			`MANY posts {
  id,
  title
}`,
			`select [1:1 - 1:18]
         table_name`,
			``,
		},
		{
			`MANY posts { id, comments: `,
			``,
			``,
		},
		{
			`MANY posts { id, comments: MANY`,
			``,
			``,
		},
	}

	log := logger.NewNoopLogger()

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			trace, err := treesql.Grammar.Parse("select", testCase.input, 0, log)
			tree := trace.ToTree().Format().String()

			if tree != testCase.output {
				t.Fatalf("EXPECTED\n\n%v\n\nGOT\n\n%v\n", testCase.output, tree)
			}
			errStr := ""
			if err != nil {
				errStr = err.Error()
			}
			if errStr != testCase.err {
				t.Fatalf("EXPECTED ERR\n\n%v\n\nGOT\n\n%v\n", testCase.err, errStr)
			}
		})
	}
}
