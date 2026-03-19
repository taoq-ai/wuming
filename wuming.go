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

	"github.com/taoq-ai/wuming/adapter/replacer"
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

// New creates a new Wuming instance with the given options.
func New(opts ...Option) *Wuming {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}

	// Default replacer.
	if cfg.replacer == nil {
		cfg.replacer = replacer.NewRedact()
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

	return &Wuming{engine: engine.New(engineOpts...)}
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
