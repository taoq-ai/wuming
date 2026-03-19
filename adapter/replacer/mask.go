package replacer

import (
	"sort"
	"strings"

	"github.com/taoq-ai/wuming/domain/model"
)

// Mask replaces characters with a mask character, preserving the last N characters.
type Mask struct {
	// Char is the masking character. Defaults to '*'.
	Char rune
	// Preserve is the number of trailing characters to keep visible. Defaults to 4.
	Preserve int
}

// NewMask creates a Mask replacer with default settings (* and preserve last 4).
func NewMask() *Mask {
	return &Mask{Char: '*', Preserve: 4}
}

func (m *Mask) Name() string { return "mask" }

func (m *Mask) Replace(text string, matches []model.Match) (string, error) {
	if len(matches) == 0 {
		return text, nil
	}

	sorted := make([]model.Match, len(matches))
	copy(sorted, matches)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Start > sorted[j].Start
	})

	result := []byte(text)
	for _, match := range sorted {
		if match.Start < 0 || match.End > len(result) || match.Start >= match.End {
			continue
		}
		value := match.Value
		runes := []rune(value)
		preserve := m.Preserve
		if preserve > len(runes) {
			preserve = len(runes)
		}
		maskLen := len(runes) - preserve
		masked := strings.Repeat(string(m.Char), maskLen) + string(runes[maskLen:])
		result = append(result[:match.Start], append([]byte(masked), result[match.End:]...)...)
	}
	return string(result), nil
}
