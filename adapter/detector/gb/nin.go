package gb

import (
	"context"
	"regexp"
	"strings"

	"github.com/taoq-ai/wuming/domain/model"
)

// NIN format: 2 letters, 6 digits, 1 letter (with optional spaces).
var ninRe = regexp.MustCompile(`(?i)\b([A-Z]{2})\s?(\d{2})\s?(\d{2})\s?(\d{2})\s?([A-D])\b`)

// invalidNINPrefixes lists two-letter prefixes that are never used.
var invalidNINPrefixes = map[string]bool{
	"BG": true,
	"GB": true,
	"NK": true,
	"KN": true,
	"TN": true,
	"NT": true,
	"ZZ": true,
}

// NINDetector detects UK National Insurance Numbers.
type NINDetector struct{}

func NewNINDetector() *NINDetector { return &NINDetector{} }

func (d *NINDetector) Name() string              { return "gb/nin" }
func (d *NINDetector) Locales() []string         { return []string{locale} }
func (d *NINDetector) PIITypes() []model.PIIType { return []model.PIIType{model.NationalID} }

func (d *NINDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := ninRe.FindAllStringSubmatchIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		full := text[loc[0]:loc[1]]
		prefix := strings.ToUpper(text[loc[2]:loc[3]])

		if !isValidNINPrefix(prefix) {
			continue
		}

		matches = append(matches, model.Match{
			Type:       model.NationalID,
			Value:      full,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.85,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}

// isValidNINPrefix checks that the two-letter prefix follows HMRC rules.
func isValidNINPrefix(prefix string) bool {
	if len(prefix) != 2 {
		return false
	}
	// First letter must not be D, F, I, Q, U, V.
	switch prefix[0] {
	case 'D', 'F', 'I', 'Q', 'U', 'V':
		return false
	}
	// Second letter must not be D, F, I, O, Q, U, V.
	switch prefix[1] {
	case 'D', 'F', 'I', 'O', 'Q', 'U', 'V':
		return false
	}
	// Certain two-letter combos are never allocated.
	if invalidNINPrefixes[prefix] {
		return false
	}
	return true
}
