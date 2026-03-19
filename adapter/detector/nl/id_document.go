package nl

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Dutch ID card: pattern like SPECI2014 — a mix of uppercase letters and digits, 9 characters.
var dutchIDRe = regexp.MustCompile(`\b[A-Z]{2,5}[A-Z0-9]{4,7}\b`)

// Dutch passport: 2 uppercase letters followed by 7 alphanumeric characters.
var dutchPassportRe = regexp.MustCompile(`\b[A-Z]{2}[A-Z0-9]{7}\b`)

// IDDocumentDetector detects Dutch ID card and passport numbers.
type IDDocumentDetector struct{}

func NewIDDocumentDetector() *IDDocumentDetector { return &IDDocumentDetector{} }

func (d *IDDocumentDetector) Name() string      { return "nl/id_document" }
func (d *IDDocumentDetector) Locales() []string { return []string{locale} }
func (d *IDDocumentDetector) PIITypes() []model.PIIType {
	return []model.PIIType{model.NationalID, model.Passport}
}

func (d *IDDocumentDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	var matches []model.Match

	// Detect Dutch ID cards (9-character alphanumeric).
	for _, loc := range dutchIDRe.FindAllStringIndex(text, -1) {
		value := text[loc[0]:loc[1]]
		if len(value) != 9 {
			continue
		}
		// Must contain both letters and digits.
		if !hasLetterAndDigit(value) {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.NationalID,
			Value:      value,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.70,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}

	// Detect Dutch passports (2 letters + 7 alphanumeric = 9 chars).
	for _, loc := range dutchPassportRe.FindAllStringIndex(text, -1) {
		value := text[loc[0]:loc[1]]
		// Must contain at least one digit to distinguish from plain words.
		if !hasDigit(value) {
			continue
		}
		// Skip if already matched as ID card (same value/position).
		if isDuplicate(matches, loc[0], loc[1]) {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.Passport,
			Value:      value,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.70,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}

	return matches, nil
}

func hasLetterAndDigit(s string) bool {
	hasLetter, hasDigitVal := false, false
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			hasLetter = true
		}
		if r >= '0' && r <= '9' {
			hasDigitVal = true
		}
	}
	return hasLetter && hasDigitVal
}

func hasDigit(s string) bool {
	for _, r := range s {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}

func isDuplicate(matches []model.Match, start, end int) bool {
	for _, m := range matches {
		if m.Start == start && m.End == end {
			return true
		}
	}
	return false
}
