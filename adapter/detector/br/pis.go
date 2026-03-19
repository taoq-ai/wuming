package br

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// PIS/PASEP: XXX.XXXXX.XX-X or 11 digits.
var pisRe = regexp.MustCompile(`\b\d{3}\.\d{5}\.\d{2}-\d\b|\b\d{11}\b`)

// PISDetector detects Brazilian PIS/PASEP numbers.
type PISDetector struct{}

func NewPISDetector() *PISDetector { return &PISDetector{} }

func (d *PISDetector) Name() string              { return "br/pis" }
func (d *PISDetector) Locales() []string         { return []string{locale} }
func (d *PISDetector) PIITypes() []model.PIIType { return []model.PIIType{model.NationalID} }

func (d *PISDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := pisRe.FindAllStringIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		raw := text[loc[0]:loc[1]]
		digits := stripNonDigits(raw)
		if len(digits) != 11 {
			continue
		}
		if !isValidPIS(digits) {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.NationalID,
			Value:      raw,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.80,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}

// isValidPIS validates a PIS/PASEP number using mod-11 with weights [3,2,9,8,7,6,5,4,3,2].
func isValidPIS(digits string) bool {
	if allSameDigit(digits) {
		return false
	}

	weights := []int{3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	d := make([]int, 11)
	for i, c := range digits {
		d[i] = int(c - '0')
	}

	sum := 0
	for i := 0; i < 10; i++ {
		sum += d[i] * weights[i]
	}
	rem := sum % 11
	check := 0
	if rem >= 2 {
		check = 11 - rem
	}
	return d[10] == check
}
