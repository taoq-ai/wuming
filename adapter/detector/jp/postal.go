package jp

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Japanese postal code: optional 〒 prefix, 3 digits, hyphen, 4 digits.
var postalRe = regexp.MustCompile(`(?:〒\s?)?\b\d{3}-\d{4}\b`)

// PostalDetector detects Japanese postal codes (郵便番号).
type PostalDetector struct{}

func NewPostalDetector() *PostalDetector { return &PostalDetector{} }

func (d *PostalDetector) Name() string              { return "jp/postal" }
func (d *PostalDetector) Locales() []string         { return []string{locale} }
func (d *PostalDetector) PIITypes() []model.PIIType { return []model.PIIType{model.PostalCode} }

func (d *PostalDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(postalRe, text, model.PostalCode, 0.75, d.Name()), nil
}
