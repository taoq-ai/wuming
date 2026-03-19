package eu

import (
	"context"
	"testing"

	"github.com/taoq-ai/wuming/domain/model"
	"github.com/taoq-ai/wuming/domain/port"
)

var ctx = context.Background()

// Verify all detectors implement port.Detector.
var (
	_ port.Detector = (*VATDetector)(nil)
	_ port.Detector = (*PassportMRZDetector)(nil)
)

func TestVATDetector(t *testing.T) {
	d := NewVATDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"VAT: ATU12345678", 1, "Austria"},
		{"VAT: BE0123456789", 1, "Belgium"},
		{"VAT: DE123456789", 1, "Germany"},
		{"VAT: FR1A123456789", 1, "France alphanumeric"},
		{"VAT: NL123456789B01", 1, "Netherlands"},
		{"VAT: PL1234567890", 1, "Poland"},
		{"VAT: EL123456789", 1, "Greece"},
		{"VAT: SE123456789012", 1, "Sweden"},
		{"VAT: IE1A12345AB", 1, "Ireland two-letter suffix"},
		{"VAT: CY12345678A", 1, "Cyprus"},
		{"VAT: ESA1234567B", 1, "Spain"},
		{"VAT: XX123456789", 0, "unknown country prefix"},
		{"no vat here", 0, "no VAT"},
		{"DE12345", 0, "too short for German VAT"},
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
			if m.Locale != "eu" {
				t.Errorf("expected locale 'eu', got %q", m.Locale)
			}
			if m.Type != model.TaxID {
				t.Errorf("expected TaxID type, got %v", m.Type)
			}
			if m.Confidence != 0.90 {
				t.Errorf("expected confidence 0.90, got %f", m.Confidence)
			}
		}
	}
}

func TestPassportMRZDetector(t *testing.T) {
	d := NewPassportMRZDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{
			input: "P<UTOERIKSSON<<ANNA<MARIA<<<<<<<<<<<<<<<<<<<\nL898902C36UTO7408122F1204159ZE184226B<<<<<10",
			want:  1,
			desc:  "valid TD3 MRZ",
		},
		{
			input: "P<GBRSMITH<<JOHN<<<<<<<<<<<<<<<<<<<<<<<<<<<<\nAB12345670GBR8501011M2501015<<<<<<<<<<<<<<00",
			want:  1,
			desc:  "UK passport MRZ",
		},
		{
			input: "P<DEUMULLER<<HANS<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\nC01X00T478DEU6408125F2010315D<<<<<<<<<<<<<<04",
			want:  0,
			desc:  "line 1 too long (>44 chars) — invalid",
		},
		{
			input: "not a passport",
			want:  0,
			desc:  "no MRZ",
		},
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
			if m.Locale != "eu" {
				t.Errorf("expected locale 'eu', got %q", m.Locale)
			}
			if m.Type != model.Passport {
				t.Errorf("expected Passport type, got %v", m.Type)
			}
			if m.Confidence != 0.95 {
				t.Errorf("expected confidence 0.95, got %f", m.Confidence)
			}
		}
	}
}
