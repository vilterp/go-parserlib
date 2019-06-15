package treesql

import (
	"fmt"

	parserlib "github.com/vilterp/go-parserlib/pkg"
	"github.com/vilterp/go-parserlib/pkg/psi"
)

func GetHighlightedElement(n psi.Node, pos parserlib.Position) *psi.HighlightedElement {
	path := psi.GetPath(n, pos)

	if path == nil {
		return nil
	}

	if path.NodeName == "Select" && path.AttrName == "table_name" {
		return &psi.HighlightedElement{
			Node: path.AttrText,
			Path: path.AttrText.Text,
		}
	}

	if path.NodeName == "Selection" && path.AttrName == "name" {
		table := findContainingTableName(path)
		return &psi.HighlightedElement{
			Path: fmt.Sprintf("%s/%s", table, path.AttrText.Text),
			Node: path.AttrText,
		}
	}

	// TODO(vilterp): don't return stuff if it doesn't exist

	return nil
}
