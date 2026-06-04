# AA-koodaus

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

### Arithmetic and strings

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

See the `examples/` folder for complete programs including a guessing game.

---

## Tests

```sh
go test ./...
```

27 golden tests in `testdata/`, plus lexer and interpreter unit tests.

---

## Build

```sh
go build -o aac .
```

No external dependencies. Go stdlib only.

---