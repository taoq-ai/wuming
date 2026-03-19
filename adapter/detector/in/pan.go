package in

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// PAN: 5 uppercase letters, 4 digits, 1 uppercase letter.
// The 4th character indicates holder type (A, B, C, F, G, H, J, L, P, T).
var panRe = regexp.MustCompile(`\b[A-Z]{3}[ABCFGHJLPT][A-Z]\d{4}[A-Z]\b`)

// PANDetector detects Indian Permanent Account Numbers.
type PANDetector struct{}

func NewPANDetector() *PANDetector { return &PANDetector{} }

func (d *PANDetector) Name() string              { return "in/pan" }
func (d *PANDetector) Locales() []string         { return []string{locale} }
func (d *PANDetector) PIITypes() []model.PIIType { return []model.PIIType{model.TaxID} }

func (d *PANDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(panRe, text, model.TaxID, 0.85, d.Name()), nil
}
