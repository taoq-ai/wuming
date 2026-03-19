# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in wuming, please report it responsibly.

**Contact:** security@taoq.ai

Alternatively, you can use [GitHub Security Advisories](https://github.com/taoq-ai/wuming/security/advisories/new) to report vulnerabilities privately.

Please **do not** open a public GitHub issue for security vulnerabilities.

## What Counts as a Security Issue

- Regular expression denial of service (ReDoS) in any detector pattern
- Incorrect redaction that leaks PII (e.g., partial replacement, off-by-one in byte offsets)
- Panics or crashes caused by crafted input
- Information leakage through error messages

## What Is Not a Security Issue

- Feature requests for new PII types or locales
- False positives or false negatives in detection (these are bugs, not security issues)
- Performance issues

## Response Timeline

- **Acknowledgment:** within 3 business days
- **Initial assessment:** within 7 business days
- **Fix or mitigation:** best effort, typically within 30 days for confirmed vulnerabilities

## Important Note

This library processes text entirely in-process. It makes **no network calls**, stores **no data**, and has **no external dependencies**. The attack surface is limited to the input text provided by the caller.
