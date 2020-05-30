package json

import p "github.com/vilterp/go-parserlib/pkg"

var Grammar *p.Grammar

func init() {
	grammar, err := p.NewGrammar(grammarRules, "value")
	if err != nil {
		panic(err)
	}
	Grammar = grammar
}

var grammarRules = map[string]p.Rule{
	"value": p.ChoiceV(
		p.Ref("object"),
		p.Ref("array"),
		p.Ref("numberLit"),
		p.Ref("stringLit"),
		p.Ref("bool"),
		p.Ref("null"),
	),
	// object
	"object": p.SeqV(
		p.Text("{"),
		p.Ref("keyValueList"),
		p.Text("}"),
	),
	"keyValueList": p.ListRule("keyValue", "keyValueList", p.Ref("sep")),
	"keyValue": p.SeqV(
		p.Ref("stringLit"),
		p.Text(":"),
		p.Ref("value"),
	),
	// array
	"array": p.SeqV(
		p.Text("["),
		p.Ref("valueList"),
		p.Text("]"),
	),
	"valueList": p.ListRule("value", "valueList", p.Ref("sep")),
	"sep":       p.Text(","),
	// literals
	"stringLit": p.ChoiceV(
		p.Text(`"foo"`),
		p.Text(`"bar"`),
		p.Text(`"baz"`),
		p.Text(`"bloop"`),
		p.Text(`"doop"`),
		p.Text(`"hello"`),
		p.Text(`"world"`),
	),
	"numberLit": p.ChoiceV(
		p.Text("1"),
		p.Text("2"),
		p.Text("3"),
		p.Text("4"),
		p.Text("5"),
		p.Text("6"),
	),
	"bool": p.ChoiceV(p.Text("true"), p.Text("false")),
	"null": p.Text("null"),
	// whitespace
	"optWhitespace": p.ListRule("whitespace", "optWhitespace", p.Succeed),
	"whitespace":    p.ChoiceV(p.Text(" "), p.Text("\n"), p.Text("\t")),
}
