package us

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Medicare Beneficiary Identifier (MBI): 11 chars, pattern C A AN N AA N AN AN.
// C=1-9, A=A-Z(excl S,L,O,I,B,Z), N=0-9, AN=alphanumeric(excl S,L,O,I,B,Z).
var medicareRe = regexp.MustCompile(`\b[1-9][A-HJKM-NP-RT-Y][A-HJKM-NP-RT-Y0-9]\d[A-HJKM-NP-RT-Y][A-HJKM-NP-RT-Y]\d[A-HJKM-NP-RT-Y][A-HJKM-NP-RT-Y0-9][A-HJKM-NP-RT-Y0-9]\d\b`)

// MedicareDetector detects US Medicare Beneficiary Identifiers.
type MedicareDetector struct{}

func NewMedicareDetector() *MedicareDetector { return &MedicareDetector{} }

func (d *MedicareDetector) Name() string              { return "us/medicare" }
func (d *MedicareDetector) Locales() []string         { return []string{locale} }
func (d *MedicareDetector) PIITypes() []model.PIIType { return []model.PIIType{model.HealthID} }

func (d *MedicareDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(medicareRe, text, model.HealthID, 0.8, d.Name()), nil
}
