package br

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// CNH: exactly 11 digits.
var cnhRe = regexp.MustCompile(`\b\d{11}\b`)

// CNHDetector detects Brazilian CNH (driver's license) numbers.
type CNHDetector struct{}

func NewCNHDetector() *CNHDetector { return &CNHDetector{} }

func (d *CNHDetector) Name() string              { return "br/cnh" }
func (d *CNHDetector) Locales() []string         { return []string{locale} }
func (d *CNHDetector) PIITypes() []model.PIIType { return []model.PIIType{model.DriversLicense} }

func (d *CNHDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := cnhRe.FindAllStringIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		raw := text[loc[0]:loc[1]]
		if !isValidCNH(raw) {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.DriversLicense,
			Value:      raw,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.75,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}

// isValidCNH validates a CNH number using the double check digit algorithm.
func isValidCNH(digits string) bool {
	if len(digits) != 11 || allSameDigit(digits) {
		return false
	}

	d := make([]int, 11)
	for i, c := range digits {
		d[i] = int(c - '0')
	}

	// First check digit.
	sum1 := 0
	for i := 0; i < 9; i++ {
		sum1 += d[i] * (9 - i)
	}
	rest1 := sum1 % 11
	check1 := 0
	if rest1 >= 2 {
		check1 = 11 - rest1
	}

	// Second check digit.
	sum2 := 0
	for i := 0; i < 9; i++ {
		sum2 += d[i] * (1 + i)
	}
	rest2 := sum2 % 11
	check2 := 0
	if rest2 >= 2 {
		check2 = 11 - rest2
	}

	return d[9] == check1 && d[10] == check2
}
