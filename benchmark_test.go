package wuming

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/taoq-ai/wuming/adapter/registry"
)

// ---------------------------------------------------------------------------
// Test data
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// Synthetic data generators
// ---------------------------------------------------------------------------

// loremChunk is a filler paragraph with zero PII.
const loremChunk = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
	"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. " +
	"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris " +
	"nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in " +
	"reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. "

// piiSamples contains PII strings from various locales for embedding in generated text.
var piiSamples = []string{
	"jane.doe@example.com",
	"(555) 123-4567",
	"4111 1111 1111 1111",
	"192.168.1.100",
	"DE89370400440532013000",
	"00:1A:2B:3C:4D:5E",
	"078-05-1120",
	"https://example.com/page",
}

// generateText builds a string of approximately size bytes with no PII.
func generateText(size int) string {
	var b strings.Builder
	b.Grow(size)
	for b.Len() < size {
		b.WriteString(loremChunk)
	}
	return b.String()[:size]
}

// generateTextWithPII builds a string of approximately size bytes with PII
// embedded at the given density (0.0 to 1.0 — fraction of lines containing PII).
func generateTextWithPII(size int, density float64) string {
	var b strings.Builder
	b.Grow(size)
	lineNum := 0
	piiIdx := 0
	for b.Len() < size {
		if density > 0 && float64(lineNum%100) < density*100 {
			b.WriteString("Contact ")
			b.WriteString(piiSamples[piiIdx%len(piiSamples)])
			b.WriteString(" for details. ")
			piiIdx++
		}
		b.WriteString(loremChunk)
		lineNum++
	}
	return b.String()[:size]
}

// cjkFiller is Chinese text with no PII.
const cjkFiller = "这是一段用于测试的中文文本。系统需要处理多种语言的文档内容。" +
	"日本語のテストテキストです。システムは多言語のドキュメントを処理する必要があります。" +
	"이것은 테스트를 위한 한국어 텍스트입니다. 시스템은 다양한 언어의 문서를 처리해야 합니다. "

// cjkPIISamples contains Asian locale PII for embedding.
var cjkPIISamples = []string{
	"身份证号码 11010519491231002X", // Chinese Resident ID
	"My Number: 123456789018",  // Japanese My Number
	"RRN: 900101-1000006",      // Korean RRN
	"电话: 13812345678",          // Chinese phone
	"Call 090-1234-5678",       // Japanese phone
	"Call 010-1234-5678",       // Korean phone
	"信用代码: 91110000710931153R", // Chinese USCC
	"法人番号: 3234567890123",      // Japanese corporate number
}

// generateCJKText builds a CJK text of approximately size bytes with Asian PII embedded.
func generateCJKText(size int) string {
	var b strings.Builder
	b.Grow(size)
	lineNum := 0
	for b.Len() < size {
		if lineNum%5 == 0 {
			b.WriteString(cjkPIISamples[lineNum%len(cjkPIISamples)])
			b.WriteString(" ")
		}
		b.WriteString(cjkFiller)
		lineNum++
	}
	return b.String()[:size]
}

// ---------------------------------------------------------------------------
// Original benchmarks
// ---------------------------------------------------------------------------

func BenchmarkDetectSmallText(b *testing.B) {
	w, bErr := New()
	if bErr != nil {
		b.Fatal(bErr)
	}
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
	w, bErr := New()
	if bErr != nil {
		b.Fatal(bErr)
	}
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
	w, bErr := New()
	if bErr != nil {
		b.Fatal(bErr)
	}
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
	w, bErr := New()
	if bErr != nil {
		b.Fatal(bErr)
	}
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
	w, bErr := New(WithLocale("us"))
	if bErr != nil {
		b.Fatal(bErr)
	}
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
	w, bErr := New()
	if bErr != nil {
		b.Fatal(bErr)
	}
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
	w, bErr := New()
	if bErr != nil {
		b.Fatal(bErr)
	}
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := w.Redact(ctx, largeText)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ---------------------------------------------------------------------------
// Per-detector benchmarks
// ---------------------------------------------------------------------------

func BenchmarkPerDetector(b *testing.B) {
	ctx := context.Background()
	detectors := registry.AllDetectors()
	text := mediumText

	for _, d := range detectors {
		d := d // capture range variable
		b.Run(d.Name(), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := d.Detect(ctx, text)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Throughput benchmarks (MB/s)
// ---------------------------------------------------------------------------

func BenchmarkThroughputSingleLocale(b *testing.B) {
	ctx := context.Background()
	text := largeText
	textBytes := int64(len(text))

	for _, locale := range registry.Locales() {
		locale := locale
		b.Run(locale, func(b *testing.B) {
			w, bErr := New(WithLocale(locale))
			if bErr != nil {
				b.Fatal(bErr)
			}
			b.SetBytes(textBytes)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := w.Detect(ctx, text)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkThroughputAllLocales(b *testing.B) {
	ctx := context.Background()
	w, bErr := New()
	if bErr != nil {
		b.Fatal(bErr)
	}
	text := largeText
	b.SetBytes(int64(len(text)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := w.Detect(ctx, text)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ---------------------------------------------------------------------------
// Latency by text size
// ---------------------------------------------------------------------------

func BenchmarkLatencyBySize(b *testing.B) {
	sizes := []struct {
		name string
		size int
	}{
		{"100B", 100},
		{"1KB", 1024},
		{"10KB", 10 * 1024},
		{"100KB", 100 * 1024},
		{"1MB", 1024 * 1024},
	}

	ctx := context.Background()
	w, bErr := New()
	if bErr != nil {
		b.Fatal(bErr)
	}

	for _, sz := range sizes {
		sz := sz
		text := generateTextWithPII(sz.size, 0.05)
		b.Run(sz.name, func(b *testing.B) {
			b.SetBytes(int64(sz.size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := w.Detect(ctx, text)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Concurrency scaling
// ---------------------------------------------------------------------------

func BenchmarkConcurrency(b *testing.B) {
	goroutineCounts := []int{1, 2, 4, 8, 16}
	ctx := context.Background()
	text := mediumText

	for _, p := range goroutineCounts {
		p := p
		b.Run(fmt.Sprintf("goroutines-%d", p), func(b *testing.B) {
			w, bErr := New()
			if bErr != nil {
				b.Fatal(bErr)
			}
			b.SetParallelism(p)
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := w.Detect(ctx, text)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

// ---------------------------------------------------------------------------
// Cold start vs warm detection
// ---------------------------------------------------------------------------

func BenchmarkColdStart(b *testing.B) {
	ctx := context.Background()
	text := mediumText

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w, newErr := New()
		if newErr != nil {
			b.Fatal(newErr)
		}
		_, err := w.Detect(ctx, text)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWarmDetection(b *testing.B) {
	ctx := context.Background()
	w, bErr := New()
	if bErr != nil {
		b.Fatal(bErr)
	}
	text := mediumText

	// Warm up: run one detection to initialize any lazy state.
	_, _ = w.Detect(ctx, text)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := w.Detect(ctx, text)
		if err != nil {
			b.Fatal(err)
		}
	}
}
