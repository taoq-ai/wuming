package common

import (
	"context"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/taoq-ai/wuming/domain/model"
)

// IBAN: 2 letter country code + 2 check digits + up to 30 alphanumeric chars.
var ibanRe = regexp.MustCompile(`\b[A-Z]{2}\d{2}[A-Z0-9]{4,30}\b`)

// IBANDetector detects International Bank Account Numbers with mod-97 validation.
type IBANDetector struct{}

func NewIBANDetector() *IBANDetector { return &IBANDetector{} }

func (d *IBANDetector) Name() string              { return "common/iban" }
func (d *IBANDetector) Locales() []string         { return nil }
func (d *IBANDetector) PIITypes() []model.PIIType { return []model.PIIType{model.IBAN} }

func (d *IBANDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	locs := ibanRe.FindAllStringIndex(text, -1)
	if len(locs) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range locs {
		val := text[loc[0]:loc[1]]
		if !validateIBANChecksum(val) {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.IBAN,
			Value:      val,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.95,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}

// validateIBANChecksum performs the ISO 13616 mod-97 check.
func validateIBANChecksum(iban string) bool {
	iban = strings.ToUpper(strings.ReplaceAll(iban, " ", ""))
	if len(iban) < 5 {
		return false
	}

	// Move the first 4 characters to the end.
	rearranged := iban[4:] + iban[:4]

	// Convert letters to numbers (A=10, B=11, ..., Z=35).
	var numeric strings.Builder
	for _, r := range rearranged {
		if r >= 'A' && r <= 'Z' {
			numeric.WriteString(strconv.Itoa(int(r - 'A' + 10)))
		} else {
			numeric.WriteRune(r)
		}
	}

	n := new(big.Int)
	n.SetString(numeric.String(), 10)
	mod := new(big.Int)
	mod.Mod(n, big.NewInt(97))

	return mod.Int64() == 1
}
