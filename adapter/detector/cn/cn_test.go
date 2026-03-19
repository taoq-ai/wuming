package cn

import (
	"context"
	"testing"

	"github.com/taoq-ai/wuming/domain/model"
	"github.com/taoq-ai/wuming/domain/port"
)

var ctx = context.Background()

// Verify all detectors implement port.Detector.
var (
	_ port.Detector = (*ResidentIDDetector)(nil)
	_ port.Detector = (*PhoneDetector)(nil)
	_ port.Detector = (*PostalDetector)(nil)
	_ port.Detector = (*PassportDetector)(nil)
	_ port.Detector = (*USCCDetector)(nil)
)

func TestResidentIDDetector(t *testing.T) {
	d := NewResidentIDDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		// 11010519491231002X is a well-known example: Beijing, 1949-12-31, check digit X.
		{"ID: 11010519491231002X", 1, "valid resident ID with X check"},
		// 110101199003070732: Beijing, 1990-03-07, check digit 2.
		{"身份证号码 110101199003070732", 1, "valid resident ID numeric check"},
		{"110105194912310021", 0, "invalid check digit"},
		{"11010519491302002X", 0, "invalid month 13"},
		{"123456", 0, "too short"},
		{"no id here", 0, "no resident ID"},
		// Future birth date should be rejected.
		{"110105209912310021", 0, "future birth date"},
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
			if m.Locale != "cn" {
				t.Errorf("expected locale 'cn', got %q", m.Locale)
			}
			if m.Type != model.NationalID {
				t.Errorf("expected NationalID type, got %v", m.Type)
			}
			if m.Confidence != 0.90 {
				t.Errorf("expected confidence 0.90, got %v", m.Confidence)
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
		{"电话: 13812345678", 1, "standard mobile"},
		{"+86 13912345678", 1, "with country code"},
		{"+86-139-1234-5678", 1, "with country code and dashes"},
		{"138 1234 5678", 1, "with spaces"},
		{"12012345678", 0, "second digit < 3 invalid"},
		{"1381234567", 0, "too short (10 digits)"},
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
		{"邮编: 100000", 1, "Beijing postal code"},
		{"518000", 1, "Shenzhen postal code"},
		{"830000", 1, "Urumqi postal code"},
		{"012345", 0, "starts with 0 invalid"},
		{"900000", 0, "starts with 9 invalid"},
		{"12345", 0, "too short (5 digits)"},
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
			if m.Confidence != 0.55 {
				t.Errorf("expected confidence 0.55, got %v", m.Confidence)
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
		{"Passport: E12345678", 1, "e-passport"},
		{"护照 G12345678", 1, "regular passport"},
		{"D1234567", 1, "diplomatic passport"},
		{"S1234567", 1, "service passport 7-digit"},
		{"S12345678", 1, "service passport 8-digit"},
		{"E1234567", 0, "e-passport too short"},
		{"A12345678", 0, "invalid prefix"},
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
			if m.Confidence != 0.75 {
				t.Errorf("expected confidence 0.75, got %v", m.Confidence)
			}
		}
	}
}

func TestUSCCDetector(t *testing.T) {
	d := NewUSCCDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		// 91110000710931153R: computed valid USCC with check char R.
		{"信用代码: 91110000710931153R", 1, "valid USCC"},
		{"91110000710931153A", 0, "invalid check digit"},
		{"12345678901234567", 0, "too short (17 chars)"},
		{"no uscc here", 0, "no USCC"},
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
			if m.Confidence != 0.85 {
				t.Errorf("expected confidence 0.85, got %v", m.Confidence)
			}
		}
	}
}

func TestAllDetectors(t *testing.T) {
	detectors := All()
	if len(detectors) != 5 {
		t.Errorf("All() returned %d detectors, want 5", len(detectors))
	}
}
