// Package testdata provides synthetic test corpus loading for PII detection benchmarks and integration tests.
package testdata

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

// TestCase represents a single test input with expected PII matches.
type TestCase struct {
	Input       string  `json:"input"`
	Locale      string  `json:"locale"`
	Expected    []Match `json:"expected"`
	Description string  `json:"description"`
}

// Match represents an expected PII detection result.
type Match struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// basePath returns the absolute path to the testdata directory.
func basePath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}

// LoadPositive loads positive test cases (texts containing PII) for the given locale file.
func LoadPositive(filename string) ([]TestCase, error) {
	return loadJSON(filepath.Join(basePath(), "positive", filename))
}

// LoadNegative loads negative test cases (false positive texts) from the given file.
func LoadNegative(filename string) ([]TestCase, error) {
	return loadJSON(filepath.Join(basePath(), "negative", filename))
}

func loadJSON(path string) ([]TestCase, error) {
	data, err := os.ReadFile(path) // #nosec G304 -- paths are constructed from compile-time test fixtures
	if err != nil {
		return nil, err
	}
	var cases []TestCase
	if err := json.Unmarshal(data, &cases); err != nil {
		return nil, err
	}
	return cases, nil
}
