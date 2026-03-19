package jp

import (
	"context"
	"testing"

	"github.com/taoq-ai/wuming/domain/model"
	"github.com/taoq-ai/wuming/domain/port"
)

var ctx = context.Background()

// Verify all detectors implement port.Detector.
var (
	_ port.Detector = (*MyNumberDetector)(nil)
	_ port.Detector = (*CorporateNumberDetector)(nil)
	_ port.Detector = (*PhoneDetector)(nil)
	_ port.Detector = (*PostalDetector)(nil)
	_ port.Detector = (*PassportDetector)(nil)
)

func TestMyNumberDetector(t *testing.T) {
	d := NewMyNumberDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"My Number: 123456789018", 1, "valid My Number"},
		{"番号: 000000000000", 1, "all zeros valid"},
		{"123456789019", 0, "invalid check digit"},
		{"12345678901", 0, "too short (11 digits)"},
		{"1234567890123", 0, "too long (13 digits)"},
		{"no number here", 0, "no My Number"},
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
			if m.Locale != "jp" {
				t.Errorf("expected locale 'jp', got %q", m.Locale)
			}
			if m.Type != model.NationalID {
				t.Errorf("expected NationalID type, got %v", m.Type)
			}
			if m.Confidence != 0.85 {
				t.Errorf("expected confidence 0.85, got %v", m.Confidence)
			}
		}
	}
}

func TestCorporateNumberDetector(t *testing.T) {
	d := NewCorporateNumberDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"法人番号: 3234567890123", 1, "valid corporate number"},
		{"9000000000000", 1, "valid all zeros body"},
		{"1234567890123", 0, "invalid check digit"},
		{"234567890123", 0, "too short (12 digits)"},
		{"no corp number", 0, "no corporate number"},
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
			if m.Confidence != 0.80 {
				t.Errorf("expected confidence 0.80, got %v", m.Confidence)
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
		{"Call 090-1234-5678", 1, "mobile 090"},
		{"080-1234-5678", 1, "mobile 080"},
		{"070-1234-5678", 1, "mobile 070"},
		{"+81 90-1234-5678", 1, "international mobile"},
		{"03-1234-5678", 1, "Tokyo landline"},
		{"06-1234-5678", 1, "Osaka landline"},
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
			if m.Confidence != 0.85 {
				t.Errorf("expected confidence 0.85, got %v", m.Confidence)
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
		{"〒100-0001", 1, "with postal mark"},
		{"100-0001", 1, "without postal mark"},
		{"〒 123-4567", 1, "postal mark with space"},
		{"1234567", 0, "no hyphen"},
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
			if m.Type != model.PostalCode {
				t.Errorf("expected PostalCode type, got %v", m.Type)
			}
			if m.Confidence != 0.75 {
				t.Errorf("expected confidence 0.75, got %v", m.Confidence)
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
		{"Passport: TK1234567", 1, "standard format"},
		{"MZ9876543", 1, "another valid format"},
		{"A1234567", 0, "only 1 letter"},
		{"TK123456", 0, "too few digits"},
		{"TK12345678", 0, "too many digits"},
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
			if m.Type != model.Passport {
				t.Errorf("expected Passport type, got %v", m.Type)
			}
			if m.Confidence != 0.70 {
				t.Errorf("expected confidence 0.70, got %v", m.Confidence)
			}
		}
	}
}

func TestAll(t *testing.T) {
	detectors := All()
	if len(detectors) != 5 {
		t.Errorf("All() returned %d detectors, want 5", len(detectors))
	}
}
