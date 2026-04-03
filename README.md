# caseorder

A Go linter that enforces consistent ordering of `switch` case statements — alphabetically for strings, numerically for integers and floats.

I built this while working on a Game Boy emulator. The CPU instruction set has 500+ opcodes, and I was implementing them in a giant switch statement. It was becoming difficult to tell at a glance whether a case was missing or duplicated. If you're building something similar: an emulator, a bytecode interpreter, a compiler, or any other state machine driven by a large switch, this linter keeps things navigable as the case list grows.

## Install

```sh
go install github.com/JoachAnts/caseorder/cmd/caseorder@latest
```

## Usage

```sh
# Lint all packages in the current module
caseorder ./...

# Lint a specific package
caseorder ./internal/handlers

# Apply suggested fixes automatically
caseorder -fix ./...
```

## What it catches

```go
// Bad — cases are out of order
switch fruit {
case "orange":
    // ...
case "apple": // want: case "apple" should come before "orange"
    // ...
case "banana": // want: case "banana" should come before "orange"
    // ...
}

// Good
switch fruit {
case "apple":
    // ...
case "banana":
    // ...
case "orange":
    // ...
}
```

It also handles integers, floats, hex literals, negative numbers, multi-value cases, and `fallthrough` chains.

## Use cases

**Code review enforcement** — run in CI to catch unordered switches before they merge:

```sh
caseorder ./... || exit 1
```

**Auto-fix on save** — pipe through `gofmt`-style tooling or wire into your editor's on-save hook:

```sh
caseorder -fix ./...
```

**Pre-commit hook** — add to `.git/hooks/pre-commit`:

```sh
#!/bin/sh
caseorder ./...
```

**Ad-hoc audit** — check a single file or subtree when reviewing unfamiliar code:

```sh
caseorder ./pkg/config/...
```

**Descending order** — for switches where largest-first is conventional (e.g. priority levels, HTTP status codes):

```sh
caseorder -order=desc ./...
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-order` | `asc` | Sort direction: `asc` or `desc` |
| `-ignore-case` | `true` | Case-insensitive string comparison |
| `-default-last` | `true` | Always place the `default` case last |
| `-autofix` | `true` | Emit suggested fixes (applied with `-fix`) |
| `-autofix-allow-fallthrough` | `false` | Also emit fixes for switches that use `fallthrough` |

## Features

### Numeric ordering

Integers, floats, hex literals, and negative numbers are compared by value:

```go
// Bad
switch code {
case 0xFF:
case 0x0A: // out of order
case -1:   // out of order
}

// Good
switch code {
case -1:
case 0x0A:
case 0xFF:
}
```

### Multi-value cases

Values within a single `case` clause are sorted too:

```go
// Bad
case 3, 1, 2:

// Good
case 1, 2, 3:
```

### Fallthrough groups

Cases connected by `fallthrough` are treated as a single unit and sorted together, preserving their internal order:

```go
// Bad
switch x {
case "zebra":
    fallthrough
case "yacht":
    fallthrough
case "venus":
case "apple": // out of order relative to the "zebra" group
}

// Good
switch x {
case "apple":
case "zebra":
    fallthrough
case "yacht":
    fallthrough
case "venus":
}
```

By default, suggested fixes are not emitted for switches containing `fallthrough`. Pass `-autofix-allow-fallthrough` to enable them.

## License

MIT
