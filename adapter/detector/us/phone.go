package us

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// NANP phone: optional +1/1 prefix, area code (2-9XX), 7-digit subscriber number.
// Supports parens, spaces, dots, and dashes as separators.
var phoneRe = regexp.MustCompile(`(?:\+?1[\s.\-]?)?\(?[2-9]\d{2}\)?[\s.\-]?\d{3}[\s.\-]?\d{4}`)

// PhoneDetector detects US phone numbers in NANP format.
type PhoneDetector struct{}

func NewPhoneDetector() *PhoneDetector { return &PhoneDetector{} }

func (d *PhoneDetector) Name() string              { return "us/phone" }
func (d *PhoneDetector) Locales() []string         { return []string{locale} }
func (d *PhoneDetector) PIITypes() []model.PIIType { return []model.PIIType{model.Phone} }

func (d *PhoneDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(phoneRe, text, model.Phone, 0.8, d.Name()), nil
}
