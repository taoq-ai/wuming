package kr

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Korean phone numbers:
// Mobile: 010-XXXX-XXXX (with optional +82 prefix)
// Seoul landline: 02-XXX(X)-XXXX
// Other landline: 0XX-XXX(X)-XXXX
var phoneRe = regexp.MustCompile(
	`(?:` +
		`\+82[\s-]?10[\s-]?\d{4}[\s-]?\d{4}` + // international mobile (+82 10-XXXX-XXXX)
		`|\+82[\s-]?2[\s-]?\d{3,4}[\s-]?\d{4}` + // international Seoul (+82 2-XXX-XXXX)
		`|\+82[\s-]?[3-6][1-9][\s-]?\d{3,4}[\s-]?\d{4}` + // international other landline
		`|010[\s-]?\d{4}[\s-]?\d{4}` + // domestic mobile
		`|02[\s-]?\d{3,4}[\s-]?\d{4}` + // domestic Seoul landline
		`|0[3-6][1-9][\s-]?\d{3,4}[\s-]?\d{4}` + // domestic other landline
		`)`)

// PhoneDetector detects South Korean phone numbers.
type PhoneDetector struct{}

func NewPhoneDetector() *PhoneDetector { return &PhoneDetector{} }

func (d *PhoneDetector) Name() string              { return "kr/phone" }
func (d *PhoneDetector) Locales() []string         { return []string{locale} }
func (d *PhoneDetector) PIITypes() []model.PIIType { return []model.PIIType{model.Phone} }

func (d *PhoneDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(phoneRe, text, model.Phone, 0.85, d.Name()), nil
}
