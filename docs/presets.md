# Compliance Presets

Presets bundle the right locales, PII types, and severity thresholds for a specific data-protection regulation. Instead of manually configuring which detectors and PII types to use, you can select a preset and wuming handles the rest.

## Usage

```go
w := wuming.New(wuming.WithPreset("gdpr"))
result, err := w.Process(ctx, text)
```

You can list all available presets programmatically:

```go
import "github.com/taoq-ai/wuming/adapter/preset"

names := preset.List()
// ["ai-act", "appi", "dpdp", "gdpr", "hipaa", "lgpd", "pci-dss", "pipa", "pipeda", "pipl", "privacy-act"]
```

## Available Presets

### GDPR

**EU General Data Protection Regulation** -- covers all personal data across EU/EEA locales.

- **Locales**: common, eu, nl, de, fr, gb
- **PII types**: all
- **Min severity**: Low

### AI Act

**EU AI Act** -- scrub training and validation data for high-risk AI systems (Articles 10, 15).

- **Locales**: all 14 (au, br, ca, cn, common, de, eu, fr, gb, in, jp, kr, nl, us)
- **PII types**: all
- **Min severity**: Low

```go
w := wuming.New(wuming.WithPreset("ai-act"))
result, err := w.Process(ctx, text)
```

### HIPAA

**US Health Insurance Portability and Accountability Act** -- protected health information.

- **Locales**: common, us
- **PII types**: National ID, Health ID, Phone, Email, Postal Code, Date of Birth
- **Min severity**: Medium

### PCI-DSS

**Payment Card Industry Data Security Standard** -- credit card data protection.

- **Locales**: common
- **PII types**: Credit Card
- **Min severity**: Critical

### LGPD

**Brazil Lei Geral de Protecao de Dados** -- covers all personal data.

- **Locales**: common, br
- **PII types**: all
- **Min severity**: Low

### APPI

**Japan Act on the Protection of Personal Information** -- covers all personal data.

- **Locales**: common, jp
- **PII types**: all
- **Min severity**: Low

### PIPL

**China Personal Information Protection Law** -- covers all personal data.

- **Locales**: common, cn
- **PII types**: all
- **Min severity**: Low

### PIPA

**South Korea Personal Information Protection Act** -- identity and contact data.

- **Locales**: common, kr
- **PII types**: National ID, Phone
- **Min severity**: Medium

### DPDP

**India Digital Personal Data Protection Act** -- covers all personal data.

- **Locales**: common, in
- **PII types**: all
- **Min severity**: Low

### PIPEDA

**Canada Personal Information Protection and Electronic Documents Act**.

- **Locales**: common, ca
- **PII types**: National ID, Health ID, Phone, Email
- **Min severity**: Medium

### Privacy Act

**Australia Privacy Act** -- tax, health, and contact data.

- **Locales**: common, au
- **PII types**: Tax ID, Health ID, Phone, Email
- **Min severity**: Medium

## How Presets Work

When you call `wuming.WithPreset("gdpr")`, wuming:

1. Looks up the preset definition in the internal registry
2. Loads detectors for each locale listed in the preset
3. Filters results to only the PII types the regulation requires
4. Common (global) detectors are always included automatically

Presets are defined in `adapter/preset/` with one file per regulation. You can inspect or extend them by adding new files that call `register()` in an `init()` function.
