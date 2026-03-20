package structured

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/taoq-ai/wuming/domain/port"
)

// CSVScanner scans CSV data cell-by-cell for PII and can produce a redacted
// copy with match annotations that include row and column context.
type CSVScanner struct {
	scanner
	// UseHeader when true treats the first row as column headers and uses
	// them in match paths (e.g. "R2:email") instead of column indices.
	UseHeader bool
}

// NewCSVScanner creates a CSVScanner that delegates text detection to the
// given Pipeline.
func NewCSVScanner(p port.Pipeline) *CSVScanner {
	return &CSVScanner{scanner: scanner{pipeline: p}}
}

// NewCSVScannerWithHeader creates a CSVScanner that uses the first row as
// column headers in match paths.
func NewCSVScannerWithHeader(p port.Pipeline) *CSVScanner {
	return &CSVScanner{
		scanner:   scanner{pipeline: p},
		UseHeader: true,
	}
}

// Scan reads all CSV records from r, detects PII in each cell, and returns
// a Result with the redacted CSV and all field matches.
func (c *CSVScanner) Scan(ctx context.Context, r io.Reader) (*Result, error) {
	reader := csv.NewReader(r)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("structured/csv: read error: %w", err)
	}

	if len(records) == 0 {
		return &Result{Data: nil}, nil
	}

	var headers []string
	startRow := 0
	if c.UseHeader && len(records) > 0 {
		headers = records[0]
		startRow = 1
	}

	var allMatches []FieldMatch
	redacted := make([][]string, len(records))

	// Copy header row as-is if using headers.
	if c.UseHeader && len(records) > 0 {
		redacted[0] = make([]string, len(records[0]))
		copy(redacted[0], records[0])
	}

	for i := startRow; i < len(records); i++ {
		row := records[i]
		redactedRow := make([]string, len(row))
		for j, cell := range row {
			path := c.cellPath(i+1, j, headers) // 1-based row numbering
			redactedCell, matches, err := c.detectAndRedact(ctx, cell, path)
			if err != nil {
				return nil, err
			}
			if len(matches) > 0 {
				redactedRow[j] = redactedCell
				allMatches = append(allMatches, matches...)
			} else {
				redactedRow[j] = cell
			}
		}
		redacted[i] = redactedRow
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	if err := writer.WriteAll(redacted); err != nil {
		return nil, fmt.Errorf("structured/csv: write error: %w", err)
	}

	return &Result{
		Data:       buf.Bytes(),
		Matches:    allMatches,
		MatchCount: len(allMatches),
	}, nil
}

// DetectCSV reads all CSV records from r, detects PII in each cell, and
// returns all field matches without modifying the data.
func (c *CSVScanner) DetectCSV(ctx context.Context, r io.Reader) ([]FieldMatch, error) {
	reader := csv.NewReader(r)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("structured/csv: read error: %w", err)
	}

	if len(records) == 0 {
		return nil, nil
	}

	var headers []string
	startRow := 0
	if c.UseHeader && len(records) > 0 {
		headers = records[0]
		startRow = 1
	}

	var allMatches []FieldMatch
	for i := startRow; i < len(records); i++ {
		for j, cell := range records[i] {
			path := c.cellPath(i+1, j, headers)
			_, matches, err := c.detectAndRedact(ctx, cell, path)
			if err != nil {
				return nil, err
			}
			allMatches = append(allMatches, matches...)
		}
	}
	return allMatches, nil
}

// cellPath builds a human-readable path for a CSV cell.
// Format: "R<row>:<header>" when headers are available, or "R<row>:C<col>" otherwise.
func (c *CSVScanner) cellPath(row, col int, headers []string) string {
	prefix := "R" + strconv.Itoa(row)
	if headers != nil && col < len(headers) {
		return prefix + ":" + headers[col]
	}
	return prefix + ":C" + strconv.Itoa(col+1) // 1-based column numbering
}
