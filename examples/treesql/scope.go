package treesql

import (
	"github.com/vilterp/go-parserlib/pkg/psi"
)

// Top level scope

type TopLevelScope struct {
	Tables map[string]*TableDesc
}

var _ psi.Scope = &TopLevelScope{}

func (t *TopLevelScope) GetItems() map[string]psi.ScopeItem {
	out := map[string]psi.ScopeItem{}
	for name := range t.Tables {
		out[name] = &TableItem{
			Name: name,
		}
	}
	return out
}

// Table scope

type TableScope struct {
	Table *TableDesc
}

var _ psi.Scope = &TableScope{}

func (t *TableScope) GetItems() map[string]psi.ScopeItem {
	// what if there are two completions with the same name
	//   that would be an issue for parsing
	//   so we'd have to namespace them
	out := map[string]psi.ScopeItem{}
	// get columns in table
	for colName := range t.Table.Columns {
		out[colName] = &ColumnItem{
			Name:  colName,
			Table: t.Table.Name,
		}
	}
	return out
}

// Column item

type ColumnItem struct {
	Name  string
	Table string
}

var _ psi.ScopeItem = &ColumnItem{}

func (c *ColumnItem) GetName() string {
	return c.Name
}

func (c *ColumnItem) GetType() string {
	return "Column"
}

// Table item

type TableItem struct {
	Name string
}

var _ psi.ScopeItem = &TableItem{}

func (t *TableItem) GetName() string {
	return t.Name
}

func (t *TableItem) GetType() string {
	return "Table"
}
