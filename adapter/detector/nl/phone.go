package nl

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Dutch phone number patterns:
// Mobile: 06-XXXXXXXX, 06 XXXXXXXX, +31 6 XXXXXXXX, 0031 6 XXXXXXXX
// Landline: 0XX-XXXXXXX, +31 XX XXXXXXX, 0031 XX XXXXXXX
var phoneRe = regexp.MustCompile(
	`(?:` +
		// International prefix +31 or 0031 followed by mobile (6) or area code (1-5, 7)
		`(?:\+31|0031)[\s\-]?(?:6[\s\-]?\d{8}|[1-57]\d[\s\-]?\d{7})` +
		`|` +
		// Domestic mobile: 06 followed by 8 digits
		`\b06[\s\-]?\d{8}` +
		`|` +
		// Domestic landline: 0 + 2-digit area code + 7-digit subscriber
		`\b0[1-57]\d[\s\-]?\d{7}` +
		`)(?:\b|$)`,
)

// PhoneDetector detects Dutch phone numbers.
type PhoneDetector struct{}

func NewPhoneDetector() *PhoneDetector { return &PhoneDetector{} }

func (d *PhoneDetector) Name() string              { return "nl/phone" }
func (d *PhoneDetector) Locales() []string         { return []string{locale} }
func (d *PhoneDetector) PIITypes() []model.PIIType { return []model.PIIType{model.Phone} }

func (d *PhoneDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(phoneRe, text, model.Phone, 0.85, d.Name()), nil
}
