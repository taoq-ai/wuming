package common

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/taoq-ai/wuming/domain/model"
)

var (
	ipv4Re = regexp.MustCompile(`\b(?:(?:25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(?:25[0-5]|2[0-4]\d|[01]?\d\d?)\b`)
	ipv6Re = regexp.MustCompile(`(?i)(?:\b(?:[0-9a-f]{1,4}:){7}[0-9a-f]{1,4}\b|(?:\b[0-9a-f]{1,4}:){1,7}:|(?:\b[0-9a-f]{1,4}:){1,6}:[0-9a-f]{1,4}\b|::(?:[0-9a-f]{1,4}:){0,5}[0-9a-f]{1,4}\b|::)`)
)

// IPDetector detects IPv4 and IPv6 addresses.
type IPDetector struct{}

func NewIPDetector() *IPDetector { return &IPDetector{} }

func (d *IPDetector) Name() string              { return "common/ip" }
func (d *IPDetector) Locales() []string         { return nil }
func (d *IPDetector) PIITypes() []model.PIIType { return []model.PIIType{model.IPAddress} }

func (d *IPDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	var matches []model.Match

	// IPv4
	for _, loc := range ipv4Re.FindAllStringIndex(text, -1) {
		val := text[loc[0]:loc[1]]
		if isValidIPv4(val) {
			matches = append(matches, model.Match{
				Type:       model.IPAddress,
				Value:      val,
				Start:      loc[0],
				End:        loc[1],
				Confidence: 0.9,
				Detector:   d.Name(),
			})
		}
	}

	// IPv6
	matches = append(matches, findAll(ipv6Re, text, model.IPAddress, 0.9, d.Name())...)

	return matches, nil
}

func isValidIPv4(s string) bool {
	parts := strings.Split(s, ".")
	if len(parts) != 4 {
		return false
	}
	for _, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 || n > 255 {
			return false
		}
	}
	return true
}
