package gb

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// UK phone patterns:
// Mobile:   07XXX XXXXXX or +44 7XXX XXXXXX
// Landline: 01XX XXX XXXX, 020 XXXX XXXX, +44 1XX XXX XXXX, +44 20 XXXX XXXX
// Separators: space, dash, dot allowed between groups.
var ukPhoneRe = regexp.MustCompile(
	`(?:(?:\+44[\s.\-]?|0)` + // +44 or leading 0
		`(?:` +
		`7\d{3}[\s.\-]?\d{3}[\s.\-]?\d{3}` + // mobile: 7XXX XXX XXX
		`|` +
		`7\d{3}[\s.\-]?\d{6}` + // mobile: 7XXX XXXXXX
		`|` +
		`20[\s.\-]?\d{4}[\s.\-]?\d{4}` + // London: 20 XXXX XXXX
		`|` +
		`1\d{2}[\s.\-]?\d{3}[\s.\-]?\d{4}` + // landline: 1XX XXX XXXX
		`|` +
		`1\d{3}[\s.\-]?\d{5,6}` + // landline: 1XXX XXXXX(X)
		`)` +
		`)`,
)

// PhoneDetector detects UK phone numbers.
type PhoneDetector struct{}

func NewPhoneDetector() *PhoneDetector { return &PhoneDetector{} }

func (d *PhoneDetector) Name() string              { return "gb/phone" }
func (d *PhoneDetector) Locales() []string         { return []string{locale} }
func (d *PhoneDetector) PIITypes() []model.PIIType { return []model.PIIType{model.Phone} }

func (d *PhoneDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(ukPhoneRe, text, model.Phone, 0.85, d.Name()), nil
}
