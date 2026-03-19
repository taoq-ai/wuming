package in

import (
	"context"
	"testing"

	"github.com/taoq-ai/wuming/domain/model"
	"github.com/taoq-ai/wuming/domain/port"
)

var ctx = context.Background()

// Verify all detectors implement port.Detector.
var (
	_ port.Detector = (*AadhaarDetector)(nil)
	_ port.Detector = (*PANDetector)(nil)
	_ port.Detector = (*PhoneDetector)(nil)
	_ port.Detector = (*PINCodeDetector)(nil)
	_ port.Detector = (*PassportDetector)(nil)
	_ port.Detector = (*GSTINDetector)(nil)
)

func TestAadhaarDetector(t *testing.T) {
	d := NewAadhaarDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"Aadhaar: 2345 6789 0124", 1, "spaced format valid Verhoeff"},
		{"234567890124", 1, "no spaces valid Verhoeff"},
		{"9876 5432 1012", 1, "starting with 9 valid Verhoeff"},
		{"5123 4567 8903", 1, "starting with 5 valid Verhoeff"},
		{"234567890120", 0, "invalid Verhoeff checksum"},
		{"0123 4567 8901", 0, "starts with 0 invalid"},
		{"1234 5678 9012", 0, "starts with 1 invalid"},
		{"no aadhaar here", 0, "no Aadhaar"},
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
			if m.Locale != "in" {
				t.Errorf("expected locale 'in', got %q", m.Locale)
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

func TestPANDetector(t *testing.T) {
	d := NewPANDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"PAN: ABCPE1234F", 1, "valid PAN person"},
		{"AAACB1234C", 1, "valid PAN company (4th=C)"},
		{"AAAHB1234C", 1, "valid PAN HUF (4th=H)"},
		{"ABCXE1234F", 0, "invalid 4th char X"},
		{"ABCP1234F", 0, "too short"},
		{"abcpe1234f", 0, "lowercase invalid"},
		{"no pan here", 0, "no PAN"},
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
		{"+91 98765 43210", 1, "with +91 prefix and spaces"},
		{"+91-9876543210", 1, "with +91 prefix and dash"},
		{"09876543210", 1, "with 0 prefix"},
		{"9876543210", 1, "bare 10-digit mobile"},
		{"6123456789", 1, "starting with 6"},
		{"5123456789", 0, "starting with 5 invalid"},
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

func TestPINCodeDetector(t *testing.T) {
	d := NewPINCodeDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"PIN: 110001", 1, "valid Delhi PIN"},
		{"400001", 1, "valid Mumbai PIN"},
		{"999999", 1, "valid max PIN"},
		{"012345", 0, "starts with 0 invalid"},
		{"12345", 0, "5 digits too short"},
		{"no pin here", 0, "no PIN"},
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

func TestPassportDetector(t *testing.T) {
	d := NewPassportDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"Passport: A1234567", 1, "valid passport"},
		{"J8765432", 1, "valid passport J prefix"},
		{"12345678", 0, "no letter prefix"},
		{"AB1234567", 0, "two letter prefix"},
		{"no passport here", 0, "no passport"},
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
		}
	}
}

func TestGSTINDetector(t *testing.T) {
	d := NewGSTINDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"GSTIN: 27AAPFU0939F1ZV", 1, "valid Maharashtra GSTIN"},
		{"06BZAHM6385P6Z2", 1, "valid state code 06"},
		{"37AAPFU0939F1ZV", 1, "valid state code 37"},
		{"00AAPFU0939F1ZV", 0, "state code 00 invalid"},
		{"38AAPFU0939F1ZV", 0, "state code 38 invalid"},
		{"27AAPFU0939F1AV", 0, "missing Z at position 13"},
		{"no gstin here", 0, "no GSTIN"},
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
			if m.Confidence != 0.90 {
				t.Errorf("expected confidence 0.90, got %f", m.Confidence)
			}
		}
	}
}

func TestAll(t *testing.T) {
	detectors := All()
	if len(detectors) != 6 {
		t.Errorf("All() returned %d detectors, want 6", len(detectors))
	}

	names := make(map[string]bool)
	for _, d := range detectors {
		if names[d.Name()] {
			t.Errorf("duplicate detector name: %s", d.Name())
		}
		names[d.Name()] = true
	}
}
