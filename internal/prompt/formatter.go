package prompt

import (
	"strings"
)

// Formatter handles section-specific formatting and limits
type Formatter struct {
	maxSectionChars int
}

func (f *Formatter) FormatFacts(facts []string) string {
	return f.limitChars(strings.Join(facts, "\n"), f.maxSectionChars)
}

func (f *Formatter) limitChars(input string, limit int) string {
	if len(input) > limit {
		return input[:limit] + "..."
	}
	return input
}
