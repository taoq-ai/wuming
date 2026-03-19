package fr

import (
	"context"
	"regexp"
	"strconv"

	"github.com/taoq-ai/wuming/domain/model"
)

// Code postal: 5 digits, first 2 = department (01-95, 97, 98).
var postalRe = regexp.MustCompile(`\b(\d{2})\d{3}\b`)

// PostalDetector detects French postal codes.
type PostalDetector struct{}

func NewPostalDetector() *PostalDetector { return &PostalDetector{} }

func (d *PostalDetector) Name() string              { return "fr/postal" }
func (d *PostalDetector) Locales() []string         { return []string{locale} }
func (d *PostalDetector) PIITypes() []model.PIIType { return []model.PIIType{model.PostalCode} }

func (d *PostalDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := postalRe.FindAllStringSubmatchIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		dept := text[loc[2]:loc[3]]
		deptNum, err := strconv.Atoi(dept)
		if err != nil {
			continue
		}

		// Valid department prefixes: 01-95, 97, 98.
		if deptNum < 1 || deptNum > 98 || deptNum == 96 {
			continue
		}

		matches = append(matches, model.Match{
			Type:       model.PostalCode,
			Value:      text[loc[0]:loc[1]],
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.60,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}
