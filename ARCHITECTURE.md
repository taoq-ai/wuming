# Architecture

Wuming uses a **hexagonal architecture** (also known as ports & adapters) to keep the core domain logic independent of concrete implementations. This makes the library easy to test, extend, and maintain.

## Layers

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

### Public API

**Package:** `wuming` (root)

The top-level package provides the user-facing API. It exposes the `Wuming` struct and functional options (`WithDetectors`, `WithReplacer`, `WithLocale`, etc.) that configure the internal engine. Users never interact with the engine directly.

Key methods:
- `Process(ctx, text)` — detect and replace, returning the full result
- `Detect(ctx, text)` — detect only, returning matches
- `Redact(ctx, text)` — detect and replace, returning just the redacted string

### Internal Engine

**Package:** `internal/engine`

The engine is the orchestrator. It receives a list of detectors and a replacer, runs detection (optionally in parallel), deduplicates overlapping matches, and applies the replacer. It implements the `port.Pipeline` interface.

The engine is internal — it cannot be imported by external packages. All access goes through the public API.

### Ports (Interfaces)

**Package:** `domain/port`

Ports define the contracts that adapters must satisfy. There are three interfaces:

- **Detector** — scans text and returns PII matches. Each detector declares its name, supported locales, and PII types.
- **Replacer** — takes text and matches, returns text with PII substituted.
- **Pipeline** — orchestrates detection and replacement (implemented by the engine).

The `Result` struct also lives here, containing the original text, redacted text, matches, and match count.

### Domain Model

**Package:** `domain/model`

Pure domain types with no dependencies:

- **PIIType** — enumeration of PII categories (Email, Phone, CreditCard, NationalID, TaxID, etc.)
- **Match** — a single detection result with type, value, byte offsets, confidence score, locale, and detector name
- **Severity** — sensitivity level (Low, Medium, High, Critical)

### Adapters

**Packages:** `adapter/detector/*`, `adapter/replacer/`

Concrete implementations of the port interfaces.

**Detectors** are organized by locale:
- `adapter/detector/common` — locale-independent patterns (email, credit card, IBAN, IP, URL, MAC)
- `adapter/detector/us` — United States (SSN, ITIN, EIN, phone, ZIP, passport, Medicare)
- `adapter/detector/nl` — Netherlands (BSN, phone, postal, KvK, ID document)
- `adapter/detector/eu` — European Union (VAT, passport MRZ)
- `adapter/detector/gb` — United Kingdom (NIN, NHS, UTR, phone, postcode)
- `adapter/detector/de` — Germany (Steuer-ID, Sozialversicherungsnummer, phone, PLZ, ID card)
- `adapter/detector/fr` — France (NIR, NIF, phone, postal, ID card)

**Replacers** provide different substitution strategies:
- `Redact` — type-based placeholder (e.g., `[EMAIL]`)
- `Mask` — character masking with preserved suffix
- `Hash` — deterministic SHA-256 hash
- `Custom` — user-defined function

## Data Flow

```
Input text
    │
    ▼
┌──────────┐
│  Engine   │──── Filters detectors by locale
│           │
│  ┌───────────────────────────────────┐
│  │  Detector 1 ──┐                   │
│  │  Detector 2 ──┼── (concurrent)    │
│  │  Detector N ──┘                   │
│  └───────────────────────────────────┘
│           │
│  Filter by confidence threshold
│  Filter by PII type
│  Deduplicate overlapping matches
│           │
│  ┌──────────────┐
│  │   Replacer   │── Substitutes matches in text
│  └──────────────┘
│           │
└──────────┘
    │
    ▼
Result { Original, Redacted, Matches, MatchCount }
```

1. The engine selects active detectors based on configured locale filters. Global detectors (those with no locale) always run.
2. Detectors run concurrently (bounded by the concurrency setting). Each returns a slice of `Match` values.
3. The engine filters matches by confidence threshold and PII type, then deduplicates overlapping matches (keeping the higher-confidence one).
4. The replacer walks the matches and substitutes each one in the original text.
5. The result is returned with the original text, redacted text, all matches, and a match count.

## Why Hexagonal Architecture?

- **Testability** — detectors and replacers can be tested in isolation with simple unit tests. The engine can be tested with mock implementations.
- **Extensibility** — adding a new locale means creating a new package under `adapter/detector/` without modifying any existing code.
- **Locale isolation** — each locale's detection logic is self-contained. A bug in the Dutch BSN detector cannot affect the US SSN detector.
- **Replacer flexibility** — swapping between redaction strategies is a single option change; no detection logic needs to change.
