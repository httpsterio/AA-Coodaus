# AA-koodaus — Roadmap

## Milestone 1 — Interpreter ✅

- [x] Lexer, parser, tree-walking interpreter
- [x] Variables, arithmetic, booleans
- [x] Sano, Kysy
- [x] Jos/Muuten, Kunnes, Tois, Loputon
- [x] Teko/Tee/Anna (functions with params and return)
- [x] Lists (literal, Lisää, Pituus, Joka, index access)
- [x] String interpolation and escape sequences
- [x] Comments, shebang
- [x] Finnish error messages
- [x] Fix: Anna at top level causes panic instead of clean error
- [x] Unit tests (lexer, interpreter)
- [x] Golden file tests (testdata/)
- [x] README
- [x] .gitignore, Makefile
- [x] examples/ folder
- [ ] GitHub release with cross-compiled binaries

---

## Milestone 2 — File I/O

- [ ] `Lue tiedosto: polku` — read file to string
- [ ] `Kirjoita tiedosto: polku sisalto` — write string to file
- [ ] `Lisää tiedosto: polku sisalto` — append to file
- [ ] `Jos Onko tiedosto: polku` — file exists check
- [ ] `Listaa hakemisto: polku` — returns list of filenames
- [ ] Cross-platform path handling (via Go stdlib)
- [ ] Tests for all file operations

---

## Milestone 3 — String operations ✅

- [x] `Pituus` for strings
- [x] `Pilko` / `Pilk` — split string by separator → list
- [x] `Korvaa` / `Korv` — replace substring → string
- [x] `Sisaltaa` / `Sis` — contains substring → boolean
- [x] `Trimmaa` / `Trim` — trim whitespace → string
- [x] `Isot` — uppercase → string
- [x] `Pienet` — lowercase → string
- [ ] Multiline string literal or heredoc syntax (TBD)
- [x] Tests

---

## Milestone 4 — Static site generator

- [ ] `Muunna markdown: sisalto` — markdown to HTML (bundles Go markdown lib, first external dep)
- [ ] `Tuo: tiedosto.aa` — import/include other `.aa` files
- [ ] Reference SSG implementation in `examples/ssg/`
- [ ] Docs: how to build a site with AA-koodaus

---

## Milestone 5 — VS Code extension

- [ ] TextMate grammar for `.aa` files
- [ ] Keyword, boolean, string, comment, number highlighting
- [ ] Interpolation highlighting inside strings (`{muuttuja}`)
- [ ] Separate repo
- [ ] Distribute as `.vsix` and publish to VS Code marketplace

---

## Backlog / maybe someday

- `Kv,ka` comma modifier system for list operations
- Local scope for functions
- `Ot-yhteys`, `Hkr`, `Blu-yhteys`, `Cr-peli` easter egg blocks
- REPL (`aac` with no arguments)
- Watch mode / dev server for SSG
- Error recovery (continue after first error instead of stopping)
- Standard library split into separate files