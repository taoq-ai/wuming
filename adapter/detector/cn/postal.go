package cn

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Chinese postal code: 6 digits, first digit 1-8.
var postalRe = regexp.MustCompile(`\b[1-8]\d{5}\b`)

// PostalDetector detects Chinese postal codes.
type PostalDetector struct{}

func NewPostalDetector() *PostalDetector { return &PostalDetector{} }

func (d *PostalDetector) Name() string              { return "cn/postal" }
func (d *PostalDetector) Locales() []string         { return []string{locale} }
func (d *PostalDetector) PIITypes() []model.PIIType { return []model.PIIType{model.PostalCode} }

func (d *PostalDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(postalRe, text, model.PostalCode, 0.55, d.Name()), nil
}
