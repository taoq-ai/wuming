# 无名 wuming

> **The Nameless** — A Go library for detecting and removing PII from text

*"The nameless (无名) is the origin of heaven and earth"* — Tao Te Ching, Chapter 1

[![Go Reference](https://pkg.go.dev/badge/github.com/taoq-ai/wuming.svg)](https://pkg.go.dev/github.com/taoq-ai/wuming)
[![Go Report Card](https://goreportcard.com/badge/github.com/taoq-ai/wuming)](https://goreportcard.com/report/github.com/taoq-ai/wuming)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/taoq-ai/wuming/actions/workflows/ci.yml/badge.svg)](https://github.com/taoq-ai/wuming/actions/workflows/ci.yml)

---

## Features

- **Hexagonal architecture** — clean separation between domain logic, ports, and adapters
- **Global coverage** — detectors for common patterns (email, credit card, IBAN) plus locale-specific PII (US SSN, Dutch BSN, UK NIN, German Steuer-ID, French NIR, EU VAT)
- **Pluggable replacers** — redact, mask, hash, or bring your own replacement strategy
- **Concurrent detection** — run multiple detectors in parallel with configurable concurrency
- **Confidence scoring** — every match includes a confidence score; filter by threshold
- **Zero dependencies** — pure Go standard library, no external modules

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/taoq-ai/wuming"
    "github.com/taoq-ai/wuming/adapter/detector/common"
    "github.com/taoq-ai/wuming/adapter/detector/nl"
)

func main() {
    w := wuming.New(
        wuming.WithDetectors(
            common.NewEmailDetector(),
            nl.NewPhoneDetector(),
        ),
    )

    result, err := w.Process(context.Background(), "Call me at 06-12345678 or email john@example.com")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(result.Redacted)
    // Output: Call me at [PHONE] or email [EMAIL]
}
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
