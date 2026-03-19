package au

import (
	"context"
	"testing"

	"github.com/taoq-ai/wuming/domain/model"
	"github.com/taoq-ai/wuming/domain/port"
)

var ctx = context.Background()

// Verify all detectors implement port.Detector.
var (
	_ port.Detector = (*TFNDetector)(nil)
	_ port.Detector = (*MedicareDetector)(nil)
	_ port.Detector = (*ABNDetector)(nil)
	_ port.Detector = (*PhoneDetector)(nil)
	_ port.Detector = (*PostcodeDetector)(nil)
)

func TestTFNDetector(t *testing.T) {
	d := NewTFNDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"TFN: 123 456 782", 1, "valid TFN with spaces"},
		{"123456782", 1, "valid TFN no separators"},
		{"987-654-303", 1, "valid TFN with dashes"},
		{"123 456 789", 0, "invalid checksum"},
		{"12345678", 0, "too few digits"},
		{"no tfn here", 0, "no TFN"},
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
			if m.Locale != "au" {
				t.Errorf("expected locale 'au', got %q", m.Locale)
			}
			if m.Type != model.TaxID {
				t.Errorf("expected TaxID type, got %v", m.Type)
			}
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
		{"Medicare: 2123 45670 1", 1, "valid Medicare with issue number"},
		{"21234567 01", 1, "valid Medicare with space before issue"},
		{"2123456701", 1, "valid Medicare no separators"},
		{"1123456789", 0, "invalid: first digit 1"},
		{"7123456789", 0, "invalid: first digit 7"},
		{"no medicare", 0, "no Medicare number"},
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

func TestABNDetector(t *testing.T) {
	d := NewABNDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"ABN: 51 824 753 556", 1, "valid ABN with spaces"},
		{"51824753556", 1, "valid ABN no separators"},
		{"53 004 085 616", 1, "valid ABN (Telstra)"},
		{"12 345 678 901", 0, "invalid checksum"},
		{"no abn here", 0, "no ABN"},
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
		{"Call 0412 345 678", 1, "mobile with spaces"},
		{"+61 412 345 678", 1, "international mobile"},
		{"(02) 1234 5678", 1, "landline with parens"},
		{"+61 2 1234 5678", 1, "international landline"},
		{"0412345678", 1, "mobile no separators"},
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
		for _, m := range matches {
			if m.Type != model.Phone {
				t.Errorf("expected Phone type, got %v", m.Type)
			}
		}
	}
}

func TestPostcodeDetector(t *testing.T) {
	d := NewPostcodeDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"Sydney 2000", 1, "NSW postcode"},
		{"Melbourne 3000", 1, "VIC postcode"},
		{"Brisbane 4000", 1, "QLD postcode"},
		{"Adelaide 5000", 1, "SA postcode"},
		{"Perth 6000", 1, "WA postcode"},
		{"Hobart 7000", 1, "TAS postcode"},
		{"Darwin 0800", 1, "NT postcode"},
		{"Canberra 2600", 1, "ACT postcode"},
		{"0001", 0, "invalid: below range"},
		{"no postcode here", 0, "no postcode"},
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
			if m.Confidence != 0.55 {
				t.Errorf("expected confidence 0.55, got %f", m.Confidence)
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
