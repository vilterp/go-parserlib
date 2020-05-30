package parserlib

func (tt *TraceTree) GetCompletions() ([]string, error) {
	rule := tt.Rule
	switch tRule := rule.(type) {
	case *ChoiceRule:
		// TODO: we sometimes want to return multiple choices here...
		// maybe only if we're on the left edge
		if tt.CursorPos == 0 {
			return tRule.Completions(tt.grammar, tt.CursorPos), nil
		}
		return tt.ChoiceTrace.GetCompletions()
	case *SeqRule:
		return tt.ItemTraces[tt.AtItemIdx].GetCompletions()
	case *TextRule:
		if tt.CursorPos == 0 {
			return []string{tRule.value}, nil
		}
		return []string{}, nil
	case *RefRule:
		return tt.RefTrace.GetCompletions()
	default:
		return []string{}, nil
	}
}

func (c *ChoiceRule) Completions(g *Grammar, cursor int) []string {
	var out []string
	for _, choice := range c.choices {
		out = append(out, choice.Completions(g, cursor)...)
	}
	return out
}

func (s *SeqRule) Completions(_ *Grammar, _ int) []string {
	// TODO: which index are we at? maybe a rule method
	// is the wrong way to do this
	return []string{}
}

func (k *TextRule) Completions(_ *Grammar, _ int) []string {
	return []string{k.value}
}

func (r *RefRule) Completions(g *Grammar, cursor int) []string {
	rule := g.rules[r.Name]
	return rule.Completions(g, cursor)
}

func (s *SucceedRule) Completions(_ *Grammar, _ int) []string {
	return []string{}
}

func (n *NamedRule) Completions(g *Grammar, cursor int) []string {
	return n.inner.Completions(g, cursor)
}
