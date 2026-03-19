package au

import (
	"context"
	"regexp"
	"strconv"

	"github.com/taoq-ai/wuming/domain/model"
)

// 4-digit postcode pattern.
var postcodeRe = regexp.MustCompile(`\b\d{4}\b`)

// PostcodeDetector detects Australian postcodes.
type PostcodeDetector struct{}

func NewPostcodeDetector() *PostcodeDetector { return &PostcodeDetector{} }

func (d *PostcodeDetector) Name() string              { return "au/postcode" }
func (d *PostcodeDetector) Locales() []string         { return []string{locale} }
func (d *PostcodeDetector) PIITypes() []model.PIIType { return []model.PIIType{model.PostalCode} }

func (d *PostcodeDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	locs := postcodeRe.FindAllStringIndex(text, -1)
	if len(locs) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range locs {
		raw := text[loc[0]:loc[1]]
		if !isValidAUPostcode(raw) {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.PostalCode,
			Value:      raw,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.55,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}

// isValidAUPostcode checks whether a 4-digit string falls within a valid
// Australian postcode range.
func isValidAUPostcode(s string) bool {
	n, err := strconv.Atoi(s)
	if err != nil {
		return false
	}

	switch {
	case n >= 800 && n <= 999: // NT
		return true
	case n >= 1000 && n <= 2999: // NSW
		return true
	case n >= 3000 && n <= 3999: // VIC
		return true
	case n >= 4000 && n <= 4999: // QLD
		return true
	case n >= 5000 && n <= 5999: // SA
		return true
	case n >= 6000 && n <= 6999: // WA
		return true
	case n >= 7000 && n <= 7999: // TAS
		return true
	default:
		return false
	}
	// Note: ACT postcodes (2600-2639, 2900-2920) are covered by the NSW
	// 1000-2999 range above, so they are valid.
}
