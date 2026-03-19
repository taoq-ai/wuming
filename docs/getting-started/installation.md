# Installation

## Requirements

- **Go 1.22** or later
- No external dependencies beyond the Go standard library

## Install

Add wuming to your Go module:

```bash
go get github.com/taoq-ai/wuming
```

## Verify

Create a simple test file to verify the installation:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/taoq-ai/wuming"
)

func main() {
    w := wuming.New()
    result, err := w.Process(context.Background(), "Email: test@example.com")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result.Redacted)
    // Output: Email: [EMAIL]
}
```

Run it:

```bash
go run main.go
```

If you see the redacted output, wuming is installed and working correctly.
