package nl

import (
	"context"
	"regexp"
	"strings"

	"github.com/taoq-ai/wuming/domain/model"
)

// KvK number: 8-digit number.
var kvkRe = regexp.MustCompile(`\b\d{8}\b`)

// kvkPrefixRe matches common KvK prefixes that increase confidence.
var kvkPrefixRe = regexp.MustCompile(`(?i)(?:kvk|kamer\s+van\s+koophandel)\s*[:.]?\s*$`)

// KvKDetector detects Dutch Kamer van Koophandel (Chamber of Commerce) numbers.
type KvKDetector struct{}

func NewKvKDetector() *KvKDetector { return &KvKDetector{} }

func (d *KvKDetector) Name() string              { return "nl/kvk" }
func (d *KvKDetector) Locales() []string         { return []string{locale} }
func (d *KvKDetector) PIITypes() []model.PIIType { return []model.PIIType{model.TaxID} }

func (d *KvKDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := kvkRe.FindAllStringIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	lower := strings.ToLower(text)
	var matches []model.Match
	for _, loc := range results {
		confidence := 0.60
		// Check for KvK prefix before the number to increase confidence.
		prefix := lower[:loc[0]]
		if kvkPrefixRe.MatchString(prefix) {
			confidence = 0.90
		}
		matches = append(matches, model.Match{
			Type:       model.TaxID,
			Value:      text[loc[0]:loc[1]],
			Start:      loc[0],
			End:        loc[1],
			Confidence: confidence,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}
