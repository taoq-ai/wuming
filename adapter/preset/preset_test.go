package preset

import (
	"testing"

	"github.com/taoq-ai/wuming/domain/model"
)

func TestListReturnsAllPresets(t *testing.T) {
	names := List()
	want := []string{
		"ai-act", "appi", "dpdp", "gdpr", "hipaa", "lgpd",
		"pci-dss", "pipa", "pipeda", "pipl", "privacy-act",
	}
	if len(names) != len(want) {
		t.Fatalf("List() returned %d presets, want %d: %v", len(names), len(want), names)
	}
	for i, name := range names {
		if name != want[i] {
			t.Errorf("List()[%d] = %q, want %q", i, name, want[i])
		}
	}
}

func TestGetKnownPreset(t *testing.T) {
	for _, name := range List() {
		p, err := Get(name)
		if err != nil {
			t.Errorf("Get(%q) returned error: %v", name, err)
			continue
		}
		if p.Name != name {
			t.Errorf("Get(%q).Name = %q", name, p.Name)
		}
		if len(p.Locales) == 0 {
			t.Errorf("Get(%q).Locales is empty", name)
		}
		if len(p.PIITypes) == 0 {
			t.Errorf("Get(%q).PIITypes is empty", name)
		}
		if p.MinSeverity < model.Low || p.MinSeverity > model.Critical {
			t.Errorf("Get(%q).MinSeverity = %d, out of range", name, p.MinSeverity)
		}
	}
}

func TestGetUnknownPresetReturnsError(t *testing.T) {
	_, err := Get("nonexistent")
	if err == nil {
		t.Error("Get(\"nonexistent\") should return an error")
	}
}

func TestAIActPreset(t *testing.T) {
	p, _ := Get("ai-act")
	assertLocales(t, p, []string{
		"au", "br", "ca", "cn", "common", "de", "eu",
		"fr", "gb", "in", "jp", "kr", "nl", "us",
	})
	assertSeverity(t, p, model.Low)
	// AI Act covers all PII types across all locales.
	if len(p.PIITypes) != len(allPIITypes()) {
		t.Errorf("ai-act: got %d PII types, want %d", len(p.PIITypes), len(allPIITypes()))
	}
}

func TestGDPRPreset(t *testing.T) {
	p, _ := Get("gdpr")
	assertLocales(t, p, []string{"common", "eu", "nl", "de", "fr", "gb"})
	assertSeverity(t, p, model.Low)
	// GDPR covers all PII types.
	if len(p.PIITypes) != len(allPIITypes()) {
		t.Errorf("gdpr: got %d PII types, want %d", len(p.PIITypes), len(allPIITypes()))
	}
}

func TestHIPAAPreset(t *testing.T) {
	p, _ := Get("hipaa")
	assertLocales(t, p, []string{"common", "us"})
	assertSeverity(t, p, model.Medium)
	assertContainsPIIType(t, p, model.NationalID)
	assertContainsPIIType(t, p, model.HealthID)
	assertContainsPIIType(t, p, model.DateOfBirth)
}

func TestPCIDSSPreset(t *testing.T) {
	p, _ := Get("pci-dss")
	assertLocales(t, p, []string{"common"})
	assertSeverity(t, p, model.Critical)
	if len(p.PIITypes) != 1 || p.PIITypes[0] != model.CreditCard {
		t.Errorf("pci-dss: PIITypes = %v, want [CreditCard]", p.PIITypes)
	}
}

func TestLGPDPreset(t *testing.T) {
	p, _ := Get("lgpd")
	assertLocales(t, p, []string{"common", "br"})
	assertSeverity(t, p, model.Low)
}

func TestAPPIPreset(t *testing.T) {
	p, _ := Get("appi")
	assertLocales(t, p, []string{"common", "jp"})
	assertSeverity(t, p, model.Low)
}

func TestPIPLPreset(t *testing.T) {
	p, _ := Get("pipl")
	assertLocales(t, p, []string{"common", "cn"})
	assertSeverity(t, p, model.Low)
}

func TestPIPAPreset(t *testing.T) {
	p, _ := Get("pipa")
	assertLocales(t, p, []string{"common", "kr"})
	assertSeverity(t, p, model.Medium)
	assertContainsPIIType(t, p, model.NationalID)
	assertContainsPIIType(t, p, model.Phone)
}

func TestDPDPPreset(t *testing.T) {
	p, _ := Get("dpdp")
	assertLocales(t, p, []string{"common", "in"})
	assertSeverity(t, p, model.Low)
}

func TestPIPEDAPreset(t *testing.T) {
	p, _ := Get("pipeda")
	assertLocales(t, p, []string{"common", "ca"})
	assertSeverity(t, p, model.Medium)
	assertContainsPIIType(t, p, model.NationalID)
	assertContainsPIIType(t, p, model.HealthID)
}

func TestPrivacyActPreset(t *testing.T) {
	p, _ := Get("privacy-act")
	assertLocales(t, p, []string{"common", "au"})
	assertSeverity(t, p, model.Medium)
	assertContainsPIIType(t, p, model.TaxID)
	assertContainsPIIType(t, p, model.HealthID)
}

// --- helpers ---

func assertLocales(t *testing.T, p Preset, want []string) {
	t.Helper()
	if len(p.Locales) != len(want) {
		t.Errorf("%s: got %d locales %v, want %v", p.Name, len(p.Locales), p.Locales, want)
		return
	}
	for i, l := range p.Locales {
		if l != want[i] {
			t.Errorf("%s: locale[%d] = %q, want %q", p.Name, i, l, want[i])
		}
	}
}

func assertSeverity(t *testing.T, p Preset, want model.Severity) {
	t.Helper()
	if p.MinSeverity != want {
		t.Errorf("%s: MinSeverity = %v, want %v", p.Name, p.MinSeverity, want)
	}
}

func assertContainsPIIType(t *testing.T, p Preset, want model.PIIType) {
	t.Helper()
	for _, pt := range p.PIITypes {
		if pt == want {
			return
		}
	}
	t.Errorf("%s: PIITypes %v does not contain %v", p.Name, p.PIITypes, want)
}
