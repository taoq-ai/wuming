package in

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// PIN code: 6 digits, first digit 1-9.
var pincodeRe = regexp.MustCompile(`\b[1-9]\d{5}\b`)

// PINCodeDetector detects Indian postal PIN codes.
type PINCodeDetector struct{}

func NewPINCodeDetector() *PINCodeDetector { return &PINCodeDetector{} }

func (d *PINCodeDetector) Name() string              { return "in/pincode" }
func (d *PINCodeDetector) Locales() []string         { return []string{locale} }
func (d *PINCodeDetector) PIITypes() []model.PIIType { return []model.PIIType{model.PostalCode} }

func (d *PINCodeDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(pincodeRe, text, model.PostalCode, 0.55, d.Name()), nil
}
