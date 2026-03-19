# Detectors Overview

Detectors are the core building blocks of wuming. Each detector scans input text for a specific type of PII and returns a list of matches with position, confidence, and metadata.

## Supported Locales and Patterns

| Locale | Detector | PII Type | Confidence | Validation |
|--------|----------|----------|------------|------------|
| Global | Email | EMAIL | 0.95 | Regex |
| Global | Credit Card | CREDIT_CARD | 0.95 | Luhn algorithm |
| Global | IBAN | IBAN | 0.95 | ISO 13616 mod-97 |
| Global | IP Address | IP_ADDRESS | 0.90 | IPv4 octet range, IPv6 regex |
| Global | URL | URL | 0.90 | Regex |
| Global | MAC Address | MAC_ADDRESS | 0.90 | Regex |
| US | SSN | NATIONAL_ID | 0.85 | Area/group/serial rules |
| US | EIN | TAX_ID | 0.85 | IRS prefix validation |
| US | ITIN | TAX_ID | 0.80 | Group range validation |
| US | Phone | PHONE | 0.80 | NANP format |
| US | Passport | PASSPORT | 0.60 | Regex |
| US | ZIP Code | POSTAL_CODE | 0.60 | Regex (5 or 5+4) |
| US | Medicare | HEALTH_ID | 0.80 | MBI pattern |
| NL | BSN | NATIONAL_ID | 0.90 | 11-proof checksum |
| NL | Phone | PHONE | 0.85 | Dutch format |
| NL | Postal Code | POSTAL_CODE | 0.90 | Format + invalid combos |
| NL | KvK | TAX_ID | 0.60-0.90 | Context-boosted |
| NL | ID Document | NATIONAL_ID / PASSPORT | 0.70 | Pattern matching |
| EU | VAT Number | TAX_ID | 0.90 | Country prefix regex |
| EU | Passport MRZ | PASSPORT | 0.95 | ICAO 9303 TD3 format |
| GB | NIN | NATIONAL_ID | 0.85 | HMRC prefix rules |
| GB | NHS Number | HEALTH_ID | 0.90 | Mod-11 check digit |
| GB | UTR | TAX_ID | 0.55-0.85 | Context-boosted |
| GB | Phone | PHONE | 0.85 | UK format |
| GB | Postcode | POSTAL_CODE | 0.90 | UK format regex |
| DE | Steuer-ID | TAX_ID | 0.85 | ISO 7064 Mod 11,10 |
| DE | ID Card | NATIONAL_ID | 0.75 | Weighted checksum (7,3,1) |
| DE | Sozialversicherung | NATIONAL_ID | 0.75 | Date validation |
| DE | Phone | PHONE | 0.85 | German format |
| DE | PLZ | POSTAL_CODE | 0.60-0.85 | Range + context boost |
| FR | NIR | NATIONAL_ID | 0.90 | Mod-97 control key |
| FR | NIF | TAX_ID | 0.70 | Regex |
| FR | ID Card | NATIONAL_ID | 0.65 | Old (12-digit) + new (9-char) |
| FR | Phone | PHONE | 0.85 | French format |
| FR | Postal Code | POSTAL_CODE | 0.60 | Department prefix validation |

## The Detector Interface

Every detector implements the `port.Detector` interface:

```go
type Detector interface {
    Detect(ctx context.Context, text string) ([]model.Match, error)
    Name() string
    Locales() []string
    PIITypes() []model.PIIType
}
```

- **`Detect`** -- Scans the input text and returns all matches found.
- **`Name`** -- Returns a unique identifier (e.g., `"nl/bsn"`, `"common/email"`).
- **`Locales`** -- Returns which locales this detector supports. An empty slice means the detector is locale-independent (global).
- **`PIITypes`** -- Returns which PII types this detector can find.

## Confidence Scoring

Each match carries a confidence score between 0.0 and 1.0:

- **0.90 -- 0.95**: High confidence, validated by checksum or strong structural rules (e.g., credit card with Luhn, IBAN with mod-97, BSN with 11-proof)
- **0.80 -- 0.89**: Good confidence, structural validation but no checksum (e.g., SSN area rules, phone formats)
- **0.60 -- 0.79**: Moderate confidence, pattern-based with limited validation (e.g., ZIP codes, passport numbers, bare KvK numbers)
- **Below 0.60**: Low confidence, context-dependent patterns

Some detectors use **context boosting** -- if keywords like "BSN", "KvK", "PLZ", or "UTR" appear near a candidate match, the confidence is increased. This reduces false positives for ambiguous patterns like bare digit sequences.

Use `WithConfidenceThreshold()` to filter out matches below a desired confidence level.
