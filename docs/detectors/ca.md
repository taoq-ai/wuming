# Canada Detectors

Canada detectors cover PII patterns specific to the Canadian context. They are activated when the `"ca"` locale is configured.

## SIN (Social Insurance Number)

| Property | Value |
|----------|-------|
| Detector | `ca/sin` |
| PII Type | `NATIONAL_ID` |
| Confidence | 0.90 |
| Validation | Luhn check digit |

Detects Canadian Social Insurance Numbers. The SIN is a 9-digit number validated using the Luhn algorithm.

**Examples:** `123 456 782`, `123-456-782`

!!! warning "Regulatory Context"
    SIN is regulated under **PIPEDA** and the Privacy Act. Collection is restricted to authorized purposes. Severity: **High**.

---

## Phone Number

| Property | Value |
|----------|-------|
| Detector | `ca/phone` |
| PII Type | `PHONE` |
| Confidence | 0.85 |
| Validation | NANP format |

Detects Canadian phone numbers following the North American Numbering Plan.

**Examples:** `(416) 123-4567`, `+1 604 123 4567`, `416-123-4567`

---

## Postal Code

| Property | Value |
|----------|-------|
| Detector | `ca/postal_code` |
| PII Type | `POSTAL_CODE` |
| Confidence | 0.60 |
| Validation | Canadian format |

Detects Canadian postal codes in the `A1A 1A1` format (alternating letter-digit-letter, space, digit-letter-digit).

**Examples:** `K1A 0B1`, `M5V 2T6`

---

## Passport

| Property | Value |
|----------|-------|
| Detector | `ca/passport` |
| PII Type | `PASSPORT` |
| Confidence | 0.60 |
| Validation | Regex |

Detects Canadian passport numbers.

**Examples:** `AB123456`
