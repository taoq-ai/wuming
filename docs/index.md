# 无名 wuming

> *The Tao that can be told is not the eternal Tao.*
> *The name that can be named is not the eternal name.*
> *The nameless is the beginning of heaven and earth.*
>
> -- Tao Te Ching, Chapter 1

---

**wuming** (无名 -- "The Nameless") is a Go library for detecting and removing Personally Identifiable Information (PII) from text. It supports global PII standards across multiple locales and provides pluggable detection and replacement strategies via a hexagonal architecture.

## Key Features

- **Multi-locale detection** -- built-in support for common (global), US, NL, EU, GB, DE, and FR PII patterns
- **Pluggable replacers** -- redact, mask, hash, or bring your own replacement strategy
- **Hexagonal architecture** -- clean separation between domain logic and adapters for easy extension
- **Concurrent processing** -- detectors run in parallel with configurable concurrency limits
- **Confidence scoring** -- each match carries a confidence score for fine-grained filtering
- **Checksum validation** -- Luhn, mod-97, 11-proof, mod-11, and other algorithms reduce false positives

## Quick Example

```go
package main

import (
    "context"
    "fmt"

    "github.com/taoq-ai/wuming"
)

func main() {
    w := wuming.New(
        wuming.WithLocale("nl"),
    )

    result, err := w.Process(context.Background(),
        "Call me at 06-12345678 or email john@example.com",
    )
    if err != nil {
        panic(err)
    }

    fmt.Println(result.Redacted)
    // Output: Call me at [PHONE] or email [EMAIL]
}
```

## Getting Started

Head over to the [Installation](getting-started/installation.md) guide to add wuming to your project, then follow the [Quick Start](getting-started/quickstart.md) to begin detecting PII in minutes.
