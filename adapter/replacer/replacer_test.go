package replacer

import (
	"testing"

	"github.com/taoq-ai/wuming/domain/model"
)

func matches() []model.Match {
	return []model.Match{
		{Type: model.Email, Value: "john@example.com", Start: 14, End: 30, Confidence: 1.0},
		{Type: model.Phone, Value: "555-1234", Start: 45, End: 53, Confidence: 0.9},
	}
}

const input = "Contact me at john@example.com or call me at 555-1234 please."

func TestRedact(t *testing.T) {
	r := NewRedact()
	got, err := r.Replace(input, matches())
	if err != nil {
		t.Fatal(err)
	}
	want := "Contact me at [EMAIL] or call me at [PHONE] please."
	if got != want {
		t.Errorf("Redact:\n got: %q\nwant: %q", got, want)
	}

	if r.Name() != "redact" {
		t.Errorf("Name() = %q, want %q", r.Name(), "redact")
	}
}

func TestRedactEmpty(t *testing.T) {
	r := NewRedact()
	got, err := r.Replace(input, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != input {
		t.Errorf("expected unchanged text for empty matches")
	}
}

func TestMask(t *testing.T) {
	m := NewMask()
	got, err := m.Replace(input, matches())
	if err != nil {
		t.Fatal(err)
	}
	// "john@example.com" (16 chars) → 12 asterisks + ".com"
	// "555-1234" (8 chars) → 4 asterisks + "1234"
	want := "Contact me at ************.com or call me at ****1234 please."
	if got != want {
		t.Errorf("Mask:\n got: %q\nwant: %q", got, want)
	}
}

func TestHash(t *testing.T) {
	h := NewHash()
	got1, err := h.Replace(input, matches())
	if err != nil {
		t.Fatal(err)
	}

	// Deterministic: same input → same output.
	got2, err := h.Replace(input, matches())
	if err != nil {
		t.Fatal(err)
	}
	if got1 != got2 {
		t.Error("Hash should be deterministic")
	}

	// Original values should not appear in output.
	if containsAny(got1, "john@example.com", "555-1234") {
		t.Error("Hash output should not contain original values")
	}
}

func TestHashWithSalt(t *testing.T) {
	h1 := &Hash{Length: 16, Salt: "salt1"}
	h2 := &Hash{Length: 16, Salt: "salt2"}

	got1, _ := h1.Replace(input, matches())
	got2, _ := h2.Replace(input, matches())

	if got1 == got2 {
		t.Error("Different salts should produce different output")
	}
}

func TestCustom(t *testing.T) {
	c := NewCustom("upper", func(m model.Match) string {
		return "REDACTED"
	})
	got, err := c.Replace(input, matches())
	if err != nil {
		t.Fatal(err)
	}
	want := "Contact me at REDACTED or call me at REDACTED please."
	if got != want {
		t.Errorf("Custom:\n got: %q\nwant: %q", got, want)
	}

	if c.Name() != "upper" {
		t.Errorf("Name() = %q, want %q", c.Name(), "upper")
	}
}

func containsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if len(sub) > 0 {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
		}
	}
	return false
}
