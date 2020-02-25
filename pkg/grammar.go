package parserlib

import (
	"fmt"
	"regexp"
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

func Sequence(items []Rule) *SeqRule {
	return &SeqRule{
		items: items,
	}
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

// keyword

type KeywordRule struct {
	value string
}

var _ Rule = &KeywordRule{}

// TODO: case insensitivity
func Keyword(value string) *KeywordRule {
	return &KeywordRule{
		value: value,
	}
}

func (k *KeywordRule) String() string {
	return fmt.Sprintf(`"%s"`, k.value)
}

func (k *KeywordRule) Validate(_ *Grammar) error {
	for _, char := range k.value {
		if char == '\n' {
			return fmt.Errorf("newlines not allowed in keywords: %v", k.value)
		}
	}
	return nil
}

func (k *KeywordRule) Children() []Rule { return []Rule{} }

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

// regex

type RegexRule struct {
	regex *regexp.Regexp
}

var _ Rule = &RegexRule{}

func Regex(re *regexp.Regexp) *RegexRule {
	return &RegexRule{
		regex: re,
	}
}

func (r *RegexRule) String() string {
	return fmt.Sprintf("/%s/", r.regex.String())
}

func (r *RegexRule) Validate(g *Grammar) error {
	return nil
}

func (r *RegexRule) Children() []Rule { return []Rule{} }

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
