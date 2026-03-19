# Australia Detectors

Australia detectors cover PII patterns specific to the Australian context. They are activated when the `"au"` locale is configured.

## TFN (Tax File Number)

| Property | Value |
|----------|-------|
| Detector | `au/tfn` |
| PII Type | `TAX_ID` |
| Confidence | 0.85 |
| Validation | Weighted checksum |

Detects Australian Tax File Numbers. The TFN is a unique 8 or 9-digit number validated using a weighted checksum algorithm.

**Examples:** `123 456 782`, `123456782`

!!! warning "Regulatory Context"
    TFN is regulated under the **Privacy Act 1988** and the Tax File Number Guidelines. Unauthorized recording or disclosure is prohibited. Severity: **High**.

---

## Medicare Number

| Property | Value |
|----------|-------|
| Detector | `au/medicare` |
| PII Type | `HEALTH_ID` |
| Confidence | 0.80 |
| Validation | Check digit |

Detects Australian Medicare card numbers (10 or 11 digits).

**Examples:** `2123 45670 1`, `21234567`

---

## ABN (Australian Business Number)

| Property | Value |
|----------|-------|
| Detector | `au/abn` |
| PII Type | `TAX_ID` |
| Confidence | 0.85 |
| Validation | Weighted checksum |

Detects Australian Business Numbers. The ABN is an 11-digit number with a weighted checksum.

**Examples:** `51 824 753 556`, `51824753556`

---

## Phone Number

| Property | Value |
|----------|-------|
| Detector | `au/phone` |
| PII Type | `PHONE` |
| Confidence | 0.85 |
| Validation | Australian format |

Detects Australian phone numbers in mobile and landline formats.

**Examples:** `0412 345 678`, `+61 2 1234 5678`, `(02) 1234 5678`

---

## Postcode

| Property | Value |
|----------|-------|
| Detector | `au/postcode` |
| PII Type | `POSTAL_CODE` |
| Confidence | 0.60 |
| Validation | Regex |

Detects Australian postcodes (4-digit numbers).

**Examples:** `2000`, `3000`, `4000`
