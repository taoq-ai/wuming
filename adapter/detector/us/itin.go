package us

import (
	"context"
	"regexp"
	"strconv"

	"github.com/taoq-ai/wuming/domain/model"
)

// ITIN: starts with 9, then 2 digits, then group, then 4 digits.
var itinRe = regexp.MustCompile(`\b(9\d{2})-?(\d{2})-?(\d{4})\b`)

// ITINDetector detects US Individual Taxpayer Identification Numbers.
type ITINDetector struct{}

func NewITINDetector() *ITINDetector { return &ITINDetector{} }

func (d *ITINDetector) Name() string              { return "us/itin" }
func (d *ITINDetector) Locales() []string         { return []string{locale} }
func (d *ITINDetector) PIITypes() []model.PIIType { return []model.PIIType{model.TaxID} }

func (d *ITINDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := itinRe.FindAllStringSubmatchIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		full := text[loc[0]:loc[1]]
		group := text[loc[4]:loc[5]]

		if !isValidITINGroup(group) {
			continue
		}

		matches = append(matches, model.Match{
			Type:       model.TaxID,
			Value:      full,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.8,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}

// isValidITINGroup checks that the group number is in valid ITIN ranges:
// 50-65, 70-88, 90-92, 94-99.
func isValidITINGroup(group string) bool {
	g, err := strconv.Atoi(group)
	if err != nil {
		return false
	}
	return (g >= 50 && g <= 65) || (g >= 70 && g <= 88) || (g >= 90 && g <= 92) || (g >= 94 && g <= 99)
}
