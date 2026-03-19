package br

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// CEP postal code: XXXXX-XXX or 8 digits.
var cepRe = regexp.MustCompile(`\b\d{5}-\d{3}\b|\b\d{8}\b`)

// CEPDetector detects Brazilian CEP postal codes.
type CEPDetector struct{}

func NewCEPDetector() *CEPDetector { return &CEPDetector{} }

func (d *CEPDetector) Name() string              { return "br/cep" }
func (d *CEPDetector) Locales() []string         { return []string{locale} }
func (d *CEPDetector) PIITypes() []model.PIIType { return []model.PIIType{model.PostalCode} }

func (d *CEPDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(cepRe, text, model.PostalCode, 0.70, d.Name()), nil
}
