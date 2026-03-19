package common

import (
	"context"
	"regexp"

	"github.com/taoq-ai/wuming/domain/model"
)

var urlRe = regexp.MustCompile(`https?://[^\s<>"{}|\\^\x60\[\]]+`)

// URLDetector detects URLs with http/https schemes.
type URLDetector struct{}

func NewURLDetector() *URLDetector { return &URLDetector{} }

func (d *URLDetector) Name() string              { return "common/url" }
func (d *URLDetector) Locales() []string         { return nil }
func (d *URLDetector) PIITypes() []model.PIIType { return []model.PIIType{model.URL} }

func (d *URLDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	return findAll(urlRe, text, model.URL, 0.9, d.Name()), nil
}
