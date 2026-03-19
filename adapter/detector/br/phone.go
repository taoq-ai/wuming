package br

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Brazilian phone patterns:
// Mobile:   (XX) 9XXXX-XXXX, +55 XX 9XXXX-XXXX, +55XX9XXXXXXXX
// Landline: (XX) XXXX-XXXX,  +55 XX XXXX-XXXX,  +55XXXXXXXXXX
var brPhoneRe = regexp.MustCompile(
	`(?:\+55[\s\-]?)?\(?[1-9]\d\)?[\s\-]?9\d{4}[\s\-]?\d{4}` + // mobile
		`|(?:\+55[\s\-]?)?\(?[1-9]\d\)?[\s\-]?[2-5]\d{3}[\s\-]?\d{4}`, // landline
)

// PhoneDetector detects Brazilian phone numbers.
type PhoneDetector struct{}

func NewPhoneDetector() *PhoneDetector { return &PhoneDetector{} }

func (d *PhoneDetector) Name() string              { return "br/phone" }
func (d *PhoneDetector) Locales() []string         { return []string{locale} }
func (d *PhoneDetector) PIITypes() []model.PIIType { return []model.PIIType{model.Phone} }

func (d *PhoneDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(brPhoneRe, text, model.Phone, 0.85, d.Name()), nil
}
