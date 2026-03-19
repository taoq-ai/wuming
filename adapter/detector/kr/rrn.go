package kr

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Korean Resident Registration Number: YYMMDD-GNNNNNN (13 digits with hyphen).
var rrnRe = regexp.MustCompile(`\b(\d{2})(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])-([1-8])(\d{6})\b`)

// RRNDetector detects South Korean Resident Registration Numbers.
type RRNDetector struct{}

func NewRRNDetector() *RRNDetector { return &RRNDetector{} }

func (d *RRNDetector) Name() string              { return "kr/rrn" }
func (d *RRNDetector) Locales() []string         { return []string{locale} }
func (d *RRNDetector) PIITypes() []model.PIIType { return []model.PIIType{model.NationalID} }

func (d *RRNDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := rrnRe.FindAllStringSubmatchIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		full := text[loc[0]:loc[1]]
		if !isValidRRN(full) {
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

// isValidRRN validates the RRN check digit.
// The 13th digit = (11 - (weighted sum mod 11)) mod 10.
// Weights: [2,3,4,5,6,7,8,9,2,3,4,5] applied to the first 12 digits.
func isValidRRN(rrn string) bool {
	// Remove the hyphen to get 13 raw digits.
	digits := make([]int, 0, 13)
	for _, c := range rrn {
		if c == '-' {
			continue
		}
		digits = append(digits, int(c-'0'))
	}
	if len(digits) != 13 {
		return false
	}

	weights := []int{2, 3, 4, 5, 6, 7, 8, 9, 2, 3, 4, 5}
	sum := 0
	for i := 0; i < 12; i++ {
		sum += digits[i] * weights[i]
	}
	check := (11 - (sum % 11)) % 10
	return check == digits[12]
}
