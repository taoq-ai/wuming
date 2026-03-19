package fr

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/taoq-ai/wuming/domain/model"
)

// NIR (Numero d'Inscription au Repertoire) / Numero de securite sociale.
// Format: X XX XX XXXXX XXX XX (15 digits, with Corsica departments 2A/2B).
// Matches with or without spaces/dashes/dots as separators.
var nirRe = regexp.MustCompile(
	`\b([12])\s?(\d{2})\s?(0[1-9]|1[0-2]|[2-9]\d)\s?` +
		`(\d{2}|2[AB])\s?(\d{3})\s?(\d{3})\s?(\d{2})\b`,
)

// NIRDetector detects French NIR (social security) numbers.
type NIRDetector struct{}

func NewNIRDetector() *NIRDetector { return &NIRDetector{} }

func (d *NIRDetector) Name() string              { return "fr/nir" }
func (d *NIRDetector) Locales() []string         { return []string{locale} }
func (d *NIRDetector) PIITypes() []model.PIIType { return []model.PIIType{model.NationalID} }

func (d *NIRDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := nirRe.FindAllStringSubmatchIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		full := text[loc[0]:loc[1]]
		dept := text[loc[8]:loc[9]] // group 4: department

		// Build the 13-digit numeric base for key validation.
		// For Corsica, replace 2A with 19 and 2B with 18 before computing.
		digits := stripSeparators(full)
		base13 := digits[:13]
		keyStr := digits[13:15]

		key, err := strconv.Atoi(keyStr)
		if err != nil {
			continue
		}

		numBase, ok := nirBaseNumber(base13, dept)
		if !ok {
			continue
		}

		// Control key = 97 - (first 13 digits mod 97).
		expectedKey := 97 - int(numBase%97)
		if key != expectedKey {
			continue
		}

		matches = append(matches, model.Match{
			Type:       model.NationalID,
			Value:      full,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.90,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}

// stripSeparators removes spaces, dashes, and dots from a NIR string,
// and replaces 2A/2B with placeholder digits for numeric parsing.
func stripSeparators(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == 'A' || r == 'a':
			b.WriteRune('A')
		case r == 'B' || r == 'b':
			b.WriteRune('B')
		}
	}
	return b.String()
}

// nirBaseNumber converts the 13-character base (which may contain A or B for
// Corsica departments) into a numeric value suitable for mod-97 computation.
// For Corsica: 2A -> subtract 1000000 from numeric value with 20 replacing 2A,
//
//	2B -> subtract 2000000 from numeric value with 20 replacing 2B.
func nirBaseNumber(base13 string, dept string) (int64, bool) {
	switch strings.ToUpper(dept) {
	case "2A":
		numeric := strings.Replace(base13, "A", "0", 1) // 2A -> 20
		n, err := strconv.ParseInt(numeric, 10, 64)
		if err != nil {
			return 0, false
		}
		return n - 1000000, true
	case "2B":
		numeric := strings.Replace(base13, "B", "0", 1) // 2B -> 20
		n, err := strconv.ParseInt(numeric, 10, 64)
		if err != nil {
			return 0, false
		}
		return n - 2000000, true
	default:
		n, err := strconv.ParseInt(base13, 10, 64)
		if err != nil {
			return 0, false
		}
		return n, true
	}
}
