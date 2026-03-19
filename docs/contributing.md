# Contributing

Thank you for your interest in contributing to wuming. This guide covers the development workflow, how to add new detectors, and code conventions.

## Development Setup

1. **Clone the repository:**

    ```bash
    git clone https://github.com/taoq-ai/wuming.git
    cd wuming
    ```

2. **Verify Go version:**

    ```bash
    go version
    # Requires Go 1.22+
    ```

3. **Run the tests:**

    ```bash
    make test
    ```

4. **Run the linter:**

    ```bash
    make lint
    ```

## How to Add a New Detector

### 1. Create the package

Create a new directory under `adapter/detector/<locale>/` for locale-specific detectors, or add to `adapter/detector/common/` for global patterns.

### 2. Implement the Detector interface

```go
package xx

import (
    "context"
    "regexp"

    "github.com/taoq-ai/wuming/domain/model"
)

const locale = "xx"

var myPatternRe = regexp.MustCompile(`...`)

type MyDetector struct{}

func NewMyDetector() *MyDetector { return &MyDetector{} }

func (d *MyDetector) Name() string              { return "xx/my_pattern" }
func (d *MyDetector) Locales() []string         { return []string{locale} }
func (d *MyDetector) PIITypes() []model.PIIType { return []model.PIIType{model.NationalID} }

func (d *MyDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
    // Your detection logic here
    return nil, nil
}
```

### 3. Add a helpers.go file

If your locale package needs shared utilities (like a `findAll` helper), add a `helpers.go` file following the convention used by existing locale packages.

### 4. Write tests

Create a `xx_test.go` file with table-driven tests covering:

- Valid matches (true positives)
- Invalid matches (true negatives)
- Edge cases (boundary values, formatting variations)
- Checksum validation (if applicable)

### 5. Register the detector

Add your detector to the appropriate registry so that it is included when users configure that locale.

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Use table-driven tests
- Keep detectors focused -- one PII type per detector
- Include doc comments on all exported types and functions
- Name detector files after the PII type they detect (e.g., `bsn.go`, `ssn.go`, `phone.go`)
- Set confidence scores consistent with existing detectors:
    - 0.90+ for checksum-validated patterns
    - 0.80-0.89 for structurally validated patterns
    - 0.60-0.79 for regex-only or context-dependent patterns

## Testing

All tests should pass before submitting a pull request:

```bash
make test
```

Run tests for a specific package:

```bash
go test ./adapter/detector/xx/...
```

## Pull Requests

- Create a feature branch from `main`
- Write a clear PR description
- Ensure all tests pass
- Reference any related issues
