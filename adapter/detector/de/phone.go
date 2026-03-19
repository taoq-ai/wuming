package de

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// German phone numbers:
// - International: +49 followed by subscriber number (without leading 0)
// - Mobile: 01XX ... (e.g. 0151, 0160, 0170, 0171, 0175, 0176, 0177, 0178, 0179)
// - Landline: 0XXX ... (area codes 2-9 followed by digits)
// Supports separators: space, dash, dot, slash.
var dePhoneRe = regexp.MustCompile(
	`(?:\+49[\s.\-/]?|00[\s.\-/]?49[\s.\-/]?)` + // international prefix +49 or 0049
		`(?:\(?\d{2,4}\)?[\s.\-/]?)` + // area code (2-4 digits, optional parens)
		`\d[\d\s.\-/]{5,10}\d` + // subscriber number
		`|` +
		`\(?0[1-9]\d{1,3}\)?[\s.\-/]?\d[\d\s.\-/]{4,9}\d`, // domestic format
)

// PhoneDetector detects German phone numbers.
type PhoneDetector struct{}

func NewPhoneDetector() *PhoneDetector { return &PhoneDetector{} }

func (d *PhoneDetector) Name() string              { return "de/phone" }
func (d *PhoneDetector) Locales() []string         { return []string{locale} }
func (d *PhoneDetector) PIITypes() []model.PIIType { return []model.PIIType{model.Phone} }

func (d *PhoneDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(dePhoneRe, text, model.Phone, 0.85, d.Name()), nil
}
