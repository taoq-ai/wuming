# United Kingdom Detectors

UK detectors cover PII patterns specific to the United Kingdom. They are activated when the `"gb"` locale is configured.

## National Insurance Number (NIN)

| Property | Value |
|----------|-------|
| Detector | `gb/nin` |
| PII Type | `NATIONAL_ID` |
| Confidence | 0.85 |
| Validation | HMRC prefix rules |

Detects UK National Insurance Numbers in the format `XX 99 99 99 X` (2 letters, 6 digits, 1 letter A-D). Validates the two-letter prefix against HMRC allocation rules:

- First letter must not be D, F, I, Q, U, or V
- Second letter must not be D, F, I, O, Q, U, or V
- Certain prefix combinations are never allocated: BG, GB, NK, KN, TN, NT, ZZ

**Examples:** `AB 12 34 56 C`, `AB123456C`

!!! warning "Regulatory Context"
    NINs are regulated by **HMRC** and protected under the **UK Data Protection Act 2018** (UK GDPR). Severity: **High**.

---

## NHS Number

| Property | Value |
|----------|-------|
| Detector | `gb/nhs` |
| PII Type | `HEALTH_ID` |
| Confidence | 0.90 |
| Validation | Mod-11 check digit |

Detects UK NHS Numbers (10 digits, optionally formatted as `XXX XXX XXXX` or `XXX-XXX-XXXX`). Validates using the **mod-11 algorithm**: the first 9 digits are weighted (10, 9, 8, ..., 2), summed, and the check digit is `11 - (sum % 11)`. If the remainder is 10, the number is invalid; if 11, the check digit is 0.

**Examples:** `943 476 5919`, `943-476-5919`, `9434765919`

!!! warning "Regulatory Context"
    NHS Numbers are protected under the **NHS Act 2006** and **UK GDPR**. They are classified as health data and subject to special category processing rules. Severity: **High**.

---

## Unique Taxpayer Reference (UTR)

| Property | Value |
|----------|-------|
| Detector | `gb/utr` |
| PII Type | `TAX_ID` |
| Confidence | 0.55 (bare) / 0.85 (with context) |
| Validation | Context-boosted confidence |

Detects UK Unique Taxpayer References (10-digit numbers). Because bare 10-digit sequences are highly ambiguous, the detector uses context boosting -- confidence is raised to 0.85 when preceded by keywords like "UTR", "tax reference", or "unique taxpayer reference".

**Examples:** `UTR: 1234567890`, `Tax reference: 1234567890`

---

## Phone Number

| Property | Value |
|----------|-------|
| Detector | `gb/phone` |
| PII Type | `PHONE` |
| Confidence | 0.85 |
| Validation | UK phone format |

Detects UK phone numbers in mobile and landline formats:

- **Mobile:** `07XXX XXXXXX`, `+44 7XXX XXXXXX`
- **London:** `020 XXXX XXXX`, `+44 20 XXXX XXXX`
- **Landline:** `01XX XXX XXXX`, `+44 1XX XXX XXXX`

Supports space, dash, and dot separators.

**Examples:** `07911 123456`, `+44 20 7946 0958`, `01onal234-567-8901`

---

## Postcode

| Property | Value |
|----------|-------|
| Detector | `gb/postcode` |
| PII Type | `POSTAL_CODE` |
| Confidence | 0.90 |
| Validation | UK postcode format regex |

Detects UK postcodes in all valid formats: `A9 9AA`, `A99 9AA`, `A9A 9AA`, `AA9 9AA`, `AA99 9AA`, `AA9A 9AA`.

**Examples:** `SW1A 1AA`, `EC1A 1BB`, `W1D 3QU`
