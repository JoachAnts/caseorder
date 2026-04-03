# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build the CLI binary
go build ./cmd/caseorder

# Run all tests
go test ./...

# Run a single test
go test -run TestAlphabetical ./...

# Run tests with verbose output
go test -v ./...

# Run the linter on a package
./caseorder ./path/to/package
```

## Architecture

`caseorder` is a Go static analysis linter that enforces consistent ordering of switch case statements. It integrates with `golang.org/x/tools/go/analysis` and can be used standalone or via golangci-lint.

**Entry points:**
- `analyzer.go` — Core logic: AST walking, grouping, sorting, and fix generation
- `cmd/caseorder/main.go` — CLI wrapper using `singlechecker.Main()`

**Analysis flow:**
1. Walk AST for `SwitchStmt` nodes
2. Collect case clauses and their constant values
3. Group cases connected by `fallthrough` into units
4. Sort values within each case clause (e.g., `case 1, 3, 2:` → `case 1, 2, 3:`)
5. Sort groups across the switch (alphabetically for strings, numerically for ints/floats)
6. Emit diagnostics with `SuggestedFix` when ordering is wrong

**Value extraction** (`getValue`): Handles `BasicLit` (strings, ints, floats, chars) and `UnaryExpr` for negatives (e.g., `-1`). Uses `constant.Value` for type-agnostic comparison.

**Testdata** (`testdata/src/`): Each subdirectory (`alphabetical`, `numbers`, `fallthru`, `multi`, `large`, `edgecases`) is an independent package used by `analysistest.Run`. The `// want` comments in those files declare expected diagnostics.

**golangci-lint integration**: The public `New()` function (added in commit `b9e36bd`) exposes the analyzer for use as a plugin.
