// Package preset provides compliance-oriented configurations for PII detection.
// Each preset bundles the locales, PII types, and severity threshold appropriate
// for a specific data-protection regulation (e.g. GDPR, HIPAA, PCI-DSS).
package preset

import (
	"fmt"
	"sort"

	"github.com/taoq-ai/wuming/domain/model"
)

// Preset describes a compliance-oriented PII detection configuration.
type Preset struct {
	Name        string
	Description string
	Locales     []string        // which locale detectors to include
	PIITypes    []model.PIIType // which PII types to detect
	MinSeverity model.Severity  // minimum severity threshold
}

// allPIITypes returns every defined PIIType for presets that cover all personal data.
func allPIITypes() []model.PIIType {
	return []model.PIIType{
		model.Email,
		model.Phone,
		model.CreditCard,
		model.IBAN,
		model.IPAddress,
		model.URL,
		model.MACAddress,
		model.NationalID,
		model.TaxID,
		model.Passport,
		model.DriversLicense,
		model.HealthID,
		model.DateOfBirth,
		model.Name,
		model.Address,
		model.PostalCode,
		model.BankAccount,
		model.SocialMedia,
	}
}

// registry maps preset names (lowercase) to their definitions.
var registry = map[string]Preset{}

// register adds a preset to the internal registry.
func register(p Preset) {
	registry[p.Name] = p
}

// Get returns the preset with the given name, or an error if it does not exist.
func Get(name string) (Preset, error) {
	p, ok := registry[name]
	if !ok {
		return Preset{}, fmt.Errorf("preset: unknown preset %q", name)
	}
	return p, nil
}

// List returns a sorted list of all registered preset names.
func List() []string {
	names := make([]string, 0, len(registry))
	for k := range registry {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}
