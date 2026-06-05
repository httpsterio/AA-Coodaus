# AA-coodaus — Roadmap

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
- [x] GitHub release with cross-compiled binaries

---

## Milestone 2 — File I/O ✅

- [x] `Lue` / `Lue tiedosto: polku` — read file to string
- [x] `Kirjoita` / `Kirjoita tiedosto: polku sisalto` — write new file (errors if exists)
- [x] `Ylikirjoita` / `Ylikirj` / `Ylikirjoita tiedosto: polku sisalto` — overwrite existing file
- [x] `Liitä` / `Liit` / `Liitä tiedosto: polku sisalto` — append to file
- [x] `Onko tiedosto: polku` / `Onko hakemisto: polku` — existence and type check
- [x] `Listaa hakemisto: polku` — returns list of filenames
- [x] `Hak: polku` — shorthand alias for `Listaa hakemisto`
- [x] `Luo hakemisto: polku` / `Luohak: polku` — create directory (mkdir -p)
- [x] Cross-platform path handling (via Go stdlib)
- [x] Tests for all file operations (13 unit tests with t.TempDir())
- [x] Shorthand/alias forms without modifier word

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

## Milestone 4 — Satunnainen ja Lopeta ✅

- [x] `Satunnainen` / `Satu` — random integer or float
- [x] `Lopeta` / `Lopt` — break out of innermost loop
- [x] Doc fixes: modulo, Pituus for strings, Satunnainen and Lopeta in reference.md

---

## Milestone 5 — Language polish

- [ ] Fix `Onko` parser precedence — `Jos Onko tiedosto polku on ei:` should parse correctly without temp variable
- [ ] Indexed list assignment — `Mu lista[i]: arvo`
- [ ] Tests for both

---

## Milestone 6 — Math

- [ ] `Tasaa arvo desimaalit` — round to N decimal places
- [ ] `Lattia arvo` — floor
- [ ] `Katto arvo` — ceil

---

## Milestone 7 — Static site generator

- [ ] `Muunna markdown: sisalto` — markdown to HTML (first external dep)
- [ ] `Tuo: tiedosto.aa` — import/include other `.aa` files
- [ ] Reference SSG implementation in `examples/ssg/`
- [ ] Docs: how to build a site with AA-coodaus

---

## Milestone 8 — VS Code extension

- [ ] TextMate grammar for `.aa` files
- [ ] Keyword, boolean, string, comment, number highlighting
- [ ] Interpolation highlighting inside strings (`{muuttuja}`)
- [ ] Same repo, `plugins/vscode-aacoodaus/` subfolder
- [ ] Distribute as `.vsix` and publish to VS Code marketplace
- [ ] GitHub Actions: build and attach `.vsix` to release

---

## Backlog / maybe someday

- `Kv,ka` comma modifier system for list operations
- Local scope for functions
- `Ot-yhteys`, `Hkr`, `Blu-yhteys`, `Cr-peli` easter egg blocks
- REPL (`aac` with no arguments)
- Watch mode / dev server for SSG
- Error recovery (continue after first error instead of stopping)
- Standard library split into separate files
- Multiline strings / heredoc syntax