# Task

Build a complete, working interpreter for a programming language called AA-koodaus. Deliver all source files ready to build.

---

## Language overview

AA-koodaus is a Finnish meme programming language. Its syntax is abbreviated Finnish with no word inflection (nominative case only), hyphen-compound commands, and comma-chained modifiers. It must feel internally consistent.

---

## Syntax rules

1. One statement per line.
2. Comments start with `--` and run to end of line.
3. Blocks are indented with 2 spaces and closed with `Loppu`.
4. Commands are hyphen-compound words: `Ot-yhteys`, `Cr-peli`.
5. Modifiers chain with commas: `Kv,ka`.
6. Values follow a colon: `Peli name: Akun suuri seikkailu`.

---

## Types

- **Number** — integers and floats. `5`, `3.14`
- **String** — bare single-word or quoted multi-word. `Aku`, `"moi maailma"`. Supports escape sequences: `\n`, `\t`, `\\`, `\"`
- **Boolean** — `kyllä` (true) and `ei` (false). Distinct type, never equal to numbers or strings.

---

## Variables

Declare and assign with `Mu`:
```
Mu nimi: Aku
Mu pisteet: 0
Mu voitto: ei
```

Reassign the same way. Types are inferred. RHS may reference the variable being assigned:
```
Mu pisteet: pisteet + 1
```

All variables are global scope.

---

## Arithmetic

Operators: `+`, `-`, `*`, `/`, `%`
```
Mu x: 5 + 3
Mu y: x * 2
```

`+` also concatenates strings.

---

## Output

```
Sano: Moi
Sano: {nimi}
Sano: "Moi {nimi} pisteet on {pisteet}"
```

`{}` interpolation works inside quoted strings and bare single-word values. Bare word that matches a variable name is printed as its value.

---

## Input

```
Kysy nimi: Mikä on nimesi?
```

Prompts user with the text after the colon, stores result in variable `nimi` as a string.

---

## Conditionals

```
Jos pisteet on 10:
  Sano: Voitit
Muuten:
  Sano: Hävisit
Loppu
```

`Muuten` is optional. Nesting is allowed. `Loppu` closes the block.

Comparison operators:
- `on` — equals
- `ei` — not equals
- `yli` — greater than
- `alle` — less than
- `yli-on` — greater than or equal
- `alle-on` — less than or equal

---

## Loops

**Counted:**
```
Tois 5:
  Sano: Olen Aku
Loppu
```

**Counted with loop variable (1-based):**
```
Tois 5 i:
  Sano: {i}
Loppu
```

**Infinite:**
```
Loputon:
  Sano: Olen Aku
Loppu
```

**Conditional (while):**
```
Kunnes pisteet yli 100:
  Mu pisteet: pisteet + 1
Loppu
```

---

## Functions

Define with `Teko`, call with `Tee`:
```
Teko tervehdys:
  Sano: Moi maailma
Loppu

Tee tervehdys
```

With parameters:
```
Teko tervehdi nimi:
  Sano: "Moi {nimi}"
Loppu

Tee tervehdi Aku
```

Return value with `Anna`:
```
Teko laske a b:
  Anna a + b
Loppu

Mu tulos: Tee laske 3 5
Sano: {tulos}
```

Calling an undefined function is a runtime error. Functions share global scope.

---

## Lists

```
Mu lista: [1 2 3 4 5]
Mu nimet: [Aku Iines Hannu]
```

Access by index (1-based):
```
Mu eka: lista[1]
```

Append with `Lisää`:
```
Lisää lista: 6
```

Length with `Pituus`:
```
Mu n: Pituus lista
```

Iterate with `Joka`:
```
Joka lista alkio:
  Sano: {alkio}
Loppu
```

---

## Error messages

All errors must be in Finnish and in the voice of the original language. Examples:

```
VIRHE rivillä 4: "Ota yhteys" ei kelpaa. Kirjoita Ot-yhteys.
VIRHE rivillä 7: Muuttujaa "pisteet" ei löydy. Muistitko Mu pisteet: 0?
VIRHE rivillä 12: Loppu puuttuu. Tarkista lohko.
VIRHE rivillä 9: Funktiota "hyppää" ei ole. Muistitko Teko hyppää?
```

---

## Full example

```aa
-- Arvausohjelma

Mu arvo: 7
Mu yritys: 0
Mu oikein: ei

Kysy arvaus: Arvaa luku 1-10:

Kunnes oikein on kyllä:
  Mu yritys: yritys + 1
  Jos arvaus on arvo:
    Mu oikein: kyllä
    Sano: "Oikein! Yrityksiä: {yritys}"
  Muuten:
    Jos arvaus alle arvo:
      Sano: Liian pieni
    Muuten:
      Sano: Liian suuri
    Loppu
    Kysy arvaus: Yritä uudestaan:
  Loppu
Loppu
```

---

## Implementation requirements

- Language: Go
- No external dependencies. Stdlib only.
- Tree-walking interpreter. No compilation step.
- Binary name: `aa`
- Usage: `aa script.aa`
- Shebang support: `#!/usr/bin/env aa`
- `go build -o aa .` must produce a working binary with no errors or warnings.

## File structure

```
aa/
  main.go          -- entry point, args, file loading
  lexer.go         -- tokenizer
  parser.go        -- AST builder
  interpreter.go   -- tree-walking evaluator
  builtins.go      -- Sano, Kysy, Lisää, Pituus
  errors.go        -- Finnish error formatting
```

Deliver all files with full implementations. No placeholders, no TODOs. The code must build and run correctly against the full example above.
