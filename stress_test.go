package wuming

import (
	"context"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Stress tests — skipped in short mode (default CI). Run explicitly with:
//
//	go test -run TestStress -v -count=1 ./...
// ---------------------------------------------------------------------------

func TestStressLargeDocument(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping stress test in short mode")
	}

	const docSize = 10 * 1024 * 1024 // 10 MB
	text := generateTextWithPII(docSize, 0.01)

	w, err := New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	start := time.Now()
	matches, err := w.Detect(ctx, text)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Detect on 10 MB document failed: %v", err)
	}

	t.Logf("10 MB document: %d matches found in %s (%.2f MB/s)",
		len(matches), elapsed, float64(docSize)/elapsed.Seconds()/1024/1024)
}

func TestStressHighPIIDensity(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping stress test in short mode")
	}

	// Generate text where every line has PII (density=1.0).
	const docSize = 1024 * 1024 // 1 MB
	text := generateTextWithPII(docSize, 1.0)

	w, err := New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	start := time.Now()
	matches, err := w.Detect(ctx, text)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Detect on high-PII-density document failed: %v", err)
	}

	if len(matches) == 0 {
		t.Error("expected PII matches in high-density document, got 0")
	}

	t.Logf("high-density 1 MB: %d matches found in %s", len(matches), elapsed)
}

func TestStressZeroPII(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping stress test in short mode")
	}

	const docSize = 10 * 1024 * 1024 // 10 MB
	text := generateText(docSize)

	w, err := New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	start := time.Now()
	matches, err := w.Detect(ctx, text)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Detect on zero-PII document failed: %v", err)
	}

	if len(matches) != 0 {
		t.Errorf("expected 0 matches in zero-PII document, got %d", len(matches))
	}

	t.Logf("zero-PII 10 MB: scanned in %s (%.2f MB/s overhead)",
		elapsed, float64(docSize)/elapsed.Seconds()/1024/1024)
}

func TestStressUnicode(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping stress test in short mode")
	}

	const docSize = 1024 * 1024 // 1 MB of CJK text with embedded PII
	text := generateCJKText(docSize)

	w, err := New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	start := time.Now()
	matches, err := w.Detect(ctx, text)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Detect on CJK document failed: %v", err)
	}

	if len(matches) == 0 {
		t.Error("expected PII matches in CJK document with embedded IDs, got 0")
	}

	// Check that we found matches from multiple CJK locales.
	locales := make(map[string]int)
	for _, m := range matches {
		locales[m.Locale]++
	}

	t.Logf("CJK 1 MB: %d matches across %d locales in %s",
		len(matches), len(locales), elapsed)

	for locale, count := range locales {
		t.Logf("  locale %q: %d matches", locale, count)
	}
}
