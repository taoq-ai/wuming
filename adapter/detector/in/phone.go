package in

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Indian mobile: optional +91 or 0 prefix, then 10 digits starting with 6-9.
// Supports spaces, dashes, and dots as separators.
var phoneRe = regexp.MustCompile(`(?:\+91[\s.\-]?|0)?[6-9]\d{4}[\s.\-]?\d{5}\b`)

// PhoneDetector detects Indian phone numbers.
type PhoneDetector struct{}

func NewPhoneDetector() *PhoneDetector { return &PhoneDetector{} }

func (d *PhoneDetector) Name() string              { return "in/phone" }
func (d *PhoneDetector) Locales() []string         { return []string{locale} }
func (d *PhoneDetector) PIITypes() []model.PIIType { return []model.PIIType{model.Phone} }

func (d *PhoneDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(phoneRe, text, model.Phone, 0.85, d.Name()), nil
}
