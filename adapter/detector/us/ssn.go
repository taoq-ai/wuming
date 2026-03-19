package us

import (
	"context"
	"regexp"
	"strings"

	"github.com/taoq-ai/wuming/domain/model"
)

// Matches 3-2-4 digit patterns with optional dashes.
var ssnRe = regexp.MustCompile(`\b(\d{3})-?(\d{2})-?(\d{4})\b`)

// SSNDetector detects US Social Security Numbers.
type SSNDetector struct{}

func NewSSNDetector() *SSNDetector { return &SSNDetector{} }

func (d *SSNDetector) Name() string              { return "us/ssn" }
func (d *SSNDetector) Locales() []string         { return []string{locale} }
func (d *SSNDetector) PIITypes() []model.PIIType { return []model.PIIType{model.NationalID} }

func (d *SSNDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := ssnRe.FindAllStringSubmatchIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		full := text[loc[0]:loc[1]]
		area := text[loc[2]:loc[3]]
		group := text[loc[4]:loc[5]]
		serial := text[loc[6]:loc[7]]

		if !isValidSSN(area, group, serial) {
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

func isValidSSN(area, group, serial string) bool {
	// Area cannot be 000, 666, or 900-999.
	if area == "000" || area == "666" || strings.HasPrefix(area, "9") {
		return false
	}
	// Group cannot be 00.
	if group == "00" {
		return false
	}
	// Serial cannot be 0000.
	if serial == "0000" {
		return false
	}
	return true
}
