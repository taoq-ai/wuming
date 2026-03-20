package engine

import (
	"context"
	"testing"

	"github.com/taoq-ai/wuming/adapter/replacer"
	"github.com/taoq-ai/wuming/domain/model"
	"github.com/taoq-ai/wuming/domain/port"
)

// stubDetector is a test detector that returns pre-configured matches.
type stubDetector struct {
	name    string
	locales []string
	types   []model.PIIType
	matches []model.Match
}

func (s *stubDetector) Detect(_ context.Context, _ string) ([]model.Match, error) {
	return s.matches, nil
}
func (s *stubDetector) Name() string              { return s.name }
func (s *stubDetector) Locales() []string         { return s.locales }
func (s *stubDetector) PIITypes() []model.PIIType { return s.types }

func TestEngineProcess(t *testing.T) {
	text := "Email me at john@example.com or call 555-1234."

	emailDetector := &stubDetector{
		name:  "email",
		types: []model.PIIType{model.Email},
		matches: []model.Match{
			{Type: model.Email, Value: "john@example.com", Start: 12, End: 28, Confidence: 1.0, Detector: "email"},
		},
	}
	phoneDetector := &stubDetector{
		name:  "phone",
		types: []model.PIIType{model.Phone},
		matches: []model.Match{
			{Type: model.Phone, Value: "555-1234", Start: 37, End: 45, Confidence: 0.9, Detector: "phone"},
		},
	}

	e := New(
		WithDetectors(emailDetector, phoneDetector),
		WithReplacer(replacer.NewRedact()),
	)

	result, err := e.Process(context.Background(), text)
	if err != nil {
		t.Fatal(err)
	}

	if result.MatchCount != 2 {
		t.Errorf("MatchCount = %d, want 2", result.MatchCount)
	}

	want := "Email me at [EMAIL] or call [PHONE]."
	if result.Redacted != want {
		t.Errorf("Redacted:\n got: %q\nwant: %q", result.Redacted, want)
	}
}

func TestEngineConfidenceThreshold(t *testing.T) {
	detector := &stubDetector{
		name: "mixed",
		matches: []model.Match{
			{Type: model.Email, Value: "a@b.com", Start: 0, End: 7, Confidence: 0.95},
			{Type: model.Phone, Value: "12345", Start: 10, End: 15, Confidence: 0.3},
		},
	}

	e := New(
		WithDetectors(detector),
		WithReplacer(replacer.NewRedact()),
		WithConfidenceThreshold(0.5),
	)

	result, err := e.Process(context.Background(), "a@b.com or 12345 maybe")
	if err != nil {
		t.Fatal(err)
	}

	if result.MatchCount != 1 {
		t.Errorf("MatchCount = %d, want 1 (low confidence should be filtered)", result.MatchCount)
	}
}

func TestEngineLocaleFiltering(t *testing.T) {
	globalDetector := &stubDetector{
		name:    "global-email",
		locales: nil, // global
		matches: []model.Match{
			{Type: model.Email, Value: "a@b.com", Start: 0, End: 7, Confidence: 1.0},
		},
	}
	nlDetector := &stubDetector{
		name:    "nl-bsn",
		locales: []string{"nl"},
		matches: []model.Match{
			{Type: model.NationalID, Value: "123456789", Start: 10, End: 19, Confidence: 0.9},
		},
	}

	e := New(
		WithDetectors(globalDetector, nlDetector),
		WithReplacer(replacer.NewRedact()),
		WithLocales("us"),
	)

	result, err := e.Process(context.Background(), "a@b.com + 123456789")
	if err != nil {
		t.Fatal(err)
	}

	// Only global detector should run; nl-bsn should be filtered out.
	if result.MatchCount != 1 {
		t.Errorf("MatchCount = %d, want 1 (NL detector should be filtered for US locale)", result.MatchCount)
	}
}

func TestEngineDedup(t *testing.T) {
	d1 := &stubDetector{
		name: "d1",
		matches: []model.Match{
			{Type: model.Email, Value: "a@b.com", Start: 0, End: 7, Confidence: 0.9},
		},
	}
	d2 := &stubDetector{
		name: "d2",
		matches: []model.Match{
			{Type: model.Email, Value: "a@b.com", Start: 0, End: 7, Confidence: 0.95},
		},
	}

	e := New(WithDetectors(d1, d2), WithReplacer(replacer.NewRedact()))

	result, err := e.Process(context.Background(), "a@b.com")
	if err != nil {
		t.Fatal(err)
	}

	if result.MatchCount != 1 {
		t.Errorf("MatchCount = %d, want 1 (overlapping matches should be deduped)", result.MatchCount)
	}
}

func TestEngineNoReplacer(t *testing.T) {
	detector := &stubDetector{
		name: "email",
		matches: []model.Match{
			{Type: model.Email, Value: "a@b.com", Start: 0, End: 7, Confidence: 1.0},
		},
	}

	e := New(WithDetectors(detector))
	result, err := e.Process(context.Background(), "a@b.com")
	if err != nil {
		t.Fatal(err)
	}

	if result.MatchCount != 1 {
		t.Errorf("MatchCount = %d, want 1", result.MatchCount)
	}
	if result.Redacted != "a@b.com" {
		t.Error("Without a replacer, redacted should equal original")
	}
}

func TestEngineAllowlist(t *testing.T) {
	detector := &stubDetector{
		name: "mixed",
		matches: []model.Match{
			{Type: model.Email, Value: "john@example.com", Start: 0, End: 16, Confidence: 1.0},
			{Type: model.URL, Value: "example.com", Start: 20, End: 31, Confidence: 0.9},
		},
	}

	e := New(
		WithDetectors(detector),
		WithReplacer(replacer.NewRedact()),
		WithAllowlist("example.com"),
	)

	result, err := e.Process(context.Background(), "john@example.com or example.com here")
	if err != nil {
		t.Fatal(err)
	}

	if result.MatchCount != 1 {
		t.Errorf("MatchCount = %d, want 1 (allowlisted value should be filtered)", result.MatchCount)
	}
	if result.Matches[0].Type != model.Email {
		t.Errorf("remaining match type = %v, want Email", result.Matches[0].Type)
	}
}

func TestEngineAllowlistCaseInsensitive(t *testing.T) {
	detector := &stubDetector{
		name: "url",
		matches: []model.Match{
			{Type: model.URL, Value: "Example.COM", Start: 0, End: 11, Confidence: 0.9},
		},
	}

	e := New(
		WithDetectors(detector),
		WithReplacer(replacer.NewRedact()),
		WithAllowlist("example.com"),
	)

	result, err := e.Process(context.Background(), "Example.COM")
	if err != nil {
		t.Fatal(err)
	}

	if result.MatchCount != 0 {
		t.Errorf("MatchCount = %d, want 0 (allowlist should be case-insensitive)", result.MatchCount)
	}
}

func TestEngineDenylist(t *testing.T) {
	e := New(
		WithDetectors(&stubDetector{name: "noop"}),
		WithReplacer(replacer.NewRedact()),
		WithDenylist(model.Name, "ACME Corp"),
	)

	text := "Contact ACME Corp for details"
	result, err := e.Process(context.Background(), text)
	if err != nil {
		t.Fatal(err)
	}

	if result.MatchCount != 1 {
		t.Fatalf("MatchCount = %d, want 1", result.MatchCount)
	}
	m := result.Matches[0]
	if m.Type != model.Name {
		t.Errorf("match type = %v, want Name", m.Type)
	}
	if m.Value != "ACME Corp" {
		t.Errorf("match value = %q, want %q", m.Value, "ACME Corp")
	}
	if m.Detector != "denylist" {
		t.Errorf("match detector = %q, want %q", m.Detector, "denylist")
	}

	want := "Contact [NAME] for details"
	if result.Redacted != want {
		t.Errorf("Redacted:\n got: %q\nwant: %q", result.Redacted, want)
	}
}

func TestEngineDenylistMultipleOccurrences(t *testing.T) {
	e := New(
		WithDetectors(&stubDetector{name: "noop"}),
		WithReplacer(replacer.NewRedact()),
		WithDenylist(model.Name, "secret"),
	)

	text := "secret and secret again"
	result, err := e.Process(context.Background(), text)
	if err != nil {
		t.Fatal(err)
	}

	if result.MatchCount != 2 {
		t.Errorf("MatchCount = %d, want 2 (denylist should match all occurrences)", result.MatchCount)
	}
}

func TestEngineDenylistCaseInsensitive(t *testing.T) {
	e := New(
		WithDetectors(&stubDetector{name: "noop"}),
		WithReplacer(replacer.NewRedact()),
		WithDenylist(model.Name, "Secret"),
	)

	text := "Found SECRET here"
	result, err := e.Process(context.Background(), text)
	if err != nil {
		t.Fatal(err)
	}

	if result.MatchCount != 1 {
		t.Errorf("MatchCount = %d, want 1 (denylist should be case-insensitive)", result.MatchCount)
	}
}

func TestEngineAllowlistAndDenylistCombined(t *testing.T) {
	detector := &stubDetector{
		name: "url",
		matches: []model.Match{
			{Type: model.URL, Value: "example.com", Start: 0, End: 11, Confidence: 0.9},
		},
	}

	e := New(
		WithDetectors(detector),
		WithReplacer(replacer.NewRedact()),
		WithAllowlist("example.com"),
		WithDenylist(model.Name, "ACME"),
	)

	text := "example.com by ACME"
	result, err := e.Process(context.Background(), text)
	if err != nil {
		t.Fatal(err)
	}

	// example.com should be filtered (allowlisted), ACME should be injected (denylisted).
	if result.MatchCount != 1 {
		t.Fatalf("MatchCount = %d, want 1", result.MatchCount)
	}
	if result.Matches[0].Value != "ACME" {
		t.Errorf("match value = %q, want %q", result.Matches[0].Value, "ACME")
	}
}

// Verify Engine implements port.Pipeline.
var _ port.Pipeline = (*Engine)(nil)
