package ca

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// NANP phone: optional +1/1 prefix, area code, 7-digit subscriber number.
var caPhoneRe = regexp.MustCompile(`(?:\+?1[\s.\-]?)?\(?([2-9]\d{2})\)?[\s.\-]?\d{3}[\s.\-]?\d{4}`)

// canadianAreaCodes contains all Canadian area codes.
var canadianAreaCodes = map[string]bool{
	"204": true, "226": true, "236": true, "249": true, "250": true,
	"289": true, "306": true, "343": true, "365": true, "387": true,
	"403": true, "416": true, "418": true, "431": true, "437": true,
	"438": true, "450": true, "506": true, "514": true, "519": true,
	"548": true, "579": true, "581": true, "587": true, "604": true,
	"613": true, "639": true, "647": true, "672": true, "705": true,
	"709": true, "742": true, "778": true, "780": true, "782": true,
	"807": true, "819": true, "825": true, "867": true, "873": true,
	"902": true, "905": true,
}

// PhoneDetector detects Canadian phone numbers in NANP format.
type PhoneDetector struct{}

func NewPhoneDetector() *PhoneDetector { return &PhoneDetector{} }

func (d *PhoneDetector) Name() string              { return "ca/phone" }
func (d *PhoneDetector) Locales() []string         { return []string{locale} }
func (d *PhoneDetector) PIITypes() []model.PIIType { return []model.PIIType{model.Phone} }

func (d *PhoneDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := caPhoneRe.FindAllStringSubmatchIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		areaCode := text[loc[2]:loc[3]]
		if !canadianAreaCodes[areaCode] {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.Phone,
			Value:      text[loc[0]:loc[1]],
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.80,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}
