package diag

import (
	"fmt"
	"sort"
)

// Severity distinguishes parser and validator outcomes.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

// Diagnostic is a stable user-facing finding.
type Diagnostic struct {
	Code     string
	Message  string
	Severity Severity
	Span     Span
}

func (d Diagnostic) String() string {
	if d.Span.IsZero() {
		return fmt.Sprintf("%s: %s (%s)", d.Severity, d.Message, d.Code)
	}
	return fmt.Sprintf(
		"%s:%d:%d: %s: %s (%s)",
		d.Span.Start.File,
		d.Span.Start.Line,
		d.Span.Start.Column,
		d.Severity,
		d.Message,
		d.Code,
	)
}

// Sort orders diagnostics deterministically for tests and CLI output.
func Sort(items []Diagnostic) {
	sort.Slice(items, func(i, j int) bool {
		left := items[i]
		right := items[j]
		if left.Span.Start.File != right.Span.Start.File {
			return left.Span.Start.File < right.Span.Start.File
		}
		if left.Span.Start.Line != right.Span.Start.Line {
			return left.Span.Start.Line < right.Span.Start.Line
		}
		if left.Span.Start.Column != right.Span.Start.Column {
			return left.Span.Start.Column < right.Span.Start.Column
		}
		if left.Code != right.Code {
			return left.Code < right.Code
		}
		return left.Message < right.Message
	})
}
