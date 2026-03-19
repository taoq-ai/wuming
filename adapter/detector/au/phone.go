package au

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Australian phone patterns:
// Mobile:   04XX XXX XXX or +61 4XX XXX XXX
// Landline: (0X) XXXX XXXX or +61 X XXXX XXXX
var auPhoneRe = regexp.MustCompile(
	`(?:` +
		`\+61[\s\-]?4\d{2}[\s\-]?\d{3}[\s\-]?\d{3}` + // +61 mobile
		`|04\d{2}[\s\-]?\d{3}[\s\-]?\d{3}` + // 04XX mobile
		`|\+61[\s\-]?[2-9][\s\-]?\d{4}[\s\-]?\d{4}` + // +61 landline
		`|\(0[2-9]\)[\s\-]?\d{4}[\s\-]?\d{4}` + // (0X) landline
		`)`,
)

// PhoneDetector detects Australian phone numbers.
type PhoneDetector struct{}

func NewPhoneDetector() *PhoneDetector { return &PhoneDetector{} }

func (d *PhoneDetector) Name() string              { return "au/phone" }
func (d *PhoneDetector) Locales() []string         { return []string{locale} }
func (d *PhoneDetector) PIITypes() []model.PIIType { return []model.PIIType{model.Phone} }

func (d *PhoneDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(auPhoneRe, text, model.Phone, 0.85, d.Name()), nil
}
