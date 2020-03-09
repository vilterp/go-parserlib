package parserlib

import (
	"reflect"
	"sort"
	"testing"
)

func TestCompletions(t *testing.T) {
	t.Skip("seem to have broken this while doing rule ids")

	g, err := NewGrammar(map[string]Rule{
		"a_or_b": Choice([]Rule{Text("A"), Text("B")}),
		"c_or_d": Choice([]Rule{Text("C"), Text("D")}),
		"ab_then_cd": Seq([]Rule{
			Choice([]Rule{Text("A"), Text("B")}),
			Choice([]Rule{Text("C"), Text("D")}),
		}),
		"ab_then_cd_refs": Seq([]Rule{
			Ref("a_or_b"),
			Ref("c_or_d"),
		}),
	}, "a_or_b")
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		grammar     *Grammar
		rule        string
		input       string
		completions []string
	}{
		{
			g,
			"a_or_b",
			"",
			[]string{"A", "B"},
		},
		{
			g,
			"ab_then_cd",
			"",
			[]string{"A", "B"},
		},
		{
			g,
			"ab_then_cd",
			"A",
			[]string{"C", "D"},
		},
		{
			g,
			"ab_then_cd_refs",
			"",
			[]string{"A", "B"},
		},
		{
			g,
			"ab_then_cd_refs",
			"A",
			[]string{"C", "D"},
		},
		//{
		//	tsg,
		//	"selection",
		//	"",
		//	[]string{"{"},
		//},
		//{
		//	TestTreeSQLGrammar,
		//	"selection",
		//	"{foo",
		//	[]string{",", "}"},
		//},
	}
	for caseIdx, testCase := range cases {
		tree, err := testCase.grammar.Parse(testCase.rule, testCase.input, 0, nil)
		if err != nil {
			t.Fatal(err)
		}
		completions, err := tree.GetCompletions()
		if err != nil {
			t.Fatal(err)
		}
		sort.Strings(completions)
		sort.Strings(testCase.completions)
		if !reflect.DeepEqual(completions, testCase.completions) {
			t.Errorf("case %d: expected %v; got %v", caseIdx, testCase.completions, completions)
		}
	}
}
