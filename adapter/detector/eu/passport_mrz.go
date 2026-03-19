package eu

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

// mrzRe matches ICAO 9303 TD3 (passport) Machine Readable Zone.
// TD3 format consists of two lines of 44 characters each.
// Line 1: document type, issuing state, name.
// Line 2: passport number, check digit, nationality, DOB, sex, expiry, optional data.
var mrzRe = regexp.MustCompile(
	`(?m)P[<A-Z][A-Z]{3}[A-Z<]{39}\n[A-Z0-9<]{9}[0-9][A-Z]{3}[0-9]{6}[0-9][MF<][0-9]{6}[0-9][A-Z0-9<]{14}[0-9][0-9]`,
)

// PassportMRZDetector detects ICAO 9303 passport Machine Readable Zones.
type PassportMRZDetector struct{}

func NewPassportMRZDetector() *PassportMRZDetector { return &PassportMRZDetector{} }

func (d *PassportMRZDetector) Name() string              { return "eu/passport_mrz" }
func (d *PassportMRZDetector) Locales() []string         { return []string{locale} }
func (d *PassportMRZDetector) PIITypes() []model.PIIType { return []model.PIIType{model.Passport} }

func (d *PassportMRZDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(mrzRe, text, model.Passport, 0.95, d.Name()), nil
}
