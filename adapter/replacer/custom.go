package replacer

import (
	"sort"

	"github.com/taoq-ai/wuming/domain/model"
)

// ReplaceFunc is a user-defined function that takes a match and returns
// the replacement string.
type ReplaceFunc func(model.Match) string

// Custom replaces each match using a user-defined function.
type Custom struct {
	fn   ReplaceFunc
	name string
}

// NewCustom creates a Custom replacer with the given name and function.
func NewCustom(name string, fn ReplaceFunc) *Custom {
	return &Custom{name: name, fn: fn}
}

func (c *Custom) Name() string { return c.name }

func (c *Custom) Replace(text string, matches []model.Match) (string, error) {
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
		replacement := c.fn(m)
		result = append(result[:m.Start], append([]byte(replacement), result[m.End:]...)...)
	}
	return string(result), nil
}
