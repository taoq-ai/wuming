package au

import (
	"context"
	"regexp"
	"strings"

	"github.com/taoq-ai/wuming/domain/model"
)

// ABN: 11 digits, optionally separated by spaces or dashes.
var abnRe = regexp.MustCompile(`\b(\d[ \-]?){10}\d\b`)

// ABNDetector detects Australian Business Numbers.
type ABNDetector struct{}

func NewABNDetector() *ABNDetector { return &ABNDetector{} }

func (d *ABNDetector) Name() string              { return "au/abn" }
func (d *ABNDetector) Locales() []string         { return []string{locale} }
func (d *ABNDetector) PIITypes() []model.PIIType { return []model.PIIType{model.TaxID} }

func (d *ABNDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := abnRe.FindAllStringIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		raw := text[loc[0]:loc[1]]
		if !isValidABN(raw) {
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

// isValidABN validates using the ABN algorithm:
// subtract 1 from first digit, apply weights, sum mod 89 must equal 0.
func isValidABN(raw string) bool {
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, raw)

	if len(digits) != 11 {
		return false
	}

	weights := [11]int{10, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19}
	sum := 0
	for i, ch := range digits {
		d := int(ch - '0')
		if i == 0 {
			d-- // subtract 1 from first digit
		}
		sum += d * weights[i]
	}
	return sum%89 == 0
}
