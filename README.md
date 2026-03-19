# 无名 wuming

> **The Nameless** — A Go library for detecting and removing PII from text

*"The nameless (无名) is the origin of heaven and earth"* — Tao Te Ching, Chapter 1

[![Go Reference](https://pkg.go.dev/badge/github.com/taoq-ai/wuming.svg)](https://pkg.go.dev/github.com/taoq-ai/wuming)
[![Go Report Card](https://goreportcard.com/badge/github.com/taoq-ai/wuming)](https://goreportcard.com/report/github.com/taoq-ai/wuming)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/taoq-ai/wuming/actions/workflows/ci.yml/badge.svg)](https://github.com/taoq-ai/wuming/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-TBD-yellowgreen)](https://github.com/taoq-ai/wuming)

---

## Features

- **Zero-config detection** — a single function call catches all PII globally, no setup required
- **Global coverage** — 14 locales, 75+ detectors spanning every major region (Americas, Europe, Asia-Pacific)
- **11 compliance presets** — preconfigured profiles for GDPR, AI Act, HIPAA, PCI-DSS, LGPD, APPI, PIPL, PIPA, DPDP, PIPEDA, and Privacy Act
- **Locale-aware registry** — filter detectors with `WithLocale()` so only relevant patterns run
- **Hexagonal architecture** — clean separation between domain logic, ports, and adapters
- **Pluggable replacers** — redact, mask, hash, or bring your own replacement strategy
- **Concurrent detection** — run multiple detectors in parallel with configurable concurrency
- **Confidence scoring** — every match includes a confidence score; filter by threshold
- **Zero dependencies** — pure Go standard library, no external modules

## Quick Start

### Zero-Config (catches everything)

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/taoq-ai/wuming"
)

func main() {
    ctx := context.Background()

    // Zero config — one line, catches everything.
    redacted, err := wuming.Redact(ctx, "SSN 123-45-6789, email john@acme.com")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(redacted)
    // Output: SSN [NATIONAL_ID], email [EMAIL]
}
```

### Compliance Preset

```go
// Configure for GDPR compliance — only EU/EEA locales and PII types.
w := wuming.New(wuming.WithPreset("gdpr"))
result, err := w.Process(ctx, "Steuer-ID 12345678911, email jan@example.de")
```

### Locale-Specific

```go
// Only Dutch + common detectors.
w := wuming.New(wuming.WithLocale("nl"))
result, err := w.Process(ctx, "BSN 123456782, call 06-12345678")
```

## Installation

```sh
go get github.com/taoq-ai/wuming
```

Requires **Go 1.22** or later.

## Supported Locales

| Locale   | Package                          | Key Patterns                                              |
|----------|----------------------------------|-----------------------------------------------------------|
| common   | `adapter/detector/common`        | Email, Credit Card, IBAN, IP Address, URL, MAC Address    |
| us       | `adapter/detector/us`            | SSN, ITIN, EIN, Phone, ZIP, Passport, Medicare            |
| nl       | `adapter/detector/nl`            | BSN, Phone, Postal Code, KvK Number, ID Document          |
| eu       | `adapter/detector/eu`            | VAT Number, Passport MRZ                                  |
| gb       | `adapter/detector/gb`            | NIN, NHS Number, UTR, Phone, Postcode                     |
| de       | `adapter/detector/de`            | Steuer-ID, Sozialversicherungsnummer, Phone, PLZ, ID Card |
| fr       | `adapter/detector/fr`            | NIR, NIF, Phone, Postal Code, ID Card                     |
| br       | `adapter/detector/br`            | CPF, CNPJ, Phone, CEP, PIS/PASEP, CNH                    |
| jp       | `adapter/detector/jp`            | My Number, Corporate Number, Phone, Postal Code, Passport |
| in       | `adapter/detector/in`            | Aadhaar, PAN, Phone, PIN Code, Passport, GSTIN            |
| cn       | `adapter/detector/cn`            | Resident ID, Phone, Postal Code, Passport, USCC           |
| kr       | `adapter/detector/kr`            | RRN, Phone, Postal Code, Passport                         |
| au       | `adapter/detector/au`            | TFN, Medicare, ABN, Phone, Postcode                       |
| ca       | `adapter/detector/ca`            | SIN, Phone, Postal Code, Passport                         |

## Compliance Presets

Presets bundle the right locales and PII types for a specific regulation. Use them to get compliant detection without manual configuration.

| Preset        | Regulation                                               | Locales                        |
|---------------|----------------------------------------------------------|--------------------------------|
| `gdpr`        | EU General Data Protection Regulation                    | common, eu, nl, de, fr, gb     |
| `ai-act`      | EU AI Act (Articles 10, 15)                              | all 14 locales                 |
| `hipaa`       | US Health Insurance Portability and Accountability Act   | common, us                     |
| `pci-dss`     | Payment Card Industry Data Security Standard             | common                         |
| `lgpd`        | Brazil Lei Geral de Protecao de Dados                    | common, br                     |
| `appi`        | Japan Act on the Protection of Personal Information      | common, jp                     |
| `pipl`        | China Personal Information Protection Law                | common, cn                     |
| `pipa`        | South Korea Personal Information Protection Act          | common, kr                     |
| `dpdp`        | India Digital Personal Data Protection Act               | common, in                     |
| `pipeda`      | Canada Personal Information Protection and Electronic Documents Act | common, ca          |
| `privacy-act` | Australia Privacy Act                                    | common, au                     |

```go
w := wuming.New(wuming.WithPreset("gdpr"))
result, err := w.Process(ctx, text)
```

## Replacer Strategies

| Strategy | Type                    | Description                                        | Example Output          |
|----------|-------------------------|----------------------------------------------------|-------------------------|
| Redact   | `replacer.NewRedact()`  | Replace with type placeholder                      | `[EMAIL]`               |
| Mask     | `replacer.NewMask()`    | Mask characters, preserve last N                   | `****5678`              |
| Hash     | `replacer.NewHash()`    | Deterministic SHA-256 hash (truncated)             | `a1b2c3d4e5f67890`     |
| Custom   | `replacer.NewCustom(…)` | User-defined replacement function                  | *(anything you want)*   |

```go
// Use a mask replacer instead of the default redact
w := wuming.New(
    wuming.WithDetectors(common.NewEmailDetector()),
    wuming.WithReplacer(replacer.NewMask()),
)
```

## Architecture

Wuming follows a hexagonal (ports & adapters) architecture:

```
┌─────────────────────────────────────┐
│         Public API (wuming.go)      │
├─────────────────────────────────────┤
│         Internal Engine             │
├──────────┬──────────────────────────┤
│  Ports   │  domain/port/            │
│          │  - Detector interface    │
│          │  - Replacer interface    │
│          │  - Pipeline interface    │
├──────────┼──────────────────────────┤
│  Domain  │  domain/model/           │
│          │  - PIIType, Match        │
│          │  - Severity              │
├──────────┼──────────────────────────┤
│ Adapters │  adapter/detector/{loc}  │
│          │  adapter/replacer/       │
│          │  adapter/preset/         │
│          │  adapter/registry/       │
└──────────┴──────────────────────────┘
```

The public API delegates to an internal engine that orchestrates detectors (through the `port.Detector` interface) and replacers (through the `port.Replacer` interface). Each locale lives in its own adapter package, making it straightforward to add new locales without touching existing code.

For a deeper dive, see [ARCHITECTURE.md](ARCHITECTURE.md).

## Contributing

Contributions are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, coding standards, and how to add new detectors.

## Security

For responsible disclosure of security issues, see [SECURITY.md](SECURITY.md).

## License

[MIT](LICENSE)
