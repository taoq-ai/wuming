package jp

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Matches 12 consecutive digits.
var myNumberRe = regexp.MustCompile(`\b\d{12}\b`)

// myNumberWeights are the weights for positions 1-11 (left to right) used in
// the check-digit calculation for the Japanese My Number (マイナンバー).
var myNumberWeights = [11]int{6, 5, 4, 3, 2, 7, 6, 5, 4, 3, 2}

// MyNumberDetector detects Japanese My Number (個人番号) identifiers.
type MyNumberDetector struct{}

func NewMyNumberDetector() *MyNumberDetector { return &MyNumberDetector{} }

func (d *MyNumberDetector) Name() string              { return "jp/my_number" }
func (d *MyNumberDetector) Locales() []string         { return []string{locale} }
func (d *MyNumberDetector) PIITypes() []model.PIIType { return []model.PIIType{model.NationalID} }

func (d *MyNumberDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := myNumberRe.FindAllStringIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		candidate := text[loc[0]:loc[1]]
		if !isValidMyNumber(candidate) {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.NationalID,
			Value:      candidate,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.85,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}

// isValidMyNumber verifies the check digit of a 12-digit My Number.
// Check digit = 11 - (weighted sum mod 11); if >= 10, check digit is 0.
func isValidMyNumber(s string) bool {
	if len(s) != 12 {
		return false
	}

	var sum int
	for i := 0; i < 11; i++ {
		digit := int(s[i] - '0')
		sum += digit * myNumberWeights[i]
	}

	remainder := sum % 11
	check := 11 - remainder
	if check >= 10 {
		check = 0
	}

	return int(s[11]-'0') == check
}
