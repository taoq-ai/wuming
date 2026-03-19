// Package au provides PII detectors for Australian-specific patterns.
package au

import (
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

const locale = "au"

// findAll returns all non-overlapping matches of re in text.
func findAll(re *regexp.Regexp, text string, piiType model.PIIType, confidence float64, detector string) []model.Match {
	locs := re.FindAllStringIndex(text, -1)
	if len(locs) == 0 {
		return nil
	}
	matches := make([]model.Match, 0, len(locs))
	for _, loc := range locs {
		matches = append(matches, model.Match{
			Type:       piiType,
			Value:      text[loc[0]:loc[1]],
			Start:      loc[0],
			End:        loc[1],
			Confidence: confidence,
			Locale:     locale,
			Detector:   detector,
		})
	}
	return matches
}
