# Netherlands Detectors

Netherlands detectors cover PII patterns specific to the Dutch context. They are activated when the `"nl"` locale is configured.

## Burgerservicenummer (BSN)

| Property | Value |
|----------|-------|
| Detector | `nl/bsn` |
| PII Type | `NATIONAL_ID` |
| Confidence | 0.90 |
| Validation | 11-proof checksum |

Detects Dutch citizen service numbers (BSN). The BSN is a 9-digit number validated using the **11-proof** (elfproef) algorithm:

`9*d1 + 8*d2 + 7*d3 + 6*d4 + 5*d5 + 4*d6 + 3*d7 + 2*d8 - 1*d9`

The result must be divisible by 11 and must not be 0.

**Examples:** `123456782`, `123.456.782`, `123 456 782`

!!! warning "Regulatory Context"
    The BSN is regulated under **Wet algemene bepalingen burgerservicenummer (Wabb)**. Its use is restricted to government agencies, healthcare providers, and certain other authorized organizations. Severity: **High**.

---

## Phone Number

| Property | Value |
|----------|-------|
| Detector | `nl/phone` |
| PII Type | `PHONE` |
| Confidence | 0.85 |
| Validation | Dutch phone format |

Detects Dutch phone numbers in mobile and landline formats:

- **Mobile:** `06-XXXXXXXX`, `+31 6 XXXXXXXX`, `0031 6 XXXXXXXX`
- **Landline:** `0XX-XXXXXXX`, `+31 XX XXXXXXX`

**Examples:** `06-12345678`, `+31 20 1234567`, `0031 6 12345678`

---

## Postal Code

| Property | Value |
|----------|-------|
| Detector | `nl/postal` |
| PII Type | `POSTAL_CODE` |
| Confidence | 0.90 |
| Validation | Format + invalid letter combinations |

Detects Dutch postal codes in the `NNNN XX` format (4 digits + 2 uppercase letters). The first digit cannot be 0, and certain letter combinations (`SA`, `SD`, `SS`) that are not used in the Dutch postal system are excluded.

**Examples:** `1234 AB`, `1012WX`

---

## KvK Number (Chamber of Commerce)

| Property | Value |
|----------|-------|
| Detector | `nl/kvk` |
| PII Type | `TAX_ID` |
| Confidence | 0.60 (bare) / 0.90 (with context) |
| Validation | Context-boosted confidence |

Detects Kamer van Koophandel (Chamber of Commerce) registration numbers. These are 8-digit numbers. Because bare 8-digit sequences are ambiguous, confidence is boosted to 0.90 when preceded by keywords like "KvK" or "Kamer van Koophandel".

**Examples:** `KvK: 12345678`, `Kamer van Koophandel: 87654321`

---

## ID Document (ID Card & Passport)

| Property | Value |
|----------|-------|
| Detector | `nl/id_document` |
| PII Type | `NATIONAL_ID` / `PASSPORT` |
| Confidence | 0.70 |
| Validation | Pattern matching |

Detects Dutch identity card and passport numbers:

- **ID Card:** 9-character alphanumeric string containing both letters and digits (e.g., `SPECI2014`)
- **Passport:** 2 uppercase letters followed by 7 alphanumeric characters containing at least one digit

**Examples:** `SPECI2014`, `NW12P4F07`

!!! info "Regulatory Context"
    Identity document numbers are protected under the **GDPR** and Dutch implementation laws (UAVG). Severity: **High**.
