package treesql

import p "github.com/vilterp/go-parserlib/pkg"

var Grammar *p.Grammar

func init() {
	grammar, err := p.NewGrammar(grammarRules, "select")
	if err != nil {
		panic(err)
	}
	Grammar = grammar
}

var grammarRules = map[string]p.Rule{
	"select": p.Seq([]p.Rule{
		p.Choice([]p.Rule{
			p.Text("ONE"),
			p.Text("MANY"),
		}),
		p.Whitespace,
		p.Named("table_name", p.Ident),
		p.Whitespace,
		p.Opt(p.Ref("where_clause")),
		p.OptWhitespace,
		p.Ref("selections"),
	}),
	"column_name": p.Ident,
	"where_clause": p.Seq([]p.Rule{
		p.Text("WHERE"),
		p.Whitespace,
		p.Ref("column_name"),
		p.OptWhitespace,
		p.Text("="),
		p.OptWhitespace,
		p.Ref("expr"),
	}),
	"selections": p.Seq([]p.Rule{
		p.Text("{"),
		p.OptWhitespaceSurround(
			p.Ref("selection_fields"),
		),
		p.Text("}"),
	}),
	// TODO: intercalate combinator (??)
	"selection_fields": p.ListRule(
		"selection_field",
		"selection_fields",
		p.Seq([]p.Rule{p.Text(","), p.OptWhitespace}),
	),
	"selection_field": p.Seq([]p.Rule{
		p.Ref("column_name"),
		p.Opt(p.Seq([]p.Rule{
			p.Text(":"),
			p.OptWhitespace,
			p.Ref("select"),
		})),
	}),
	"expr": p.Choice([]p.Rule{
		p.Ident,
		p.StringLit,
		p.SignedIntLit,
	}),
}
