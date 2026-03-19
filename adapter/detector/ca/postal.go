package ca

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Canadian postal code: letter-digit-letter space? digit-letter-digit.
// First letter excludes D, F, I, O, Q, U, W, Z.
var postalRe = regexp.MustCompile(`\b[ABCEGHJKLMNPRSTVXY]\d[A-Z]\s?\d[A-Z]\d\b`)

// PostalCodeDetector detects Canadian postal codes.
type PostalCodeDetector struct{}

func NewPostalCodeDetector() *PostalCodeDetector { return &PostalCodeDetector{} }

func (d *PostalCodeDetector) Name() string              { return "ca/postal_code" }
func (d *PostalCodeDetector) Locales() []string         { return []string{locale} }
func (d *PostalCodeDetector) PIITypes() []model.PIIType { return []model.PIIType{model.PostalCode} }

func (d *PostalCodeDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(postalRe, text, model.PostalCode, 0.90, d.Name()), nil
}
