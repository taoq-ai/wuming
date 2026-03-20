# Wuming (无名) — AI Agent Integration Guide

## Library Overview

Wuming ("The Nameless") is a Go library for detecting and removing Personally Identifiable Information (PII) from text. It supports global PII standards across 14 locales and provides pluggable detection and replacement strategies.

**Module path:** `github.com/taoq-ai/wuming`
**Go version:** 1.22+
**Architecture:** Hexagonal (ports & adapters)

```
domain/model/   — Core types: Match, PIIType, Severity
domain/port/    — Interfaces: Detector, Replacer, Pipeline, Result
adapter/        — Implementations: detectors, replacers, presets, registry, structured
internal/engine — Orchestrator wiring detectors + replacers (implements Pipeline)
wuming.go       — Public API facade
```

---

## Quick Start

### Zero-config (all detectors, all locales)

```go
ctx := context.Background()

// Redact all PII
redacted, err := wuming.Redact(ctx, "Call me at 06-12345678 or email john@example.com")
// "Call me at [PHONE] or email [EMAIL]"

// Detect only (no modification)
matches, err := wuming.Detect(ctx, "SSN: 123-45-6789")

// Full result (original, redacted, matches)
result, err := wuming.Process(ctx, text)
// result.Original, result.Redacted, result.Matches, result.MatchCount
```

### Configured instance

```go
w, err := wuming.New(
    wuming.WithLocale("nl"),
    wuming.WithReplacer(replacer.NewMask()),
    wuming.WithConfidenceThreshold(0.8),
    wuming.WithAllowlist("example.com"),
    wuming.WithDenylist(model.Custom, "ProjectX"),
)
if err != nil {
    log.Fatal(err)
}

result, err := w.Process(ctx, text)
```

### Compliance preset

```go
w, err := wuming.New(wuming.WithPreset("gdpr"))
```

---

## Public API Reference

### Top-level package functions (use default instance with all detectors)

```go
func Redact(ctx context.Context, text string) (string, error)
func Detect(ctx context.Context, text string) ([]model.Match, error)
func Process(ctx context.Context, text string) (*port.Result, error)
func RedactJSON(ctx context.Context, data []byte) (*structured.Result, error)
func DetectJSON(ctx context.Context, data []byte) ([]structured.FieldMatch, error)
func RedactCSV(ctx context.Context, r io.Reader) (*structured.Result, error)
func DetectCSV(ctx context.Context, r io.Reader) ([]structured.FieldMatch, error)
```

### Constructor and instance methods

```go
func New(opts ...Option) (*Wuming, error)

func (w *Wuming) Redact(ctx context.Context, text string) (string, error)
func (w *Wuming) Detect(ctx context.Context, text string) ([]model.Match, error)
func (w *Wuming) Process(ctx context.Context, text string) (*port.Result, error)
func (w *Wuming) RedactJSON(ctx context.Context, data []byte) (*structured.Result, error)
func (w *Wuming) DetectJSON(ctx context.Context, data []byte) ([]structured.FieldMatch, error)
func (w *Wuming) RedactCSV(ctx context.Context, r io.Reader) (*structured.Result, error)
func (w *Wuming) DetectCSV(ctx context.Context, r io.Reader) ([]structured.FieldMatch, error)
```

### Options

| Option | Signature | Description |
|--------|-----------|-------------|
| `WithLocale` | `WithLocale(locale string) Option` | Filter detectors to a specific locale. Global detectors always run. |
| `WithDetectors` | `WithDetectors(d ...port.Detector) Option` | Add explicit detectors (bypasses registry). |
| `WithReplacer` | `WithReplacer(r port.Replacer) Option` | Set replacement strategy. Default: `Redact` (`[EMAIL]`). |
| `WithPIITypes` | `WithPIITypes(types ...model.PIIType) Option` | Filter results to specific PII types only. |
| `WithConcurrency` | `WithConcurrency(n int) Option` | Max parallel detectors. Default: unlimited. |
| `WithConfidenceThreshold` | `WithConfidenceThreshold(f float64) Option` | Drop matches below this score (0.0-1.0). |
| `WithAllowlist` | `WithAllowlist(values ...string) Option` | Values that should never be flagged. Case-insensitive. |
| `WithDenylist` | `WithDenylist(piiType model.PIIType, values ...string) Option` | Values that should always be flagged as the given type. |
| `WithPreset` | `WithPreset(name string) Option` | Load a compliance preset (e.g. "gdpr", "hipaa"). |

---

## Structured Data

### JSON scanning

```go
// Scans each string value in the JSON document
result, err := w.RedactJSON(ctx, jsonBytes)
// result.Data     — redacted JSON as []byte
// result.Matches  — []structured.FieldMatch with Path like "user.email", "contacts[0].phone"

matches, err := w.DetectJSON(ctx, jsonBytes) // detect only, no modification
```

### CSV scanning

```go
// First row is treated as column headers
result, err := w.RedactCSV(ctx, reader)
// result.Data     — redacted CSV as []byte
// result.Matches  — []structured.FieldMatch with Path like "R2:email", "R3:phone"

matches, err := w.DetectCSV(ctx, reader) // detect only
```

### Structured types

```go
// structured.FieldMatch extends model.Match with:
type FieldMatch struct {
    model.Match
    Path string  // e.g. "user.email" (JSON) or "R2:C3" (CSV)
}

type Result struct {
    Data       []byte
    Matches    []FieldMatch
    MatchCount int
}
```

---

## Architecture Guide

```
                    wuming.go (public facade)
                         |
                   internal/engine
                   /             \
          domain/port             domain/model
         (interfaces)              (types)
              |
         adapter/
        /    |    \       \
  detector/ replacer/ preset/ structured/
      |
  us/ nl/ common/ de/ fr/ gb/ ...
```

**Flow:** `wuming.New()` -> creates `engine.Engine` -> `Engine.Process()` runs active detectors concurrently -> merges/deduplicates matches -> filters by confidence/PIIType/allowlist -> injects denylist -> applies replacer -> returns `port.Result`.

**Key interfaces (domain/port/):**

```go
type Detector interface {
    Detect(ctx context.Context, text string) ([]model.Match, error)
    Name() string                    // e.g. "us/ssn", "common/email"
    Locales() []string               // empty = global/locale-independent
    PIITypes() []model.PIIType
}

type Replacer interface {
    Replace(text string, matches []model.Match) (string, error)
    Name() string
}

type Pipeline interface {
    Process(ctx context.Context, text string) (*Result, error)
}
```

**Key types (domain/model/):**

```go
type Match struct {
    Type       PIIType
    Value      string
    Start      int      // byte offset
    End        int      // byte offset (exclusive)
    Confidence float64  // 0.0 to 1.0
    Locale     string   // e.g. "nl", "us", "" for global
    Detector   string   // e.g. "us/ssn"
}

type PIIType int  // Email, Phone, CreditCard, IBAN, IPAddress, URL, MACAddress,
                  // NationalID, TaxID, Passport, DriversLicense, HealthID,
                  // DateOfBirth, Name, Address, PostalCode, BankAccount,
                  // SocialMedia, Custom

type Severity int // Low, Medium, High, Critical
```

---

## Adding a New Locale/Detector

### Step 1: Create the locale package

```
adapter/detector/xx/
    helpers.go   — shared constants and helper functions
    all.go       — All() function returning all detectors
    foo.go       — individual detector
    xx_test.go   — tests
```

### Step 2: Implement helpers.go

```go
package xx

import (
    "regexp"
    "github.com/taoq-ai/wuming/domain/model"
)

const locale = "xx"

func findAll(re *regexp.Regexp, text string, piiType model.PIIType, confidence float64, detector string) []model.Match {
    locs := re.FindAllStringIndex(text, -1)
    if len(locs) == 0 {
        return nil
    }
    matches := make([]model.Match, 0, len(locs))
    for _, loc := range locs {
        matches = append(matches, model.Match{
            Type:       piiType,
            Value:      text[loc[0]:loc[1]],
            Start:      loc[0],
            End:        loc[1],
            Confidence: confidence,
            Locale:     locale,
            Detector:   detector,
        })
    }
    return matches
}
```

### Step 3: Implement a detector

```go
package xx

import (
    "context"
    "regexp"
    "github.com/taoq-ai/wuming/domain/model"
)

var fooRe = regexp.MustCompile(`\bPATTERN\b`)

type FooDetector struct{}

func NewFooDetector() *FooDetector { return &FooDetector{} }

func (d *FooDetector) Name() string              { return "xx/foo" }
func (d *FooDetector) Locales() []string         { return []string{locale} }
func (d *FooDetector) PIITypes() []model.PIIType { return []model.PIIType{model.NationalID} }

func (d *FooDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
    return findAll(fooRe, text, model.NationalID, 0.9, d.Name()), nil
}
```

### Step 4: Create all.go

```go
package xx

import "github.com/taoq-ai/wuming/domain/port"

func All() []port.Detector {
    return []port.Detector{
        NewFooDetector(),
    }
}
```

### Step 5: Register in adapter/registry/registry.go

Add the import and map entry:

```go
import "github.com/taoq-ai/wuming/adapter/detector/xx"

// In localeProviders:
"xx": xx.All,
```

### Step 6: Write tests

Follow the pattern in existing `xx_test.go` files. Test both positive matches and non-matches.

---

## Adding a New Preset

Create a file in `adapter/preset/` (e.g. `my_regulation.go`):

```go
package preset

import "github.com/taoq-ai/wuming/domain/model"

func init() {
    register(Preset{
        Name:        "my-regulation",
        Description: "Description of the regulation",
        Locales:     []string{"common", "xx"},
        PIITypes: []model.PIIType{
            model.NationalID,
            model.Phone,
            model.Email,
        },
        MinSeverity: model.Medium,
    })
}
```

The preset is auto-registered via `init()` and becomes available through `wuming.WithPreset("my-regulation")`.

---

## Available Locales

| Locale | Package | Detectors |
|--------|---------|-----------|
| `common` | `adapter/detector/common` | Email, CreditCard, IP, URL, IBAN, MAC |
| `au` | `adapter/detector/au` | TFN, Medicare, ABN, Phone, Postcode |
| `br` | `adapter/detector/br` | CPF, CNPJ, Phone, CEP, PIS, CNH |
| `ca` | `adapter/detector/ca` | SIN, Phone, PostalCode, Passport |
| `cn` | `adapter/detector/cn` | ResidentID, Phone, Postal, Passport, USCC |
| `de` | `adapter/detector/de` | SteuerID, Sozialversicherung, Phone, PLZ, IDCard |
| `eu` | `adapter/detector/eu` | PassportMRZ, VAT |
| `fr` | `adapter/detector/fr` | NIR, NIF, Phone, Postal, IDCard |
| `gb` | `adapter/detector/gb` | NIN, NHS, UTR, Phone, Postcode |
| `in` | `adapter/detector/in` | Aadhaar, PAN, Phone, PINCode, Passport, GSTIN |
| `jp` | `adapter/detector/jp` | MyNumber, CorporateNumber, Phone, Postal, Passport |
| `kr` | `adapter/detector/kr` | RRN, Phone, Postal, Passport |
| `nl` | `adapter/detector/nl` | BSN, Phone, Postal, KvK, IDDocument |
| `us` | `adapter/detector/us` | SSN, ITIN, EIN, Phone, ZIP, Passport, Medicare |

---

## Available Presets

| Preset | Description | Locales |
|--------|-------------|---------|
| `ai-act` | EU AI Act — scrub training data for high-risk AI systems | All 14 locales |
| `appi` | Japan Act on Protection of Personal Information | common, jp |
| `dpdp` | India Digital Personal Data Protection Act | common, in |
| `gdpr` | EU General Data Protection Regulation | common, eu, nl, de, fr, gb |
| `hipaa` | US Health Insurance Portability and Accountability Act | common, us |
| `lgpd` | Brazil Lei Geral de Protecao de Dados | common, br |
| `pci-dss` | Payment Card Industry Data Security Standard | common |
| `pipa` | South Korea Personal Information Protection Act | common, kr |
| `pipeda` | Canada Personal Information Protection and Electronic Documents Act | common, ca |
| `pipl` | China Personal Information Protection Law | common, cn |
| `privacy-act` | Australia Privacy Act | common, au |

---

## Replacer Strategies

All replacers implement `port.Replacer`. Import from `github.com/taoq-ai/wuming/adapter/replacer`.

### Redact (default)

Replaces each match with a type-based placeholder.

```go
r := replacer.NewRedact()
// "john@example.com" -> "[EMAIL]"
// Custom format:
r = &replacer.Redact{Format: "<%s>"}
// "john@example.com" -> "<EMAIL>"
```

### Mask

Replaces characters with a mask character, preserving trailing characters.

```go
r := replacer.NewMask()
// "john@example.com" -> "************.com"
// Custom settings:
r = &replacer.Mask{Char: '#', Preserve: 2}
```

### Hash

Replaces with a deterministic SHA-256 hash (truncated).

```go
r := replacer.NewHash()
// Same input always produces same hash
// Custom settings:
r = &replacer.Hash{Length: 16, Salt: "my-salt"}
```

### Custom

User-defined replacement function.

```go
r := replacer.NewCustom("my-replacer", func(m model.Match) string {
    return fmt.Sprintf("<%s:%d>", m.Type, m.Start)
})
```

### Consistent

Wraps any replacer to ensure the same PII value always maps to the same replacement. With `Redact`, produces numbered placeholders.

```go
r := replacer.NewConsistent(replacer.NewRedact())
// First "john@example.com" -> "[EMAIL_1]"
// Second "john@example.com" -> "[EMAIL_1]" (same)
// "jane@example.com" -> "[EMAIL_2]"

// Reset mapping between unrelated texts:
r.Reset()
```

---

## Development Conventions

### Commit messages

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add XX locale detector
fix: correct SSN validation for area 900+
docs: update SKILL.md with new locale
test: add benchmarks for JSON scanner
chore: bump Go version to 1.22
```

### Branch workflow

- Never push directly to `main`
- Create feature branches: `feat/XX-description`, `fix/XX-description`, `docs/XX-description`
- Open a pull request, wait for CI, then merge

### Testing

- Every detector must have tests covering positive matches and non-matches
- Run tests: `go test ./...`
- Run tests with race detector: `go test -race ./...`
- Follow existing test patterns in `*_test.go` files

### Project structure rules

- Detectors go in `adapter/detector/<locale>/`
- Every locale package needs `helpers.go`, `all.go`, individual detector files, and `<locale>_test.go`
- Presets go in `adapter/preset/<name>.go` using `init()` + `register()`
- Replacers go in `adapter/replacer/`
- The `domain/` layer has no external dependencies
- The `internal/` layer is not importable by consumers
