package au

import (
	"context"
	"regexp"
	"strings"

	"github.com/taoq-ai/wuming/domain/model"
)

// TFN: 9 digits, optionally separated by spaces or dashes.
var tfnRe = regexp.MustCompile(`\b(\d[ \-]?){8}\d\b`)

// TFNDetector detects Australian Tax File Numbers.
type TFNDetector struct{}

func NewTFNDetector() *TFNDetector { return &TFNDetector{} }

func (d *TFNDetector) Name() string              { return "au/tfn" }
func (d *TFNDetector) Locales() []string         { return []string{locale} }
func (d *TFNDetector) PIITypes() []model.PIIType { return []model.PIIType{model.TaxID} }

func (d *TFNDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := tfnRe.FindAllStringIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		raw := text[loc[0]:loc[1]]
		if !isValidTFN(raw) {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.TaxID,
			Value:      raw,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.85,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}

// isValidTFN checks the TFN weighted checksum.
// Weights: [1,4,3,7,5,8,6,9,10], sum mod 11 must equal 0.
func isValidTFN(raw string) bool {
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, raw)

	if len(digits) != 9 {
		return false
	}

	weights := [9]int{1, 4, 3, 7, 5, 8, 6, 9, 10}
	sum := 0
	for i, ch := range digits {
		sum += int(ch-'0') * weights[i]
	}
	return sum%11 == 0
}
