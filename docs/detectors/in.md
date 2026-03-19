# India Detectors

India detectors cover PII patterns specific to the Indian context. They are activated when the `"in"` locale is configured.

## Aadhaar

| Property | Value |
|----------|-------|
| Detector | `in/aadhaar` |
| PII Type | `NATIONAL_ID` |
| Confidence | 0.90 |
| Validation | Verhoeff checksum |

Detects Indian Aadhaar numbers. Aadhaar is a 12-digit unique identity number issued by UIDAI, validated using the Verhoeff checksum algorithm.

**Examples:** `2345 6789 0123`, `234567890123`

!!! warning "Regulatory Context"
    Aadhaar is regulated under the **Aadhaar Act, 2016** and the **DPDP Act**. Unauthorized collection or storage is prohibited. Severity: **Critical**.

---

## PAN (Permanent Account Number)

| Property | Value |
|----------|-------|
| Detector | `in/pan` |
| PII Type | `TAX_ID` |
| Confidence | 0.85 |
| Validation | Pattern matching |

Detects Indian Permanent Account Numbers issued by the Income Tax Department. PAN follows the format `AAAAA9999A` (5 letters, 4 digits, 1 letter).

**Examples:** `ABCDE1234F`

---

## Phone Number

| Property | Value |
|----------|-------|
| Detector | `in/phone` |
| PII Type | `PHONE` |
| Confidence | 0.85 |
| Validation | Indian format |

Detects Indian phone numbers in mobile and landline formats.

**Examples:** `+91 98765 43210`, `098765 43210`

---

## PIN Code

| Property | Value |
|----------|-------|
| Detector | `in/pin_code` |
| PII Type | `POSTAL_CODE` |
| Confidence | 0.60 |
| Validation | Regex |

Detects Indian postal PIN codes (6-digit numbers).

**Examples:** `110001`, `400001`

---

## Passport

| Property | Value |
|----------|-------|
| Detector | `in/passport` |
| PII Type | `PASSPORT` |
| Confidence | 0.60 |
| Validation | Regex |

Detects Indian passport numbers.

**Examples:** `A1234567`, `Z9876543`

---

## GSTIN

| Property | Value |
|----------|-------|
| Detector | `in/gstin` |
| PII Type | `TAX_ID` |
| Confidence | 0.85 |
| Validation | Check digit |

Detects Indian Goods and Services Tax Identification Numbers (GSTIN). GSTIN is a 15-character alphanumeric identifier.

**Examples:** `22AAAAA0000A1Z5`
