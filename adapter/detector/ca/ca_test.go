package ca

import (
	"context"
	"testing"

	"github.com/taoq-ai/wuming/domain/model"
	"github.com/taoq-ai/wuming/domain/port"
)

var ctx = context.Background()

// Verify all detectors implement port.Detector.
var (
	_ port.Detector = (*SINDetector)(nil)
	_ port.Detector = (*PhoneDetector)(nil)
	_ port.Detector = (*PostalCodeDetector)(nil)
	_ port.Detector = (*PassportDetector)(nil)
)

func TestSINDetector(t *testing.T) {
	d := NewSINDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"SIN: 046 454 286", 1, "valid SIN with spaces"},
		{"SIN: 046-454-286", 1, "valid SIN with dashes"},
		{"046454286", 1, "valid SIN no separators"},
		{"000 000 000", 0, "all zeros invalid"},
		{"123 456 789", 0, "fails Luhn check"},
		{"12 345 678", 0, "too few digits in first group"},
		{"no sin here", 0, "no SIN"},
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
			if m.Locale != "ca" {
				t.Errorf("expected locale 'ca', got %q", m.Locale)
			}
			if m.Type != model.NationalID {
				t.Errorf("expected NationalID type, got %v", m.Type)
			}
			if m.Confidence != 0.85 {
				t.Errorf("expected confidence 0.85, got %f", m.Confidence)
			}
		}
	}
}

func TestSINLuhn(t *testing.T) {
	tests := []struct {
		digits string
		valid  bool
	}{
		{"046454286", true},
		{"123456789", false},
		{"000000000", false},
	}
	for _, tt := range tests {
		if got := luhnValid(tt.digits); got != tt.valid {
			t.Errorf("luhnValid(%q) = %v, want %v", tt.digits, got, tt.valid)
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
		{"Call (416) 555-1234", 1, "Toronto area code with parens"},
		{"+1-604-555-1234", 1, "Vancouver with +1 prefix"},
		{"514.555.1234", 1, "Montreal dot-separated"},
		{"9055551234", 1, "Hamilton no separators"},
		{"1-867-555-1234", 1, "Northern territories"},
		{"(555) 123-4567", 0, "non-Canadian area code 555"},
		{"(212) 555-1234", 0, "US area code 212"},
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
			if m.Locale != "ca" {
				t.Errorf("expected locale 'ca', got %q", m.Locale)
			}
			if m.Confidence != 0.80 {
				t.Errorf("expected confidence 0.80, got %f", m.Confidence)
			}
		}
	}
}

func TestPostalCodeDetector(t *testing.T) {
	d := NewPostalCodeDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"Postal: K1A 0B1", 1, "Ottawa with space"},
		{"V6B3K9", 1, "Vancouver no space"},
		{"M5V 2T6", 1, "Toronto"},
		{"D1A 0B1", 0, "starts with D (invalid)"},
		{"F1A 0B1", 0, "starts with F (invalid)"},
		{"W1A 0B1", 0, "starts with W (invalid)"},
		{"Z1A 0B1", 0, "starts with Z (invalid)"},
		{"no postal here", 0, "no postal code"},
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
			if m.Locale != "ca" {
				t.Errorf("expected locale 'ca', got %q", m.Locale)
			}
			if m.Type != model.PostalCode {
				t.Errorf("expected PostalCode type, got %v", m.Type)
			}
			if m.Confidence != 0.90 {
				t.Errorf("expected confidence 0.90, got %f", m.Confidence)
			}
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
		{"Passport: AB123456", 1, "valid Canadian passport"},
		{"GK987654", 1, "another valid passport"},
		{"A1234567", 0, "only one letter"},
		{"ABC12345", 0, "three letters"},
		{"AB12345", 0, "only 5 digits"},
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
		for _, m := range matches {
			if m.Locale != "ca" {
				t.Errorf("expected locale 'ca', got %q", m.Locale)
			}
			if m.Type != model.Passport {
				t.Errorf("expected Passport type, got %v", m.Type)
			}
			if m.Confidence != 0.65 {
				t.Errorf("expected confidence 0.65, got %f", m.Confidence)
			}
		}
	}
}

func TestAll(t *testing.T) {
	detectors := All()
	if len(detectors) != 4 {
		t.Errorf("All() returned %d detectors, want 4", len(detectors))
	}
}
