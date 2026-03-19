package de

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Personalausweisnummer: 10 alphanumeric characters.
// First char is a letter from the set {L, M, N, P, R, T, V, W, X, Y},
// followed by 8 alphanumeric chars and 1 check digit.
var idCardRe = regexp.MustCompile(`\b[LMNPRTVWXY][0-9A-Z]{8}[0-9]\b`)

// IDCardDetector detects German Personalausweisnummer (ID card numbers).
type IDCardDetector struct{}

func NewIDCardDetector() *IDCardDetector { return &IDCardDetector{} }

func (d *IDCardDetector) Name() string              { return "de/id_card" }
func (d *IDCardDetector) Locales() []string         { return []string{locale} }
func (d *IDCardDetector) PIITypes() []model.PIIType { return []model.PIIType{model.NationalID} }

func (d *IDCardDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := idCardRe.FindAllStringIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		value := text[loc[0]:loc[1]]
		if !validIDCardCheckDigit(value) {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.NationalID,
			Value:      value,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.75,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}

// charValue returns the numeric value of a character for the check digit calculation.
// Digits 0-9 map to 0-9; letters A-Z map to 10-35.
func charValue(c byte) int {
	if c >= '0' && c <= '9' {
		return int(c - '0')
	}
	return int(c-'A') + 10
}

// validIDCardCheckDigit verifies the weighted checksum (weights 7,3,1 cycling).
// The check digit is the last character of the 10-character ID.
func validIDCardCheckDigit(id string) bool {
	if len(id) != 10 {
		return false
	}
	weights := []int{7, 3, 1}
	sum := 0
	for i := 0; i < 9; i++ {
		sum += charValue(id[i]) * weights[i%3]
	}
	return (sum % 10) == int(id[9]-'0')
}
