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
		p.OptWhitespace,
		p.Ref("keyValueList"),
		p.OptWhitespace,
		p.Text("}"),
	),
	"keyValueList": p.ListRule("keyValue", "keyValueList", p.Ref("sep")),
	"keyValue": p.SeqV(
		p.Ref("stringLit"),
		p.OptWhitespace,
		p.Text(":"),
		p.OptWhitespace,
		p.Ref("value"),
	),
	// array
	"array": p.SeqV(
		p.Text("["),
		p.OptWhitespace,
		p.Ref("valueList"),
		p.OptWhitespace,
		p.Text("]"),
	),
	"valueList": p.ListRule("value", "valueList", p.Ref("sep")),
	"sep":       p.SeqV(p.OptWhitespace, p.Text(","), p.OptWhitespace),
	// literals
	"stringLit": p.StringLit,
	"numberLit": p.SignedFloatLit,
	"bool":      p.ChoiceV(p.Text("true"), p.Text("false")),
	"null":      p.Text("null"),
}
