package us

import (
	"context"
	"testing"

	"github.com/taoq-ai/wuming/domain/model"
	"github.com/taoq-ai/wuming/domain/port"
)

var ctx = context.Background()

// Verify all detectors implement port.Detector.
var (
	_ port.Detector = (*SSNDetector)(nil)
	_ port.Detector = (*EINDetector)(nil)
	_ port.Detector = (*PhoneDetector)(nil)
	_ port.Detector = (*ZIPDetector)(nil)
	_ port.Detector = (*PassportDetector)(nil)
	_ port.Detector = (*ITINDetector)(nil)
	_ port.Detector = (*MedicareDetector)(nil)
)

func TestSSNDetector(t *testing.T) {
	d := NewSSNDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"SSN: 123-45-6789", 1, "standard format"},
		{"123456789", 1, "no dashes"},
		{"000-12-3456", 0, "area 000 invalid"},
		{"666-12-3456", 0, "area 666 invalid"},
		{"900-12-3456", 0, "area 900+ invalid"},
		{"123-00-6789", 0, "group 00 invalid"},
		{"123-45-0000", 0, "serial 0000 invalid"},
		{"no ssn here", 0, "no SSN"},
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
			if m.Locale != "us" {
				t.Errorf("expected locale 'us', got %q", m.Locale)
			}
		}
	}
}

func TestEINDetector(t *testing.T) {
	d := NewEINDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"EIN: 12-3456789", 1, "valid EIN"},
		{"01-1234567", 1, "valid prefix 01"},
		{"00-1234567", 0, "invalid prefix 00"},
		{"99-1234567", 0, "invalid prefix 99"},
		{"no ein here", 0, "no EIN"},
	}

	for _, tt := range tests {
		matches, err := d.Detect(ctx, tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if len(matches) != tt.want {
			t.Errorf("%s: Detect(%q) got %d matches, want %d", tt.desc, tt.input, len(matches), tt.want)
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
		{"Call (555) 123-4567", 1, "parenthesized area code"},
		{"+1-555-123-4567", 1, "international format"},
		{"555.123.4567", 1, "dot-separated"},
		{"5551234567", 1, "no separators"},
		{"1-800-555-1234", 1, "toll-free"},
		{"123-456", 0, "too short"},
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
	}
}

func TestZIPDetector(t *testing.T) {
	d := NewZIPDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"ZIP: 90210", 1, "5-digit ZIP"},
		{"90210-1234", 1, "ZIP+4"},
		{"1234", 0, "too short"},
		{"no zip here", 0, "no ZIP"},
	}

	for _, tt := range tests {
		matches, err := d.Detect(ctx, tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if len(matches) != tt.want {
			t.Errorf("%s: Detect(%q) got %d matches, want %d", tt.desc, tt.input, len(matches), tt.want)
		}
	}
}

func TestPassportDetector(t *testing.T) {
	d := NewPassportDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"Passport: 123456789", 1, "9-digit"},
		{"C12345678", 1, "letter prefix"},
		{"1234567", 0, "too short"},
		{"no passport", 0, "no passport"},
	}

	for _, tt := range tests {
		matches, err := d.Detect(ctx, tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if len(matches) != tt.want {
			t.Errorf("%s: Detect(%q) got %d matches, want %d", tt.desc, tt.input, len(matches), tt.want)
		}
	}
}

func TestITINDetector(t *testing.T) {
	d := NewITINDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"ITIN: 912-50-1234", 1, "valid ITIN with group 50"},
		{"970-70-1234", 1, "valid ITIN with group 70"},
		{"912501234", 1, "no dashes"},
		{"no itin here", 0, "no ITIN"},
	}

	for _, tt := range tests {
		matches, err := d.Detect(ctx, tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if len(matches) != tt.want {
			t.Errorf("%s: Detect(%q) got %d matches, want %d", tt.desc, tt.input, len(matches), tt.want)
		}
	}
}

func TestMedicareDetector(t *testing.T) {
	d := NewMedicareDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"MBI: 1EG4TE5MK72", 1, "valid MBI"},
		{"no medicare", 0, "no MBI"},
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
			if m.Type != model.HealthID {
				t.Errorf("expected HealthID type, got %v", m.Type)
			}
		}
	}
}
