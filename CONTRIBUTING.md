# Contributing to wuming

Thank you for your interest in contributing to wuming. This guide covers everything you need to get started.

## Development Setup

### Prerequisites

- **Go 1.22** or later
- **golangci-lint** (for linting)

### Getting Started

```sh
git clone https://github.com/taoq-ai/wuming.git
cd wuming
make test
```

### Useful Make Targets

| Command      | Description                       |
|--------------|-----------------------------------|
| `make build` | Build all packages                |
| `make test`  | Run tests with race detection     |
| `make lint`  | Run golangci-lint                 |
| `make fmt`   | Format code with gofmt            |
| `make vet`   | Run go vet                        |
| `make check` | Run fmt, vet, test, and lint      |

## How to Add a New Locale Detector

Each locale lives in its own package under `adapter/detector/`. Follow these steps to add a new locale (for example, `jp` for Japan):

### 1. Create the package directory

```sh
mkdir -p adapter/detector/jp
```

### 2. Create a helpers file

Create `adapter/detector/jp/helpers.go` with shared utilities:

```go
package jp

import (
    "context"
    "regexp"

    "github.com/taoq-ai/wuming/domain/model"
)

const locale = "jp"

func findAll(re *regexp.Regexp, text string, piiType model.PIIType, confidence float64, detector string) []model.Match {
    results := re.FindAllStringIndex(text, -1)
    if len(results) == 0 {
        return nil
    }

    var matches []model.Match
    for _, loc := range results {
        matches = append(matches, model.Match{
            Type:       piiType,
            Value:      text[loc[0]:loc[1]],
            Start:      loc[0],
            End:        loc[1],
            Confidence: confidence,
            Locale:     locale,
            Detector:   detector,
        })
    }
    return matches
}
```

### 3. Create a detector

Create `adapter/detector/jp/my_number.go`:

```go
package jp

import (
    "context"
    "regexp"

    "github.com/taoq-ai/wuming/domain/model"
)

var myNumberRe = regexp.MustCompile(`\b\d{4}\s?\d{4}\s?\d{4}\b`)

// MyNumberDetector detects Japanese My Number (Individual Number).
type MyNumberDetector struct{}

func NewMyNumberDetector() *MyNumberDetector { return &MyNumberDetector{} }

func (d *MyNumberDetector) Name() string              { return "jp/my_number" }
func (d *MyNumberDetector) Locales() []string         { return []string{locale} }
func (d *MyNumberDetector) PIITypes() []model.PIIType { return []model.PIIType{model.NationalID} }

func (d *MyNumberDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
    return findAll(myNumberRe, text, model.NationalID, 0.85, d.Name()), nil
}
```

### 4. Write tests

Create `adapter/detector/jp/jp_test.go` using table-driven tests:

```go
package jp

import (
    "context"
    "testing"

    "github.com/taoq-ai/wuming/domain/model"
)

func TestMyNumberDetector(t *testing.T) {
    d := NewMyNumberDetector()

    tests := []struct {
        name  string
        input string
        want  int
    }{
        {"valid my number", "My number is 1234 5678 9012", 1},
        {"no match", "Hello world", 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            matches, err := d.Detect(context.Background(), tt.input)
            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
            if len(matches) != tt.want {
                t.Errorf("got %d matches, want %d", len(matches), tt.want)
            }
        })
    }
}
```

### 5. Verify the interface

Ensure your detector satisfies the `port.Detector` interface by adding a compile-time check at the top of your detector file:

```go
var _ port.Detector = (*MyNumberDetector)(nil)
```

## How to Add a New Replacer

Replacers live in `adapter/replacer/`. To add one:

1. Create a new file in `adapter/replacer/` (e.g., `anonymize.go`).
2. Implement the `port.Replacer` interface:
   - `Replace(text string, matches []model.Match) (string, error)`
   - `Name() string`
3. Add a constructor function (e.g., `NewAnonymize()`).
4. Write tests in `adapter/replacer/replacer_test.go` or a new test file.
5. Add a compile-time interface check: `var _ port.Replacer = (*Anonymize)(nil)`

## Code Style

- Run `gofmt` on all code before committing.
- Run `golangci-lint run ./...` and fix any issues.
- Follow standard Go conventions: exported names have GoDoc comments, packages have doc comments.
- Keep functions focused and small.

## PR Workflow

1. **Branch from `main`** — create a feature branch with a descriptive name (e.g., `feat/jp-detectors`).
2. **One feature per PR** — keep pull requests focused on a single change.
3. **Use conventional commits** — prefix commit messages with a type:
   - `feat:` for new features
   - `fix:` for bug fixes
   - `docs:` for documentation
   - `test:` for test additions or changes
   - `refactor:` for refactoring
   - `chore:` for maintenance tasks
4. **Ensure CI passes** — `make check` should succeed before opening a PR.
5. **Write a clear PR description** — explain what the change does and why.

## Testing Expectations

- Use **table-driven tests** for detector and replacer tests.
- Include both positive matches and negative (no-match) cases.
- Add **interface compliance checks** at compile time:
  ```go
  var _ port.Detector = (*MyDetector)(nil)
  ```
- Run tests with race detection: `go test -race ./...`
- Aim for meaningful coverage of edge cases, not just line coverage.
