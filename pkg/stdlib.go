package parserlib

import (
	"regexp"

	pp "github.com/vilterp/go-pretty-print"
)

func ListRule1(ruleName string, listName string, sep Rule) Rule {
	return Choice([]Rule{
		Sequence([]Rule{
			Ref(ruleName),
			sep,
			Ref(listName),
		}),
		Ref(ruleName),
	})
}

func ListRule(ruleName string, listName string, sep Rule) Rule {
	return Opt(ListRule1(ruleName, listName, sep))
}

func Opt(r Rule) Rule {
	return &ChoiceRule{
		choices: []Rule{
			r,
			Succeed,
		},
	}
}

var OptWhitespace = Opt(Whitespace)

func WhitespaceSeq(items []Rule) Rule {
	// hoo, a generic intercalate function sure would be nice
	var outItems []Rule
	for idx, item := range items {
		if idx > 0 {
			outItems = append(outItems, Whitespace)
		}
		outItems = append(outItems, item)
	}
	return &SeqRule{
		items: outItems,
	}
}

func OptWhitespaceSurround(r Rule) Rule {
	return Sequence([]Rule{
		OptWhitespace,
		r,
		OptWhitespace,
	})
}

func Block(start string, doc pp.Doc, end string) pp.Doc {
	return SeqV(
		pp.Text(start), pp.Newline,
		pp.Indent(2, doc),
		pp.Newline, pp.Text(end),
	)
}

func SeqV(docs ...pp.Doc) pp.Doc {
	return pp.Seq(docs)
}

var Whitespace = Regex(regexp.MustCompile("\\s+"))

var CommaOptWhitespace = Sequence([]Rule{Keyword(","), OptWhitespace})

var UnsignedIntLit = Regex(regexp.MustCompile("[0-9]+"))

var SignedIntLit = Regex(regexp.MustCompile("-?[0-9]+"))

// Thank you https://stackoverflow.com/a/2039820
var StringLit = Regex(regexp.MustCompile(`\"(\\.|[^"\\])*\"`))

var Ident = Regex(regexp.MustCompile("[a-zA-Z_][a-zA-Z0-9_]*"))
