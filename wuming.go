// Package wuming (无名 — "The Nameless") is a Go library for detecting and
// removing Personally Identifiable Information (PII) from text.
//
// It supports global PII standards across multiple locales and provides
// pluggable detection and replacement strategies via a hexagonal architecture.
//
// Quick start:
//
//	// Simple one-liner — detect and redact all PII.
//	redacted, err := wuming.Redact("Call me at 06-12345678 or email john@example.com")
//	// → "Call me at [PHONE] or email [EMAIL]"
//
//	// Configured for a specific locale.
//	w := wuming.New(wuming.WithLocale("nl"))
//	result, err := w.Process(ctx, text)
package wuming

import (
	"context"
	"io"
	"sync"

	"github.com/taoq-ai/wuming/adapter/preset"
	"github.com/taoq-ai/wuming/adapter/registry"
	"github.com/taoq-ai/wuming/adapter/replacer"
	"github.com/taoq-ai/wuming/adapter/structured"
	"github.com/taoq-ai/wuming/domain/model"
	"github.com/taoq-ai/wuming/domain/port"
	"github.com/taoq-ai/wuming/internal/engine"
)

// Wuming is the main entry point for PII detection and redaction.
type Wuming struct {
	engine *engine.Engine
}

// Option configures a Wuming instance.
type Option func(*config)

type config struct {
	detectors           []port.Detector
	replacer            port.Replacer
	locales             []string
	piiTypes            []model.PIIType
	concurrency         int
	confidenceThreshold float64
	allowlist           []string
	denylist            []engine.DenylistEntry
	err                 error
	consistentRedaction bool
}

// WithDetectors adds PII detectors.
func WithDetectors(d ...port.Detector) Option {
	return func(c *config) {
		c.detectors = append(c.detectors, d...)
	}
}

// WithReplacer sets the replacement strategy.
// Defaults to Redact (e.g. "[EMAIL]") if not specified.
func WithReplacer(r port.Replacer) Option {
	return func(c *config) {
		c.replacer = r
	}
}

// WithLocale filters detectors to those supporting the given locale.
// Global detectors always run regardless of locale.
func WithLocale(locale string) Option {
	return func(c *config) {
		c.locales = append(c.locales, locale)
	}
}

// WithPIITypes filters results to only the specified PII types.
func WithPIITypes(types ...model.PIIType) Option {
	return func(c *config) {
		c.piiTypes = append(c.piiTypes, types...)
	}
}

// WithConcurrency sets the maximum number of detectors to run in parallel.
func WithConcurrency(n int) Option {
	return func(c *config) {
		c.concurrency = n
	}
}

// WithConfidenceThreshold filters out matches below this confidence score (0.0–1.0).
func WithConfidenceThreshold(f float64) Option {
	return func(c *config) {
		c.confidenceThreshold = f
	}
}

// WithAllowlist specifies values that should never be flagged as PII.
// For example, a company domain like "example.com" can be allowlisted so
// it is not reported as a URL match. Matching is case-insensitive.
func WithAllowlist(values ...string) Option {
	return func(c *config) {
		c.allowlist = append(c.allowlist, values...)
	}
}

// WithDenylist specifies values that should always be flagged as the given PII
// type, even if no detector would normally find them. The engine injects
// synthetic matches with confidence 1.0 for every occurrence of each value.
func WithDenylist(piiType model.PIIType, values ...string) Option {
	return func(c *config) {
		for _, v := range values {
			c.denylist = append(c.denylist, engine.DenylistEntry{Value: v, PIIType: piiType})
		}
	}
}

// WithConsistentRedaction ensures that the same PII value always maps to
// the same replacement string. For example, if "john@example.com" appears
// three times the Redact replacer will produce "[EMAIL_1]" in every position.
// This works with any replacer strategy (Redact, Hash, Mask, Custom).
func WithConsistentRedaction() Option {
	return func(c *config) {
		c.consistentRedaction = true
	}
}

// WithPreset configures the instance for a specific compliance regulation
// (e.g. "gdpr", "hipaa", "pci-dss"). It sets the appropriate locales, PII
// types, and detectors based on the preset definition. If the preset name
// is unknown, New will return a non-nil error.
func WithPreset(name string) Option {
	return func(c *config) {
		p, err := preset.Get(name)
		if err != nil {
			c.err = err
			return
		}

		// Collect detectors for each locale in the preset.
		seen := make(map[string]bool)
		for _, locale := range p.Locales {
			if seen[locale] {
				continue
			}
			seen[locale] = true
			c.detectors = append(c.detectors, registry.DetectorsForLocale(locale)...)
			c.locales = append(c.locales, locale)
		}

		c.piiTypes = append(c.piiTypes, p.PIITypes...)
	}
}

// New creates a new Wuming instance with the given options.
// It returns an error if any option fails (e.g. an unknown preset name).
func New(opts ...Option) (*Wuming, error) {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.err != nil {
		return nil, cfg.err
	}

	// Default to all detectors if none explicitly provided.
	if len(cfg.detectors) == 0 {
		cfg.detectors = registry.AllDetectors()
	}

	// Default replacer.
	if cfg.replacer == nil {
		cfg.replacer = replacer.NewRedact()
	}

	// Wrap with consistent-replacement behavior when requested.
	if cfg.consistentRedaction {
		cfg.replacer = replacer.NewConsistent(cfg.replacer)
	}

	var engineOpts []engine.Option
	if len(cfg.detectors) > 0 {
		engineOpts = append(engineOpts, engine.WithDetectors(cfg.detectors...))
	}
	engineOpts = append(engineOpts, engine.WithReplacer(cfg.replacer))
	if len(cfg.locales) > 0 {
		engineOpts = append(engineOpts, engine.WithLocales(cfg.locales...))
	}
	if len(cfg.piiTypes) > 0 {
		engineOpts = append(engineOpts, engine.WithPIITypes(cfg.piiTypes...))
	}
	if cfg.concurrency > 0 {
		engineOpts = append(engineOpts, engine.WithConcurrency(cfg.concurrency))
	}
	if cfg.confidenceThreshold > 0 {
		engineOpts = append(engineOpts, engine.WithConfidenceThreshold(cfg.confidenceThreshold))
	}
	if len(cfg.allowlist) > 0 {
		engineOpts = append(engineOpts, engine.WithAllowlist(cfg.allowlist...))
	}
	if len(cfg.denylist) > 0 {
		for _, entry := range cfg.denylist {
			engineOpts = append(engineOpts, engine.WithDenylist(entry.PIIType, entry.Value))
		}
	}

	return &Wuming{engine: engine.New(engineOpts...)}, nil
}

// Process runs PII detection and replacement, returning the full result.
func (w *Wuming) Process(ctx context.Context, text string) (*port.Result, error) {
	return w.engine.Process(ctx, text)
}

// Detect runs PII detection only, returning all matches without modifying text.
func (w *Wuming) Detect(ctx context.Context, text string) ([]model.Match, error) {
	result, err := w.engine.Process(ctx, text)
	if err != nil {
		return nil, err
	}
	return result.Matches, nil
}

// Redact runs PII detection and returns the redacted text.
func (w *Wuming) Redact(ctx context.Context, text string) (string, error) {
	result, err := w.engine.Process(ctx, text)
	if err != nil {
		return "", err
	}
	return result.Redacted, nil
}

// --- Structured data support (JSON, CSV) ---

// RedactJSON detects and redacts PII in a JSON document, scanning each string
// value individually. It returns a structured result containing the redacted
// JSON bytes and matches annotated with their JSON path (e.g. "user.email").
func (w *Wuming) RedactJSON(ctx context.Context, data []byte) (*structured.Result, error) {
	return structured.NewJSONScanner(w.engine).Scan(ctx, data)
}

// DetectJSON detects PII in a JSON document without modifying it. Each match
// is annotated with the JSON path where it was found.
func (w *Wuming) DetectJSON(ctx context.Context, data []byte) ([]structured.FieldMatch, error) {
	return structured.NewJSONScanner(w.engine).DetectJSON(ctx, data)
}

// RedactCSV detects and redacts PII in CSV data read from r. The first row is
// treated as column headers and used in match paths (e.g. "R2:email").
// Returns a structured result containing the redacted CSV bytes and matches.
func (w *Wuming) RedactCSV(ctx context.Context, r io.Reader) (*structured.Result, error) {
	return structured.NewCSVScannerWithHeader(w.engine).Scan(ctx, r)
}

// DetectCSV detects PII in CSV data without modifying it. The first row is
// treated as column headers. Each match is annotated with its row/column path.
func (w *Wuming) DetectCSV(ctx context.Context, r io.Reader) ([]structured.FieldMatch, error) {
	return structured.NewCSVScannerWithHeader(w.engine).DetectCSV(ctx, r)
}

// defaultInstance is lazily initialized with all detectors.
var (
	defaultOnce     sync.Once
	defaultInstance *Wuming
)

func getDefault() *Wuming {
	defaultOnce.Do(func() {
		w, err := New()
		if err != nil {
			panic("wuming: failed to create default instance: " + err.Error())
		}
		defaultInstance = w
	})
	return defaultInstance
}

// Redact detects and redacts all PII from text using all available detectors.
func Redact(ctx context.Context, text string) (string, error) {
	return getDefault().Redact(ctx, text)
}

// Detect finds all PII matches in text using all available detectors.
func Detect(ctx context.Context, text string) ([]model.Match, error) {
	return getDefault().Detect(ctx, text)
}

// Process runs full PII detection and replacement using all available detectors.
func Process(ctx context.Context, text string) (*port.Result, error) {
	return getDefault().Process(ctx, text)
}

// RedactJSON detects and redacts PII in a JSON document using all available detectors.
func RedactJSON(ctx context.Context, data []byte) (*structured.Result, error) {
	return getDefault().RedactJSON(ctx, data)
}

// DetectJSON detects PII in a JSON document using all available detectors.
func DetectJSON(ctx context.Context, data []byte) ([]structured.FieldMatch, error) {
	return getDefault().DetectJSON(ctx, data)
}

// RedactCSV detects and redacts PII in CSV data using all available detectors.
func RedactCSV(ctx context.Context, r io.Reader) (*structured.Result, error) {
	return getDefault().RedactCSV(ctx, r)
}

// DetectCSV detects PII in CSV data using all available detectors.
func DetectCSV(ctx context.Context, r io.Reader) ([]structured.FieldMatch, error) {
	return getDefault().DetectCSV(ctx, r)
}
