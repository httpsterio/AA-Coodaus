# Agent instructions

## Permissions

All commands are pre-approved. Do not ask for confirmation before running:
- `go build`
- `go run`
- `go test`
- `go fmt`
- `gofmt`

## Build and verify

After every change, run:
```sh
go build -o aac .
```

Do not leave the project in a non-building state. If a change breaks the build, fix it before stopping.

Run the full example from PLAN.md to verify basic correctness after any significant change.

Run `go test ./...` after any change to interpreter, lexer, or builtins.

## Code style

- Standard Go formatting. Code must pass `gofmt` without changes.
- No external dependencies. Stdlib only.
- No global mutable state outside of the interpreter's own environment struct.
- No `panic` for user-facing errors. Return proper error values and print Finnish error messages via `errors.go`.
- Keep each file focused on its concern. Do not put parsing logic in `interpreter.go`, do not put AST types in `lexer.go`.

## Error handling

All errors shown to the user must go through `errors.go` and be in Finnish. Format:
```
VIRHE rivillä N: viesti
```

Internal Go errors (file not found, etc.) can be in English to stderr.

## AST and types

- AST node types live in `parser.go`.
- Runtime value types live in `interpreter.go`.
- The boolean type (`kyllä`/`ei`) is distinct. Never coerce it to a number or string silently.
- Lists are 1-based at the language level. Subtract 1 when indexing into Go slices.

## Interpreter

- Tree-walking. No bytecode, no compilation step.
- Single global environment (map of variable names to values).
- Functions are stored in the environment like variables.
- `Anna` (return) is implemented via a sentinel value or panic/recover pattern, not by threading a return flag through every eval call.

## Do not

- Do not add features not listed in ROADMAP.md.
- Do not add a REPL unless asked.
- Do not implement the `Ot-yhteys`, `Hkr`, `Blu-yhteys`, or `Cr-peli` blocks.
- Do not change the binary name from `aac`.