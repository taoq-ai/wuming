package common

import (
	"context"
	"regexp"
	"strings"

	"github.com/taoq-ai/wuming/domain/model"
)

// Matches 13-19 digit sequences optionally separated by spaces or hyphens.
var creditCardRe = regexp.MustCompile(`\b(?:\d[ \-]*?){13,19}\b`)

// CreditCardDetector detects credit card numbers using pattern matching and Luhn validation.
type CreditCardDetector struct{}

func NewCreditCardDetector() *CreditCardDetector { return &CreditCardDetector{} }

func (d *CreditCardDetector) Name() string              { return "common/creditcard" }
func (d *CreditCardDetector) Locales() []string         { return nil }
func (d *CreditCardDetector) PIITypes() []model.PIIType { return []model.PIIType{model.CreditCard} }

func (d *CreditCardDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	locs := creditCardRe.FindAllStringIndex(text, -1)
	if len(locs) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range locs {
		raw := text[loc[0]:loc[1]]
		digits := stripNonDigits(raw)
		if len(digits) < 13 || len(digits) > 19 {
			continue
		}
		if !luhn(digits) {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.CreditCard,
			Value:      raw,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.95,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}

func stripNonDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// luhn validates a number string using the Luhn algorithm.
func luhn(digits string) bool {
	if len(digits) == 0 {
		return false
	}
	sum := 0
	alt := false
	for i := len(digits) - 1; i >= 0; i-- {
		n := int(digits[i] - '0')
		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}
	return sum%10 == 0
}
