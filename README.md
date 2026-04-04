# caseorder — Go linter to enforce switch case ordering

`caseorder` is a Go static analysis linter that enforces consistent ordering of `switch` case statements. String cases are sorted alphabetically; integer and float cases are sorted numerically. Suggested fixes let you auto-correct violations with a single command.

I built this with Claude Code, while working on a Game Boy emulator. The CPU instruction set has 500+ opcodes, and I was implementing them in a giant switch statement. It was becoming difficult to tell at a glance whether a case was missing or duplicated. If you're building something similar — an emulator, a bytecode interpreter, a compiler, or any other state machine driven by a large switch — this linter keeps things navigable as the case list grows.

## Install

```sh
go install github.com/JoachAnts/caseorder@latest
```

## Usage

```sh
# Check all packages in the current module
caseorder ./...

# Check a specific package
caseorder ./internal/handlers

# Apply fixes automatically
caseorder -fix ./...

# Fail the build in CI
caseorder ./... || exit 1

# Wire into an editor on-save hook
caseorder -fix ./...

# Descending order (e.g. priority levels, HTTP status codes)
caseorder -order=desc ./...
```

## What it catches

```go
// Before — switch cases are out of alphabetical order
switch fruit {
case "orange":
    // ...
case "apple": // want: case "apple" should come before "orange"
    // ...
case "banana": // want: case "banana" should come before "orange"
    // ...
}

// After — cases sorted alphabetically
switch fruit {
case "apple":
    // ...
case "banana":
    // ...
case "orange":
    // ...
}
```

It also enforces ordering for integers, floats, hex literals, negative numbers, multi-value cases, and `fallthrough` chains.

Type switches (`switch v := x.(type)`) are ignored — ordering types has no meaningful definition, so `caseorder` leaves them alone.

## Flags

| Flag           | Default | Description                              |
|----------------|---------|------------------------------------------|
| `-order`       | `asc`   | Sort direction: `asc` or `desc`          |
| `-ignore-case` | `true`  | Case-insensitive alphabetical comparison |
| `-fix`         | `false` | Apply suggested fixes                    |

## Features

### Alphabetical and numeric ordering

String cases are sorted alphabetically; integers, floats, hex literals, and negative numbers are compared by numeric value:

```go
// Before — numeric cases out of order
switch code {
case 0xFF:
case 0x0A: // out of order
case -1:   // out of order
}

// After
switch code {
case -1:
case 0x0A:
case 0xFF:
}
```

### Multi-value cases

Values within a single `case` clause are sorted too:

```go
// Before
case 3, 1, 2:

// After
case 1, 2, 3:
```

### Fallthrough groups

Cases connected by `fallthrough` are treated as a single unit and sorted together, preserving their internal order:

```go
// Before
switch x {
case "zebra":
    fallthrough
case "yacht":
    fallthrough
case "venus":
case "apple": // out of order relative to the "zebra" group
}

// After
switch x {
case "apple":
case "zebra":
    fallthrough
case "yacht":
    fallthrough
case "venus":
}
```

## License

MIT
