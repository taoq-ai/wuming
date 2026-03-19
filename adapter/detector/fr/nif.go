package fr

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// NIF (Numero d'Identification Fiscale) — 13 digits, first digit 0-3.
// Matches bare 13-digit strings or formatted XX XX XXX XXX XXX.
var nifRe = regexp.MustCompile(`\b[0-3]\d[\s.-]?\d{2}[\s.-]?\d{3}[\s.-]?\d{3}[\s.-]?\d{3}\b`)

// NIFDetector detects French tax identification numbers.
type NIFDetector struct{}

func NewNIFDetector() *NIFDetector { return &NIFDetector{} }

func (d *NIFDetector) Name() string              { return "fr/nif" }
func (d *NIFDetector) Locales() []string         { return []string{locale} }
func (d *NIFDetector) PIITypes() []model.PIIType { return []model.PIIType{model.TaxID} }

func (d *NIFDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(nifRe, text, model.TaxID, 0.70, d.Name()), nil
}
