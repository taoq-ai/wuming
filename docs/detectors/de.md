# Germany Detectors

German detectors cover PII patterns specific to Germany. They are activated when the `"de"` locale is configured.

## Steuerliche Identifikationsnummer (Steuer-ID)

| Property | Value |
|----------|-------|
| Detector | `de/steuer_id` |
| PII Type | `TAX_ID` |
| Confidence | 0.85 |
| Validation | Digit distribution + ISO 7064 Mod 11,10 check digit |

Detects German tax identification numbers (11 digits). Validation is two-fold:

1. **Digit distribution check** (first 10 digits): exactly one digit must appear twice, exactly one digit must not appear at all, and the remaining 8 digits each appear exactly once. The first digit cannot be 0.
2. **Check digit** (11th digit): validated using the **ISO 7064 Mod 11,10** algorithm.

**Examples:** `12345679812` (valid structure with correct check digit)

!!! warning "Regulatory Context"
    The Steuer-ID is regulated under the **Abgabenordnung (AO)** and is a lifelong identifier. Protected under **GDPR** and **BDSG** (Bundesdatenschutzgesetz). Severity: **High**.

---

## Personalausweisnummer (ID Card)

| Property | Value |
|----------|-------|
| Detector | `de/id_card` |
| PII Type | `NATIONAL_ID` |
| Confidence | 0.75 |
| Validation | Weighted checksum (7, 3, 1 cycling) |

Detects German national ID card numbers (10 alphanumeric characters). The first character is a letter from the set {L, M, N, P, R, T, V, W, X, Y}, followed by 8 alphanumeric characters and 1 check digit. The check digit is validated using a **weighted checksum** with cycling weights 7, 3, 1.

**Examples:** `L01X00T471`

---

## Sozialversicherungsnummer (Social Security Number)

| Property | Value |
|----------|-------|
| Detector | `de/sozialversicherung` |
| PII Type | `NATIONAL_ID` |
| Confidence | 0.75 |
| Validation | Date component validation |

Detects German social security numbers (12 characters: `AADDMMYYXNNN` -- area number, birth date, birth name initial, serial + check). Validates that the date components (day 01-31, month 01-12) and area number (non-zero) are within valid ranges.

**Examples:** `12 010290 A 123`, `12010290A123`

!!! info "Regulatory Context"
    The Sozialversicherungsnummer is issued by the **Deutsche Rentenversicherung** and is protected under **SGB IV** (Sozialgesetzbuch). Severity: **High**.

---

## Phone Number

| Property | Value |
|----------|-------|
| Detector | `de/phone` |
| PII Type | `PHONE` |
| Confidence | 0.85 |
| Validation | German phone format |

Detects German phone numbers in international and domestic formats:

- **International:** `+49 XXX XXXXXXX`, `0049 XXX XXXXXXX`
- **Mobile:** `01XX XXXXXXXX` (e.g., 0151, 0160, 0170, etc.)
- **Landline:** `0XXX XXXXXXX`

Supports separators: space, dash, dot, slash.

**Examples:** `+49 30 12345678`, `0151-12345678`, `030/1234567`

---

## Postleitzahl (PLZ)

| Property | Value |
|----------|-------|
| Detector | `de/plz` |
| PII Type | `POSTAL_CODE` |
| Confidence | 0.60 (bare) / 0.80-0.85 (with context) |
| Validation | Range check + context boost |

Detects German postal codes (5 digits, range 01001-99998). Because 5-digit sequences are common, the detector uses **context boosting** -- confidence is increased when keywords like "PLZ", "Postleitzahl", "Str.", "Stadt", "Ort", "Adresse", or "Anschrift" appear in the surrounding text.

**Examples:** `PLZ 10115`, `Adresse: Berliner Str. 1, 10115 Berlin`
