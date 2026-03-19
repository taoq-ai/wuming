package jp

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Matches 13 consecutive digits.
var corporateNumberRe = regexp.MustCompile(`\b\d{13}\b`)

// CorporateNumberDetector detects Japanese Corporate Numbers (法人番号).
type CorporateNumberDetector struct{}

func NewCorporateNumberDetector() *CorporateNumberDetector { return &CorporateNumberDetector{} }

func (d *CorporateNumberDetector) Name() string              { return "jp/corporate_number" }
func (d *CorporateNumberDetector) Locales() []string         { return []string{locale} }
func (d *CorporateNumberDetector) PIITypes() []model.PIIType { return []model.PIIType{model.TaxID} }

func (d *CorporateNumberDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := corporateNumberRe.FindAllStringIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		candidate := text[loc[0]:loc[1]]
		if !isValidCorporateNumber(candidate) {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.TaxID,
			Value:      candidate,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.80,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}

// isValidCorporateNumber verifies the check digit of a 13-digit Corporate Number.
// The first digit is the check digit: 9 - (weighted sum of digits 2-13 mod 9).
// Weights alternate 1, 2, 1, 2, ... for positions 2-13 (left to right).
func isValidCorporateNumber(s string) bool {
	if len(s) != 13 {
		return false
	}

	var sum int
	for i := 1; i < 13; i++ {
		digit := int(s[i] - '0')
		// Positions 2-13 (1-indexed): odd positions get weight 1, even get weight 2.
		if i%2 == 1 {
			sum += digit * 1
		} else {
			sum += digit * 2
		}
	}

	check := 9 - (sum % 9)
	return int(s[0]-'0') == check
}
