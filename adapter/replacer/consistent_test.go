package replacer

import (
	"strings"
	"testing"

	"github.com/taoq-ai/wuming/domain/model"
)

func TestConsistentRedactSameValueSameReplacement(t *testing.T) {
	// Same email appears twice -> both should get [EMAIL_1].
	text := "Email john@example.com and again john@example.com"
	ms := []model.Match{
		{Type: model.Email, Value: "john@example.com", Start: 6, End: 22, Confidence: 1.0},
		{Type: model.Email, Value: "john@example.com", Start: 33, End: 49, Confidence: 1.0},
	}

	c := NewConsistent(NewRedact())
	got, err := c.Replace(text, ms)
	if err != nil {
		t.Fatal(err)
	}
	want := "Email [EMAIL_1] and again [EMAIL_1]"
	if got != want {
		t.Errorf("ConsistentRedact same value:\n got: %q\nwant: %q", got, want)
	}
}

func TestConsistentRedactDifferentValuesDifferentNumbers(t *testing.T) {
	// Two different emails -> [EMAIL_1] and [EMAIL_2].
	text := "Email john@example.com and jane@example.com"
	ms := []model.Match{
		{Type: model.Email, Value: "john@example.com", Start: 6, End: 22, Confidence: 1.0},
		{Type: model.Email, Value: "jane@example.com", Start: 27, End: 43, Confidence: 1.0},
	}

	c := NewConsistent(NewRedact())
	got, err := c.Replace(text, ms)
	if err != nil {
		t.Fatal(err)
	}
	want := "Email [EMAIL_1] and [EMAIL_2]"
	if got != want {
		t.Errorf("ConsistentRedact different values:\n got: %q\nwant: %q", got, want)
	}
}

func TestConsistentRedactMixedTypes(t *testing.T) {
	// Email and phone: counters are per-type.
	text := "Email john@example.com call 555-1234 email john@example.com"
	ms := []model.Match{
		{Type: model.Email, Value: "john@example.com", Start: 6, End: 22, Confidence: 1.0},
		{Type: model.Phone, Value: "555-1234", Start: 28, End: 36, Confidence: 0.9},
		{Type: model.Email, Value: "john@example.com", Start: 43, End: 59, Confidence: 1.0},
	}

	c := NewConsistent(NewRedact())
	got, err := c.Replace(text, ms)
	if err != nil {
		t.Fatal(err)
	}
	want := "Email [EMAIL_1] call [PHONE_1] email [EMAIL_1]"
	if got != want {
		t.Errorf("ConsistentRedact mixed types:\n got: %q\nwant: %q", got, want)
	}
}

func TestConsistentRedactEmpty(t *testing.T) {
	c := NewConsistent(NewRedact())
	got, err := c.Replace("no pii here", nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "no pii here" {
		t.Errorf("expected unchanged text for empty matches, got: %q", got)
	}
}

func TestConsistentName(t *testing.T) {
	c := NewConsistent(NewRedact())
	want := "consistent(redact)"
	if c.Name() != want {
		t.Errorf("Name() = %q, want %q", c.Name(), want)
	}
}

func TestConsistentWithHash(t *testing.T) {
	// Same value with hash replacer should produce same hash output.
	text := "Email john@example.com and again john@example.com"
	ms := []model.Match{
		{Type: model.Email, Value: "john@example.com", Start: 6, End: 22, Confidence: 1.0},
		{Type: model.Email, Value: "john@example.com", Start: 33, End: 49, Confidence: 1.0},
	}

	c := NewConsistent(NewHash())
	got, err := c.Replace(text, ms)
	if err != nil {
		t.Fatal(err)
	}

	// The two replacements should be identical.
	parts := strings.SplitN(got, " and again ", 2)
	if len(parts) != 2 {
		t.Fatalf("unexpected output format: %q", got)
	}
	first := strings.TrimPrefix(parts[0], "Email ")
	second := parts[1]
	if first != second {
		t.Errorf("hash replacements differ: %q vs %q", first, second)
	}

	// Original should not appear.
	if strings.Contains(got, "john@example.com") {
		t.Error("original value should not appear in output")
	}
}

func TestConsistentWithMask(t *testing.T) {
	text := "Call 555-1234 and again 555-1234"
	ms := []model.Match{
		{Type: model.Phone, Value: "555-1234", Start: 5, End: 13, Confidence: 0.9},
		{Type: model.Phone, Value: "555-1234", Start: 24, End: 32, Confidence: 0.9},
	}

	c := NewConsistent(NewMask())
	got, err := c.Replace(text, ms)
	if err != nil {
		t.Fatal(err)
	}
	want := "Call ****1234 and again ****1234"
	if got != want {
		t.Errorf("ConsistentMask:\n got: %q\nwant: %q", got, want)
	}
}

func TestConsistentReset(t *testing.T) {
	c := NewConsistent(NewRedact())

	text1 := "Email john@example.com"
	ms1 := []model.Match{
		{Type: model.Email, Value: "john@example.com", Start: 6, End: 22, Confidence: 1.0},
	}
	got1, err := c.Replace(text1, ms1)
	if err != nil {
		t.Fatal(err)
	}
	if got1 != "Email [EMAIL_1]" {
		t.Errorf("before reset: got %q, want %q", got1, "Email [EMAIL_1]")
	}

	// Without reset, counter continues.
	text2 := "Email jane@example.com"
	ms2 := []model.Match{
		{Type: model.Email, Value: "jane@example.com", Start: 6, End: 22, Confidence: 1.0},
	}
	got2, err := c.Replace(text2, ms2)
	if err != nil {
		t.Fatal(err)
	}
	if got2 != "Email [EMAIL_2]" {
		t.Errorf("without reset: got %q, want %q", got2, "Email [EMAIL_2]")
	}

	// After reset, counter restarts.
	c.Reset()
	got3, err := c.Replace(text2, ms2)
	if err != nil {
		t.Fatal(err)
	}
	if got3 != "Email [EMAIL_1]" {
		t.Errorf("after reset: got %q, want %q", got3, "Email [EMAIL_1]")
	}
}

func TestConsistentThreeOccurrences(t *testing.T) {
	// Three occurrences of the same value -> all get [EMAIL_1].
	text := "a john@example.com b john@example.com c john@example.com"
	ms := []model.Match{
		{Type: model.Email, Value: "john@example.com", Start: 2, End: 18, Confidence: 1.0},
		{Type: model.Email, Value: "john@example.com", Start: 21, End: 37, Confidence: 1.0},
		{Type: model.Email, Value: "john@example.com", Start: 40, End: 56, Confidence: 1.0},
	}

	c := NewConsistent(NewRedact())
	got, err := c.Replace(text, ms)
	if err != nil {
		t.Fatal(err)
	}
	want := "a [EMAIL_1] b [EMAIL_1] c [EMAIL_1]"
	if got != want {
		t.Errorf("three occurrences:\n got: %q\nwant: %q", got, want)
	}
}
