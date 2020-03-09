package parserlib

func ExtractList(tree *RuleNode, listRule string, singularRule string) ([]*RuleNode, error) {
	return extractListRecur(nil, tree, listRule, singularRule)
}

func extractListRecur(
	out []*RuleNode, cur *RuleNode, listRule string, singularRule string,
) ([]*RuleNode, error) {
	child := cur.GetChildWithName(singularRule)
	if child == nil {
		return out, nil
	}
	out = append(out, child)
	rest := cur.GetChildWithName(listRule)
	if rest == nil {
		return out, nil
	}
	return extractListRecur(out, rest, listRule, singularRule)
}
