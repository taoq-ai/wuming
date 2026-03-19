package kr

import (
	"context"
	"testing"

	"github.com/taoq-ai/wuming/domain/model"
	"github.com/taoq-ai/wuming/domain/port"
)

var ctx = context.Background()

// Verify all detectors implement port.Detector.
var (
	_ port.Detector = (*RRNDetector)(nil)
	_ port.Detector = (*PhoneDetector)(nil)
	_ port.Detector = (*PostalDetector)(nil)
	_ port.Detector = (*PassportDetector)(nil)
)

func TestRRNDetector(t *testing.T) {
	d := NewRRNDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"RRN: 900101-1000006", 1, "valid male 1900s"},
		{"850515-2123451", 1, "valid female 1900s"},
		{"000101-3000008", 1, "valid male 2000s"},
		{"900101-1000007", 0, "invalid check digit"},
		{"000000-1234567", 0, "invalid date 000000"},
		{"901301-1000001", 0, "invalid month 13"},
		{"900100-1000001", 0, "invalid day 00"},
		{"900101-9000001", 0, "invalid gender digit 9"},
		{"no rrn here", 0, "no RRN"},
		{"12345-1234567", 0, "only 5 digit date prefix"},
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
			if m.Locale != "kr" {
				t.Errorf("expected locale 'kr', got %q", m.Locale)
			}
			if m.Type != model.NationalID {
				t.Errorf("expected NationalID type, got %v", m.Type)
			}
			if m.Confidence != 0.90 {
				t.Errorf("expected confidence 0.90, got %f", m.Confidence)
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
		{"Call 010-1234-5678", 1, "mobile with dashes"},
		{"010 1234 5678", 1, "mobile with spaces"},
		{"+82 10-1234-5678", 1, "mobile international prefix"},
		{"+82-10-1234-5678", 1, "mobile international with dash"},
		{"02-123-4567", 1, "Seoul landline 7-digit"},
		{"02-1234-5678", 1, "Seoul landline 8-digit"},
		{"031-123-4567", 1, "Gyeonggi landline"},
		{"051-1234-5678", 1, "Busan landline"},
		{"123-4567", 0, "too short, no area code"},
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
				t.Errorf("expected confidence 0.85, got %f", m.Confidence)
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
		{"Postal: 06236", 1, "valid Seoul postal code"},
		{"48058", 1, "valid Busan postal code"},
		{"1234", 0, "too short"},
		{"123456", 0, "too long"},
		{"no postal code", 0, "no postal code"},
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
			if m.Confidence != 0.50 {
				t.Errorf("expected confidence 0.50, got %f", m.Confidence)
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
		{"Passport: M12345678", 1, "regular passport"},
		{"R12345678", 1, "residence passport letter"},
		{"S98765432", 1, "other letter prefix"},
		{"12345678", 0, "no letter prefix"},
		{"AB12345678", 0, "two letter prefix"},
		{"M1234567", 0, "too few digits"},
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
			if m.Confidence != 0.70 {
				t.Errorf("expected confidence 0.70, got %f", m.Confidence)
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
