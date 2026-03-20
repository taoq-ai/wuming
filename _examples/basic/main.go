package main

import (
	"context"
	"fmt"
	"log"

	"github.com/taoq-ai/wuming"
)

func main() {
	ctx := context.Background()

	// Zero-config: one line catches all PII across every locale.
	redacted, err := wuming.Redact(ctx, "Email john@acme.com, SSN 123-45-6789, BSN 111222333, call 06-12345678")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("=== Zero-config Redact ===")
	fmt.Println("Redacted:", redacted)
	fmt.Println()

	// Configured instance: restrict to Dutch locale.
	w, err := wuming.New(
		wuming.WithLocale("nl"),
	)
	if err != nil {
		log.Fatal(err)
	}

	text := "Email john@acme.com, SSN 123-45-6789, BSN 111222333, call 06-12345678"
	result, err := w.Process(ctx, text)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Locale-specific (nl) ===")
	fmt.Println("Original:", result.Original)
	fmt.Println("Redacted:", result.Redacted)
	fmt.Printf("Found %d PII matches\n", result.MatchCount)
	for _, m := range result.Matches {
		fmt.Printf("  [%s] %q (confidence: %.0f%%)\n", m.Type, m.Value, m.Confidence*100)
	}
}
