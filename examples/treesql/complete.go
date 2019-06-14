package treesql

import (
	parserlib "github.com/vilterp/go-parserlib/pkg"
	"github.com/vilterp/go-parserlib/pkg/psi"
)

func Complete(sel *Select, schema *SchemaDesc, pos parserlib.Position) psi.Completions {
	path := psi.GetPath(sel, pos)

	if path == nil {
		return nil
	}

	if path.NodeName == "Select" && path.AttrName == "table_name" {
		return completeTableName(schema, path.AttrText.Text)
	}

	if path.NodeName == "Selection" && path.AttrName == "name" {
		return completeColumnName(schema, path)
	}

	return nil
}

func completeTableName(schema *SchemaDesc, text string) psi.Completions {
	var out psi.Completions
	for name := range schema.Tables {
		if psi.PrefixMatch(name, text) {
			out = append(out, &psi.Completion{
				Kind:    "table",
				Content: name,
			})
		}
	}
	return out
}

func completeColumnName(schema *SchemaDesc, path *psi.Path) psi.Completions {
	tableName := findContainingTableName(path)
	if tableName == "" {
		// TODO(vilterp): log or something? panicking seems harsh.
		return nil
	}
	tableDesc, ok := schema.Tables[tableName]
	if !ok {
		// can't autocomplete columns for a nonexistent table
		return nil
	}
	var out psi.Completions
	for name := range tableDesc.Columns {
		// TODO: fuzzy in-order-contains instead of has prefix
		//   i.e. foo => /f.*o.*o/
		if psi.PrefixMatch(name, path.AttrText.Text) {
			out = append(out, &psi.Completion{
				Kind:    "column",
				Content: name,
			})
		}
	}
	return out
}

func findContainingTableName(path *psi.Path) string {
	// iterate backwards up the chain til we find a Select
	for i := len(path.Nodes) - 1; i >= 0; i-- {
		node := path.Nodes[i]
		if node.TypeName() == "Select" {
			return node.AttrNodes()["table_name"].Text
		}
	}
	return ""
}
