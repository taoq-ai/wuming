# Brazil Detectors

Brazil detectors cover PII patterns specific to the Brazilian context. They are activated when the `"br"` locale is configured.

## CPF (Cadastro de Pessoas Fisicas)

| Property | Value |
|----------|-------|
| Detector | `br/cpf` |
| PII Type | `NATIONAL_ID` |
| Confidence | 0.90 |
| Validation | Check digits |

Detects Brazilian individual taxpayer registration numbers (CPF). The CPF is an 11-digit number with two check digits validated using a weighted modular arithmetic algorithm.

**Examples:** `123.456.789-09`, `12345678909`

---

## CNPJ (Cadastro Nacional da Pessoa Juridica)

| Property | Value |
|----------|-------|
| Detector | `br/cnpj` |
| PII Type | `TAX_ID` |
| Confidence | 0.90 |
| Validation | Check digits |

Detects Brazilian corporate taxpayer registration numbers (CNPJ). The CNPJ is a 14-digit number with two check digits.

**Examples:** `11.222.333/0001-81`, `11222333000181`

---

## Phone Number

| Property | Value |
|----------|-------|
| Detector | `br/phone` |
| PII Type | `PHONE` |
| Confidence | 0.85 |
| Validation | Brazilian format |

Detects Brazilian phone numbers in mobile and landline formats.

**Examples:** `(11) 91234-5678`, `+55 11 91234-5678`

---

## CEP (Codigo de Enderecamento Postal)

| Property | Value |
|----------|-------|
| Detector | `br/cep` |
| PII Type | `POSTAL_CODE` |
| Confidence | 0.60 |
| Validation | Regex |

Detects Brazilian postal codes in the `NNNNN-NNN` format.

**Examples:** `01001-000`, `12345-678`

---

## PIS/PASEP

| Property | Value |
|----------|-------|
| Detector | `br/pis` |
| PII Type | `NATIONAL_ID` |
| Confidence | 0.80 |
| Validation | Check digit |

Detects Brazilian social integration program numbers (PIS/PASEP), used for social security and employment purposes.

**Examples:** `123.45678.90-1`

---

## CNH (Carteira Nacional de Habilitacao)

| Property | Value |
|----------|-------|
| Detector | `br/cnh` |
| PII Type | `DRIVERS_LICENSE` |
| Confidence | 0.75 |
| Validation | Pattern matching |

Detects Brazilian driver's license numbers (CNH).

**Examples:** `12345678901`
