package port

import (
	"context"

	"github.com/taoq-ai/wuming/domain/model"
)

// Result holds the outcome of a full detection + replacement pipeline run.
type Result struct {
	// Original is the input text before any redaction.
	Original string
	// Redacted is the text after PII has been replaced.
	Redacted string
	// Matches contains all PII detections found.
	Matches []model.Match
	// MatchCount is the total number of matches found.
	MatchCount int
}

// Pipeline orchestrates detection and replacement in sequence.
type Pipeline interface {
	// Process runs detection + replacement and returns the result.
	Process(ctx context.Context, text string) (*Result, error)
}
