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

type choice struct {
	choices []Rule
}

var _ Rule = &choice{}

func Choice(choices []Rule) *choice {
	return &choice{
		choices: choices,
	}
}

func (c *choice) String() string {
	choicesStrs := make([]string, len(c.choices))
	for idx, choice := range c.choices {
		choicesStrs[idx] = choice.String()
	}
	return fmt.Sprintf("(%s)", strings.Join(choicesStrs, " | "))
}

func (c *choice) Validate(g *Grammar) error {
	for idx, choice := range c.choices {
		if err := choice.Validate(g); err != nil {
			return fmt.Errorf("in choice %d: %v", idx, err)
		}
	}
	return nil
}

func (c *choice) Children() []Rule {
	return c.choices
}

// sequence

type sequence struct {
	items []Rule
}

var _ Rule = &sequence{}

func Sequence(items []Rule) *sequence {
	return &sequence{
		items: items,
	}
}

func (s *sequence) String() string {
	itemsStrs := make([]string, len(s.items))
	for idx, item := range s.items {
		itemsStrs[idx] = item.String()
	}
	return fmt.Sprintf("[%s]", strings.Join(itemsStrs, ", "))
}

func (s *sequence) Validate(g *Grammar) error {
	for idx, item := range s.items {
		if err := item.Validate(g); err != nil {
			return fmt.Errorf("in seq item %d: %v", idx, err)
		}
	}
	return nil
}

func (s *sequence) Children() []Rule {
	return s.items
}

// keyword

type keyword struct {
	value string
}

var _ Rule = &keyword{}

// TODO: case insensitivity
func Keyword(value string) *keyword {
	return &keyword{
		value: value,
	}
}

func (k *keyword) String() string {
	return fmt.Sprintf(`"%s"`, k.value)
}

func (k *keyword) Validate(_ *Grammar) error {
	for _, char := range k.value {
		if char == '\n' {
			return fmt.Errorf("newlines not allowed in keywords: %v", k.value)
		}
	}
	return nil
}

func (k *keyword) Children() []Rule { return []Rule{} }

// Rule ref

type ref struct {
	name string
}

var _ Rule = &ref{}

func Ref(name string) *ref {
	return &ref{
		name: name,
	}
}

func (r *ref) String() string {
	return string(r.name)
}

func (r *ref) Validate(g *Grammar) error {
	if _, ok := g.rules[r.name]; !ok {
		return fmt.Errorf(`ref not found: "%s"`, r.name)
	}
	return nil
}

func (r *ref) Children() []Rule { return []Rule{} }

// regex

type regex struct {
	regex *regexp.Regexp
}

var _ Rule = &regex{}

func Regex(re *regexp.Regexp) *regex {
	return &regex{
		regex: re,
	}
}

func (r *regex) String() string {
	return fmt.Sprintf("/%s/", r.regex.String())
}

func (r *regex) Validate(g *Grammar) error {
	return nil
}

func (r *regex) Children() []Rule { return []Rule{} }

// Succeed

var Succeed = &succeed{}

type succeed struct{}

var _ Rule = &succeed{}

func (s *succeed) String() string {
	return "<succeed>"
}

func (s *succeed) Validate(g *Grammar) error {
	return nil
}

func (s *succeed) Children() []Rule { return []Rule{} }
