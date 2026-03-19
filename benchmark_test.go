package wuming

import (
	"context"
	"strings"
	"testing"
)

// smallText is ~100 characters with a couple of PII items.
var smallText = "Contact jane.doe@example.com or call (555) 123-4567 for details about the quarterly report update."

// mediumText is ~1000 characters with PII scattered throughout.
var mediumText = buildMediumText()

// largeText is ~10000 characters with PII scattered throughout.
var largeText = buildLargeText()

func buildMediumText() string {
	var b strings.Builder
	paragraphs := []string{
		"Dear customer, your account has been updated. Please contact support at help@example.com if you have questions. ",
		"Your reference number is 42 and the case was opened on Monday. The office is located at 123 Main Street. ",
		"For billing inquiries, reach us at (555) 234-5678 during business hours from 9 AM to 5 PM Eastern Time. ",
		"Payment was processed via card ending in 4111 1111 1111 1111 on the third of the month successfully. ",
		"Server logs indicate connections from 192.168.1.100 and 10.0.0.1 during the maintenance window last night. ",
		"Please visit https://example.com/account/settings to update your preferences and notification settings today. ",
		"The IBAN for wire transfers is DE89370400440532013000, please include your invoice number as reference. ",
		"Network device MAC 00:1A:2B:3C:4D:5E was registered to the access point in building C floor three. ",
		"Your SSN 078-05-1120 is on file for identity verification purposes and will be kept confidential. ",
		"Ship orders to ZIP 90210-1234, attention warehouse B, dock seven, for overnight delivery options. ",
	}
	for _, p := range paragraphs {
		b.WriteString(p)
	}
	return b.String()
}

func buildLargeText() string {
	var b strings.Builder
	// Repeat medium text ~10 times with slight variation to reach ~10000 chars.
	base := buildMediumText()
	for i := 0; i < 10; i++ {
		b.WriteString(base)
		b.WriteString("\n")
	}
	return b.String()
}

func BenchmarkDetectSmallText(b *testing.B) {
	w := New()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := w.Detect(ctx, smallText)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDetectMediumText(b *testing.B) {
	w := New()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := w.Detect(ctx, mediumText)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDetectLargeText(b *testing.B) {
	w := New()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := w.Detect(ctx, largeText)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDetectAllLocales(b *testing.B) {
	// Zero-config uses all registered detectors across all locales.
	w := New()
	ctx := context.Background()
	text := mediumText
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := w.Detect(ctx, text)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDetectSingleLocale(b *testing.B) {
	// Restricted to a single locale (US) plus common detectors.
	w := New(WithLocale("us"))
	ctx := context.Background()
	text := mediumText
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := w.Detect(ctx, text)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRedactSmallText(b *testing.B) {
	w := New()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := w.Redact(ctx, smallText)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRedactLargeText(b *testing.B) {
	w := New()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := w.Redact(ctx, largeText)
		if err != nil {
			b.Fatal(err)
		}
	}
}
