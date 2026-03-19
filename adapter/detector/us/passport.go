package us

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// US passport: 9 digits (may also have letter prefix in newer format).
var passportRe = regexp.MustCompile(`\b[A-Z]?\d{8,9}\b`)

// PassportDetector detects US passport numbers.
type PassportDetector struct{}

func NewPassportDetector() *PassportDetector { return &PassportDetector{} }

func (d *PassportDetector) Name() string              { return "us/passport" }
func (d *PassportDetector) Locales() []string         { return []string{locale} }
func (d *PassportDetector) PIITypes() []model.PIIType { return []model.PIIType{model.Passport} }

func (d *PassportDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(passportRe, text, model.Passport, 0.6, d.Name()), nil
}
