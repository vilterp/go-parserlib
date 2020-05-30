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
		p.Ref("optWhitespace"),
		p.Ref("keyValueList"),
		p.Ref("optWhitespace"),
		p.Text("}"),
	),
	"keyValueList": p.ListRule("keyValue", "keyValueList", p.Ref("sep")),
	"keyValue": p.SeqV(
		p.Ref("stringLit"),
		p.Ref("optWhitespace"),
		p.Text(":"),
		p.Ref("optWhitespace"),
		p.Ref("value"),
	),
	// array
	"array": p.SeqV(
		p.Text("["),
		p.Ref("optWhitespace"),
		p.Ref("valueList"),
		p.Ref("optWhitespace"),
		p.Text("]"),
	),
	"valueList": p.ListRule("value", "valueList", p.Ref("sep")),
	"sep":       p.SeqV(p.Ref("optWhitespace"), p.Text(","), p.Ref("optWhitespace")),
	// literals
	"stringLit": p.StringLit,
	"numberLit": p.SignedFloatLit,
	"bool":      p.ChoiceV(p.Text("true"), p.Text("false")),
	"null":      p.Text("null"),
	// whitespace
	"optWhitespace": p.ListRule("whitespace", "optWhitespace", p.Succeed),
	"whitespace":    p.ChoiceV(p.Text(" "), p.Text("\n"), p.Text("\t")),
}
