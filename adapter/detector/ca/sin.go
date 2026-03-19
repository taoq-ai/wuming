package ca

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Matches 9-digit patterns: XXX XXX XXX or XXX-XXX-XXX or XXXXXXXXX.
var sinRe = regexp.MustCompile(`\b(\d{3})[\s-]?(\d{3})[\s-]?(\d{3})\b`)

// SINDetector detects Canadian Social Insurance Numbers.
type SINDetector struct{}

func NewSINDetector() *SINDetector { return &SINDetector{} }

func (d *SINDetector) Name() string              { return "ca/sin" }
func (d *SINDetector) Locales() []string         { return []string{locale} }
func (d *SINDetector) PIITypes() []model.PIIType { return []model.PIIType{model.NationalID} }

func (d *SINDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := sinRe.FindAllStringSubmatchIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		full := text[loc[0]:loc[1]]
		g1 := text[loc[2]:loc[3]]
		g2 := text[loc[4]:loc[5]]
		g3 := text[loc[6]:loc[7]]

		digits := g1 + g2 + g3
		if !luhnValid(digits) {
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

// luhnValid checks whether digits satisfies the Luhn algorithm.
func luhnValid(digits string) bool {
	if len(digits) != 9 {
		return false
	}
	// Reject all-zero SINs.
	if digits == "000000000" {
		return false
	}

	sum := 0
	for i, ch := range digits {
		n := int(ch - '0')
		// Double every second digit (0-indexed: positions 1, 3, 5, 7).
		if i%2 == 1 {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
	}
	return sum%10 == 0
}
