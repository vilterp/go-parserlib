package treesql

import (
	"strings"

	parserlib "github.com/vilterp/go-parserlib/pkg"
	"github.com/vilterp/go-parserlib/pkg/psi"
)

func Complete(sel *Select, schema *SchemaDesc, pos parserlib.Position) psi.Completions {
	path := psi.GetPath(sel, pos)

	if path == nil {
		return nil
	}

	if path.AttrName == "table_name" {
		return completeTableName(schema, path.AttrText.Text)
	}

	return nil
}

func completeTableName(schema *SchemaDesc, text string) psi.Completions {
	var out psi.Completions
	for name := range schema.Tables {
		if strings.HasPrefix(name, text) {
			out = append(out, psi.Completion(name))
		}
	}
	return out
}
