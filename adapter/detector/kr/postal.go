package kr

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Korean postal code: 5 digits (new system since 2015).
var postalRe = regexp.MustCompile(`\b\d{5}\b`)

// PostalDetector detects South Korean postal codes.
type PostalDetector struct{}

func NewPostalDetector() *PostalDetector { return &PostalDetector{} }

func (d *PostalDetector) Name() string              { return "kr/postal" }
func (d *PostalDetector) Locales() []string         { return []string{locale} }
func (d *PostalDetector) PIITypes() []model.PIIType { return []model.PIIType{model.PostalCode} }

func (d *PostalDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(postalRe, text, model.PostalCode, 0.50, d.Name()), nil
}
