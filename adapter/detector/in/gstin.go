package in

import (
	"context"
	"regexp"
	"strconv"

	"github.com/taoq-ai/wuming/domain/model"
)

// GSTIN: 2-digit state code (01-37) + PAN (5 letters + 4 digits + 1 letter) + 1 alphanumeric + Z + 1 alphanumeric.
var gstinRe = regexp.MustCompile(`\b\d{2}[A-Z]{5}\d{4}[A-Z][0-9A-Z]Z[0-9A-Z]\b`)

// GSTINDetector detects Indian Goods and Services Tax Identification Numbers.
type GSTINDetector struct{}

func NewGSTINDetector() *GSTINDetector { return &GSTINDetector{} }

func (d *GSTINDetector) Name() string              { return "in/gstin" }
func (d *GSTINDetector) Locales() []string         { return []string{locale} }
func (d *GSTINDetector) PIITypes() []model.PIIType { return []model.PIIType{model.TaxID} }

func (d *GSTINDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := gstinRe.FindAllStringIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		raw := text[loc[0]:loc[1]]
		stateCode, _ := strconv.Atoi(raw[:2])
		if stateCode < 1 || stateCode > 37 {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.TaxID,
			Value:      raw,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.90,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}
