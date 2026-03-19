package gb

import (
	"context"
	"strings"
	"testing"

	"github.com/taoq-ai/wuming/domain/model"
	"github.com/taoq-ai/wuming/domain/port"
)

var ctx = context.Background()

// Verify all detectors implement port.Detector.
var (
	_ port.Detector = (*NINDetector)(nil)
	_ port.Detector = (*NHSDetector)(nil)
	_ port.Detector = (*PhoneDetector)(nil)
	_ port.Detector = (*PostcodeDetector)(nil)
	_ port.Detector = (*UTRDetector)(nil)
)

func TestNINDetector(t *testing.T) {
	d := NewNINDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"NIN: AB 12 34 56 A", 1, "standard format with spaces"},
		{"AB123456A", 1, "no spaces"},
		{"ab 12 34 56 c", 1, "lowercase"},
		{"AB 12 34 56 D", 1, "suffix D"},
		// Invalid prefixes: first letter D, F, I, Q, U, V.
		{"DA 12 34 56 A", 0, "first letter D invalid"},
		{"FA 12 34 56 A", 0, "first letter F invalid"},
		{"IA 12 34 56 A", 0, "first letter I invalid"},
		{"QA 12 34 56 A", 0, "first letter Q invalid"},
		{"UA 12 34 56 A", 0, "first letter U invalid"},
		{"VA 12 34 56 A", 0, "first letter V invalid"},
		// Invalid prefixes: second letter D, F, I, O, Q, U, V.
		{"AD 12 34 56 A", 0, "second letter D invalid"},
		{"AF 12 34 56 A", 0, "second letter F invalid"},
		{"AI 12 34 56 A", 0, "second letter I invalid"},
		{"AO 12 34 56 A", 0, "second letter O invalid"},
		{"AQ 12 34 56 A", 0, "second letter Q invalid"},
		{"AU 12 34 56 A", 0, "second letter U invalid"},
		{"AV 12 34 56 A", 0, "second letter V invalid"},
		// Forbidden two-letter prefixes.
		{"BG 12 34 56 A", 0, "prefix BG invalid"},
		{"GB 12 34 56 A", 0, "prefix GB invalid"},
		{"NK 12 34 56 A", 0, "prefix NK invalid"},
		{"KN 12 34 56 A", 0, "prefix KN invalid"},
		{"TN 12 34 56 A", 0, "prefix TN invalid"},
		{"NT 12 34 56 A", 0, "prefix NT invalid"},
		{"ZZ 12 34 56 A", 0, "prefix ZZ invalid"},
		// Invalid suffix.
		{"AB 12 34 56 E", 0, "suffix E invalid"},
		{"AB 12 34 56 Z", 0, "suffix Z invalid"},
		// No NIN.
		{"no nin here", 0, "no NIN"},
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
			if m.Locale != "gb" {
				t.Errorf("expected locale 'gb', got %q", m.Locale)
			}
			if m.Type != model.NationalID {
				t.Errorf("expected NationalID type, got %v", m.Type)
			}
		}
	}
}

func TestNHSDetector(t *testing.T) {
	d := NewNHSDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		// 943 476 5919: weights 10*9+9*4+8*3+7*4+6*7+5*6+4*5+3*9+2*1 = 90+36+24+28+42+30+20+27+2 = 299
		// 299 mod 11 = 2; check = 11-2 = 9. Last digit is 9. Valid.
		{"NHS: 943 476 5919", 1, "valid with spaces"},
		{"9434765919", 1, "valid no spaces"},
		// Invalid check digit: change last digit.
		{"943 476 5910", 0, "invalid check digit"},
		// Remainder == 10 means invalid.
		// Need to craft: sum mod 11 == 1 => remainder = 10 => invalid.
		// 0000000000: sum = 0, 11 - 0 = 11 => check = 0. Last = 0. Valid actually.
		{"no nhs here", 0, "no NHS"},
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

func TestNHSMod11Validation(t *testing.T) {
	// Directly test the mod-11 logic with known values.
	tests := []struct {
		digits string
		valid  bool
		desc   string
	}{
		{"9434765919", true, "valid NHS number"},
		{"9434765910", false, "wrong check digit"},
		{"4505577104", true, "another valid NHS number"},
		// 1234567881: weights 10*1+9*2+8*3+7*4+6*5+5*6+4*7+3*8+2*8 = 10+18+24+28+30+30+28+24+16 = 208
		// 208 mod 11 = 10; remainder = 11 - 10 = 1. Last digit = 1. Valid.
		{"1234567881", true, "computed valid"},
	}

	for _, tt := range tests {
		got := isValidNHS(tt.digits)
		if got != tt.valid {
			t.Errorf("%s: isValidNHS(%q) = %v, want %v", tt.desc, tt.digits, got, tt.valid)
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
		{"Call 07911 123456", 1, "mobile standard"},
		{"+44 7911 123456", 1, "mobile international"},
		{"07911-123-456", 1, "mobile with dashes"},
		{"07911.123.456", 1, "mobile with dots"},
		{"020 7946 0958", 1, "London landline"},
		{"+44 20 7946 0958", 1, "London international"},
		{"0113 496 0123", 1, "landline 3-digit area code"},
		{"+44 113 496 0123", 1, "landline international"},
		{"01onal call", 0, "not a phone number"},
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
		{"M1 1AA", 1, "A9 9AA format"},
		{"M60 1NW", 1, "A99 9AA format"},
		{"W1A 1HQ", 1, "A9A 9AA format"},
		{"CR2 6XH", 1, "AA9 9AA format"},
		{"DN55 1PT", 1, "AA99 9AA format"},
		{"EC1A 1BB", 1, "AA9A 9AA format"},
		{"ec1a 1bb", 1, "lowercase"},
		{"SW1A 2AA", 1, "Downing Street"},
		{"not a postcode", 0, "no postcode"},
		{"12345", 0, "US ZIP not matched"},
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
			// Value should be normalized to uppercase.
			if m.Value != strings.ToUpper(m.Value) {
				t.Errorf("value should be uppercase, got %q", m.Value)
			}
		}
	}
}

func TestPostcodeNormalization(t *testing.T) {
	d := NewPostcodeDetector()
	matches, err := d.Detect(ctx, "ec1a 1bb")
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].Value != "EC1A 1BB" {
		t.Errorf("expected normalized value 'EC1A 1BB', got %q", matches[0].Value)
	}
}

func TestUTRDetector(t *testing.T) {
	d := NewUTRDetector()

	tests := []struct {
		input      string
		want       int
		desc       string
		confidence float64
	}{
		{"UTR: 1234567890", 1, "with UTR context", 0.85},
		{"tax reference: 1234567890", 1, "with tax reference context", 0.85},
		{"Unique Taxpayer Reference 1234567890", 1, "with full context", 0.85},
		{"1234567890", 1, "bare 10-digit number", 0.55},
		{"123456789", 0, "9 digits too short", 0},
		{"12345678901", 0, "11 digits too long", 0},
		{"no utr here", 0, "no UTR", 0},
	}

	for _, tt := range tests {
		matches, err := d.Detect(ctx, tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if len(matches) != tt.want {
			t.Errorf("%s: Detect(%q) got %d matches, want %d", tt.desc, tt.input, len(matches), tt.want)
			continue
		}
		for _, m := range matches {
			if m.Type != model.TaxID {
				t.Errorf("expected TaxID type, got %v", m.Type)
			}
			if m.Confidence != tt.confidence {
				t.Errorf("%s: expected confidence %v, got %v", tt.desc, tt.confidence, m.Confidence)
			}
		}
	}
}
