package fr

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// French phone numbers:
//   - Local format: 0X XX XX XX XX (X = 1-7 for landline/mobile)
//   - International: +33 X XX XX XX XX
//
// Separators: space, dot, dash, or none.
var phoneRe = regexp.MustCompile(
	`(?:` +
		`(?:\+33[\s.-]?|0)[1-7][\s.-]?\d{2}[\s.-]?\d{2}[\s.-]?\d{2}[\s.-]?\d{2}` +
		`)\b`,
)

// PhoneDetector detects French phone numbers.
type PhoneDetector struct{}

func NewPhoneDetector() *PhoneDetector { return &PhoneDetector{} }

func (d *PhoneDetector) Name() string              { return "fr/phone" }
func (d *PhoneDetector) Locales() []string         { return []string{locale} }
func (d *PhoneDetector) PIITypes() []model.PIIType { return []model.PIIType{model.Phone} }

func (d *PhoneDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(phoneRe, text, model.Phone, 0.85, d.Name()), nil
}
