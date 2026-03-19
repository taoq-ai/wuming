package registry

import (
	"slices"
	"testing"
)

func TestAllDetectors(t *testing.T) {
	detectors := AllDetectors()
	if len(detectors) == 0 {
		t.Fatal("AllDetectors() returned 0 detectors, expected > 0")
	}
}

func TestAllDetectorsNoDuplicates(t *testing.T) {
	detectors := AllDetectors()
	seen := make(map[string]bool)
	for _, d := range detectors {
		if seen[d.Name()] {
			t.Errorf("duplicate detector: %s", d.Name())
		}
		seen[d.Name()] = true
	}
}

func TestDetectorsForLocale(t *testing.T) {
	nlDetectors := DetectorsForLocale("nl")
	if len(nlDetectors) == 0 {
		t.Fatal("DetectorsForLocale(\"nl\") returned 0 detectors")
	}

	// Should include common detectors.
	commonDetectors := DetectorsForLocale("common")
	if len(nlDetectors) <= len(commonDetectors) {
		t.Errorf("nl detectors (%d) should be more than common-only (%d)",
			len(nlDetectors), len(commonDetectors))
	}

	// Verify common detectors are present in the nl set.
	nlNames := make(map[string]bool)
	for _, d := range nlDetectors {
		nlNames[d.Name()] = true
	}
	for _, d := range commonDetectors {
		if !nlNames[d.Name()] {
			t.Errorf("common detector %q missing from nl locale", d.Name())
		}
	}
}

func TestDetectorsForLocaleUnknown(t *testing.T) {
	detectors := DetectorsForLocale("xx")
	commonDetectors := DetectorsForLocale("common")
	if len(detectors) != len(commonDetectors) {
		t.Errorf("unknown locale should return only common detectors: got %d, want %d",
			len(detectors), len(commonDetectors))
	}
}

func TestLocales(t *testing.T) {
	locales := Locales()
	expected := []string{"au", "cn", "common", "de", "eu", "fr", "gb", "jp", "kr", "nl", "us"}
	if len(locales) != len(expected) {
		t.Fatalf("Locales() = %v, want %v", locales, expected)
	}
	for _, e := range expected {
		if !slices.Contains(locales, e) {
			t.Errorf("Locales() missing %q", e)
		}
	}
}
