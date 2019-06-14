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
			`select [1:1 - 3:2]
  table_name [1:6 - 1:11]
  selections [1:12 - 3:2]
    selection_fields [2:3 - 2:5]
      selection_field [2:3 - 2:5]
        column_name [2:3 - 2:5]`,
			``,
		},
		{
			`MANY posts {
  id,
  title
}`,
			`select [1:1 - 4:2]
  table_name [1:6 - 1:11]
  selections [1:12 - 4:2]
    selection_fields [2:3 - 3:8]
      selection_field [2:3 - 2:5]
        column_name [2:3 - 2:5]
      selection_fields [3:3 - 3:8]
        selection_field [3:3 - 3:8]
          column_name [3:3 - 3:8]`,
			``,
		},
		{
			`MANY posts { id, comments: `,
			`select [1:1 - 1:28]
  table_name [1:6 - 1:11]
  selections [1:12 - 1:28]
    selection_fields [1:14 - 1:28]
      selection_field [1:14 - 1:16]
        column_name [1:14 - 1:16]
      selection_fields [1:18 - 1:28]
        selection_field [1:18 - 1:28]
          column_name [1:18 - 1:26]
          select [1:28 - 1:28]`,
			"1:1: no match for rule \"select\": 1:12: no match for sequence item 6: 1:12: no match for rule \"selections\": 1:13: no match for sequence item 1: 1:14: no match for sequence item 1: 1:14: no match for rule \"selection_fields\": 1:14: no match for rule `(([selection_field, [\",\", (/\\s+/ | <succeed>)], selection_fields] | selection_field) | <succeed>)`",
		},
		{
			`MANY posts { id, comments: MANY c`,
			`select [1:1 - 1:34]
  table_name [1:6 - 1:11]
  selections [1:12 - 1:34]
    selection_fields [1:14 - 1:34]
      selection_field [1:14 - 1:16]
        column_name [1:14 - 1:16]
      selection_fields [1:18 - 1:34]
        selection_field [1:18 - 1:34]
          column_name [1:18 - 1:26]
          select [1:28 - 1:34]
            table_name [1:33 - 1:34]`,
			"1:1: no match for rule \"select\": 1:12: no match for sequence item 6: 1:12: no match for rule \"selections\": 1:13: no match for sequence item 1: 1:14: no match for sequence item 1: 1:14: no match for rule \"selection_fields\": 1:14: no match for rule `(([selection_field, [\",\", (/\\s+/ | <succeed>)], selection_fields] | selection_field) | <succeed>)`",
		},
	}

	log := logger.NewNoopLogger()

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			trace, err := treesql.Grammar.Parse("select", testCase.input, 0, log)
			tree := trace.ToRuleTree().Format().String()

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
