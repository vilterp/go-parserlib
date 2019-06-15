package treesql

import (
	"fmt"

	parserlib "github.com/vilterp/go-parserlib/pkg"
	"github.com/vilterp/go-parserlib/pkg/psi"
)

func GetHighlightedElement(n psi.Node, pos parserlib.Position, schema *SchemaDesc) *psi.HighlightedElement {
	path := psi.GetPath(n, pos)

	if path == nil {
		return nil
	}

	if path.NodeName == "Select" && path.AttrName == "table_name" {
		tableName := path.AttrText.Text
		if _, ok := schema.Tables[tableName]; !ok {
			return nil
		}
		return &psi.HighlightedElement{
			Node: path.AttrText,
			Path: tableName,
		}
	}

	if path.NodeName == "Selection" && path.AttrName == "name" {
		colName := path.AttrText.Text
		table := findContainingTableName(path)
		tableDesc := schema.Tables[table]
		if _, ok := tableDesc.Columns[colName]; !ok {
			return nil
		}
		return &psi.HighlightedElement{
			Path: fmt.Sprintf("%s/%s", table, colName),
			Node: path.AttrText,
		}
	}

	return nil
}
