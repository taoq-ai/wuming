package au

import (
	"context"
	"regexp"
	"strings"

	"github.com/taoq-ai/wuming/domain/model"
)

// Medicare number: 10 base digits + optional 1-digit issue number, with optional separators.
var medicareRe = regexp.MustCompile(`\b(\d[ \-]?){9}\d(\d)?\b`)

// MedicareDetector detects Australian Medicare card numbers.
type MedicareDetector struct{}

func NewMedicareDetector() *MedicareDetector { return &MedicareDetector{} }

func (d *MedicareDetector) Name() string              { return "au/medicare" }
func (d *MedicareDetector) Locales() []string         { return []string{locale} }
func (d *MedicareDetector) PIITypes() []model.PIIType { return []model.PIIType{model.HealthID} }

func (d *MedicareDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := medicareRe.FindAllStringIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		raw := text[loc[0]:loc[1]]
		if !isValidMedicare(raw) {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.HealthID,
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

// isValidMedicare checks the Medicare check digit.
// Weights [1,3,7,9,1,3,7,9] on first 8 digits; sum mod 10 must equal the 9th digit.
func isValidMedicare(raw string) bool {
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, raw)

	if len(digits) < 10 || len(digits) > 11 {
		return false
	}

	// First digit must be 2-6 (valid Medicare card range).
	if digits[0] < '2' || digits[0] > '6' {
		return false
	}

	weights := [8]int{1, 3, 7, 9, 1, 3, 7, 9}
	sum := 0
	for i := 0; i < 8; i++ {
		sum += int(digits[i]-'0') * weights[i]
	}
	return sum%10 == int(digits[8]-'0')
}
