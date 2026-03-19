// Package port defines the hexagonal architecture port interfaces.
// These interfaces decouple the domain logic from concrete adapter implementations.
package port

import (
	"context"

	"github.com/taoq-ai/wuming/domain/model"
)

// Detector scans text and returns all PII matches found.
type Detector interface {
	// Detect scans text and returns all PII matches found.
	Detect(ctx context.Context, text string) ([]model.Match, error)
	// Name returns a unique identifier for this detector.
	Name() string
	// Locales returns which locales this detector supports.
	// An empty slice means the detector is locale-independent (global).
	Locales() []string
	// PIITypes returns which PII types this detector can find.
	PIITypes() []model.PIIType
}
