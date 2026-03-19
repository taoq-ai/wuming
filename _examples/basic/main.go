package main

import (
	"context"
	"fmt"
	"log"

	"github.com/taoq-ai/wuming"
	"github.com/taoq-ai/wuming/adapter/detector/common"
	"github.com/taoq-ai/wuming/adapter/detector/nl"
	"github.com/taoq-ai/wuming/adapter/detector/us"
)

func main() {
	// Create a wuming instance with detectors from multiple locales.
	w := wuming.New(
		wuming.WithDetectors(
			common.NewEmailDetector(),
			common.NewCreditCardDetector(),
			us.NewSSNDetector(),
			nl.NewBSNDetector(),
			nl.NewPhoneDetector(),
		),
	)

	text := "Email john@acme.com, SSN 123-45-6789, BSN 111222333, call 06-12345678"
	result, err := w.Process(context.Background(), text)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Original:", result.Original)
	fmt.Println("Redacted:", result.Redacted)
	fmt.Printf("Found %d PII matches\n", result.MatchCount)
	for _, m := range result.Matches {
		fmt.Printf("  [%s] %q (confidence: %.0f%%)\n", m.Type, m.Value, m.Confidence*100)
	}
}
