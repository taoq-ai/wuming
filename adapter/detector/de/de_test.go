package de

import (
	"context"
	"testing"

	"github.com/taoq-ai/wuming/domain/model"
	"github.com/taoq-ai/wuming/domain/port"
)

var ctx = context.Background()

// Verify all detectors implement port.Detector.
var (
	_ port.Detector = (*IDCardDetector)(nil)
	_ port.Detector = (*SteuerIDDetector)(nil)
	_ port.Detector = (*SozialversicherungDetector)(nil)
	_ port.Detector = (*PhoneDetector)(nil)
	_ port.Detector = (*PLZDetector)(nil)
)

func TestIDCardDetector(t *testing.T) {
	d := NewIDCardDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"ID: T220001293", 1, "valid Personalausweisnummer"},
		{"L01X00T471", 1, "valid with mixed alphanumeric"},
		{"A220001293", 0, "invalid first char (A not in allowed set)"},
		{"T22000129X", 0, "check digit is not a digit"},
		{"T220001290", 0, "wrong check digit"},
		{"no id here", 0, "no ID card number"},
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
			if m.Locale != "de" {
				t.Errorf("expected locale 'de', got %q", m.Locale)
			}
			if m.Type != model.NationalID {
				t.Errorf("expected NationalID type, got %v", m.Type)
			}
		}
	}
}

func TestSteuerIDDetector(t *testing.T) {
	d := NewSteuerIDDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"Steuer-ID: 86095742719", 1, "valid Steuer-ID with correct check digit"},
		{"Tax: 65491328003", 1, "valid Steuer-ID"},
		{"11234567890", 1, "valid digit distribution and check digit"},
		{"12345678901", 0, "all digits unique (no double)"},
		{"11111111111", 0, "all same digit"},
		{"01234567890", 0, "starts with 0"},
		{"no tax id", 0, "no Steuer-ID"},
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
		}
	}
}

func TestSozialversicherungDetector(t *testing.T) {
	d := NewSozialversicherungDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"SVN: 12150175A123", 1, "valid compact format"},
		{"SVN: 12 150175 A 123", 1, "valid with spaces"},
		{"00150175A123", 0, "invalid area 00"},
		{"12130075A123", 0, "invalid month 00"},
		{"12153275A123", 0, "invalid day 32"},
		{"no svn here", 0, "no SVN"},
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
			if m.Type != model.NationalID {
				t.Errorf("expected NationalID type, got %v", m.Type)
			}
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
		{"Mobil: 0151 12345678", 1, "mobile with space"},
		{"+49 151 12345678", 1, "international mobile"},
		{"030 12345678", 1, "Berlin landline"},
		{"+49 30 12345678", 1, "international landline"},
		{"0221-1234567", 1, "Cologne landline with dash"},
		{"089/12345678", 1, "Munich landline with slash"},
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

func TestPLZDetector(t *testing.T) {
	d := NewPLZDetector()

	tests := []struct {
		input   string
		want    int
		desc    string
		minConf float64
	}{
		{"PLZ 10115", 1, "valid Berlin PLZ with context", 0.80},
		{"80331", 1, "valid Munich PLZ without context", 0.55},
		{"Straße 5, 50667 Köln", 1, "valid Cologne PLZ with address context", 0.75},
		{"00999", 0, "below valid range", 0.0},
		{"00000", 0, "all zeros", 0.0},
		{"no plz here", 0, "no PLZ", 0.0},
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
			if m.Confidence < tt.minConf {
				t.Errorf("%s: expected confidence >= %.2f, got %.2f", tt.desc, tt.minConf, m.Confidence)
			}
		}
	}
}
