package cn

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Chinese mobile phone: 11 digits starting with 1, second digit 3-9.
// Supports optional +86 country code, optional spaces/dashes between groups.
var phoneRe = regexp.MustCompile(`(?:\+86[\s-]?)?1[3-9]\d[\s-]?\d{4}[\s-]?\d{4}\b`)

// PhoneDetector detects Chinese phone numbers.
type PhoneDetector struct{}

func NewPhoneDetector() *PhoneDetector { return &PhoneDetector{} }

func (d *PhoneDetector) Name() string              { return "cn/phone" }
func (d *PhoneDetector) Locales() []string         { return []string{locale} }
func (d *PhoneDetector) PIITypes() []model.PIIType { return []model.PIIType{model.Phone} }

func (d *PhoneDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(phoneRe, text, model.Phone, 0.85, d.Name()), nil
}
