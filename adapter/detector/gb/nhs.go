package gb

import (
	"context"
	"regexp"
	"strings"

	"github.com/taoq-ai/wuming/domain/model"
)

// NHS number: 10 digits, optionally formatted as XXX XXX XXXX or XXX-XXX-XXXX.
var nhsRe = regexp.MustCompile(`\b(\d{3})[\s\-]?(\d{3})[\s\-]?(\d{4})\b`)

// NHSDetector detects UK NHS Numbers with mod-11 check digit validation.
type NHSDetector struct{}

func NewNHSDetector() *NHSDetector { return &NHSDetector{} }

func (d *NHSDetector) Name() string              { return "gb/nhs" }
func (d *NHSDetector) Locales() []string         { return []string{locale} }
func (d *NHSDetector) PIITypes() []model.PIIType { return []model.PIIType{model.HealthID} }

func (d *NHSDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := nhsRe.FindAllStringSubmatchIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		full := text[loc[0]:loc[1]]
		// Extract the 10 raw digits by concatenating the 3 capture groups.
		digits := text[loc[2]:loc[3]] + text[loc[4]:loc[5]] + text[loc[6]:loc[7]]

		if !isValidNHS(digits) {
			continue
		}

		matches = append(matches, model.Match{
			Type:       model.HealthID,
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

// isValidNHS validates the NHS number using the mod-11 algorithm.
func isValidNHS(digits string) bool {
	digits = strings.ReplaceAll(digits, " ", "")
	if len(digits) != 10 {
		return false
	}

	weights := [9]int{10, 9, 8, 7, 6, 5, 4, 3, 2}
	sum := 0
	for i := 0; i < 9; i++ {
		d := int(digits[i] - '0')
		sum += d * weights[i]
	}

	remainder := 11 - (sum % 11)
	if remainder == 11 {
		remainder = 0
	}
	// If remainder is 10, the number is invalid.
	if remainder == 10 {
		return false
	}

	checkDigit := int(digits[9] - '0')
	return checkDigit == remainder
}
