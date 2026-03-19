package nl

import (
	"context"
	"testing"

	"github.com/taoq-ai/wuming/domain/model"
	"github.com/taoq-ai/wuming/domain/port"
)

var ctx = context.Background()

// Verify all detectors implement port.Detector.
var (
	_ port.Detector = (*BSNDetector)(nil)
	_ port.Detector = (*PhoneDetector)(nil)
	_ port.Detector = (*PostalDetector)(nil)
	_ port.Detector = (*KvKDetector)(nil)
	_ port.Detector = (*IDDocumentDetector)(nil)
)

func TestBSNDetector(t *testing.T) {
	d := NewBSNDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"BSN: 111222333", 1, "valid BSN plain"},
		{"BSN: 111.222.333", 1, "valid BSN with dots"},
		{"BSN: 111 222 333", 1, "valid BSN with spaces"},
		{"123456782", 1, "valid 11-proof number"},
		{"000000000", 0, "all zeros fails 11-proof (sum is 0)"},
		{"123456789", 0, "invalid 11-proof"},
		{"12345678", 0, "too few digits"},
		{"no bsn here", 0, "no BSN"},
	}

	for _, tt := range tests {
		matches, err := d.Detect(ctx, tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if len(matches) != tt.want {
			t.Errorf("%s: Detect(%q) got %d matches, want %d", tt.desc, tt.input, len(matches), tt.want)
		}
		for _, m := range matches {
			if m.Locale != "nl" {
				t.Errorf("expected locale 'nl', got %q", m.Locale)
			}
			if m.Type != model.NationalID {
				t.Errorf("expected NationalID type, got %v", m.Type)
			}
		}
	}
}

func TestBSN11Proof(t *testing.T) {
	tests := []struct {
		digits string
		valid  bool
		desc   string
	}{
		{"111222333", true, "known valid BSN"},
		{"123456782", true, "another valid BSN"},
		{"000000000", false, "all zeros (sum is 0)"},
		{"123456789", false, "invalid checksum"},
	}

	for _, tt := range tests {
		got := isValid11Proof(tt.digits)
		if got != tt.valid {
			t.Errorf("%s: isValid11Proof(%q) = %v, want %v", tt.desc, tt.digits, got, tt.valid)
		}
	}
}

func TestPhoneDetector(t *testing.T) {
	d := NewPhoneDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"Bel 0612345678", 1, "mobile no separator"},
		{"Bel 06-12345678", 1, "mobile with dash"},
		{"Bel 06 12345678", 1, "mobile with space"},
		{"+31 6 12345678", 1, "international mobile +31"},
		{"0031 6 12345678", 1, "international mobile 0031"},
		{"+31612345678", 1, "international mobile compact"},
		{"020-1234567", 1, "landline Amsterdam"},
		{"+31 20 1234567", 1, "international landline"},
		{"0031 20 1234567", 1, "international landline 0031"},
		{"12345", 0, "too short"},
		{"no phone here", 0, "no phone"},
	}

	for _, tt := range tests {
		matches, err := d.Detect(ctx, tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if len(matches) != tt.want {
			t.Errorf("%s: Detect(%q) got %d matches, want %d", tt.desc, tt.input, len(matches), tt.want)
		}
		for _, m := range matches {
			if m.Type != model.Phone {
				t.Errorf("expected Phone type, got %v", m.Type)
			}
		}
	}
}

func TestPostalDetector(t *testing.T) {
	d := NewPostalDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"1234 AB", 1, "standard format with space"},
		{"1234AB", 1, "no space"},
		{"9999 ZZ", 1, "high number valid"},
		{"1000 SA", 0, "SA excluded"},
		{"1000 SD", 0, "SD excluded"},
		{"1000 SS", 0, "SS excluded"},
		{"0123 AB", 0, "starts with 0 invalid"},
		{"1234 ab", 0, "lowercase invalid"},
		{"no postal", 0, "no postal code"},
	}

	for _, tt := range tests {
		matches, err := d.Detect(ctx, tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if len(matches) != tt.want {
			t.Errorf("%s: Detect(%q) got %d matches, want %d", tt.desc, tt.input, len(matches), tt.want)
		}
		for _, m := range matches {
			if m.Type != model.PostalCode {
				t.Errorf("expected PostalCode type, got %v", m.Type)
			}
			if m.Confidence != 0.90 {
				t.Errorf("expected confidence 0.90, got %f", m.Confidence)
			}
		}
	}
}

func TestKvKDetector(t *testing.T) {
	d := NewKvKDetector()

	tests := []struct {
		input      string
		want       int
		desc       string
		confidence float64
	}{
		{"12345678", 1, "bare 8-digit number", 0.60},
		{"KvK 12345678", 1, "KvK prefix", 0.90},
		{"kvk: 12345678", 1, "kvk lowercase with colon", 0.90},
		{"Kamer van Koophandel 12345678", 1, "full name prefix", 0.90},
		{"1234567", 0, "too few digits", 0},
		{"123456789", 0, "too many digits", 0},
		{"no kvk here", 0, "no KvK", 0},
	}

	for _, tt := range tests {
		matches, err := d.Detect(ctx, tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if len(matches) != tt.want {
			t.Errorf("%s: Detect(%q) got %d matches, want %d", tt.desc, tt.input, len(matches), tt.want)
		}
		for _, m := range matches {
			if m.Type != model.TaxID {
				t.Errorf("expected TaxID type, got %v", m.Type)
			}
			if tt.confidence > 0 && m.Confidence != tt.confidence {
				t.Errorf("%s: expected confidence %f, got %f", tt.desc, tt.confidence, m.Confidence)
			}
		}
	}
}

func TestIDDocumentDetector(t *testing.T) {
	d := NewIDDocumentDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"ID: SPECI2014", 1, "Dutch ID card format"},
		{"ID: SPEC12014", 1, "Dutch ID alternate"},
		{"PP: NR1234567", 1, "Dutch passport format"},
		{"PP: AB1234567", 1, "Dutch passport alternate"},
		{"ABCDEFGHI", 0, "all letters no digits"},
		{"no id here", 0, "no ID document"},
	}

	for _, tt := range tests {
		matches, err := d.Detect(ctx, tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if len(matches) != tt.want {
			t.Errorf("%s: Detect(%q) got %d matches, want %d", tt.desc, tt.input, len(matches), tt.want)
		}
		for _, m := range matches {
			if m.Locale != "nl" {
				t.Errorf("expected locale 'nl', got %q", m.Locale)
			}
		}
	}
}
