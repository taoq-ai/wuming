# Japan Detectors

Japan detectors cover PII patterns specific to the Japanese context. They are activated when the `"jp"` locale is configured.

## My Number (Individual Number)

| Property | Value |
|----------|-------|
| Detector | `jp/my_number` |
| PII Type | `NATIONAL_ID` |
| Confidence | 0.90 |
| Validation | Check digit |

Detects Japanese Individual Number (My Number) identifiers. The My Number is a 12-digit number assigned to residents, with check digit validation.

**Examples:** `123456789012`

!!! warning "Regulatory Context"
    My Number is regulated under **APPI** and the My Number Act. Unauthorized use or collection is prohibited. Severity: **High**.

---

## Corporate Number

| Property | Value |
|----------|-------|
| Detector | `jp/corporate_number` |
| PII Type | `TAX_ID` |
| Confidence | 0.90 |
| Validation | Check digit |

Detects Japanese corporate numbers (13-digit identifiers assigned to corporations).

**Examples:** `1234567890123`

---

## Phone Number

| Property | Value |
|----------|-------|
| Detector | `jp/phone` |
| PII Type | `PHONE` |
| Confidence | 0.85 |
| Validation | Japanese format |

Detects Japanese phone numbers in fixed-line and mobile formats.

**Examples:** `03-1234-5678`, `090-1234-5678`, `+81 3 1234 5678`

---

## Postal Code

| Property | Value |
|----------|-------|
| Detector | `jp/postal` |
| PII Type | `POSTAL_CODE` |
| Confidence | 0.60 |
| Validation | Regex |

Detects Japanese postal codes in the `NNN-NNNN` format.

**Examples:** `100-0001`, `160-0023`

---

## Passport

| Property | Value |
|----------|-------|
| Detector | `jp/passport` |
| PII Type | `PASSPORT` |
| Confidence | 0.60 |
| Validation | Regex |

Detects Japanese passport numbers.

**Examples:** `TK1234567`
