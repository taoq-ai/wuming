package in

import (
	"context"
	"regexp"
	"strings"

	"github.com/taoq-ai/wuming/domain/model"
)

// Aadhaar: 12 digits, first digit 2-9, with optional spaces every 4 digits.
var aadhaarRe = regexp.MustCompile(`\b[2-9]\d{3}\s?\d{4}\s?\d{4}\b`)

// Verhoeff algorithm tables.
var (
	verhoeffMultiply = [10][10]int{
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		{1, 2, 3, 4, 0, 6, 7, 8, 9, 5},
		{2, 3, 4, 0, 1, 7, 8, 9, 5, 6},
		{3, 4, 0, 1, 2, 8, 9, 5, 6, 7},
		{4, 0, 1, 2, 3, 9, 5, 6, 7, 8},
		{5, 9, 8, 7, 6, 0, 4, 3, 2, 1},
		{6, 5, 9, 8, 7, 1, 0, 4, 3, 2},
		{7, 6, 5, 9, 8, 2, 1, 0, 4, 3},
		{8, 7, 6, 5, 9, 3, 2, 1, 0, 4},
		{9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
	}

	verhoeffPermute = [8][10]int{
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		{1, 5, 7, 6, 2, 8, 3, 0, 9, 4},
		{5, 8, 0, 3, 7, 9, 6, 1, 4, 2},
		{8, 9, 1, 6, 0, 4, 3, 5, 2, 7},
		{9, 4, 5, 3, 1, 2, 6, 8, 7, 0},
		{4, 2, 8, 6, 5, 7, 3, 9, 0, 1},
		{2, 7, 9, 3, 8, 0, 6, 4, 1, 5},
		{7, 0, 4, 6, 9, 1, 3, 2, 5, 8},
	}
)

// verhoeffChecksum returns true if the digit string passes the Verhoeff check.
func verhoeffChecksum(s string) bool {
	c := 0
	digits := []int{}
	for i := len(s) - 1; i >= 0; i-- {
		digits = append(digits, int(s[i]-'0'))
	}
	for i, d := range digits {
		c = verhoeffMultiply[c][verhoeffPermute[i%8][d]]
	}
	return c == 0
}

// AadhaarDetector detects Indian Aadhaar numbers.
type AadhaarDetector struct{}

func NewAadhaarDetector() *AadhaarDetector { return &AadhaarDetector{} }

func (d *AadhaarDetector) Name() string              { return "in/aadhaar" }
func (d *AadhaarDetector) Locales() []string         { return []string{locale} }
func (d *AadhaarDetector) PIITypes() []model.PIIType { return []model.PIIType{model.NationalID} }

func (d *AadhaarDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := aadhaarRe.FindAllStringIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		raw := text[loc[0]:loc[1]]
		digits := strings.ReplaceAll(raw, " ", "")
		if len(digits) != 12 {
			continue
		}
		if !verhoeffChecksum(digits) {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.NationalID,
			Value:      raw,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.85,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}
