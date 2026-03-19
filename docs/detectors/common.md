# Common (Global) Detectors

Common detectors are locale-independent and run regardless of which locale is configured. They cover PII patterns that are universal across geographies.

## Email Address

| Property | Value |
|----------|-------|
| Detector | `common/email` |
| PII Type | `EMAIL` |
| Confidence | 0.95 |
| Validation | Regex pattern matching |

Detects email addresses in standard `user@domain.tld` format.

**Examples:** `john@example.com`, `user.name+tag@company.co.uk`

---

## Credit Card

| Property | Value |
|----------|-------|
| Detector | `common/creditcard` |
| PII Type | `CREDIT_CARD` |
| Confidence | 0.95 |
| Validation | Luhn algorithm |

Detects credit card numbers (13-19 digits) with optional separators (spaces or hyphens). Every candidate is validated using the **Luhn checksum algorithm** to eliminate false positives.

**Examples:** `4111 1111 1111 1111`, `5500-0000-0000-0004`

!!! info "Regulatory Context"
    Credit card numbers are regulated under **PCI-DSS** (Payment Card Industry Data Security Standard) globally. Severity: **Critical**.

---

## IBAN

| Property | Value |
|----------|-------|
| Detector | `common/iban` |
| PII Type | `IBAN` |
| Confidence | 0.95 |
| Validation | ISO 13616 mod-97 checksum |

Detects International Bank Account Numbers. The format is a 2-letter country code, 2 check digits, and up to 30 alphanumeric characters. Every candidate is validated using the **mod-97 algorithm** (ISO 13616).

**Examples:** `NL91ABNA0417164300`, `DE89370400440532013000`

---

## IP Address

| Property | Value |
|----------|-------|
| Detector | `common/ip` |
| PII Type | `IP_ADDRESS` |
| Confidence | 0.90 |
| Validation | IPv4 octet range (0-255), IPv6 regex |

Detects both IPv4 and IPv6 addresses. IPv4 addresses are validated to ensure each octet is within the 0-255 range.

**Examples:** `192.168.1.1`, `2001:0db8:85a3:0000:0000:8a2e:0370:7334`

!!! info "Regulatory Context"
    IP addresses are considered personal data under **GDPR** (EU) when they can be linked to an individual.

---

## URL

| Property | Value |
|----------|-------|
| Detector | `common/url` |
| PII Type | `URL` |
| Confidence | 0.90 |
| Validation | Regex pattern matching |

Detects URLs with `http://` or `https://` schemes.

**Examples:** `https://example.com/profile/12345`, `http://intranet.company.com/users/jdoe`

---

## MAC Address

| Property | Value |
|----------|-------|
| Detector | `common/mac` |
| PII Type | `MAC_ADDRESS` |
| Confidence | 0.90 |
| Validation | Regex pattern matching |

Detects MAC addresses in colon-separated, hyphen-separated, and dot-separated formats.

**Examples:** `00:1A:2B:3C:4D:5E`, `00-1A-2B-3C-4D-5E`, `001A.2B3C.4D5E`
