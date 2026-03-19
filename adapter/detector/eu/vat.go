package eu

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// vatRe matches EU VAT identification numbers for all 27 member states.
// Each country has a two-letter prefix followed by a country-specific format.
var vatRe = regexp.MustCompile(
	`\b(?:` +
		`ATU\d{8}` +
		`|BE[01]\d{9}` +
		`|BG\d{9,10}` +
		`|CY\d{8}[A-Z]` +
		`|CZ\d{8,10}` +
		`|DE\d{9}` +
		`|DK\d{8}` +
		`|EE\d{9}` +
		`|EL\d{9}` +
		`|ES[A-Z0-9]\d{7}[A-Z0-9]` +
		`|FI\d{8}` +
		`|FR[A-Z0-9]{2}\d{9}` +
		`|HR\d{11}` +
		`|HU\d{8}` +
		`|IE\d[A-Z0-9]\d{5}[A-Z]{1,2}` +
		`|IT\d{11}` +
		`|LT\d{9,12}` +
		`|LU\d{8}` +
		`|LV\d{11}` +
		`|MT\d{8}` +
		`|NL\d{9}B\d{2}` +
		`|PL\d{10}` +
		`|PT\d{9}` +
		`|RO\d{2,10}` +
		`|SE\d{12}` +
		`|SI\d{8}` +
		`|SK\d{10}` +
		`)` + `\b`,
)

// VATDetector detects EU VAT identification numbers.
type VATDetector struct{}

func NewVATDetector() *VATDetector { return &VATDetector{} }

func (d *VATDetector) Name() string              { return "eu/vat" }
func (d *VATDetector) Locales() []string         { return []string{locale} }
func (d *VATDetector) PIITypes() []model.PIIType { return []model.PIIType{model.TaxID} }

func (d *VATDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(vatRe, text, model.TaxID, 0.90, d.Name()), nil
}
