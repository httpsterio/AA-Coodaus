# AA-coodaus

Someone had to do it. AA-coodaus is a real, working programming language with abbreviated
Finnish syntax, hyphen-compound commands, and zero regrets.

It is a complete, working interpreter written in Go.

---

## Quick start

```sh
git clone https://github.com/httpsterio/AA-Coodaus
cd AA-Coodaus
go build -o aac .
```

Hello world:

```aa
Sano: Moi maailma
```

```sh
./aac script.aa
```

Or with shebang:

```aa
#!/usr/bin/env aac
Sano: Moi maailma
```

---

## Syntax

### Variables

```aa
Mu nimi: Aku
Mu pisteet: 0
Mu voitto: ei
```

Types are inferred: numbers, strings, booleans (`kyllä` / `ei`). All variables are global. RHS can reference itself:

```aa
Mu pisteet: pisteet + 1
```

### Output and input

```aa
Sano: "Moi {nimi}, pisteet: {pisteet}"
Kysy nimi: Mikä on nimesi?
```

### Conditionals

```aa
Jos pisteet yli 10:
  Sano: Voitit
Muuten:
  Sano: Hävisit
Loppu
```

Operators: `on`, `ei`, `yli`, `alle`, `yli-on`, `alle-on`

### Loops

```aa
Tois 5:
  Sano: Olen Aku
Loppu

Tois 5 i:
  Sano: {i}
Loppu

Kunnes pisteet yli 100:
  Mu pisteet: pisteet + 1
Loppu

Loputon:
  Sano: apua
Loppu
```

### Functions

```aa
Teko tervehdi nimi:
  Sano: "Moi {nimi}"
Loppu

Tee tervehdi Aku

Teko laske a b:
  Anna a + b
Loppu

Mu tulos: Tee laske 3 5
```

### Lists

```aa
Mu lista: [1 2 3 4 5]
Mu eka: lista[1]
Lisää lista: 6
Mu n: Pituus lista

Joka lista alkio:
  Sano: {alkio}
Loppu
```

### String operations

String literals can be bare single-word values or quoted multi-word strings. Escape sequences: `\n`, `\t`, `\\`, `\"`. Interpolation with `{muuttuja}` works inside quoted strings.

| Keyword | Alias | What it does |
|---|---|---|
| `Pilko` | `Pilk` | Split string by separator → list |
| `Korvaa` | `Korv` | Replace substring → string |
| `Sisaltaa` | `Sis` | Check if string contains substring → boolean |
| `Trimmaa` | `Trim` | Trim whitespace → string |
| `Isot` | — | Uppercase → string |
| `Pienet` | — | Lowercase → string |
| `Pituus` | — | Length of string (also works on lists) |

```aa
Mu osat: Pilk "yksi kaksi kolme" " "
Sano: "{Pituus osat}"   -- 3
Mu korvattu: Korv "aabbcc" "bb" "XX"
Sano: "{korvattu}"       -- aaXXcc
Mu onko: Sis "moi maailma" "maa"
Sano: "{onko}"           -- kyllä
Mu iso: Isot "ääkköset"
Sano: "{iso}"            -- ÄÄKKÖSET
```

### File I/O

| Keyword | Alias | Syntax | Returns | Notes |
|---|---|---|---|---|
| `Lue tiedosto` | `Lue` | `Lue tiedosto "polku"` | string | Read file. Trailing newline stripped. |
| `Kirjoita tiedosto` | `Kirj` | `Kirjoita tiedosto "polku" sisalto` | — | Create file. Fails if exists. |
| `Ylikirjoita tiedosto` | `Ylikirj` | `Ylikirjoita tiedosto "polku" sisalto` | — | Overwrite file. Creates if missing. |
| `Liitä tiedosto` | `Liit` | `Liitä tiedosto "polku" sisalto` | — | Append to file. Fails if missing. |
| `Listaa hakemisto` | `Hak` | `Listaa hakemisto "polku"` | list | List directory contents (flat). |
| `Luo hakemisto` | `Luohak` | `Luo hakemisto "polku"` | — | Create directory (mkdir -p). |
| `Onko tiedosto` | — | `Onko tiedosto "polku"` | boolean | File exists. |
| `Onko hakemisto` | — | `Onko hakemisto "polku"` | boolean | Directory exists. |

The modifier words (`tiedosto`, `hakemisto`) are optional — short forms work without them:

```aa
Mu sisalto: Lue "tiedosto.txt"
Kirjoita "uusi.txt" "terve"
Ylikirj "vanha.txt" "uusisisalto"
Liit "loki.txt" "uusi rivi"
Hak "/tmp"
Luohak "kansio"
```

### Arithmetic

```aa
Mu x: 5 + 3
Mu y: x * 2
Mu z: "Moi " + nimi
```

### Comments

```aa
-- tämä on kommentti
```

---

## Examples

See the `examples/` folder for complete programs:

| File | What it does |
|---|---|
| `hello.aa` | Hello world with variables, conditionals, and a loop |
| `lista.aa` | List operations: create, index, append, iterate |
| `arvaus.aa` | Number guessing game using Kysy, Jos, loops |
| `csv.aa` | Parse CSV-style data and compute results |
| `merkkijono.aa` | String operations: Trim, Isot, Pienet, Pilko, Korvaa |
| `tiedostot-ja-hakemistot.aa` | File I/O: read, write, list, create directories |

---

## Tests

```sh
go test ./...
```

39 golden tests in `testdata/`, plus lexer and interpreter unit tests.

---

## Build

```sh
go build -o aac .
```

No external dependencies. Go stdlib only.

---