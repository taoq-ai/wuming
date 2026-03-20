package replacer

import (
	"fmt"
	"sort"

	"github.com/taoq-ai/wuming/domain/model"
	"github.com/taoq-ai/wuming/domain/port"
)

// Consistent wraps any Replacer to ensure the same PII value always maps to
// the same replacement string within a single Replace call. For example, if
// "john@example.com" appears three times it will be replaced with "[EMAIL_1]"
// in all three positions (when wrapping a Redact replacer).
//
// The mapping is built per Replace call—callers who need cross-call consistency
// should reuse the same Consistent instance and call [Consistent.Reset] between
// unrelated texts (or simply create a new one).
type Consistent struct {
	inner port.Replacer
	// seen maps "TYPE:value" → replacement string across calls.
	seen map[string]string
	// counters tracks the next index per PIIType for Redact-style numbering.
	counters map[model.PIIType]int
}

// NewConsistent wraps the given replacer with consistent-replacement behavior.
func NewConsistent(inner port.Replacer) *Consistent {
	return &Consistent{
		inner:    inner,
		seen:     make(map[string]string),
		counters: make(map[model.PIIType]int),
	}
}

func (c *Consistent) Name() string { return "consistent(" + c.inner.Name() + ")" }

// Reset clears the value→replacement mapping so subsequent Replace calls start
// fresh. Useful when processing unrelated texts with the same instance.
func (c *Consistent) Reset() {
	c.seen = make(map[string]string)
	c.counters = make(map[model.PIIType]int)
}

func (c *Consistent) Replace(text string, matches []model.Match) (string, error) {
	if len(matches) == 0 {
		return text, nil
	}

	// Sort by start position ascending so we assign stable counter values
	// (first occurrence gets _1, second unique value gets _2, etc.).
	sorted := make([]model.Match, len(matches))
	copy(sorted, matches)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Start < sorted[j].Start
	})

	// Build the replacement for each unique value. We delegate to the inner
	// replacer for a single synthetic match to obtain its replacement string,
	// then cache it.
	for _, m := range sorted {
		key := cacheKey(m)
		if _, ok := c.seen[key]; ok {
			continue
		}
		rep, err := c.replacementFor(m)
		if err != nil {
			return "", err
		}
		c.seen[key] = rep
	}

	// Apply replacements in reverse order to preserve byte offsets.
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Start > sorted[j].Start
	})

	result := []byte(text)
	for _, m := range sorted {
		if m.Start < 0 || m.End > len(result) || m.Start >= m.End {
			continue
		}
		rep := c.seen[cacheKey(m)]
		result = append(result[:m.Start], append([]byte(rep), result[m.End:]...)...)
	}
	return string(result), nil
}

// replacementFor generates the replacement string for a match. For the
// built-in Redact replacer it produces numbered placeholders like [EMAIL_1].
// For all other replacers it delegates to the inner replacer with a synthetic
// single-match call.
func (c *Consistent) replacementFor(m model.Match) (string, error) {
	if r, ok := c.inner.(*Redact); ok {
		c.counters[m.Type]++
		return fmt.Sprintf(
			numberFormat(r.Format),
			m.Type.String(),
			c.counters[m.Type],
		), nil
	}

	// Generic path: ask the inner replacer to replace a tiny synthetic text.
	synthetic := m.Value
	syntheticMatch := model.Match{
		Type:       m.Type,
		Value:      m.Value,
		Start:      0,
		End:        len(synthetic),
		Confidence: m.Confidence,
		Locale:     m.Locale,
		Detector:   m.Detector,
	}
	return c.inner.Replace(synthetic, []model.Match{syntheticMatch})
}

// cacheKey returns a string that uniquely identifies a PII value within a type.
func cacheKey(m model.Match) string {
	return m.Type.String() + ":" + m.Value
}

// numberFormat converts a format like "[%s]" to "[%s_%d]".
func numberFormat(f string) string {
	if len(f) == 0 {
		return "[%s_%d]"
	}
	// Insert _%d before the closing character (e.g. "]").
	return f[:len(f)-1] + "_%d" + f[len(f)-1:]
}
