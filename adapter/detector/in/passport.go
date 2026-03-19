package in

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Indian passport: 1 uppercase letter followed by 7 digits.
var passportRe = regexp.MustCompile(`\b[A-Z]\d{7}\b`)

// PassportDetector detects Indian passport numbers.
type PassportDetector struct{}

func NewPassportDetector() *PassportDetector { return &PassportDetector{} }

func (d *PassportDetector) Name() string              { return "in/passport" }
func (d *PassportDetector) Locales() []string         { return []string{locale} }
func (d *PassportDetector) PIITypes() []model.PIIType { return []model.PIIType{model.Passport} }

func (d *PassportDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(passportRe, text, model.Passport, 0.65, d.Name()), nil
}
