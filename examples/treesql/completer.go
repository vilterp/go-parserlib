package treesql

import (
	"fmt"

	parserlib "github.com/vilterp/go-parserlib/pkg"
	pp "github.com/vilterp/go-pretty-print"
)

type Completion struct {
	// TODO(vilterp): how does this represent completion partway into a token?
	Pos  parserlib.Position
	Text string
}

func Complete(t *parserlib.Node, cursorPos parserlib.Position) []*Completion {
	if !t.Span.Contains(cursorPos) {
		return nil
	}
	return nil
}

// TODO(vilterp): make these implement a PSINode interface??
//   trace tree => node => PSI node?

// can do completion on this representation
// and annotate it with errors
type Select struct {
	Many       bool
	TableName  *parserlib.TextNode
	Selections []*Selection
}

type Selection struct {
	Name      *parserlib.TextNode
	SubSelect *Select
}

func ToSelect(t *parserlib.Node) *Select {
	if t.Name != "select" {
		panic(fmt.Sprintf("expecting `select` got %s", t.Name))
	}

	tn := t.MustGetChildWithName("table_name")
	selectionFieldsNode := t.
		MustGetChildWithName("selections").
		MustGetChildWithName("selection_fields")

	return &Select{
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
		Name: columnName.Text(),
	}
	selects := n.GetChildrenWithName("select")
	if len(selects) == 1 {
		out.SubSelect = ToSelect(selects[0])
	}
	return out
}

func (s *Select) Format() pp.Doc {
	var selDocs []pp.Doc
	for _, selection := range s.Selections {
		selDocs = append(selDocs, selection.Format())
	}
	return pp.Seq([]pp.Doc{
		pp.Text("MANY"), // TODO(vilterp): fix
		pp.Text(" "),
		pp.Text(s.TableName.String()),
		pp.Text(" {"),
		pp.Newline,
		pp.Indent(2, pp.Join(selDocs, pp.Newline)),
		pp.Newline,
		pp.Text("}"),
	})
}

func (s *Selection) Format() pp.Doc {
	if s.SubSelect == nil {
		return pp.Seq([]pp.Doc{
			pp.Text(s.Name.String()),
		})
	}
	return pp.Seq([]pp.Doc{
		pp.Text(s.Name.String()),
		pp.Text(": "),
		s.SubSelect.Format(),
	})
}
