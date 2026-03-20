package structured

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/taoq-ai/wuming/domain/model"
	"github.com/taoq-ai/wuming/domain/port"
)

// stubPipeline is a test double that detects a hardcoded set of PII patterns.
type stubPipeline struct{}

func (s *stubPipeline) Process(_ context.Context, text string) (*port.Result, error) {
	var matches []model.Match
	// Detect "john@example.com" as EMAIL.
	if idx := indexOf(text, "john@example.com"); idx >= 0 {
		matches = append(matches, model.Match{
			Type: model.Email, Value: "john@example.com",
			Start: idx, End: idx + 16, Confidence: 1.0,
		})
	}
	// Detect "555-1234" as PHONE.
	if idx := indexOf(text, "555-1234"); idx >= 0 {
		matches = append(matches, model.Match{
			Type: model.Phone, Value: "555-1234",
			Start: idx, End: idx + 8, Confidence: 0.9,
		})
	}
	// Detect "078-05-1120" as NATIONAL_ID.
	if idx := indexOf(text, "078-05-1120"); idx >= 0 {
		matches = append(matches, model.Match{
			Type: model.NationalID, Value: "078-05-1120",
			Start: idx, End: idx + 11, Confidence: 1.0,
		})
	}

	redacted := text
	if len(matches) > 0 {
		redacted = replaceAll(text, matches)
	}
	return &port.Result{
		Original:   text,
		Redacted:   redacted,
		Matches:    matches,
		MatchCount: len(matches),
	}, nil
}

// indexOf returns the byte index of substr in s, or -1.
func indexOf(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// replaceAll replaces all matches in text with their type placeholder.
func replaceAll(text string, matches []model.Match) string {
	result := []byte(text)
	// Process from end to start to keep offsets valid.
	for i := len(matches) - 1; i >= 0; i-- {
		m := matches[i]
		placeholder := "[" + m.Type.String() + "]"
		result = append(result[:m.Start], append([]byte(placeholder), result[m.End:]...)...)
	}
	return string(result)
}

func TestJSONScanner_FlatObject(t *testing.T) {
	scanner := NewJSONScanner(&stubPipeline{})
	input := `{"name":"Alice","email":"john@example.com","age":30}`

	result, err := scanner.Scan(context.Background(), []byte(input))
	if err != nil {
		t.Fatal(err)
	}

	if result.MatchCount != 1 {
		t.Fatalf("got %d matches, want 1", result.MatchCount)
	}
	if result.Matches[0].Path != "email" {
		t.Errorf("path = %q, want %q", result.Matches[0].Path, "email")
	}
	if result.Matches[0].Type != model.Email {
		t.Errorf("type = %v, want EMAIL", result.Matches[0].Type)
	}

	// Verify redacted output contains placeholder.
	var out map[string]interface{}
	if err := json.Unmarshal(result.Data, &out); err != nil {
		t.Fatal(err)
	}
	if out["email"] != "[EMAIL]" {
		t.Errorf("redacted email = %q, want %q", out["email"], "[EMAIL]")
	}
	// Non-PII fields should be unchanged.
	if out["name"] != "Alice" {
		t.Errorf("name = %q, want %q", out["name"], "Alice")
	}
}

func TestJSONScanner_NestedObject(t *testing.T) {
	scanner := NewJSONScanner(&stubPipeline{})
	input := `{"user":{"contact":{"email":"john@example.com","phone":"555-1234"}}}`

	result, err := scanner.Scan(context.Background(), []byte(input))
	if err != nil {
		t.Fatal(err)
	}

	if result.MatchCount != 2 {
		t.Fatalf("got %d matches, want 2", result.MatchCount)
	}

	paths := map[string]bool{}
	for _, m := range result.Matches {
		paths[m.Path] = true
	}
	if !paths["user.contact.email"] {
		t.Error("missing match at path user.contact.email")
	}
	if !paths["user.contact.phone"] {
		t.Error("missing match at path user.contact.phone")
	}
}

func TestJSONScanner_Array(t *testing.T) {
	scanner := NewJSONScanner(&stubPipeline{})
	input := `{"emails":["john@example.com","safe@example.org"]}`

	result, err := scanner.Scan(context.Background(), []byte(input))
	if err != nil {
		t.Fatal(err)
	}

	if result.MatchCount != 1 {
		t.Fatalf("got %d matches, want 1", result.MatchCount)
	}
	if result.Matches[0].Path != "emails[0]" {
		t.Errorf("path = %q, want %q", result.Matches[0].Path, "emails[0]")
	}
}

func TestJSONScanner_TopLevelArray(t *testing.T) {
	scanner := NewJSONScanner(&stubPipeline{})
	input := `[{"email":"john@example.com"},{"email":"safe@example.org"}]`

	result, err := scanner.Scan(context.Background(), []byte(input))
	if err != nil {
		t.Fatal(err)
	}

	if result.MatchCount != 1 {
		t.Fatalf("got %d matches, want 1", result.MatchCount)
	}
	if result.Matches[0].Path != "[0].email" {
		t.Errorf("path = %q, want %q", result.Matches[0].Path, "[0].email")
	}
}

func TestJSONScanner_NoPII(t *testing.T) {
	scanner := NewJSONScanner(&stubPipeline{})
	input := `{"name":"Alice","city":"Amsterdam"}`

	result, err := scanner.Scan(context.Background(), []byte(input))
	if err != nil {
		t.Fatal(err)
	}

	if result.MatchCount != 0 {
		t.Errorf("got %d matches, want 0", result.MatchCount)
	}
}

func TestJSONScanner_InvalidJSON(t *testing.T) {
	scanner := NewJSONScanner(&stubPipeline{})
	_, err := scanner.Scan(context.Background(), []byte(`{invalid`))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestJSONScanner_NonStringValues(t *testing.T) {
	scanner := NewJSONScanner(&stubPipeline{})
	input := `{"count":42,"active":true,"data":null}`

	result, err := scanner.Scan(context.Background(), []byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if result.MatchCount != 0 {
		t.Errorf("got %d matches for non-string values, want 0", result.MatchCount)
	}
}

func TestJSONScanner_DetectJSON(t *testing.T) {
	scanner := NewJSONScanner(&stubPipeline{})
	input := `{"email":"john@example.com","phone":"555-1234"}`

	matches, err := scanner.DetectJSON(context.Background(), []byte(input))
	if err != nil {
		t.Fatal(err)
	}

	if len(matches) != 2 {
		t.Fatalf("got %d matches, want 2", len(matches))
	}

	paths := map[string]model.PIIType{}
	for _, m := range matches {
		paths[m.Path] = m.Type
	}
	if paths["email"] != model.Email {
		t.Errorf("email path type = %v, want EMAIL", paths["email"])
	}
	if paths["phone"] != model.Phone {
		t.Errorf("phone path type = %v, want PHONE", paths["phone"])
	}
}

func TestJSONScanner_MultiplePIIInSingleField(t *testing.T) {
	scanner := NewJSONScanner(&stubPipeline{})
	input := `{"info":"Contact john@example.com or call 555-1234"}`

	result, err := scanner.Scan(context.Background(), []byte(input))
	if err != nil {
		t.Fatal(err)
	}

	if result.MatchCount != 2 {
		t.Fatalf("got %d matches, want 2", result.MatchCount)
	}
	// Both matches should have the same path.
	for _, m := range result.Matches {
		if m.Path != "info" {
			t.Errorf("path = %q, want %q", m.Path, "info")
		}
	}
}

func TestJSONScanner_EmptyObject(t *testing.T) {
	scanner := NewJSONScanner(&stubPipeline{})
	result, err := scanner.Scan(context.Background(), []byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	if result.MatchCount != 0 {
		t.Errorf("got %d matches for empty object, want 0", result.MatchCount)
	}
}

func TestJSONScanner_EmptyArray(t *testing.T) {
	scanner := NewJSONScanner(&stubPipeline{})
	result, err := scanner.Scan(context.Background(), []byte(`[]`))
	if err != nil {
		t.Fatal(err)
	}
	if result.MatchCount != 0 {
		t.Errorf("got %d matches for empty array, want 0", result.MatchCount)
	}
}
