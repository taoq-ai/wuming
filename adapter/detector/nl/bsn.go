package nl

import (
	"context"
	"regexp"
	"strings"

	"github.com/taoq-ai/wuming/domain/model"
)

// Matches 9-digit sequences, optionally separated by dots or spaces (e.g. 123456789, 123.456.789, 123 456 789).
var bsnRe = regexp.MustCompile(`\b\d{3}[.\s]?\d{3}[.\s]?\d{3}\b`)

// BSNDetector detects Dutch Burgerservicenummer (citizen service numbers).
type BSNDetector struct{}

func NewBSNDetector() *BSNDetector { return &BSNDetector{} }

func (d *BSNDetector) Name() string              { return "nl/bsn" }
func (d *BSNDetector) Locales() []string         { return []string{locale} }
func (d *BSNDetector) PIITypes() []model.PIIType { return []model.PIIType{model.NationalID} }

func (d *BSNDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := bsnRe.FindAllStringIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		raw := text[loc[0]:loc[1]]
		digits := stripNonDigits(raw)
		if len(digits) != 9 {
			continue
		}
		if !isValid11Proof(digits) {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.NationalID,
			Value:      raw,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.90,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}

// stripNonDigits removes all non-digit characters from s.
func stripNonDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// isValid11Proof checks the BSN 11-proof checksum.
// Formula: 9*d1 + 8*d2 + 7*d3 + 6*d4 + 5*d5 + 4*d6 + 3*d7 + 2*d8 - 1*d9
// Result must be divisible by 11 and must not be 0.
func isValid11Proof(digits string) bool {
	if len(digits) != 9 {
		return false
	}
	weights := []int{9, 8, 7, 6, 5, 4, 3, 2, -1}
	sum := 0
	for i, w := range weights {
		d := int(digits[i] - '0')
		sum += w * d
	}
	return sum != 0 && sum%11 == 0
}
