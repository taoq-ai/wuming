package common

import (
	"context"
	"testing"

	"github.com/taoq-ai/wuming/domain/model"
	"github.com/taoq-ai/wuming/domain/port"
)

var ctx = context.Background()

// Verify all detectors implement port.Detector.
var (
	_ port.Detector = (*EmailDetector)(nil)
	_ port.Detector = (*CreditCardDetector)(nil)
	_ port.Detector = (*IPDetector)(nil)
	_ port.Detector = (*URLDetector)(nil)
	_ port.Detector = (*IBANDetector)(nil)
	_ port.Detector = (*MACDetector)(nil)
)

func TestEmailDetector(t *testing.T) {
	d := NewEmailDetector()

	tests := []struct {
		input string
		want  int
	}{
		{"email john@example.com here", 1},
		{"user.name+tag@domain.co.uk", 1},
		{"a@b.cd and x@y.ef", 2},
		{"no email here", 0},
		{"@@invalid", 0},
		{"test@", 0},
	}

	for _, tt := range tests {
		matches, err := d.Detect(ctx, tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if len(matches) != tt.want {
			t.Errorf("Detect(%q) got %d matches, want %d", tt.input, len(matches), tt.want)
		}
		for _, m := range matches {
			if m.Type != model.Email {
				t.Errorf("expected Email type, got %v", m.Type)
			}
		}
	}
}

func TestCreditCardDetector(t *testing.T) {
	d := NewCreditCardDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"4111111111111111", 1, "Visa"},
		{"4111 1111 1111 1111", 1, "Visa with spaces"},
		{"4111-1111-1111-1111", 1, "Visa with dashes"},
		{"5500000000000004", 1, "Mastercard"},
		{"378282246310005", 1, "Amex"},
		{"6011111111111117", 1, "Discover"},
		{"1234567890123456", 0, "invalid Luhn"},
		{"123", 0, "too short"},
		{"no card here", 0, "no digits"},
	}

	for _, tt := range tests {
		matches, err := d.Detect(ctx, tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if len(matches) != tt.want {
			t.Errorf("%s: Detect(%q) got %d matches, want %d", tt.desc, tt.input, len(matches), tt.want)
		}
	}
}

func TestLuhn(t *testing.T) {
	if !luhn("4111111111111111") {
		t.Error("Visa test number should pass Luhn")
	}
	if luhn("4111111111111112") {
		t.Error("Invalid number should fail Luhn")
	}
	if luhn("") {
		t.Error("Empty string should fail Luhn")
	}
}

func TestIPDetector(t *testing.T) {
	d := NewIPDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"IP is 192.168.1.1 here", 1, "simple IPv4"},
		{"10.0.0.1 and 172.16.0.1", 2, "two IPv4"},
		{"2001:0db8:85a3:0000:0000:8a2e:0370:7334", 1, "full IPv6"},
		{"::1", 1, "loopback IPv6"},
		{"999.999.999.999", 0, "invalid octets"},
		{"no ip here", 0, "no IP"},
	}

	for _, tt := range tests {
		matches, err := d.Detect(ctx, tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if len(matches) != tt.want {
			t.Errorf("%s: Detect(%q) got %d matches, want %d", tt.desc, tt.input, len(matches), tt.want)
		}
	}
}

func TestURLDetector(t *testing.T) {
	d := NewURLDetector()

	tests := []struct {
		input string
		want  int
	}{
		{"visit https://example.com/path?q=1 now", 1},
		{"http://localhost:8080", 1},
		{"https://a.com and https://b.com", 2},
		{"ftp://notmatched.com", 0},
		{"no url here", 0},
	}

	for _, tt := range tests {
		matches, err := d.Detect(ctx, tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if len(matches) != tt.want {
			t.Errorf("Detect(%q) got %d matches, want %d", tt.input, len(matches), tt.want)
		}
	}
}

func TestIBANDetector(t *testing.T) {
	d := NewIBANDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"IBAN: NL91ABNA0417164300", 1, "valid NL IBAN"},
		{"DE89370400440532013000", 1, "valid DE IBAN"},
		{"GB29NWBK60161331926819", 1, "valid GB IBAN"},
		{"NL00ABNA0000000000", 0, "invalid checksum"},
		{"no iban here", 0, "no IBAN"},
	}

	for _, tt := range tests {
		matches, err := d.Detect(ctx, tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if len(matches) != tt.want {
			t.Errorf("%s: Detect(%q) got %d matches, want %d", tt.desc, tt.input, len(matches), tt.want)
		}
	}
}

func TestMACDetector(t *testing.T) {
	d := NewMACDetector()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"MAC: 00:1A:2B:3C:4D:5E", 1, "colon-separated"},
		{"00-1A-2B-3C-4D-5E", 1, "hyphen-separated"},
		{"001A.2B3C.4D5E", 1, "dot-separated"},
		{"not a mac address", 0, "no MAC"},
		{"00:1A:2B", 0, "too short"},
	}

	for _, tt := range tests {
		matches, err := d.Detect(ctx, tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if len(matches) != tt.want {
			t.Errorf("%s: Detect(%q) got %d matches, want %d", tt.desc, tt.input, len(matches), tt.want)
		}
	}
}
