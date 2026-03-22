package quasi

import (
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

// Profile describes one style-specific quasi surface-validation contract.
type Profile struct {
	Version    int             `yaml:"version"`
	Style      StyleMeta       `yaml:"style"`
	Layout     LayoutRules     `yaml:"layout"`
	Keywords   KeywordGroups   `yaml:"keywords"`
	Slots      map[string]Slot `yaml:"slots"`
	Rules      []Rule          `yaml:"rules"`
	Validation Validation      `yaml:"validation"`
}

// StyleMeta names the profile and describes its provenance.
type StyleMeta struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	DerivedFrom []string `yaml:"derivedFrom"`
}

// LayoutRules constrain indentation and blank-line behavior.
type LayoutRules struct {
	Indentation         Indentation `yaml:"indentation"`
	AllowBlankLines     bool        `yaml:"allowBlankLines"`
	AllowNarrativeLines bool        `yaml:"allowNarrativeLines"`
}

// Indentation constrains how raw block indentation is interpreted.
type Indentation struct {
	Mode  string `yaml:"mode"`
	Width int    `yaml:"width"`
}

// KeywordGroups provides cheap first-token filtering before regex matching.
type KeywordGroups struct {
	BlockOpeners       []string `yaml:"blockOpeners"`
	BlockContinuations []string `yaml:"blockContinuations"`
	SimpleStatements   []string `yaml:"simpleStatements"`
}

// Slot defines an optionally named sub-pattern used by rule captures.
type Slot struct {
	Pattern string         `yaml:"pattern"`
	Regexp  *regexp.Regexp `yaml:"-"`
}

// RuleKind identifies how a rule participates in structural validation.
type RuleKind string

const (
	RuleKindStatement    RuleKind = "statement"
	RuleKindBlock        RuleKind = "block"
	RuleKindContinuation RuleKind = "continuation"
)

// Rule matches one complete quasi line shape.
type Rule struct {
	ID                  string            `yaml:"id"`
	Kind                RuleKind          `yaml:"kind"`
	Keyword             string            `yaml:"keyword"`
	Pattern             string            `yaml:"pattern"`
	Captures            map[string]string `yaml:"captures"`
	OpensBlock          bool              `yaml:"opensBlock"`
	ChildPolicy         string            `yaml:"childPolicy"`
	Continuations       []string          `yaml:"continuations"`
	AttachesTo          []string          `yaml:"attachesTo"`
	MustAlignWithParent bool              `yaml:"mustAlignWithParent"`
	Regexp              *regexp.Regexp    `yaml:"-"`
}

// Validation controls style-level acceptance and rejection behavior.
type Validation struct {
	FirstTokenMustBeKeyword                        bool   `yaml:"firstTokenMustBeKeyword"`
	UnknownKeywordPolicy                           string `yaml:"unknownKeywordPolicy"`
	NonMatchingLinePolicy                          string `yaml:"nonMatchingLinePolicy"`
	RequireContinuationImmediatelyAfterParentBlock bool   `yaml:"requireContinuationImmediatelyAfterParentBlock"`
	RequireConsistentIndentation                   bool   `yaml:"requireConsistentIndentation"`
	RequireIndentAfterBlockOpener                  bool   `yaml:"requireIndentAfterBlockOpener"`
	RequireReturnExpression                        bool   `yaml:"requireReturnExpression"`
}

// Compile prepares slot and rule regular expressions for reuse.
func (p *Profile) Compile() error {
	if p.Slots == nil {
		p.Slots = map[string]Slot{}
	}
	for name, slot := range p.Slots {
		re, err := regexp.Compile(slot.Pattern)
		if err != nil {
			return &CompileError{Kind: "slot", Name: name, Pattern: slot.Pattern, Err: err}
		}
		slot.Regexp = re
		p.Slots[name] = slot
	}

	for i := range p.Rules {
		re, err := regexp.Compile(p.Rules[i].Pattern)
		if err != nil {
			return &CompileError{Kind: "rule", Name: p.Rules[i].ID, Pattern: p.Rules[i].Pattern, Err: err}
		}
		p.Rules[i].Regexp = re
	}

	return nil
}

// LoadProfile reads and compiles a YAML style profile from disk.
func LoadProfile(path string) (Profile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Profile{}, err
	}

	var profile Profile
	if err := yaml.Unmarshal(data, &profile); err != nil {
		return Profile{}, err
	}
	if err := profile.Compile(); err != nil {
		return Profile{}, err
	}
	return profile, nil
}

// RulesForKeyword returns the ordered rule subset that can match the keyword.
func (p *Profile) RulesForKeyword(keyword string) []Rule {
	matches := make([]Rule, 0, len(p.Rules))
	for _, rule := range p.Rules {
		if rule.Keyword == keyword {
			matches = append(matches, rule)
		}
	}
	return matches
}

// CompileError reports YAML profile regex issues before validation starts.
type CompileError struct {
	Kind    string
	Name    string
	Pattern string
	Err     error
}

func (e *CompileError) Error() string {
	return "quasi profile compile error: " + e.Kind + " " + e.Name + " pattern " + e.Pattern + ": " + e.Err.Error()
}
