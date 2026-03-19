# EU-wide Detectors

EU-wide detectors cover PII patterns that apply across European Union member states. They are activated when the `"eu"` locale is configured.

## VAT Identification Number

| Property | Value |
|----------|-------|
| Detector | `eu/vat` |
| PII Type | `TAX_ID` |
| Confidence | 0.90 |
| Validation | Country prefix + format regex |

Detects EU VAT identification numbers for all 27 member states. Each number begins with a two-letter country prefix followed by a country-specific format.

**Supported countries and formats:**

| Prefix | Country | Format |
|--------|---------|--------|
| AT | Austria | `ATU` + 8 digits |
| BE | Belgium | `BE0` or `BE1` + 9 digits |
| BG | Bulgaria | 9-10 digits |
| CY | Cyprus | 8 digits + 1 letter |
| CZ | Czech Republic | 8-10 digits |
| DE | Germany | 9 digits |
| DK | Denmark | 8 digits |
| EE | Estonia | 9 digits |
| EL | Greece | 9 digits |
| ES | Spain | 1 alphanumeric + 7 digits + 1 alphanumeric |
| FI | Finland | 8 digits |
| FR | France | 2 alphanumeric + 9 digits |
| HR | Croatia | 11 digits |
| HU | Hungary | 8 digits |
| IE | Ireland | 1 digit + 1 alphanumeric + 5 digits + 1-2 letters |
| IT | Italy | 11 digits |
| LT | Lithuania | 9 or 12 digits |
| LU | Luxembourg | 8 digits |
| LV | Latvia | 11 digits |
| MT | Malta | 8 digits |
| NL | Netherlands | 9 digits + `B` + 2 digits |
| PL | Poland | 10 digits |
| PT | Portugal | 9 digits |
| RO | Romania | 2-10 digits |
| SE | Sweden | 12 digits |
| SI | Slovenia | 8 digits |
| SK | Slovakia | 10 digits |

**Examples:** `NL123456789B01`, `DE123456789`, `ATU12345678`

!!! info "Regulatory Context"
    VAT numbers are business identifiers regulated under EU VAT Directive (2006/112/EC). While they are primarily organizational, they can identify sole traders and are therefore considered PII under the GDPR.

---

## Passport MRZ (Machine Readable Zone)

| Property | Value |
|----------|-------|
| Detector | `eu/passport_mrz` |
| PII Type | `PASSPORT` |
| Confidence | 0.95 |
| Validation | ICAO 9303 TD3 format |

Detects ICAO 9303 TD3 (passport) Machine Readable Zones. The MRZ consists of two lines of 44 characters each, containing the document type, issuing state, holder's name, passport number, nationality, date of birth, sex, expiry date, and check digits.

This detector has a high confidence score because the MRZ format is highly structured and unlikely to produce false positives.

!!! warning "Regulatory Context"
    Passport MRZ data contains multiple PII elements (name, nationality, date of birth, passport number) and is protected under GDPR and national identity document laws. Severity: **Critical**.
