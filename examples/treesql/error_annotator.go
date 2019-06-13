package treesql

import (
	"fmt"

	"github.com/vilterp/go-parserlib/pkg/psi"
)

type SchemaDesc struct {
	Tables map[string]*TableDesc
}

type TableDesc struct {
	Columns map[string]*ColDesc
}

type ColDesc struct {
	// TODO(vilterp): type
	// TODO(vilterp): foreign key
}

func Annotate(schema *SchemaDesc, n *Select) []*psi.ErrorAnnotation {
	var out []*psi.ErrorAnnotation
	tableDesc, ok := schema.Tables[n.TableName.Text]
	if !ok {
		out = append(out, &psi.ErrorAnnotation{
			Span:    n.TableName.Span,
			Message: fmt.Sprintf("no such table"),
		})
		return out
	}
	// TODO(vilterp): check foreign key existence...
	for _, sel := range n.Selections {
		if sel.SubSelect == nil {
			_, ok := tableDesc.Columns[sel.Name.Text]
			if !ok {
				out = append(out, &psi.ErrorAnnotation{
					Span:    sel.Name.Span,
					Message: fmt.Sprintf("no column `%s` in table `%s`", sel.Name.Text, n.TableName.Text),
				})
			}
		} else {
			out = append(out, Annotate(schema, sel.SubSelect)...)
		}
		// TODO(vilterp): check for duplicate names
	}
	return out
}
