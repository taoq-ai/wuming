// Package common provides PII detectors for locale-independent patterns.
package common

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

var emailRe = regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)

// EmailDetector detects email addresses.
type EmailDetector struct{}

func NewEmailDetector() *EmailDetector { return &EmailDetector{} }

func (d *EmailDetector) Name() string              { return "common/email" }
func (d *EmailDetector) Locales() []string         { return nil }
func (d *EmailDetector) PIITypes() []model.PIIType { return []model.PIIType{model.Email} }

func (d *EmailDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(emailRe, text, model.Email, 0.95, d.Name()), nil
}
