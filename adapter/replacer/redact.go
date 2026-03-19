// Package replacer provides concrete Replacer implementations for different
// PII substitution strategies.
package replacer

import (
	"fmt"
	"sort"

	"github.com/taoq-ai/wuming/domain/model"
)

// Redact replaces each match with a placeholder like [EMAIL] or [NATIONAL_ID].
type Redact struct {
	// Format is the format string for the placeholder. It receives the PIIType
	// string as its argument. Defaults to "[%s]".
	Format string
}

// NewRedact creates a Redact replacer with the default "[%s]" format.
func NewRedact() *Redact {
	return &Redact{Format: "[%s]"}
}

func (r *Redact) Name() string { return "redact" }

func (r *Redact) Replace(text string, matches []model.Match) (string, error) {
	if len(matches) == 0 {
		return text, nil
	}

	sorted := make([]model.Match, len(matches))
	copy(sorted, matches)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Start > sorted[j].Start
	})

	result := []byte(text)
	for _, m := range sorted {
		if m.Start < 0 || m.End > len(result) || m.Start >= m.End {
			continue
		}
		placeholder := fmt.Sprintf(r.Format, m.Type.String())
		result = append(result[:m.Start], append([]byte(placeholder), result[m.End:]...)...)
	}
	return string(result), nil
}
