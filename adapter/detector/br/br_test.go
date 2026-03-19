package br

import (
	"context"
	"testing"

	"github.com/taoq-ai/wuming/domain/model"
	"github.com/taoq-ai/wuming/domain/port"
)

var ctx = context.Background()

// Verify all detectors implement port.Detector.
var (
	_ port.Detector = (*CPFDetector)(nil)
	_ port.Detector = (*CNPJDetector)(nil)
	_ port.Detector = (*PhoneDetector)(nil)
	_ port.Detector = (*CEPDetector)(nil)
	_ port.Detector = (*PISDetector)(nil)
	_ port.Detector = (*CNHDetector)(nil)
)

func TestCPFDetector(t *testing.T) {
	d := NewCPFDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"CPF: 529.982.247-25", 1, "valid formatted CPF"},
		{"52998224725", 1, "valid unformatted CPF"},
		{"111.111.111-11", 0, "all same digits invalid"},
		{"000.000.000-00", 0, "all zeros invalid"},
		{"529.982.247-00", 0, "wrong check digits"},
		{"123.456.789-00", 0, "invalid check digits"},
		{"no cpf here", 0, "no CPF"},
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
			if m.Locale != "br" {
				t.Errorf("expected locale 'br', got %q", m.Locale)
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

func TestCNPJDetector(t *testing.T) {
	d := NewCNPJDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"CNPJ: 11.222.333/0001-81", 1, "valid formatted CNPJ"},
		{"11222333000181", 1, "valid unformatted CNPJ"},
		{"11.111.111/1111-11", 0, "all same digits invalid"},
		{"11.222.333/0001-00", 0, "wrong check digits"},
		{"no cnpj here", 0, "no CNPJ"},
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
			if m.Locale != "br" {
				t.Errorf("expected locale 'br', got %q", m.Locale)
			}
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
		{"(11) 91234-5678", 1, "mobile with parens"},
		{"+55 11 91234-5678", 1, "mobile with country code"},
		{"(21) 2345-6789", 1, "landline with parens"},
		{"+55 21 2345-6789", 1, "landline with country code"},
		{"11 91234-5678", 1, "mobile without parens"},
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

func TestCEPDetector(t *testing.T) {
	d := NewCEPDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"CEP: 01001-000", 1, "formatted CEP"},
		{"01001000", 1, "unformatted CEP"},
		{"1234567", 0, "too short"},
		{"no cep here", 0, "no CEP"},
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

func TestPISDetector(t *testing.T) {
	d := NewPISDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"PIS: 123.45678.90-0", 1, "valid formatted PIS"},
		{"12345678900", 1, "valid unformatted PIS"},
		{"11111111111", 0, "all same digits invalid"},
		{"no pis here", 0, "no PIS"},
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

func TestCNHDetector(t *testing.T) {
	d := NewCNHDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"CNH: 50680454028", 1, "valid CNH"},
		{"11111111111", 0, "all same digits invalid"},
		{"12345678999", 0, "invalid check digits"},
		{"no cnh here", 0, "no CNH"},
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
			if m.Type != model.DriversLicense {
				t.Errorf("expected DriversLicense type, got %v", m.Type)
			}
		}
	}
}

func TestAll(t *testing.T) {
	detectors := All()
	if len(detectors) != 6 {
		t.Errorf("All() returned %d detectors, want 6", len(detectors))
	}
}
