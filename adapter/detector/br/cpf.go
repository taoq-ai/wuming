package br

import (
	"context"
	"regexp"
	"strings"

	"github.com/taoq-ai/wuming/domain/model"
)

// Matches CPF in XXX.XXX.XXX-XX or 11-digit unformatted form.
var cpfRe = regexp.MustCompile(`\b(\d{3})\.(\d{3})\.(\d{3})-(\d{2})\b|\b(\d{11})\b`)

// CPFDetector detects Brazilian CPF numbers.
type CPFDetector struct{}

func NewCPFDetector() *CPFDetector { return &CPFDetector{} }

func (d *CPFDetector) Name() string              { return "br/cpf" }
func (d *CPFDetector) Locales() []string         { return []string{locale} }
func (d *CPFDetector) PIITypes() []model.PIIType { return []model.PIIType{model.TaxID} }

func (d *CPFDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := cpfRe.FindAllStringIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		raw := text[loc[0]:loc[1]]
		digits := stripNonDigits(raw)
		if len(digits) != 11 {
			continue
		}
		if !isValidCPF(digits) {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.TaxID,
			Value:      raw,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.90,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}

// isValidCPF validates a CPF using the mod-11 double check digit algorithm.
func isValidCPF(digits string) bool {
	// Reject all-same-digit sequences (e.g. 111.111.111-11).
	if allSameDigit(digits) {
		return false
	}

	d := make([]int, 11)
	for i, c := range digits {
		d[i] = int(c - '0')
	}

	// First check digit.
	sum := 0
	for i := 0; i < 9; i++ {
		sum += d[i] * (10 - i)
	}
	rem := sum % 11
	check1 := 0
	if rem >= 2 {
		check1 = 11 - rem
	}
	if d[9] != check1 {
		return false
	}

	// Second check digit.
	sum = 0
	for i := 0; i < 10; i++ {
		sum += d[i] * (11 - i)
	}
	rem = sum % 11
	check2 := 0
	if rem >= 2 {
		check2 = 11 - rem
	}
	return d[10] == check2
}

func stripNonDigits(s string) string {
	var b strings.Builder
	for _, c := range s {
		if c >= '0' && c <= '9' {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func allSameDigit(s string) bool {
	if len(s) == 0 {
		return true
	}
	first := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] != first {
			return false
		}
	}
	return true
}
