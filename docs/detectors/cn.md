# China Detectors

China detectors cover PII patterns specific to the Chinese context. They are activated when the `"cn"` locale is configured.

## Resident Identity Card Number

| Property | Value |
|----------|-------|
| Detector | `cn/resident_id` |
| PII Type | `NATIONAL_ID` |
| Confidence | 0.90 |
| Validation | Check digit (GB 11643) |

Detects Chinese resident identity card numbers. The 18-digit number encodes region, birth date, sequence, and a check digit calculated per the GB 11643 standard.

**Examples:** `110101199003077593`

!!! warning "Regulatory Context"
    Resident ID numbers are regulated under **PIPL** (Personal Information Protection Law). Severity: **Critical**.

---

## Phone Number

| Property | Value |
|----------|-------|
| Detector | `cn/phone` |
| PII Type | `PHONE` |
| Confidence | 0.85 |
| Validation | Chinese format |

Detects Chinese phone numbers in mobile and landline formats.

**Examples:** `13812345678`, `+86 138 1234 5678`

---

## Postal Code

| Property | Value |
|----------|-------|
| Detector | `cn/postal` |
| PII Type | `POSTAL_CODE` |
| Confidence | 0.60 |
| Validation | Regex |

Detects Chinese postal codes (6-digit numbers).

**Examples:** `100000`, `200001`

---

## Passport

| Property | Value |
|----------|-------|
| Detector | `cn/passport` |
| PII Type | `PASSPORT` |
| Confidence | 0.60 |
| Validation | Regex |

Detects Chinese passport numbers.

**Examples:** `E12345678`, `G87654321`

---

## USCC (Unified Social Credit Code)

| Property | Value |
|----------|-------|
| Detector | `cn/uscc` |
| PII Type | `TAX_ID` |
| Confidence | 0.85 |
| Validation | Pattern matching |

Detects Chinese Unified Social Credit Codes, which are 18-character alphanumeric identifiers assigned to organizations.

**Examples:** `91110000MA001AA123`
