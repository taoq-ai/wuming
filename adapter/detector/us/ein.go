package us

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// EIN format: XX-XXXXXXX. Valid prefixes per IRS campus assignment.
var einRe = regexp.MustCompile(`\b(?:0[1-6]|1[0-6]|2[0-7]|3[0-9]|4[0-8]|5[0-9]|6[0-8]|7[1-7]|8[0-5]|9[0-5|8])-\d{7}\b`)

// EINDetector detects US Employer Identification Numbers.
type EINDetector struct{}

func NewEINDetector() *EINDetector { return &EINDetector{} }

func (d *EINDetector) Name() string              { return "us/ein" }
func (d *EINDetector) Locales() []string         { return []string{locale} }
func (d *EINDetector) PIITypes() []model.PIIType { return []model.PIIType{model.TaxID} }

func (d *EINDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(einRe, text, model.TaxID, 0.85, d.Name()), nil
}
