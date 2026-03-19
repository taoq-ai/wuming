package de

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Steuerliche Identifikationsnummer: exactly 11 digits.
var steuerIDRe = regexp.MustCompile(`\b\d{11}\b`)

// SteuerIDDetector detects German tax identification numbers (Steuer-ID).
type SteuerIDDetector struct{}

func NewSteuerIDDetector() *SteuerIDDetector { return &SteuerIDDetector{} }

func (d *SteuerIDDetector) Name() string              { return "de/steuer_id" }
func (d *SteuerIDDetector) Locales() []string         { return []string{locale} }
func (d *SteuerIDDetector) PIITypes() []model.PIIType { return []model.PIIType{model.TaxID} }

func (d *SteuerIDDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := steuerIDRe.FindAllStringIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		value := text[loc[0]:loc[1]]
		if !validSteuerID(value) {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.TaxID,
			Value:      value,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.85,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}

// validSteuerID checks the digit distribution and ISO 7064 Mod 11,10 check digit.
// Rules for the first 10 digits:
//   - Exactly one digit (0-9) must appear twice
//   - Exactly one digit must not appear at all
//   - The remaining 8 digits each appear exactly once
//
// The 11th digit is a check digit based on ISO 7064 Mod 11,10.
func validSteuerID(id string) bool {
	if len(id) != 11 {
		return false
	}
	// First digit must not be 0.
	if id[0] == '0' {
		return false
	}

	// Check digit distribution among first 10 digits.
	var freq [10]int
	for i := 0; i < 10; i++ {
		freq[id[i]-'0']++
	}
	doubles := 0
	zeros := 0
	ones := 0
	for _, f := range freq {
		switch f {
		case 0:
			zeros++
		case 1:
			ones++
		case 2:
			doubles++
		default:
			return false // no digit may appear 3+ times
		}
	}
	if doubles != 1 || zeros != 1 || ones != 8 {
		return false
	}

	// ISO 7064 Mod 11,10 check digit.
	product := 10
	for i := 0; i < 10; i++ {
		sum := (int(id[i]-'0') + product) % 10
		if sum == 0 {
			sum = 10
		}
		product = (sum * 2) % 11
	}
	check := 11 - product
	if check == 10 {
		check = 0
	}
	return check == int(id[10]-'0')
}
