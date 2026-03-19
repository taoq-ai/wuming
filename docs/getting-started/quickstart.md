# Quick Start

This guide walks through the most common usage patterns for wuming.

## Basic Usage

The simplest way to use wuming is with the default configuration, which loads all available detectors and uses the redact replacer:

```go
w := wuming.New()
result, err := w.Process(ctx, "My SSN is 123-45-6789")
// result.Redacted: "My SSN is [NATIONAL_ID]"
```

## Configuring Locales

To restrict detection to a specific locale (global detectors always run regardless):

```go
w := wuming.New(
    wuming.WithLocale("nl"),
)
result, err := w.Process(ctx, "BSN: 123456782, SSN: 123-45-6789")
// Only the BSN is detected (SSN is US-specific and filtered out)
```

You can combine multiple locales:

```go
w := wuming.New(
    wuming.WithLocale("nl"),
    wuming.WithLocale("de"),
)
```

## Choosing a Replacer Strategy

wuming ships with four replacer strategies:

=== "Redact"

    Replaces PII with a type label (default):

    ```go
    import "github.com/taoq-ai/wuming/adapter/replacer"

    w := wuming.New(
        wuming.WithReplacer(replacer.NewRedact()),
    )
    // "john@example.com" -> "[EMAIL]"
    ```

=== "Mask"

    Replaces characters with asterisks, preserving the last 4:

    ```go
    w := wuming.New(
        wuming.WithReplacer(replacer.NewMask()),
    )
    // "john@example.com" -> "************.com"
    ```

=== "Hash"

    Replaces PII with a deterministic SHA-256 hash (truncated to 16 hex chars):

    ```go
    w := wuming.New(
        wuming.WithReplacer(replacer.NewHash()),
    )
    // "john@example.com" -> "a8cfcd74832004e0"
    ```

=== "Custom"

    Provide your own replacement function:

    ```go
    w := wuming.New(
        wuming.WithReplacer(replacer.NewCustom("my-replacer", func(m model.Match) string {
            return fmt.Sprintf("<%s:REDACTED>", m.Type)
        })),
    )
    ```

## Filtering by PII Type

Restrict detection to specific PII types:

```go
w := wuming.New(
    wuming.WithPIITypes(model.Email, model.Phone),
)
```

## Confidence Threshold

Filter out low-confidence matches:

```go
w := wuming.New(
    wuming.WithConfidenceThreshold(0.8),
)
```

## Concurrency Control

Limit the number of detectors running in parallel:

```go
w := wuming.New(
    wuming.WithConcurrency(4),
)
```

## Full Working Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/taoq-ai/wuming"
    "github.com/taoq-ai/wuming/adapter/replacer"
    "github.com/taoq-ai/wuming/domain/model"
)

func main() {
    w := wuming.New(
        wuming.WithLocale("nl"),
        wuming.WithReplacer(replacer.NewMask()),
        wuming.WithConfidenceThreshold(0.7),
        wuming.WithPIITypes(model.Email, model.Phone, model.NationalID),
        wuming.WithConcurrency(4),
    )

    text := "Contact Jan at jan@bedrijf.nl or 06-12345678. BSN: 123456782."

    result, err := w.Process(context.Background(), text)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Original:", result.Original)
    fmt.Println("Redacted:", result.Redacted)
    fmt.Printf("Matches:  %d\n", result.MatchCount)

    for _, m := range result.Matches {
        fmt.Printf("  - %s: %q (confidence: %.2f, detector: %s)\n",
            m.Type, m.Value, m.Confidence, m.Detector)
    }
}
```
