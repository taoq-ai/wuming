# Replacers

Replacers determine how detected PII is substituted in the output text. wuming ships with four built-in strategies, and you can create your own.

## Built-in Replacers

### Redact

Replaces each match with a type-labeled placeholder. This is the **default** replacer.

```go
import "github.com/taoq-ai/wuming/adapter/replacer"

r := replacer.NewRedact()
// "john@example.com" -> "[EMAIL]"
// "123-45-6789"      -> "[NATIONAL_ID]"
```

The placeholder format defaults to `[%s]` where `%s` is the PII type string. You can customize it:

```go
r := &replacer.Redact{Format: "<%s>"}
// "john@example.com" -> "<EMAIL>"
```

---

### Mask

Replaces characters with a mask character, preserving the last N characters.

```go
r := replacer.NewMask()
// Default: mask with '*', preserve last 4 characters
// "john@example.com" -> "************.com"
// "123-45-6789"      -> "*******6789"
```

Customize the mask character and number of preserved characters:

```go
r := &replacer.Mask{Char: '#', Preserve: 2}
// "john@example.com" -> "##############om"
```

---

### Hash

Replaces each match with a deterministic SHA-256 hash, truncated to a configurable length. The same input always produces the same hash, which is useful for de-identification while preserving the ability to detect duplicates.

```go
r := replacer.NewHash()
// Default: 16 hex characters, no salt
// "john@example.com" -> "a8cfcd74832004e0"
```

Add a salt for extra security:

```go
r := &replacer.Hash{Length: 16, Salt: "my-secret-salt"}
```

---

### Custom

Provide your own replacement function for full control:

```go
import (
    "fmt"
    "github.com/taoq-ai/wuming/adapter/replacer"
    "github.com/taoq-ai/wuming/domain/model"
)

r := replacer.NewCustom("my-replacer", func(m model.Match) string {
    return fmt.Sprintf("[%s:%s:%.0f%%]", m.Locale, m.Type, m.Confidence*100)
})
// "123456782" (Dutch BSN) -> "[nl:NATIONAL_ID:90%]"
```

## Creating a Custom Replacer

You can also implement the `Replacer` interface directly:

```go
type Replacer interface {
    Replace(text string, matches []model.Match) (string, error)
    Name() string
}
```

When implementing `Replace`, process matches in **reverse order** (from the end of the text to the beginning) to preserve byte offsets as you substitute text. The built-in replacers demonstrate this pattern -- they sort matches by descending start position before applying replacements.

```go
type MyReplacer struct{}

func (r *MyReplacer) Name() string { return "my-replacer" }

func (r *MyReplacer) Replace(text string, matches []model.Match) (string, error) {
    // Sort matches by start position descending
    sorted := make([]model.Match, len(matches))
    copy(sorted, matches)
    sort.Slice(sorted, func(i, j int) bool {
        return sorted[i].Start > sorted[j].Start
    })

    result := []byte(text)
    for _, m := range sorted {
        replacement := []byte("***")
        result = append(result[:m.Start], append(replacement, result[m.End:]...)...)
    }
    return string(result), nil
}
```

Then use it:

```go
w := wuming.New(
    wuming.WithReplacer(&MyReplacer{}),
)
```
