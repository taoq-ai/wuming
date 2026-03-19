# South Korea Detectors

South Korea detectors cover PII patterns specific to the Korean context. They are activated when the `"kr"` locale is configured.

## RRN (Resident Registration Number)

| Property | Value |
|----------|-------|
| Detector | `kr/rrn` |
| PII Type | `NATIONAL_ID` |
| Confidence | 0.90 |
| Validation | Check digit |

Detects Korean Resident Registration Numbers. The RRN is a 13-digit number encoding birth date, gender, and a check digit.

**Examples:** `901231-1234567`

!!! warning "Regulatory Context"
    RRN is regulated under **PIPA** (Personal Information Protection Act). Collection requires explicit consent. Severity: **Critical**.

---

## Phone Number

| Property | Value |
|----------|-------|
| Detector | `kr/phone` |
| PII Type | `PHONE` |
| Confidence | 0.85 |
| Validation | Korean format |

Detects Korean phone numbers in mobile and landline formats.

**Examples:** `010-1234-5678`, `+82 10 1234 5678`, `02-123-4567`

---

## Postal Code

| Property | Value |
|----------|-------|
| Detector | `kr/postal` |
| PII Type | `POSTAL_CODE` |
| Confidence | 0.60 |
| Validation | Regex |

Detects Korean postal codes (5-digit numbers).

**Examples:** `06236`, `12345`

---

## Passport

| Property | Value |
|----------|-------|
| Detector | `kr/passport` |
| PII Type | `PASSPORT` |
| Confidence | 0.60 |
| Validation | Regex |

Detects Korean passport numbers.

**Examples:** `M12345678`, `R98765432`
