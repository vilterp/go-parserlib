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
	base := &Select{
		BaseNode: psi.BaseNode{Span: n.Span},
	}

	tn := n.GetChildWithName("table_name")
	if tn == nil {
		return base
	}
	base.TableName = tn.Text()
	selectionsNode := n.GetChildWithName("selections")
	if selectionsNode == nil {
		return base
	}
	selectionFieldsNode := selectionsNode.GetChildWithName("selection_fields")
	if selectionFieldsNode == nil {
		return base
	}

	// TODO: keep modifying base instead...
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

	selectionField := n.GetChildWithName("selection_field")
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
	columnName := n.GetChildWithName("column_name")
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
