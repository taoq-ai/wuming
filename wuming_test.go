package wuming

import (
	"context"
	"testing"

	"github.com/taoq-ai/wuming/domain/model"
	"github.com/taoq-ai/wuming/domain/port"
)

type stubDetector struct {
	matches []model.Match
}

func (s *stubDetector) Detect(_ context.Context, _ string) ([]model.Match, error) {
	return s.matches, nil
}
func (s *stubDetector) Name() string              { return "stub" }
func (s *stubDetector) Locales() []string         { return nil }
func (s *stubDetector) PIITypes() []model.PIIType { return nil }

func TestWumingRedact(t *testing.T) {
	det := &stubDetector{
		matches: []model.Match{
			{Type: model.Email, Value: "john@example.com", Start: 10, End: 26, Confidence: 1.0},
		},
	}

	w := New(WithDetectors(det))
	got, err := w.Redact(context.Background(), "Email me: john@example.com please")
	if err != nil {
		t.Fatal(err)
	}

	want := "Email me: [EMAIL] please"
	if got != want {
		t.Errorf("Redact:\n got: %q\nwant: %q", got, want)
	}
}

func TestWumingDetect(t *testing.T) {
	det := &stubDetector{
		matches: []model.Match{
			{Type: model.Phone, Value: "555-1234", Start: 5, End: 13, Confidence: 0.9},
		},
	}

	w := New(WithDetectors(det))
	matches, err := w.Detect(context.Background(), "Call 555-1234")
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatalf("got %d matches, want 1", len(matches))
	}
	if matches[0].Type != model.Phone {
		t.Errorf("match type = %v, want Phone", matches[0].Type)
	}
}

func TestWumingProcess(t *testing.T) {
	det := &stubDetector{
		matches: []model.Match{
			{Type: model.Email, Value: "a@b.com", Start: 0, End: 7, Confidence: 1.0},
		},
	}

	w := New(WithDetectors(det))
	result, err := w.Process(context.Background(), "a@b.com")
	if err != nil {
		t.Fatal(err)
	}

	if result.Original != "a@b.com" {
		t.Errorf("Original = %q, want %q", result.Original, "a@b.com")
	}
	if result.Redacted != "[EMAIL]" {
		t.Errorf("Redacted = %q, want %q", result.Redacted, "[EMAIL]")
	}
	if result.MatchCount != 1 {
		t.Errorf("MatchCount = %d, want 1", result.MatchCount)
	}
}

// Verify stub implements port.Detector.
var _ port.Detector = (*stubDetector)(nil)
