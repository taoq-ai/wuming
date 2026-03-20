package wuming

import (
	"context"
	"strings"
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

	w, err := New(WithDetectors(det))
	if err != nil {
		t.Fatal(err)
	}
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

	w, err := New(WithDetectors(det))
	if err != nil {
		t.Fatal(err)
	}
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

	w, err := New(WithDetectors(det))
	if err != nil {
		t.Fatal(err)
	}
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

func TestZeroConfigNew(t *testing.T) {
	w, err := New()
	if err != nil {
		t.Fatal(err)
	}
	text := "Email john@example.com and SSN 078-05-1120"
	result, err := w.Process(context.Background(), text)
	if err != nil {
		t.Fatal(err)
	}
	if result.MatchCount == 0 {
		t.Error("zero-config New() should detect PII, but found no matches")
	}
}

func TestPackageLevelRedact(t *testing.T) {
	got, err := Redact(context.Background(), "Email john@example.com please")
	if err != nil {
		t.Fatal(err)
	}
	if got == "Email john@example.com please" {
		t.Error("package-level Redact() did not redact anything")
	}
}

func TestPackageLevelDetect(t *testing.T) {
	matches, err := Detect(context.Background(), "Email john@example.com please")
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) == 0 {
		t.Error("package-level Detect() found no matches")
	}
}

func TestWithPresetGDPR(t *testing.T) {
	w, err := New(WithPreset("gdpr"))
	if err != nil {
		t.Fatal(err)
	}
	// GDPR should detect email (common locale).
	text := "Email john@example.com"
	result, err := w.Process(context.Background(), text)
	if err != nil {
		t.Fatal(err)
	}
	if result.MatchCount == 0 {
		t.Error("WithPreset(\"gdpr\") should detect PII, but found no matches")
	}
}

func TestWithPresetHIPAA(t *testing.T) {
	w, err := New(WithPreset("hipaa"))
	if err != nil {
		t.Fatal(err)
	}
	// HIPAA should detect SSN (us locale, NationalID type).
	text := "SSN: 078-05-1120"
	result, err := w.Process(context.Background(), text)
	if err != nil {
		t.Fatal(err)
	}
	hasNationalID := false
	for _, m := range result.Matches {
		if m.Type == model.NationalID {
			hasNationalID = true
		}
	}
	if !hasNationalID {
		t.Error("WithPreset(\"hipaa\") should detect NationalID (SSN)")
	}
}

func TestWithPresetPCIDSS(t *testing.T) {
	w, err := New(WithPreset("pci-dss"))
	if err != nil {
		t.Fatal(err)
	}
	// PCI-DSS should detect credit cards.
	text := "Card: 4111 1111 1111 1111"
	result, err := w.Process(context.Background(), text)
	if err != nil {
		t.Fatal(err)
	}
	hasCreditCard := false
	for _, m := range result.Matches {
		if m.Type == model.CreditCard {
			hasCreditCard = true
		}
	}
	if !hasCreditCard {
		t.Error("WithPreset(\"pci-dss\") should detect CreditCard")
	}
}

func TestWithPresetUnknownReturnsError(t *testing.T) {
	w, err := New(WithPreset("nonexistent"))
	if err == nil {
		t.Error("WithPreset(\"nonexistent\") should return an error, got nil")
	}
	if w != nil {
		t.Error("WithPreset(\"nonexistent\") should return nil Wuming, got non-nil")
	}
	if err != nil && !strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("error should mention the unknown preset name, got: %v", err)
	}
}

func TestWithAllowlist(t *testing.T) {
	det := &stubDetector{
		matches: []model.Match{
			{Type: model.Email, Value: "john@example.com", Start: 0, End: 16, Confidence: 1.0},
			{Type: model.URL, Value: "example.com", Start: 20, End: 31, Confidence: 0.9},
		},
	}

	w, err := New(WithDetectors(det), WithAllowlist("example.com"))
	if err != nil {
		t.Fatal(err)
	}
	matches, err := w.Detect(context.Background(), "john@example.com or example.com here")
	if err != nil {
		t.Fatal(err)
	}

	if len(matches) != 1 {
		t.Fatalf("got %d matches, want 1", len(matches))
	}
	if matches[0].Type != model.Email {
		t.Errorf("match type = %v, want Email", matches[0].Type)
	}
}

func TestWithDenylist(t *testing.T) {
	det := &stubDetector{} // no matches from detector

	w, err := New(WithDetectors(det), WithDenylist(model.Name, "ACME Corp"))
	if err != nil {
		t.Fatal(err)
	}
	result, err := w.Process(context.Background(), "Contact ACME Corp for details")
	if err != nil {
		t.Fatal(err)
	}

	if result.MatchCount != 1 {
		t.Fatalf("MatchCount = %d, want 1", result.MatchCount)
	}
	if result.Matches[0].Value != "ACME Corp" {
		t.Errorf("match value = %q, want %q", result.Matches[0].Value, "ACME Corp")
	}
	want := "Contact [NAME] for details"
	if result.Redacted != want {
		t.Errorf("Redacted = %q, want %q", result.Redacted, want)
	}
}

func TestWithAllowlistAndDenylist(t *testing.T) {
	det := &stubDetector{
		matches: []model.Match{
			{Type: model.URL, Value: "example.com", Start: 0, End: 11, Confidence: 0.9},
		},
	}

	w, err := New(
		WithDetectors(det),
		WithAllowlist("example.com"),
		WithDenylist(model.Name, "ACME"),
	)
	if err != nil {
		t.Fatal(err)
	}

	result, err := w.Process(context.Background(), "example.com by ACME")
	if err != nil {
		t.Fatal(err)
	}

	if result.MatchCount != 1 {
		t.Fatalf("MatchCount = %d, want 1", result.MatchCount)
	}
	if result.Matches[0].Value != "ACME" {
		t.Errorf("match value = %q, want %q", result.Matches[0].Value, "ACME")
	}
}

func TestWithLocaleFilters(t *testing.T) {
	w, err := New(WithLocale("nl"))
	if err != nil {
		t.Fatal(err)
	}
	// BSN should be detected, SSN should not (US-specific).
	text := "BSN: 123456782, SSN: 078-05-1120"
	result, err := w.Process(context.Background(), text)
	if err != nil {
		t.Fatal(err)
	}

	hasNL := false
	hasUS := false
	for _, m := range result.Matches {
		if strings.HasPrefix(m.Detector, "nl/") {
			hasNL = true
		}
		if strings.HasPrefix(m.Detector, "us/") {
			hasUS = true
		}
	}

	if !hasNL {
		t.Error("WithLocale(\"nl\") should detect NL PII")
	}
	if hasUS {
		t.Error("WithLocale(\"nl\") should not detect US-specific PII")
	}
}
