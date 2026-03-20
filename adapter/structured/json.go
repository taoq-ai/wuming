package structured

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/taoq-ai/wuming/domain/port"
)

// JSONScanner scans JSON documents field-by-field for PII and can produce
// a redacted copy with match annotations that include the JSON path.
type JSONScanner struct {
	scanner
}

// NewJSONScanner creates a JSONScanner that delegates text detection to
// the given Pipeline (typically an engine.Engine or a wuming.Wuming
// instance accessed through its Process method).
func NewJSONScanner(p port.Pipeline) *JSONScanner {
	return &JSONScanner{scanner: scanner{pipeline: p}}
}

// Scan parses the JSON data, walks every string value, detects PII in each,
// and returns a Result with the redacted JSON and all field matches.
func (j *JSONScanner) Scan(ctx context.Context, data []byte) (*Result, error) {
	var root interface{}
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("structured/json: invalid JSON: %w", err)
	}

	var allMatches []FieldMatch
	redacted, matches, err := j.walkValue(ctx, root, "")
	if err != nil {
		return nil, err
	}
	allMatches = append(allMatches, matches...)

	out, err := json.Marshal(redacted)
	if err != nil {
		return nil, fmt.Errorf("structured/json: marshal error: %w", err)
	}

	return &Result{
		Data:       out,
		Matches:    allMatches,
		MatchCount: len(allMatches),
	}, nil
}

// DetectJSON parses the JSON data, walks every string value, and returns
// all PII field matches without modifying the data.
func (j *JSONScanner) DetectJSON(ctx context.Context, data []byte) ([]FieldMatch, error) {
	var root interface{}
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("structured/json: invalid JSON: %w", err)
	}

	_, matches, err := j.walkValue(ctx, root, "")
	if err != nil {
		return nil, err
	}
	return matches, nil
}

// walkValue recursively processes a JSON value and returns its redacted form
// plus any PII matches found.
func (j *JSONScanner) walkValue(ctx context.Context, v interface{}, path string) (interface{}, []FieldMatch, error) {
	switch val := v.(type) {
	case map[string]interface{}:
		return j.walkObject(ctx, val, path)
	case []interface{}:
		return j.walkArray(ctx, val, path)
	case string:
		redacted, matches, err := j.detectAndRedact(ctx, val, path)
		if err != nil {
			return nil, nil, err
		}
		if len(matches) > 0 {
			return redacted, matches, nil
		}
		return val, nil, nil
	default:
		// Numbers, booleans, null — no PII to detect.
		return val, nil, nil
	}
}

// walkObject processes a JSON object. Keys are iterated in sorted order for
// deterministic output.
func (j *JSONScanner) walkObject(ctx context.Context, obj map[string]interface{}, path string) (interface{}, []FieldMatch, error) {
	var allMatches []FieldMatch
	result := make(map[string]interface{}, len(obj))

	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		childPath := k
		if path != "" {
			childPath = path + "." + k
		}
		redacted, matches, err := j.walkValue(ctx, obj[k], childPath)
		if err != nil {
			return nil, nil, err
		}
		result[k] = redacted
		allMatches = append(allMatches, matches...)
	}
	return result, allMatches, nil
}

// walkArray processes a JSON array.
func (j *JSONScanner) walkArray(ctx context.Context, arr []interface{}, path string) (interface{}, []FieldMatch, error) {
	var allMatches []FieldMatch
	result := make([]interface{}, len(arr))
	for i, elem := range arr {
		childPath := path + "[" + strconv.Itoa(i) + "]"
		redacted, matches, err := j.walkValue(ctx, elem, childPath)
		if err != nil {
			return nil, nil, err
		}
		result[i] = redacted
		allMatches = append(allMatches, matches...)
	}
	return result, allMatches, nil
}
