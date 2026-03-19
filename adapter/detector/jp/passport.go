package jp

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Japanese passport: 2 uppercase letters followed by 7 digits (e.g. TK1234567).
var jpPassportRe = regexp.MustCompile(`\b[A-Z]{2}\d{7}\b`)

// PassportDetector detects Japanese passport numbers.
type PassportDetector struct{}

func NewPassportDetector() *PassportDetector { return &PassportDetector{} }

func (d *PassportDetector) Name() string              { return "jp/passport" }
func (d *PassportDetector) Locales() []string         { return []string{locale} }
func (d *PassportDetector) PIITypes() []model.PIIType { return []model.PIIType{model.Passport} }

func (d *PassportDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(jpPassportRe, text, model.Passport, 0.70, d.Name()), nil
}
