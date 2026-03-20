// Package structured provides scanners for detecting and redacting PII in
// structured data formats such as JSON and CSV. Each scanner walks the
// structure field-by-field, delegates text-level detection to the core
// engine, and annotates matches with structural context (JSON path,
// row/column index).
package structured

import (
	"context"

	"github.com/taoq-ai/wuming/domain/model"
	"github.com/taoq-ai/wuming/domain/port"
)

// FieldMatch extends a model.Match with the structural path where it was found.
type FieldMatch struct {
	model.Match
	// Path is the structural location of the match, e.g. "user.email" for
	// JSON or "R2:C3" for CSV.
	Path string
}

// Result holds the outcome of scanning a structured document.
type Result struct {
	// Data is the redacted output (JSON bytes or CSV string).
	Data []byte
	// Matches contains all PII detections annotated with field paths.
	Matches []FieldMatch
	// MatchCount is the total number of PII detections.
	MatchCount int
}

// scanner is shared logic used by both JSON and CSV scanners.
type scanner struct {
	pipeline port.Pipeline
}

// detectAndRedact runs PII detection on a single text value, returning the
// redacted text and annotated field matches.
func (s *scanner) detectAndRedact(ctx context.Context, text, path string) (string, []FieldMatch, error) {
	result, err := s.pipeline.Process(ctx, text)
	if err != nil {
		return "", nil, err
	}

	var fm []FieldMatch
	for _, m := range result.Matches {
		fm = append(fm, FieldMatch{Match: m, Path: path})
	}
	return result.Redacted, fm, nil
}
