package fr

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// CNI (Carte Nationale d'Identite):
//   - Old format (before 2021): 12 digits
//   - New format (since 2021): 9 alphanumeric characters
var (
	cniOldRe = regexp.MustCompile(`\b\d{12}\b`)
	cniNewRe = regexp.MustCompile(`\b[A-Z0-9]{9}\b`)
)

// IDCardDetector detects French CNI (national identity card) numbers.
type IDCardDetector struct{}

func NewIDCardDetector() *IDCardDetector { return &IDCardDetector{} }

func (d *IDCardDetector) Name() string              { return "fr/id_card" }
func (d *IDCardDetector) Locales() []string         { return []string{locale} }
func (d *IDCardDetector) PIITypes() []model.PIIType { return []model.PIIType{model.NationalID} }

func (d *IDCardDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	var matches []model.Match
	matches = append(matches, findAll(cniOldRe, text, model.NationalID, 0.65, d.Name())...)
	matches = append(matches, findAll(cniNewRe, text, model.NationalID, 0.65, d.Name())...)

	// Deduplicate: old-format 12-digit matches may overlap with new-format 9-char matches.
	// Since old-format is strictly digits and 12 chars, and new is 9 chars, no real overlap.
	return matches, nil
}
