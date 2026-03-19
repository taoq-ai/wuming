# United States Detectors

US detectors cover PII patterns specific to the United States. They are activated when the `"us"` locale is configured.

## Social Security Number (SSN)

| Property | Value |
|----------|-------|
| Detector | `us/ssn` |
| PII Type | `NATIONAL_ID` |
| Confidence | 0.85 |
| Validation | Area/group/serial rules |

Detects US Social Security Numbers in the `XXX-XX-XXXX` format (with or without dashes).

**Validation rules:**

- Area number (first 3 digits) cannot be `000`, `666`, or `900-999`
- Group number (middle 2 digits) cannot be `00`
- Serial number (last 4 digits) cannot be `0000`

**Examples:** `123-45-6789`, `123456789`

!!! warning "Regulatory Context"
    SSNs are protected under various US federal and state privacy laws. Severity: **High**.

---

## Employer Identification Number (EIN)

| Property | Value |
|----------|-------|
| Detector | `us/ein` |
| PII Type | `TAX_ID` |
| Confidence | 0.85 |
| Validation | IRS campus prefix validation |

Detects US EINs in the `XX-XXXXXXX` format. The two-digit prefix is validated against the IRS campus assignment table.

**Examples:** `12-3456789`, `95-1234567`

---

## Individual Taxpayer Identification Number (ITIN)

| Property | Value |
|----------|-------|
| Detector | `us/itin` |
| PII Type | `TAX_ID` |
| Confidence | 0.80 |
| Validation | Group range validation |

Detects US ITINs (assigned to individuals who are not eligible for an SSN). Starts with `9`, followed by a group number in valid ITIN ranges (50-65, 70-88, 90-92, 94-99).

**Examples:** `9XX-70-XXXX`, `912-78-1234`

---

## Phone Number

| Property | Value |
|----------|-------|
| Detector | `us/phone` |
| PII Type | `PHONE` |
| Confidence | 0.80 |
| Validation | NANP format |

Detects US phone numbers in NANP (North American Numbering Plan) format. Supports optional `+1` prefix, parenthesized area codes, and various separators (space, dash, dot).

**Examples:** `(555) 123-4567`, `+1 555.123.4567`, `5551234567`

---

## Passport Number

| Property | Value |
|----------|-------|
| Detector | `us/passport` |
| PII Type | `PASSPORT` |
| Confidence | 0.60 |
| Validation | Regex pattern |

Detects US passport numbers (8-9 digits, optionally prefixed with a letter for newer formats).

**Examples:** `123456789`, `A12345678`

!!! note
    This detector has a lower confidence score because bare digit sequences can produce false positives. Use `WithConfidenceThreshold()` to filter as needed.

---

## ZIP Code

| Property | Value |
|----------|-------|
| Detector | `us/zip` |
| PII Type | `POSTAL_CODE` |
| Confidence | 0.60 |
| Validation | Regex pattern |

Detects US ZIP codes in 5-digit format, optionally with the ZIP+4 extension.

**Examples:** `90210`, `10001-1234`

---

## Medicare Beneficiary Identifier (MBI)

| Property | Value |
|----------|-------|
| Detector | `us/medicare` |
| PII Type | `HEALTH_ID` |
| Confidence | 0.80 |
| Validation | MBI character pattern |

Detects US Medicare Beneficiary Identifiers. The 11-character MBI follows a strict character pattern: `C A AN N AA N AN AN` where C=1-9, A=letter (excluding S,L,O,I,B,Z), N=digit, AN=alphanumeric (with same letter exclusions).

**Examples:** `1EG4-TE5-MK72` (without formatting: `1EG4TE5MK72`)

!!! warning "Regulatory Context"
    MBIs are protected under **HIPAA** (Health Insurance Portability and Accountability Act). Severity: **High**.
