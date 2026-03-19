package port

import "github.com/taoq-ai/wuming/domain/model"

// Replacer substitutes detected PII matches with redacted values.
type Replacer interface {
	// Replace takes the original text and a set of matches, and returns
	// the text with all matches substituted.
	Replace(text string, matches []model.Match) (string, error)
	// Name returns a unique identifier for this replacer strategy.
	Name() string
}
