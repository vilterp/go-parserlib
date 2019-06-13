package treesql

import (
	"fmt"

	parserlib "github.com/vilterp/go-parserlib/pkg"
	"github.com/vilterp/go-parserlib/pkg/psi"
)

func ToSelect(n *parserlib.Node) *Select {
	if n.Name != "select" {
		panic(fmt.Sprintf("expecting `select` got %s", n.Name))
	}

	tn := n.MustGetChildWithName("table_name")
	selectionFieldsNode := n.
		MustGetChildWithName("selections").
		MustGetChildWithName("selection_fields")

	return &Select{
		BaseNode:   psi.BaseNode{Span: n.Span},
		Many:       true, // TODO: guess I have to name this in the grammar
		TableName:  tn.Text(),
		Selections: getSubSelections(selectionFieldsNode),
	}
}

func getSubSelections(n *parserlib.Node) []*Selection {
	if n.Name != "selection_fields" {
		panic("expecting `selections`")
	}
	var out []*Selection
	if len(n.Children) == 0 {
		return nil
	}

	selectionField := n.MustGetChildWithName("selection_field")
	out = append(out, getSelection(selectionField))

	nextSelectionFieldss := n.GetChildrenWithName("selection_fields")
	if len(nextSelectionFieldss) == 1 {
		nextSelectionFields := nextSelectionFieldss[0]
		out = append(out, getSubSelections(nextSelectionFields)...)
	}
	return out
}

func getSelection(n *parserlib.Node) *Selection {
	if n.Name != "selection_field" {
		panic(fmt.Sprintf("expecting `selection_field`; got %s", n.Name))
	}
	columnName := n.MustGetChildWithName("column_name")
	out := &Selection{
		BaseNode: psi.BaseNode{Span: n.Span},
		Name:     columnName.Text(),
	}
	selects := n.GetChildrenWithName("select")
	if len(selects) == 1 {
		out.SubSelect = ToSelect(selects[0])
	}
	return out
}
