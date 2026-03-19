package fr

import (
	"context"
	"testing"

	"github.com/taoq-ai/wuming/domain/model"
	"github.com/taoq-ai/wuming/domain/port"
)

var ctx = context.Background()

// Verify all detectors implement port.Detector.
var (
	_ port.Detector = (*NIRDetector)(nil)
	_ port.Detector = (*NIFDetector)(nil)
	_ port.Detector = (*PhoneDetector)(nil)
	_ port.Detector = (*PostalDetector)(nil)
	_ port.Detector = (*IDCardDetector)(nil)
)

func TestNIRDetector(t *testing.T) {
	d := NewNIRDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		// key = 97 - (2840578006084 % 97) = 90.
		{"NIR: 2 84 05 78 006 084 90", 1, "valid NIR with spaces"},
		{"284057800608490", 1, "valid NIR without spaces"},
		// Invalid control key.
		{"2 84 05 78 006 084 99", 0, "invalid control key"},
		// Corsica 2A: key = 97 - ((1850120033084 - 1000000) % 97) = 34.
		{"1 85 01 2A 033 084 34", 1, "valid NIR Corsica 2A"},
		// Corsica 2B: key = 97 - ((2900320001002 - 2000000) % 97) = 64.
		{"2 90 03 2B 001 002 64", 1, "valid NIR Corsica 2B"},
		{"no nir here", 0, "no NIR"},
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
			if m.Locale != "fr" {
				t.Errorf("expected locale 'fr', got %q", m.Locale)
			}
			if m.Type != model.NationalID {
				t.Errorf("expected NationalID type, got %v", m.Type)
			}
		}
	}
}

func TestNIFDetector(t *testing.T) {
	d := NewNIFDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"NIF: 0123456789012", 1, "valid NIF starting with 0"},
		{"1234567890123", 1, "valid NIF starting with 1"},
		{"32 99 123 456 789", 1, "formatted NIF"},
		{"4123456789012", 0, "invalid first digit 4"},
		{"9123456789012", 0, "invalid first digit 9"},
		{"no nif here", 0, "no NIF"},
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

func TestPhoneDetector(t *testing.T) {
	d := NewPhoneDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"Tel: 06 12 34 56 78", 1, "mobile 06 with spaces"},
		{"07.12.34.56.78", 1, "mobile 07 with dots"},
		{"01 23 45 67 89", 1, "landline 01"},
		{"02-34-56-78-90", 1, "landline 02 with dashes"},
		{"0345678901", 1, "landline 03 no separators"},
		{"04 56 78 90 12", 1, "landline 04"},
		{"05 67 89 01 23", 1, "landline 05"},
		{"+33 6 12 34 56 78", 1, "international mobile +33 6"},
		{"+33 7 12 34 56 78", 1, "international mobile +33 7"},
		{"+33 1 23 45 67 89", 1, "international landline +33 1"},
		{"+33.5.67.89.01.23", 1, "international with dots"},
		{"08 12 34 56 78", 0, "08 prefix not matched"},
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
		{"Code postal: 75001", 1, "Paris 1er"},
		{"13100", 1, "Aix-en-Provence"},
		{"97400", 1, "DOM department 97"},
		{"98000", 1, "department 98"},
		{"00123", 0, "invalid department 00"},
		{"96000", 0, "invalid department 96"},
		{"99000", 0, "invalid department 99"},
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
		}
	}
}

func TestIDCardDetector(t *testing.T) {
	d := NewIDCardDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"CNI: 880123456789", 1, "old format 12 digits"},
		{"CNI: F4R7X2K9Q", 1, "new format 9 alphanumeric"},
		{"AB12CD345", 1, "new format mixed"},
		{"no cni here", 0, "no CNI"},
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
