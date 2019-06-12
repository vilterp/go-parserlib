package treesql

import "testing"

func TestTreeSQL(t *testing.T) {
	tree, err := Grammar.Parse("select", "MANY posts { id }", 0)
	if err != nil {
		t.Fatal(t)
	}
	t.Log(tree.Format())
}
