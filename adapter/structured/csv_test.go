package structured

import (
	"context"
	"strings"
	"testing"

	"github.com/taoq-ai/wuming/domain/model"
)

func TestCSVScanner_WithHeaders(t *testing.T) {
	scanner := NewCSVScannerWithHeader(&stubPipeline{})
	input := "name,email,phone\nAlice,john@example.com,555-1234\nBob,safe@example.org,123-4567\n"

	result, err := scanner.Scan(context.Background(), strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}

	if result.MatchCount != 2 {
		t.Fatalf("got %d matches, want 2", result.MatchCount)
	}

	paths := map[string]model.PIIType{}
	for _, m := range result.Matches {
		paths[m.Path] = m.Type
	}
	if paths["R2:email"] != model.Email {
		t.Errorf("expected EMAIL at R2:email, got %v", paths["R2:email"])
	}
	if paths["R2:phone"] != model.Phone {
		t.Errorf("expected PHONE at R2:phone, got %v", paths["R2:phone"])
	}

	// Verify redacted CSV.
	csv := string(result.Data)
	if !strings.Contains(csv, "[EMAIL]") {
		t.Error("redacted CSV should contain [EMAIL]")
	}
	if !strings.Contains(csv, "[PHONE]") {
		t.Error("redacted CSV should contain [PHONE]")
	}
	// Header row should be unchanged.
	lines := strings.Split(strings.TrimSpace(csv), "\n")
	if lines[0] != "name,email,phone" {
		t.Errorf("header = %q, want %q", lines[0], "name,email,phone")
	}
}

func TestCSVScanner_WithoutHeaders(t *testing.T) {
	scanner := NewCSVScanner(&stubPipeline{})
	input := "Alice,john@example.com,555-1234\nBob,safe@example.org,123-4567\n"

	result, err := scanner.Scan(context.Background(), strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}

	if result.MatchCount != 2 {
		t.Fatalf("got %d matches, want 2", result.MatchCount)
	}

	paths := map[string]model.PIIType{}
	for _, m := range result.Matches {
		paths[m.Path] = m.Type
	}
	// Without headers, paths use column indices.
	if paths["R1:C2"] != model.Email {
		t.Errorf("expected EMAIL at R1:C2, got %v", paths["R1:C2"])
	}
	if paths["R1:C3"] != model.Phone {
		t.Errorf("expected PHONE at R1:C3, got %v", paths["R1:C3"])
	}
}

func TestCSVScanner_NoPII(t *testing.T) {
	scanner := NewCSVScannerWithHeader(&stubPipeline{})
	input := "name,city\nAlice,Amsterdam\nBob,Berlin\n"

	result, err := scanner.Scan(context.Background(), strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}

	if result.MatchCount != 0 {
		t.Errorf("got %d matches, want 0", result.MatchCount)
	}
}

func TestCSVScanner_EmptyInput(t *testing.T) {
	scanner := NewCSVScannerWithHeader(&stubPipeline{})
	result, err := scanner.Scan(context.Background(), strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	if result.MatchCount != 0 {
		t.Errorf("got %d matches for empty input, want 0", result.MatchCount)
	}
}

func TestCSVScanner_MultipleRows(t *testing.T) {
	scanner := NewCSVScannerWithHeader(&stubPipeline{})
	input := "name,email\nAlice,john@example.com\nBob,john@example.com\n"

	result, err := scanner.Scan(context.Background(), strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}

	if result.MatchCount != 2 {
		t.Fatalf("got %d matches, want 2", result.MatchCount)
	}

	// Both rows should have PII detected.
	paths := map[string]bool{}
	for _, m := range result.Matches {
		paths[m.Path] = true
	}
	if !paths["R2:email"] {
		t.Error("missing match at R2:email")
	}
	if !paths["R3:email"] {
		t.Error("missing match at R3:email")
	}
}

func TestCSVScanner_DetectCSV(t *testing.T) {
	scanner := NewCSVScannerWithHeader(&stubPipeline{})
	input := "name,email\nAlice,john@example.com\n"

	matches, err := scanner.DetectCSV(context.Background(), strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}

	if len(matches) != 1 {
		t.Fatalf("got %d matches, want 1", len(matches))
	}
	if matches[0].Path != "R2:email" {
		t.Errorf("path = %q, want %q", matches[0].Path, "R2:email")
	}
	if matches[0].Type != model.Email {
		t.Errorf("type = %v, want EMAIL", matches[0].Type)
	}
}

func TestCSVScanner_InvalidCSV(t *testing.T) {
	scanner := NewCSVScannerWithHeader(&stubPipeline{})
	// Uneven columns trigger a CSV parse error.
	input := "a,b\n1,2,3\n"
	_, err := scanner.Scan(context.Background(), strings.NewReader(input))
	if err == nil {
		t.Error("expected error for malformed CSV")
	}
}

func TestCSVScanner_SingleColumn(t *testing.T) {
	scanner := NewCSVScannerWithHeader(&stubPipeline{})
	input := "email\njohn@example.com\n"

	result, err := scanner.Scan(context.Background(), strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}

	if result.MatchCount != 1 {
		t.Fatalf("got %d matches, want 1", result.MatchCount)
	}
	if result.Matches[0].Path != "R2:email" {
		t.Errorf("path = %q, want %q", result.Matches[0].Path, "R2:email")
	}
}
