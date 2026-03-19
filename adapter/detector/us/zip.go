package us

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// ZIP code: 5 digits, optionally followed by -XXXX (ZIP+4).
var zipRe = regexp.MustCompile(`\b\d{5}(?:-\d{4})?\b`)

// ZIPDetector detects US ZIP codes.
type ZIPDetector struct{}

func NewZIPDetector() *ZIPDetector { return &ZIPDetector{} }

func (d *ZIPDetector) Name() string              { return "us/zip" }
func (d *ZIPDetector) Locales() []string         { return []string{locale} }
func (d *ZIPDetector) PIITypes() []model.PIIType { return []model.PIIType{model.PostalCode} }

func (d *ZIPDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(zipRe, text, model.PostalCode, 0.6, d.Name()), nil
}
