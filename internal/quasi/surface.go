package quasi

import "strings"

// SourcePosition identifies a 1-based source location.
type SourcePosition struct {
	File   string
	Line   int
	Column int
}

// SourceSpan covers a half-open interval in one source file.
type SourceSpan struct {
	Start SourcePosition
	End   SourcePosition
}

// RawLine is the parser-side representation of one captured quasi line.
type RawLine struct {
	Text         string
	TrimmedText  string
	IndentColumn int
	Span         SourceSpan
}

// Block is a raw quasi body captured by the host parser.
type Block struct {
	StyleName string
	Lines     []RawLine
	Span      SourceSpan
}

// Diagnostic is a lightweight validation finding for quasi surface checks.
type Diagnostic struct {
	Code    string
	Message string
	Span    SourceSpan
}

// Result contains the outcome of surface validation.
type Result struct {
	Diagnostics []Diagnostic
}

// OK reports whether the validation succeeded without diagnostics.
func (r Result) OK() bool {
	return len(r.Diagnostics) == 0
}

// Validator validates raw quasi blocks against a compiled profile.
type Validator struct {
	Profile Profile
}

// Validate performs surface validation over a raw captured quasi block.
// This is intentionally lexical and structural; it does not interpret meaning.
func (v Validator) Validate(block Block) Result {
	var result Result
	var stack []frame
	var pendingChild *frame
	var recentClosed []frame

	for i := range block.Lines {
		line := block.Lines[i]
		if line.TrimmedText == "" {
			if !v.Profile.Layout.AllowBlankLines {
				result.Diagnostics = append(result.Diagnostics, Diagnostic{
					Code:    "quasi.blank-line-forbidden",
					Message: "blank lines are not allowed in this quasi style",
					Span:    line.Span,
				})
			}
			continue
		}

		keyword := leadingKeyword(line.TrimmedText)
		if keyword == "" && v.Profile.Validation.FirstTokenMustBeKeyword {
			result.Diagnostics = append(result.Diagnostics, Diagnostic{
				Code:    "quasi.keyword.missing",
				Message: "quasi line must start with a recognized keyword",
				Span:    line.Span,
			})
			continue
		}

		rules := v.Profile.RulesForKeyword(keyword)
		if len(rules) == 0 {
			result.Diagnostics = append(result.Diagnostics, Diagnostic{
				Code:    "quasi.keyword.unknown",
				Message: "keyword is not allowed by the resolved quasi style: " + keyword,
				Span:    line.Span,
			})
			continue
		}

		matched, rule := matchRule(rules, line.TrimmedText)
		if !matched {
			result.Diagnostics = append(result.Diagnostics, Diagnostic{
				Code:    "quasi.rule.no-match",
				Message: "line does not match any rule for keyword: " + keyword,
				Span:    line.Span,
			})
			continue
		}

		stack, recentClosed = closeFrames(stack, line.IndentColumn)
		result.Diagnostics = append(result.Diagnostics, validateIndentAndContinuation(v.Profile, stack, recentClosed, pendingChild, line, rule)...)
		pendingChild = nil
		if rule.OpensBlock {
			opened := frame{RuleID: rule.ID, IndentColumn: line.IndentColumn, LineIndex: i}
			stack = append(stack, opened)
			pendingChild = &opened
			recentClosed = nil
			continue
		}
		recentClosed = nil
	}

	return result
}

type frame struct {
	RuleID       string
	IndentColumn int
	LineIndex    int
}

func closeFrames(stack []frame, indentColumn int) ([]frame, []frame) {
	closed := make([]frame, 0)
	for len(stack) > 0 && indentColumn <= stack[len(stack)-1].IndentColumn {
		closed = append(closed, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	return stack, closed
}

func validateIndentAndContinuation(profile Profile, stack []frame, recentClosed []frame, pendingChild *frame, line RawLine, rule Rule) []Diagnostic {
	var diagnostics []Diagnostic

	if profile.Validation.RequireConsistentIndentation && profile.Layout.Indentation.Width > 0 {
		if line.IndentColumn%profile.Layout.Indentation.Width != 0 {
			diagnostics = append(diagnostics, Diagnostic{
				Code:    "quasi.indent.invalid",
				Message: "line indentation does not match style indentation width",
				Span:    line.Span,
			})
		}
	}

	if pendingChild != nil && profile.Validation.RequireIndentAfterBlockOpener {
		if line.IndentColumn <= pendingChild.IndentColumn {
			diagnostics = append(diagnostics, Diagnostic{
				Code:    "quasi.block.child-indent",
				Message: "block opener must be followed by a more-indented child line",
				Span:    line.Span,
			})
		}
	}

	if rule.Kind == RuleKindContinuation {
		parent, ok := continuationParent(stack, recentClosed, line.IndentColumn)
		if !ok {
			diagnostics = append(diagnostics, Diagnostic{
				Code:    "quasi.continuation.orphan",
				Message: "continuation line has no compatible parent block",
				Span:    line.Span,
			})
			return diagnostics
		}
		if rule.MustAlignWithParent && line.IndentColumn != parent.IndentColumn {
			diagnostics = append(diagnostics, Diagnostic{
				Code:    "quasi.continuation.misaligned",
				Message: "continuation line must align with its parent block",
				Span:    line.Span,
			})
		}
		if profile.Validation.RequireContinuationImmediatelyAfterParentBlock && len(recentClosed) == 0 {
			diagnostics = append(diagnostics, Diagnostic{
				Code:    "quasi.continuation.not-immediate",
				Message: "continuation line must appear immediately after its parent block",
				Span:    line.Span,
			})
		}
	}

	return diagnostics
}

func continuationParent(stack []frame, recentClosed []frame, indentColumn int) (frame, bool) {
	for _, parent := range recentClosed {
		if parent.IndentColumn == indentColumn {
			return parent, true
		}
	}
	for i := len(stack) - 1; i >= 0; i-- {
		if stack[i].IndentColumn == indentColumn {
			return stack[i], true
		}
	}
	return frame{}, false
}

func leadingKeyword(line string) string {
	line = strings.TrimSpace(line)
	if line == "" {
		return ""
	}

	if strings.HasPrefix(line, "else if ") || line == "else if:" {
		return "else if"
	}
	if strings.HasPrefix(line, "else:") {
		return "else"
	}

	parts := strings.Fields(line)
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

func matchRule(rules []Rule, line string) (bool, Rule) {
	for _, rule := range rules {
		if rule.Regexp != nil && rule.Regexp.MatchString(line) {
			return true, rule
		}
	}
	return false, Rule{}
}
