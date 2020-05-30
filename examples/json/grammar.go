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
		p.RepSep(p.Ref("keyValue"), p.Ref("sep")),
		p.Text("}"),
	),
	"keyValue": p.SeqV(
		p.Ref("stringLit"),
		p.Text(":"),
		p.Ref("value"),
	),
	// array
	"array": p.SeqV(
		p.Text("["),
		p.RepSep(p.Ref("value"), p.Ref("sep")),
		p.Text("]"),
	),
	// TODO: put whitespace back in
	"sep": p.Text(","),
	// literals
	"stringLit": p.SeqV(
		p.Text(`"`),
		// TODO: escaping, more chars than alphanum
		// regex: \"(\\.|[^"\\])*\"
		p.RepSep(p.Ref("alphaNum"), p.Succeed),
		p.Text(`"`),
	),
	"numberLit": p.RepSep(p.Ref("digit"), p.Succeed),
	"bool":      p.ChoiceV(p.Text("true"), p.Text("false")),
	"null":      p.Text("null"),
	// whitespace
	"optWhitespace": p.ListRule("whitespace", "optWhitespace", p.Succeed),
	"whitespace":    p.ChoiceV(p.Text(" "), p.Text("\n"), p.Text("\t")),
	"digit":         p.RuneRange('0', '9'),
	"alphaNum": p.ChoiceV(
		p.RuneRange('a', 'z'),
		p.RuneRange('A', 'Z'),
		p.Ref("digit"),
	),
}
