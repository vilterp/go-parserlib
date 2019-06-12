package treesql

import (
	"testing"

	"github.com/vilterp/go-parserlib/pkg/logger"
)

func TestTreeSQL(t *testing.T) {
	queries := []string{
		"MANY posts { id }",
		"MANY posts { id, comments: ",
		"MANY posts { id, comments: MANY",
	}

	for i, query := range queries {
		trace, err := Grammar.Parse("select", query, 0, logger.NewStdoutLogger())
		if err != nil {
			t.Fatalf("case %d: %v\n\n%v\n\n%v", i, err, trace.Format(), trace.ToTree("select").Format().String())
		}
		t.Logf("tree:\n\n%v", trace.ToTree("select").Format())
	}
}
