package parserlib

import (
	"fmt"
	"strconv"
	"strings"

	pp "github.com/vilterp/go-pretty-print"
)

type RuleID int

type Grammar struct {
	StartRule string
	rules     map[string]Rule

	idForRule  map[Rule]RuleID
	ruleForID  map[RuleID]Rule
	nameForID  map[RuleID]string
	nextRuleID RuleID
}

func NewGrammar(rules map[string]Rule, startRule string) (*Grammar, error) {
	g := &Grammar{
		StartRule: startRule,
		rules:     rules,
		idForRule: make(map[Rule]RuleID),
		ruleForID: make(map[RuleID]Rule),
		nameForID: make(map[RuleID]string),
		// prevent zero value from accidentally making things work that shouldn't
		nextRuleID: 1,
	}
	if err := g.validate(); err != nil {
		return nil, err
	}
	for name, rule := range rules {
		id := g.assignRuleIDs(rule)
		g.nameForID[id] = name
	}
	return g, nil
}

func (g *Grammar) assignRuleIDs(r Rule) RuleID {
	id := g.nextRuleID
	g.idForRule[r] = id
	g.ruleForID[id] = r
	g.nextRuleID++
	for _, child := range r.Children() {
		g.assignRuleIDs(child)
	}
	return id
}

func (g *Grammar) validate() error {
	for ruleName, rule := range g.rules {
		if err := rule.Validate(g); err != nil {
			return fmt.Errorf(`in rule "%s": %v`, ruleName, err)
		}
	}
	if _, ok := g.rules[g.StartRule]; !ok {
		return fmt.Errorf("start rule %s not found", g.StartRule)
	}
	return nil
}

func (g *Grammar) Format() pp.Doc {
	var ruleDocs []pp.Doc
	for name, rule := range g.rules {
		ruleDocs = append(ruleDocs, pp.Textf("%s: %s", name, rule.String()))
	}
	return pp.Join(ruleDocs, pp.Newline)
}

func (g *Grammar) NameForID(id RuleID) string {
	return g.nameForID[id]
}

type Rule interface {
	String() string
	Validate(g *Grammar) error
	Completions(g *Grammar, cursor int) []string
	Children() []Rule
	Serialize(g *Grammar) SerializedRule
}

// choice

type ChoiceRule struct {
	choices []Rule
}

var _ Rule = &ChoiceRule{}

func Choice(choices []Rule) *ChoiceRule {
	return &ChoiceRule{
		choices: choices,
	}
}

func ChoiceV(choices ...Rule) *ChoiceRule {
	return Choice(choices)
}

func (c *ChoiceRule) String() string {
	choicesStrs := make([]string, len(c.choices))
	for idx, choice := range c.choices {
		choicesStrs[idx] = choice.String()
	}
	return fmt.Sprintf("(%s)", strings.Join(choicesStrs, " | "))
}

func (c *ChoiceRule) Validate(g *Grammar) error {
	for idx, choice := range c.choices {
		if err := choice.Validate(g); err != nil {
			return fmt.Errorf("in choice %d: %v", idx, err)
		}
	}
	return nil
}

func (c *ChoiceRule) Children() []Rule {
	return c.choices
}

// sequence

type SeqRule struct {
	items []Rule
}

var _ Rule = &SeqRule{}

func Seq(items []Rule) *SeqRule {
	return &SeqRule{
		items: items,
	}
}

func SeqV(items ...Rule) *SeqRule {
	return Seq(items)
}

func (s *SeqRule) String() string {
	itemsStrs := make([]string, len(s.items))
	for idx, item := range s.items {
		itemsStrs[idx] = item.String()
	}
	return fmt.Sprintf("[%s]", strings.Join(itemsStrs, ", "))
}

func (s *SeqRule) Validate(g *Grammar) error {
	for idx, item := range s.items {
		if err := item.Validate(g); err != nil {
			return fmt.Errorf("in seq item %d: %v", idx, err)
		}
	}
	return nil
}

func (s *SeqRule) Children() []Rule {
	return s.items
}

// Text

type TextRule struct {
	value string
}

var _ Rule = &TextRule{}

// TODO: case insensitivity
func Text(value string) *TextRule {
	return &TextRule{
		value: value,
	}
}

func (k *TextRule) String() string {
	return strconv.Quote(k.value)
}

func (k *TextRule) Validate(_ *Grammar) error {
	// TODO: I forget why this was necessary
	//for _, char := range k.value {
	//	if char == '\n' {
	//		return fmt.Errorf("newlines not allowed in texts: %v", k.value)
	//	}
	//}
	return nil
}

func (k *TextRule) Children() []Rule { return []Rule{} }

// Rule ref

type RefRule struct {
	Name string
}

var _ Rule = &RefRule{}

func Ref(name string) *RefRule {
	return &RefRule{
		Name: name,
	}
}

func (r *RefRule) String() string {
	return r.Name
}

func (r *RefRule) Validate(g *Grammar) error {
	if _, ok := g.rules[r.Name]; !ok {
		return fmt.Errorf(`ref not found: "%s"`, r.Name)
	}
	return nil
}

func (r *RefRule) Children() []Rule { return []Rule{} }

// named

type NamedRule struct {
	name  string
	inner Rule
}

var _ Rule = &NamedRule{}

func Named(name string, rule Rule) Rule {
	return &NamedRule{
		name:  name,
		inner: rule,
	}
}

func (n *NamedRule) String() string {
	return fmt.Sprintf("NAMED(%s, %s)", n.name, n.inner.String())
}

func (n *NamedRule) Validate(g *Grammar) error {
	return n.inner.Validate(g)
}

func (n *NamedRule) Children() []Rule {
	return []Rule{n.inner}
}

// repsep

type RepSepRule struct {
	Rep Rule
	Sep Rule
}

var _ Rule = &RepSepRule{}

func RepSep(rep Rule, sep Rule) *RepSepRule {
	return &RepSepRule{
		Rep: rep,
		Sep: sep,
	}
}

func (r *RepSepRule) String() string {
	return fmt.Sprintf("RepSep(%v, %v)", r.Rep.String(), r.Sep.String())
}

func (r *RepSepRule) Validate(g *Grammar) error {
	if err := r.Rep.Validate(g); err != nil {
		return err
	}
	if err := r.Sep.Validate(g); err != nil {
		return err
	}
	return nil
}

func (r *RepSepRule) Completions(g *Grammar, cursor int) []string {
	return nil
}

func (r *RepSepRule) Children() []Rule {
	return []Rule{r.Rep, r.Sep}
}

func (r *RepSepRule) Serialize(g *Grammar) SerializedRule {
	panic("implement me")
}

// rune range

type RuneRangeRule struct {
	from rune
	to   rune
}

var _ Rule = &RuneRangeRule{}

func RuneRange(from rune, to rune) *RuneRangeRule {
	return &RuneRangeRule{
		from: from,
		to:   to,
	}
}

func (r *RuneRangeRule) String() string {
	return fmt.Sprintf("[%c-%c]", r.from, r.to)
}

func (r *RuneRangeRule) Validate(g *Grammar) error {
	if r.from >= r.to {
		return fmt.Errorf(
			"from %s not less than to %s",
			strconv.QuoteRune(r.from), strconv.QuoteRune(r.to),
		)
	}
	return nil
}

func (r *RuneRangeRule) Completions(g *Grammar, cursor int) []string {
	return nil
}

func (r *RuneRangeRule) Children() []Rule {
	return nil
}

func (r *RuneRangeRule) Serialize(g *Grammar) SerializedRule {
	panic("implement me")
}

// Succeed

var Succeed = &SucceedRule{}

type SucceedRule struct{}

var _ Rule = &SucceedRule{}

func (s *SucceedRule) String() string {
	return "<succeed>"
}

func (s *SucceedRule) Validate(g *Grammar) error {
	return nil
}

func (s *SucceedRule) Children() []Rule { return []Rule{} }
