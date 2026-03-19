package gb

import (
	"context"
	"regexp"
	"strings"

	"github.com/taoq-ai/wuming/domain/model"
)

// UTR: exactly 10 digits.
var utrRe = regexp.MustCompile(`\b\d{10}\b`)

// utrContextRe detects context keywords preceding a 10-digit number.
var utrContextRe = regexp.MustCompile(`(?i)(?:UTR|tax\s+reference|unique\s+taxpayer\s+reference)\s*:?\s*(\d{10})\b`)

// UTRDetector detects UK Unique Taxpayer References.
type UTRDetector struct{}

func NewUTRDetector() *UTRDetector { return &UTRDetector{} }

func (d *UTRDetector) Name() string              { return "gb/utr" }
func (d *UTRDetector) Locales() []string         { return []string{locale} }
func (d *UTRDetector) PIITypes() []model.PIIType { return []model.PIIType{model.TaxID} }

func (d *UTRDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	// First, find high-confidence matches with context keywords.
	contextMatches := make(map[string]bool)
	contextResults := utrContextRe.FindAllStringSubmatchIndex(text, -1)
	var matches []model.Match

	for _, loc := range contextResults {
		// loc[0]:loc[1] is the full match; loc[2]:loc[3] is the 10-digit group.
		digits := text[loc[2]:loc[3]]
		contextMatches[digits] = true
		matches = append(matches, model.Match{
			Type:       model.TaxID,
			Value:      strings.TrimSpace(text[loc[0]:loc[1]]),
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.85,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}

	// Then find bare 10-digit numbers that were not already matched with context.
	bareLocs := utrRe.FindAllStringIndex(text, -1)
	for _, loc := range bareLocs {
		digits := text[loc[0]:loc[1]]
		if contextMatches[digits] {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.TaxID,
			Value:      digits,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.55,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}

	return matches, nil
}
