package cn

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Chinese passport formats:
// E-passport:   E followed by 8 digits
// Regular:      G followed by 8 digits
// Diplomatic:   D followed by 7 digits
// Service:      S followed by 7 or 8 digits
var passportRe = regexp.MustCompile(`\b(?:[EG]\d{8}|D\d{7}|S\d{7,8})\b`)

// PassportDetector detects Chinese passport numbers.
type PassportDetector struct{}

func NewPassportDetector() *PassportDetector { return &PassportDetector{} }

func (d *PassportDetector) Name() string              { return "cn/passport" }
func (d *PassportDetector) Locales() []string         { return []string{locale} }
func (d *PassportDetector) PIITypes() []model.PIIType { return []model.PIIType{model.Passport} }

func (d *PassportDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(passportRe, text, model.Passport, 0.75, d.Name()), nil
}
