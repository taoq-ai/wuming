package nl

import (
	"context"
	"regexp"
	"strings"

	"github.com/taoq-ai/wuming/domain/model"
)

// Dutch postal code: 4 digits + optional space + 2 uppercase letters.
// First digit cannot be 0.
var postalRe = regexp.MustCompile(`\b[1-9]\d{3}\s?[A-Z]{2}\b`)

// invalidLetterCombos contains letter combinations not used in Dutch postal codes.
var invalidLetterCombos = map[string]bool{
	"SA": true,
	"SD": true,
	"SS": true,
}

// PostalDetector detects Dutch postal codes.
type PostalDetector struct{}

func NewPostalDetector() *PostalDetector { return &PostalDetector{} }

func (d *PostalDetector) Name() string              { return "nl/postal" }
func (d *PostalDetector) Locales() []string         { return []string{locale} }
func (d *PostalDetector) PIITypes() []model.PIIType { return []model.PIIType{model.PostalCode} }

func (d *PostalDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := postalRe.FindAllStringIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		value := text[loc[0]:loc[1]]
		// Extract the 2-letter suffix.
		letters := strings.TrimSpace(value)
		letterPart := letters[len(letters)-2:]
		if invalidLetterCombos[letterPart] {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.PostalCode,
			Value:      value,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.90,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}
