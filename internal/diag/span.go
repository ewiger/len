package diag

// Position identifies a 1-based source location.
type Position struct {
	File   string
	Line   int
	Column int
}

// Span covers a half-open source interval.
type Span struct {
	Start Position
	End   Position
}

// IsZero reports whether the span is unset.
func (s Span) IsZero() bool {
	return s.Start == (Position{}) && s.End == (Position{})
}
