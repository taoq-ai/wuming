package de

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/taoq-ai/wuming/domain/model"
)

// PLZ: German 5-digit postal code.
var plzRe = regexp.MustCompile(`\b\d{5}\b`)

// Context keywords that increase confidence when found near a PLZ.
var plzContextRe = regexp.MustCompile(`(?i)\b(?:PLZ|Postleitzahl|Straße|Str\.|Stadt|Ort|Adresse|Anschrift)\b`)

// PLZDetector detects German postal codes (Postleitzahl).
type PLZDetector struct{}

func NewPLZDetector() *PLZDetector { return &PLZDetector{} }

func (d *PLZDetector) Name() string              { return "de/plz" }
func (d *PLZDetector) Locales() []string         { return []string{locale} }
func (d *PLZDetector) PIITypes() []model.PIIType { return []model.PIIType{model.PostalCode} }

func (d *PLZDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := plzRe.FindAllStringIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	// Pre-check if any context keyword exists in the text.
	hasContext := plzContextRe.MatchString(text)

	var matches []model.Match
	for _, loc := range results {
		value := text[loc[0]:loc[1]]
		n, _ := strconv.Atoi(value)
		if n < 1001 || n > 99998 {
			continue
		}

		confidence := 0.60
		if hasContext {
			confidence = 0.80
		}
		// Also boost if "PLZ" or similar appears right before the number.
		if loc[0] >= 4 {
			prefix := strings.ToUpper(text[loc[0]-4 : loc[0]])
			if strings.Contains(prefix, "PLZ") {
				confidence = 0.85
			}
		}

		matches = append(matches, model.Match{
			Type:       model.PostalCode,
			Value:      value,
			Start:      loc[0],
			End:        loc[1],
			Confidence: confidence,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}
