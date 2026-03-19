package br

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Matches CNPJ in XX.XXX.XXX/XXXX-XX or 14-digit unformatted form.
var cnpjRe = regexp.MustCompile(`\b\d{2}\.\d{3}\.\d{3}/\d{4}-\d{2}\b|\b\d{14}\b`)

// CNPJDetector detects Brazilian CNPJ numbers.
type CNPJDetector struct{}

func NewCNPJDetector() *CNPJDetector { return &CNPJDetector{} }

func (d *CNPJDetector) Name() string              { return "br/cnpj" }
func (d *CNPJDetector) Locales() []string         { return []string{locale} }
func (d *CNPJDetector) PIITypes() []model.PIIType { return []model.PIIType{model.TaxID} }

func (d *CNPJDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := cnpjRe.FindAllStringIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		raw := text[loc[0]:loc[1]]
		digits := stripNonDigits(raw)
		if len(digits) != 14 {
			continue
		}
		if !isValidCNPJ(digits) {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.TaxID,
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

// isValidCNPJ validates a CNPJ using the mod-11 algorithm with standard weights.
func isValidCNPJ(digits string) bool {
	if allSameDigit(digits) {
		return false
	}

	d := make([]int, 14)
	for i, c := range digits {
		d[i] = int(c - '0')
	}

	// First check digit: weights [5,4,3,2,9,8,7,6,5,4,3,2].
	w1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum := 0
	for i := 0; i < 12; i++ {
		sum += d[i] * w1[i]
	}
	rem := sum % 11
	check1 := 0
	if rem >= 2 {
		check1 = 11 - rem
	}
	if d[12] != check1 {
		return false
	}

	// Second check digit: weights [6,5,4,3,2,9,8,7,6,5,4,3,2].
	w2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum = 0
	for i := 0; i < 13; i++ {
		sum += d[i] * w2[i]
	}
	rem = sum % 11
	check2 := 0
	if rem >= 2 {
		check2 = 11 - rem
	}
	return d[13] == check2
}
