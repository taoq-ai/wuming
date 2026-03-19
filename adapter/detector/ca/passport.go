package ca

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Canadian passport: 2 uppercase letters followed by 6 digits.
var caPassportRe = regexp.MustCompile(`\b[A-Z]{2}\d{6}\b`)

// PassportDetector detects Canadian passport numbers.
type PassportDetector struct{}

func NewPassportDetector() *PassportDetector { return &PassportDetector{} }

func (d *PassportDetector) Name() string              { return "ca/passport" }
func (d *PassportDetector) Locales() []string         { return []string{locale} }
func (d *PassportDetector) PIITypes() []model.PIIType { return []model.PIIType{model.Passport} }

func (d *PassportDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(caPassportRe, text, model.Passport, 0.65, d.Name()), nil
}
