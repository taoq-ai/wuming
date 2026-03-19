package de

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Sozialversicherungsnummer: DDMMYYXNNNN where DDMMYY is birth date,
// X is first letter of birth name (encoded), NNNN is serial + check.
// 12 chars total: 8 digits, 1 letter, 3 digits.
// Also matches with space separators (e.g. "12 010290 A 123").
var svnRe = regexp.MustCompile(`\b\d{2}\s?\d{6}\s?[A-Z]\s?\d{3}\b`)

// SozialversicherungDetector detects German social security numbers.
type SozialversicherungDetector struct{}

func NewSozialversicherungDetector() *SozialversicherungDetector {
	return &SozialversicherungDetector{}
}

func (d *SozialversicherungDetector) Name() string      { return "de/sozialversicherung" }
func (d *SozialversicherungDetector) Locales() []string { return []string{locale} }
func (d *SozialversicherungDetector) PIITypes() []model.PIIType {
	return []model.PIIType{model.NationalID}
}

func (d *SozialversicherungDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := svnRe.FindAllStringIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		value := text[loc[0]:loc[1]]
		if !validSVN(value) {
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

// validSVN performs basic validation on the Sozialversicherungsnummer.
// It strips spaces and checks that the date portion has valid day (01-31)
// and month (01-12) values.
func validSVN(svn string) bool {
	// Strip spaces to get the 12-char canonical form.
	clean := make([]byte, 0, 12)
	for i := 0; i < len(svn); i++ {
		if svn[i] != ' ' {
			clean = append(clean, svn[i])
		}
	}
	if len(clean) != 12 {
		return false
	}

	// First two digits: area number (02-99, but we just check it's not 00).
	area := (int(clean[0]-'0') * 10) + int(clean[1]-'0')
	if area == 0 {
		return false
	}

	// Digits 2-7 are DDMMYY (birth date).
	day := (int(clean[2]-'0') * 10) + int(clean[3]-'0')
	month := (int(clean[4]-'0') * 10) + int(clean[5]-'0')
	if day < 1 || day > 31 {
		return false
	}
	if month < 1 || month > 12 {
		return false
	}

	return true
}
