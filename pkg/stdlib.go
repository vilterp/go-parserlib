package parserlib

func ListRule1(ruleName string, listName string, sep Rule) Rule {
	return Choice([]Rule{
		Seq([]Rule{
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

var Whitespace = ChoiceV(Text(" "), Text("\n"), Text("\t"))

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

func Block(start string, inner Rule, end string) Rule {
	return Seq([]Rule{
		Text(start),
		OptWhitespace,
		inner,
		Seq([]Rule{
			OptWhitespace,
			Text(end),
		}),
	})
}

func OptWhitespaceSurround(r Rule) Rule {
	return Seq([]Rule{
		OptWhitespace,
		r,
		OptWhitespace,
	})
}

var CommaOptWhitespace = Seq([]Rule{Text(","), OptWhitespace})
