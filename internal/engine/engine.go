// Package engine provides the core orchestrator that wires detectors and
// replacers together, implementing the Pipeline port interface.
package engine

import (
	"context"
	"sort"
	"sync"

	"github.com/taoq-ai/wuming/domain/model"
	"github.com/taoq-ai/wuming/domain/port"
)

// Engine orchestrates PII detection and replacement.
type Engine struct {
	detectors           []port.Detector
	replacer            port.Replacer
	locales             map[string]bool
	piiTypes            map[model.PIIType]bool
	concurrency         int
	confidenceThreshold float64
}

// Option configures an Engine.
type Option func(*Engine)

// WithDetectors adds detectors to the engine.
func WithDetectors(d ...port.Detector) Option {
	return func(e *Engine) {
		e.detectors = append(e.detectors, d...)
	}
}

// WithReplacer sets the replacement strategy.
func WithReplacer(r port.Replacer) Option {
	return func(e *Engine) {
		e.replacer = r
	}
}

// WithLocales filters detectors to only those supporting the given locales.
func WithLocales(locales ...string) Option {
	return func(e *Engine) {
		for _, l := range locales {
			e.locales[l] = true
		}
	}
}

// WithPIITypes filters results to only the specified PII types.
func WithPIITypes(types ...model.PIIType) Option {
	return func(e *Engine) {
		for _, t := range types {
			e.piiTypes[t] = true
		}
	}
}

// WithConcurrency sets the maximum number of detectors to run in parallel.
// Defaults to 0 (unlimited).
func WithConcurrency(n int) Option {
	return func(e *Engine) {
		e.concurrency = n
	}
}

// WithConfidenceThreshold filters out matches below this confidence score.
func WithConfidenceThreshold(f float64) Option {
	return func(e *Engine) {
		e.confidenceThreshold = f
	}
}

// New creates a new Engine with the given options.
func New(opts ...Option) *Engine {
	e := &Engine{
		locales:  make(map[string]bool),
		piiTypes: make(map[model.PIIType]bool),
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Process runs all configured detectors, merges matches, and applies the replacer.
func (e *Engine) Process(ctx context.Context, text string) (*port.Result, error) {
	detectors := e.activeDetectors()
	allMatches, err := e.runDetectors(ctx, text, detectors)
	if err != nil {
		return nil, err
	}

	allMatches = e.filterByConfidence(allMatches)
	allMatches = e.filterByPIIType(allMatches)
	allMatches = dedup(allMatches)

	redacted := text
	if e.replacer != nil && len(allMatches) > 0 {
		redacted, err = e.replacer.Replace(text, allMatches)
		if err != nil {
			return nil, err
		}
	}

	return &port.Result{
		Original:   text,
		Redacted:   redacted,
		Matches:    allMatches,
		MatchCount: len(allMatches),
	}, nil
}

// activeDetectors returns detectors matching the configured locale filters.
func (e *Engine) activeDetectors() []port.Detector {
	if len(e.locales) == 0 {
		return e.detectors
	}

	var filtered []port.Detector
	for _, d := range e.detectors {
		locales := d.Locales()
		if len(locales) == 0 {
			// Global detectors always run.
			filtered = append(filtered, d)
			continue
		}
		for _, l := range locales {
			if e.locales[l] {
				filtered = append(filtered, d)
				break
			}
		}
	}
	return filtered
}

func (e *Engine) runDetectors(ctx context.Context, text string, detectors []port.Detector) ([]model.Match, error) {
	if len(detectors) == 0 {
		return nil, nil
	}

	type result struct {
		matches []model.Match
		err     error
	}

	results := make([]result, len(detectors))
	var wg sync.WaitGroup

	sem := make(chan struct{}, e.effectiveConcurrency(len(detectors)))

	for i, d := range detectors {
		wg.Add(1)
		go func(idx int, det port.Detector) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			m, err := det.Detect(ctx, text)
			results[idx] = result{matches: m, err: err}
		}(i, d)
	}
	wg.Wait()

	var all []model.Match
	for _, r := range results {
		if r.err != nil {
			return nil, r.err
		}
		all = append(all, r.matches...)
	}
	return all, nil
}

func (e *Engine) effectiveConcurrency(n int) int {
	if e.concurrency > 0 && e.concurrency < n {
		return e.concurrency
	}
	return n
}

func (e *Engine) filterByConfidence(matches []model.Match) []model.Match {
	if e.confidenceThreshold <= 0 {
		return matches
	}
	var filtered []model.Match
	for _, m := range matches {
		if m.Confidence >= e.confidenceThreshold {
			filtered = append(filtered, m)
		}
	}
	return filtered
}

func (e *Engine) filterByPIIType(matches []model.Match) []model.Match {
	if len(e.piiTypes) == 0 {
		return matches
	}
	var filtered []model.Match
	for _, m := range matches {
		if e.piiTypes[m.Type] {
			filtered = append(filtered, m)
		}
	}
	return filtered
}

// dedup removes overlapping matches, preferring higher confidence.
func dedup(matches []model.Match) []model.Match {
	if len(matches) <= 1 {
		return matches
	}

	// Sort by start position, then by confidence descending for ties.
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].Start != matches[j].Start {
			return matches[i].Start < matches[j].Start
		}
		return matches[i].Confidence > matches[j].Confidence
	})

	var result []model.Match
	result = append(result, matches[0])

	for i := 1; i < len(matches); i++ {
		last := result[len(result)-1]
		cur := matches[i]
		// If current match overlaps with the last kept match, skip it.
		if cur.Start < last.End {
			// Keep the one with higher confidence (already kept).
			continue
		}
		result = append(result, cur)
	}
	return result
}
