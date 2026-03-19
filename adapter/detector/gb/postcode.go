package gb

import (
	"context"
	"regexp"
	"strings"

	"github.com/taoq-ai/wuming/domain/model"
)

// UK postcode formats:
// A9 9AA, A99 9AA, A9A 9AA, AA9 9AA, AA99 9AA, AA9A 9AA
var postcodeRe = regexp.MustCompile(`(?i)\b[A-Z]{1,2}\d[A-Z\d]?\s*\d[A-Z]{2}\b`)

// PostcodeDetector detects UK postcodes.
type PostcodeDetector struct{}

func NewPostcodeDetector() *PostcodeDetector { return &PostcodeDetector{} }

func (d *PostcodeDetector) Name() string              { return "gb/postcode" }
func (d *PostcodeDetector) Locales() []string         { return []string{locale} }
func (d *PostcodeDetector) PIITypes() []model.PIIType { return []model.PIIType{model.PostalCode} }

func (d *PostcodeDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	locs := postcodeRe.FindAllStringIndex(text, -1)
	if len(locs) == 0 {
		return nil, nil
	}
	matches := make([]model.Match, 0, len(locs))
	for _, loc := range locs {
		matches = append(matches, model.Match{
			Type:       model.PostalCode,
			Value:      strings.ToUpper(text[loc[0]:loc[1]]),
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.90,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}
