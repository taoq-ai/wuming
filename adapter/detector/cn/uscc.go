package cn

import (
	"context"
	"regexp"
	"strings"

	"github.com/taoq-ai/wuming/domain/model"
)

// Unified Social Credit Code (统一社会信用代码): 18 characters.
// 2-char registration management department + 6-digit area code + 9-char organization code + 1 check char.
// Valid characters: digits and uppercase letters excluding I, O, S, V, Z.
var usccRe = regexp.MustCompile(`\b[0-9A-HJ-NP-RTUW-Y]{2}\d{6}[0-9A-HJ-NP-RTUW-Y]{10}\b`)

// usccCharset maps valid characters to their numeric values for mod-31.
const usccChars = "0123456789ABCDEFGHJKLMNPQRTUWXY"

// usccWeights are the positional weights for the mod-31 check.
var usccWeights = [17]int{1, 3, 9, 27, 19, 26, 16, 17, 20, 29, 25, 13, 8, 24, 10, 30, 28}

// USCCDetector detects Chinese Unified Social Credit Codes.
type USCCDetector struct{}

func NewUSCCDetector() *USCCDetector { return &USCCDetector{} }

func (d *USCCDetector) Name() string              { return "cn/uscc" }
func (d *USCCDetector) Locales() []string         { return []string{locale} }
func (d *USCCDetector) PIITypes() []model.PIIType { return []model.PIIType{model.TaxID} }

func (d *USCCDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := usccRe.FindAllStringIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		code := text[loc[0]:loc[1]]
		if !isValidUSCC(code) {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.TaxID,
			Value:      code,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.85,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}

func isValidUSCC(code string) bool {
	upper := strings.ToUpper(code)

	sum := 0
	for i := 0; i < 17; i++ {
		idx := strings.IndexByte(usccChars, upper[i])
		if idx < 0 {
			return false
		}
		sum += idx * usccWeights[i]
	}
	remainder := sum % 31
	checkValue := (31 - remainder) % 31

	// Map check value back to character.
	if checkValue >= len(usccChars) {
		return false
	}
	expected := usccChars[checkValue]

	return upper[17] == expected
}
