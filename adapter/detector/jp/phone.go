package jp

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Japanese phone numbers:
// - Mobile: 090/080/070-XXXX-XXXX (with optional +81 prefix)
// - Landline: 0X-XXXX-XXXX or 0XX-XXX-XXXX (common area codes)
// - International: +81 X0-XXXX-XXXX
var jpPhoneRe = regexp.MustCompile(
	`(?:\+81[\s-]?[0-9]{1,2}[\s-]?[0-9]{4}[\s-]?[0-9]{4})` + // +81 format
		`|(?:0[789]0[\s-]?[0-9]{4}[\s-]?[0-9]{4})` + // mobile 090/080/070
		`|(?:0[1-9][0-9]{0,3}[\s-]?[0-9]{2,4}[\s-]?[0-9]{4})`, // landline
)

// PhoneDetector detects Japanese phone numbers.
type PhoneDetector struct{}

func NewPhoneDetector() *PhoneDetector { return &PhoneDetector{} }

func (d *PhoneDetector) Name() string              { return "jp/phone" }
func (d *PhoneDetector) Locales() []string         { return []string{locale} }
func (d *PhoneDetector) PIITypes() []model.PIIType { return []model.PIIType{model.Phone} }

func (d *PhoneDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(jpPhoneRe, text, model.Phone, 0.85, d.Name()), nil
}
