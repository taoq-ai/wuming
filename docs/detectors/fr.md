# France Detectors

French detectors cover PII patterns specific to France. They are activated when the `"fr"` locale is configured.

## NIR (Social Security Number)

| Property | Value |
|----------|-------|
| Detector | `fr/nir` |
| PII Type | `NATIONAL_ID` |
| Confidence | 0.90 |
| Validation | Mod-97 control key |

Detects the French NIR (Numero d'Inscription au Repertoire), also known as the numero de securite sociale. The NIR is a 15-digit number with the structure `X XX XX XXXXX XXX XX`:

- Digit 1: sex (1=male, 2=female)
- Digits 2-3: year of birth
- Digits 4-5: month of birth (01-12)
- Digits 6-7: department (01-95, or 2A/2B for Corsica)
- Digits 8-10: commune code
- Digits 11-13: order number
- Digits 14-15: control key

The control key is validated as `97 - (first 13 digits mod 97)`. Special handling is applied for Corsica departments (2A and 2B) which contain letters in the department field.

**Examples:** `1 85 05 78 006 084 36` (with spaces), `185057800608436` (without spaces)

!!! warning "Regulatory Context"
    The NIR is regulated under French law and the **CNIL** (Commission nationale de l'informatique et des libertes). Its use is strictly controlled and protected under **GDPR**. Severity: **High**.

---

## NIF (Tax Identification Number)

| Property | Value |
|----------|-------|
| Detector | `fr/nif` |
| PII Type | `TAX_ID` |
| Confidence | 0.70 |
| Validation | Regex pattern |

Detects the French Numero d'Identification Fiscale (13 digits, first digit 0-3). Supports formatted versions with spaces, dashes, or dots as separators.

**Examples:** `01 23 456 789 012`, `0123456789012`

---

## CNI (National Identity Card)

| Property | Value |
|----------|-------|
| Detector | `fr/id_card` |
| PII Type | `NATIONAL_ID` |
| Confidence | 0.65 |
| Validation | Regex pattern (old + new format) |

Detects French Carte Nationale d'Identite numbers in both formats:

- **Old format** (before 2021): 12 digits
- **New format** (since 2021): 9 alphanumeric characters

**Examples:** `123456789012` (old), `ABC123DEF` (new)

!!! note
    This detector has a lower confidence score because both formats can match non-CNI strings. Use `WithConfidenceThreshold()` to filter as needed.

---

## Phone Number

| Property | Value |
|----------|-------|
| Detector | `fr/phone` |
| PII Type | `PHONE` |
| Confidence | 0.85 |
| Validation | French phone format |

Detects French phone numbers in local and international formats:

- **Local:** `0X XX XX XX XX` (where X = 1-7 for landline/mobile)
- **International:** `+33 X XX XX XX XX`

Supports separators: space, dot, dash.

**Examples:** `01 23 45 67 89`, `+33 6 12 34 56 78`, `06.12.34.56.78`

---

## Postal Code (Code Postal)

| Property | Value |
|----------|-------|
| Detector | `fr/postal` |
| PII Type | `POSTAL_CODE` |
| Confidence | 0.60 |
| Validation | Department prefix validation |

Detects French postal codes (5 digits). The first two digits represent the department number and are validated to be within valid ranges (01-95, 97, 98 -- department 96 does not exist).

**Examples:** `75001` (Paris), `13001` (Marseille), `97100` (overseas territory)
