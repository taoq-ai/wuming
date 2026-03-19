package kr

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Korean passport: single uppercase letter followed by 8 digits.
// Regular passports start with M; other types use various letters.
var passportRe = regexp.MustCompile(`\b[A-Z]\d{8}\b`)

// PassportDetector detects South Korean passport numbers.
type PassportDetector struct{}

func NewPassportDetector() *PassportDetector { return &PassportDetector{} }

func (d *PassportDetector) Name() string              { return "kr/passport" }
func (d *PassportDetector) Locales() []string         { return []string{locale} }
func (d *PassportDetector) PIITypes() []model.PIIType { return []model.PIIType{model.Passport} }

func (d *PassportDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(passportRe, text, model.Passport, 0.70, d.Name()), nil
}
