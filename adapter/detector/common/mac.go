package common

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// Matches colon-separated, hyphen-separated, and dot-separated MAC formats.
var macRe = regexp.MustCompile(`(?i)\b(?:[0-9a-f]{2}[:\-]){5}[0-9a-f]{2}\b|(?i)\b(?:[0-9a-f]{4}\.){2}[0-9a-f]{4}\b`)

// MACDetector detects MAC addresses.
type MACDetector struct{}

func NewMACDetector() *MACDetector { return &MACDetector{} }

func (d *MACDetector) Name() string              { return "common/mac" }
func (d *MACDetector) Locales() []string         { return nil }
func (d *MACDetector) PIITypes() []model.PIIType { return []model.PIIType{model.MACAddress} }

func (d *MACDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(macRe, text, model.MACAddress, 0.9, d.Name()), nil
}
