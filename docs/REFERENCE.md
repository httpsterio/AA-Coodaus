# AA-koodaus — Keyword reference

Complete keyword reference. All keywords are case-sensitive, first letter capitalized.

**Syntax notation:** `Muuttuja` / `arvo` / `polku` etc. are placeholders for actual values.
A trailing `:` starts a block (closed with `Loppu`).
Modifier words like `tiedosto` / `hakemisto` are optional — commands work with or without them.

---

## Variables

| Keyword | Alias | Syntax | Notes |
|---|---|---|---|
| `Mu` | — | `Mu nimi: arvo` | Declare or reassign. RHS may reference same variable. |

---

## Output and input

| Keyword | Alias | Syntax | Notes |
|---|---|---|---|
| `Sano` | — | `Sano: arvo` | Print to stdout. Supports `{muuttuja}` interpolation. |
| `Kysy` | — | `Kysy nimi: kysymys` | Prompt user, store result in variable. |

---

## Conditionals

| Keyword | Alias | Syntax | Notes |
|---|---|---|---|
| `Jos` | — | `Jos a operaattori b:` | If. Closed with `Loppu`. |
| `Muuten` | — | `Muuten:` | Else. Optional. Inside `Jos` block. |
| `Loppu` | — | `Loppu` | Closes any block. |

### Comparison operators

| Operator | Meaning |
|---|---|
| `on` | equals |
| `ei` | not equals |
| `yli` | greater than |
| `alle` | less than |
| `yli-on` | greater than or equal |
| `alle-on` | less than or equal |

---

## Loops

| Keyword | Alias | Syntax | Notes |
|---|---|---|---|
| `Tois` | — | `Tois n:` or `Tois n i:` | Counted loop. `i` is 1-based loop variable. |
| `Kunnes` | — | `Kunnes ehto:` | While loop. |
| `Loputon` | — | `Loputon:` | Infinite loop. |
| `Joka` | — | `Joka lista alkio:` | Iterate over list. |

---

## Functions

| Keyword | Alias | Syntax | Notes |
|---|---|---|---|
| `Teko` | — | `Teko nimi params:` | Define function. |
| `Tee` | — | `Tee nimi args` | Call function. |
| `Anna` | — | `Anna arvo` | Return value. Only valid inside `Teko`. |

---

## Lists

| Keyword | Alias | Syntax | Returns | Notes |
|---|---|---|---|---|
| — | — | `Mu lista: [1 2 3]` | list | List literal. Space-separated. |
| — | — | `lista[n]` | value | Index access. 1-based. |
| `Lisää` | — | `Lisää lista: arvo` | — | Append to list. |
| `Pituus` | — | `Pituus lista` | number | Length of list or string. |

---

## Strings

| Keyword | Alias | Syntax | Returns |
|---|---|---|---|
| `Pilko` | `Pilk` | `Pilko str sep` | list |
| `Korvaa` | `Korv` | `Korvaa str vanha uusi` | string |
| `Sisaltaa` | `Sis` | `Sisaltaa str osa` | boolean |
| `Trimmaa` | `Trim` | `Trimmaa str` | string |
| `Isot` | — | `Isot str` | string |
| `Pienet` | — | `Pienet str` | string |

String literals: bare single-word or quoted multi-word. Escape sequences: `\n`, `\t`, `\\`, `\"`.
Interpolation: `{muuttuja}` inside strings.

---

## File I/O

| Keyword | Alias | Syntax | Returns | Notes |
|---|---|---|---|---|
| `Lue tiedosto` | `Lue` | `Lue tiedosto "polku"` | string | Read file. Trailing newline stripped. |
| `Kirjoita tiedosto` | `Kirj` | `Kirjoita tiedosto "polku" sisalto` | — | Create file. Fails if exists. |
| `Ylikirjoita tiedosto` | `Ylikirj` | `Ylikirjoita tiedosto "polku" sisalto` | — | Overwrite file. Creates if not exists. |
| `Liitä tiedosto` | `Liit` | `Liitä tiedosto "polku" sisalto` | — | Append to file. Fails if not exists. |
| `Listaa hakemisto` | `Hak` | `Listaa hakemisto "polku"` | list | List directory. Flat, filenames only. |
| `Luo hakemisto` | `Luohak` | `Luo hakemisto "polku"` | — | Create directory. No-op if exists. |
| `Onko tiedosto` | — | `Onko tiedosto "polku"` | boolean | File exists check. |
| `Onko hakemisto` | — | `Onko hakemisto "polku"` | boolean | Directory exists check. |

---

## Types

| Type | Values | Notes |
|---|---|---|
| number | `5`, `3.14` | Integer or float. |
| string | `Aku`, `"moi maailma"` | Bare single-word or quoted. |
| boolean | `kyllä`, `ei` | Distinct type. Never equal to numbers or strings. |
| list | `[1 2 3]` | 1-based indexing. |