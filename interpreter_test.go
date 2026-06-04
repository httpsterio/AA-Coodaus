package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func runProgram(t testing.TB, src string) string {
	t.Helper()

	if strings.HasPrefix(src, "#!") {
		idx := strings.Index(src, "\n")
		if idx >= 0 {
			src = "--" + src[2:idx] + src[idx:]
		}
	}

	PreCheckSource(src)

	l := NewLexer(src)
	p := NewParser(l)
	program := p.ParseProgram()

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	env := NewEnv()
	ev := NewEvaluator(env)
	ev.Evaluate(program)

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = old

	return string(out)
}

func runProgramWithInput(t testing.TB, src, input string) string {
	t.Helper()

	if strings.HasPrefix(src, "#!") {
		idx := strings.Index(src, "\n")
		if idx >= 0 {
			src = "--" + src[2:idx] + src[idx:]
		}
	}

	PreCheckSource(src)

	l := NewLexer(src)
	p := NewParser(l)
	program := p.ParseProgram()

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	env := NewEnv()
	ev := NewEvaluator(env)
	ev.reader = bufio.NewReader(readerFromString(input))
	ev.Evaluate(program)

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = old

	return string(out)
}

func readerFromString(s string) *bytes.Reader {
	return bytes.NewReader([]byte(s))
}

func TestVarDeclAndPrint(t *testing.T) {
	out := runProgram(t, `Mu nimi: Aku
Sano: {nimi}`)
	if out != "Aku\n" {
		t.Errorf("got %q, want %q", out, "Aku\n")
	}
}

func TestArithmetic(t *testing.T) {
	out := runProgram(t, `Mu x: 5 + 3
Mu y: x * 2
Mu z: y - 4
Mu w: z / 3
Sano: {x}
Sano: {y}
Sano: {z}
Sano: {w}`)
	want := "8\n16\n12\n4\n"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestModulo(t *testing.T) {
	out := runProgram(t, `Mu a: 10
Mu b: 3
Mu c: a % b
Sano: {c}`)
	if out != "1\n" {
		t.Errorf("got %q, want %q", out, "1\n")
	}
}

func TestNegativeNumbers(t *testing.T) {
	out := runProgram(t, `Mu a: -5
Sano: {a}
Mu b: 10 - -3
Sano: {b}`)
	want := "-5\n13\n"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestBooleanLiteral(t *testing.T) {
	out := runProgram(t, `Mu oikein: kyllä
Mu vaarin: ei
Sano: {oikein}
Sano: {vaarin}`)
	want := "kyllä\nei\n"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestStringConcat(t *testing.T) {
	out := runProgram(t, `Mu etu: Aku
Mu suku: Ankka
Mu koko: etu + " " + suku
Sano: {koko}`)
	if out != "Aku Ankka\n" {
		t.Errorf("got %q, want %q", out, "Aku Ankka\n")
	}
}

func TestStringInterpolation(t *testing.T) {
	out := runProgram(t, `Mu nimi: Aku
Sano: "Hei {nimi}"`)
	if out != "Hei Aku\n" {
		t.Errorf("got %q, want %q", out, "Hei Aku\n")
	}
}

func TestSanoUnquoted(t *testing.T) {
	out := runProgram(t, `Sano: Moi
Sano: "Moi maailma"`)
	want := "Moi\nMoi maailma\n"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestSanoBareVar(t *testing.T) {
	// Sano with a bare word that matches a variable name
	out := runProgram(t, `Mu terve: Moi
Sano: {terve}`)
	if out != "Moi\n" {
		t.Errorf("got %q, want %q", out, "Moi\n")
	}
}

func TestConditionalTrue(t *testing.T) {
	out := runProgram(t, `Mu pisteet: 10
Jos pisteet on 10:
  Sano: Voitit
Loppu`)
	if out != "Voitit\n" {
		t.Errorf("got %q, want %q", out, "Voitit\n")
	}
}

func TestConditionalFalse(t *testing.T) {
	out := runProgram(t, `Mu pisteet: 5
Jos pisteet on 10:
  Sano: Voitit
Loppu`)
	if out != "" {
		t.Errorf("got %q, want empty", out)
	}
}

func TestConditionalElse(t *testing.T) {
	out := runProgram(t, `Mu pisteet: 5
Jos pisteet on 10:
  Sano: Voitit
Muuten:
  Sano: Hävisit
Loppu`)
	if out != "Hävisit\n" {
		t.Errorf("got %q, want %q", out, "Hävisit\n")
	}
}

func TestNestedConditional(t *testing.T) {
	out := runProgram(t, `Mu x: 5
Jos x yli 0:
  Jos x on 5:
    Sano: Viisi
  Muuten:
    Sano: Ei viisi
  Loppu
Loppu`)
	if out != "Viisi\n" {
		t.Errorf("got %q, want %q", out, "Viisi\n")
	}
}

func TestComparisonOn(t *testing.T) {
	out := runProgram(t, `Mu a: 5
Mu b: 5
Jos a on b:
  Sano: Yhtä
Loppu`)
	if out != "Yhtä\n" {
		t.Errorf("got %q, want %q", out, "Yhtä\n")
	}
}

func TestComparisonEi(t *testing.T) {
	out := runProgram(t, `Mu a: 5
Mu b: 3
Jos a ei b:
  Sano: Eri
Loppu`)
	if out != "Eri\n" {
		t.Errorf("got %q, want %q", out, "Eri\n")
	}
}

func TestComparisonYli(t *testing.T) {
	out := runProgram(t, `Mu a: 5
Mu b: 3
Jos a yli b:
  Sano: Yli
Loppu`)
	if out != "Yli\n" {
		t.Errorf("got %q, want %q", out, "Yli\n")
	}
}

func TestComparisonAlle(t *testing.T) {
	out := runProgram(t, `Mu a: 3
Mu b: 5
Jos a alle b:
  Sano: Alle
Loppu`)
	if out != "Alle\n" {
		t.Errorf("got %q, want %q", out, "Alle\n")
	}
}

func TestComparisonYliOn(t *testing.T) {
	out := runProgram(t, `Mu a: 5
Jos a yli-on 5:
  Sano: Yli tai yhtä
Loppu
Jos a yli-on 4:
  Sano: Yli neljän
Loppu`)
	want := "Yli tai yhtä\nYli neljän\n"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestComparisonAlleOn(t *testing.T) {
	out := runProgram(t, `Mu a: 5
Jos a alle-on 5:
  Sano: Alle tai yhtä
Loppu
Jos a alle-on 6:
  Sano: Alle kuuden
Loppu`)
	want := "Alle tai yhtä\nAlle kuuden\n"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestToisLoop(t *testing.T) {
	out := runProgram(t, `Tois 3 i:
  Sano: {i}
Loppu`)
	want := "1\n2\n3\n"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestToisLoopNoVar(t *testing.T) {
	out := runProgram(t, `Tois 3:
  Sano: Olen Aku
Loppu`)
	want := "Olen Aku\nOlen Aku\nOlen Aku\n"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestToisLoopZero(t *testing.T) {
	out := runProgram(t, `Tois 0:
  Sano: Ei näy
Loppu`)
	if out != "" {
		t.Errorf("got %q, want empty", out)
	}
}

func TestKunnesLoop(t *testing.T) {
	out := runProgram(t, `Mu luku: 1
Kunnes luku yli 5:
  Sano: {luku}
  Mu luku: luku + 1
Loppu`)
	want := "1\n2\n3\n4\n5\n"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestKunnesLoopImmediate(t *testing.T) {
	// Condition is immediately true, body never executes
	out := runProgram(t, `Mu x: kyllä
Kunnes x:
  Sano: Ei näy
Loppu`)
	if out != "" {
		t.Errorf("got %q, want empty", out)
	}
}

func TestFunctionCall(t *testing.T) {
	out := runProgram(t, `Teko tervehdi nimi:
  Sano: "Moi {nimi}"
Loppu

Tee tervehdi Aku`)
	if out != "Moi Aku\n" {
		t.Errorf("got %q, want %q", out, "Moi Aku\n")
	}
}

func TestFunctionReturn(t *testing.T) {
	out := runProgram(t, `Teko laske a b:
  Anna a + b
Loppu

Mu tulos: Tee laske 3 5
Sano: {tulos}
Sano: "Tulos on {tulos}"`)
	want := "8\nTulos on 8\n"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestFunctionDefaultReturn(t *testing.T) {
	// Function without Anna returns ei (false)
	out := runProgram(t, `Teko foo:
  Sano: Moi
Loppu

Mu t: Tee foo
Sano: {t}`)
	want := "Moi\nei\n"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestListLiteral(t *testing.T) {
	out := runProgram(t, `Mu lista: [1 2 3]
Mu eka: lista[1]
Sano: {eka}`)
	if out != "1\n" {
		t.Errorf("got %q, want %q", out, "1\n")
	}
}

func TestListIndex(t *testing.T) {
	out := runProgram(t, `Mu lista: [10 20 30]
Mu tok: lista[2]
Sano: {tok}`)
	if out != "20\n" {
		t.Errorf("got %q, want %q", out, "20\n")
	}
}

func TestListAppend(t *testing.T) {
	out := runProgram(t, `Mu lista: [1 2 3]
Lisää lista: 4
Mu n: Pituus lista
Sano: {n}`)
	if out != "4\n" {
		t.Errorf("got %q, want %q", out, "4\n")
	}
}

func TestListLength(t *testing.T) {
	out := runProgram(t, `Mu lista: [1 2 3 4 5]
Mu n: Pituus lista
Sano: {n}`)
	if out != "5\n" {
		t.Errorf("got %q, want %q", out, "5\n")
	}
}

func TestJokaLoop(t *testing.T) {
	out := runProgram(t, `Mu nimet: [Aku Iines Hannu]
Joka nimet nimi:
  Sano: {nimi}
Loppu`)
	want := "Aku\nIines\nHannu\n"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestComments(t *testing.T) {
	out := runProgram(t, `-- tämä on kommentti
Sano: Moi -- kommentti rivin lopussa`)
	if out != "Moi\n" {
		t.Errorf("got %q, want %q", out, "Moi\n")
	}
}

func TestEmptyProgram(t *testing.T) {
	out := runProgram(t, ``)
	if out != "" {
		t.Errorf("got %q, want empty", out)
	}
}

func TestReassignVariable(t *testing.T) {
	out := runProgram(t, `Mu x: 5
Sano: {x}
Mu x: x + 1
Sano: {x}`)
	want := "5\n6\n"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestSelfReference(t *testing.T) {
	out := runProgram(t, `Mu luku: 1
Mu luku: luku + 1
Sano: {luku}`)
	if out != "2\n" {
		t.Errorf("got %q, want %q", out, "2\n")
	}
}

func TestKysyInput(t *testing.T) {
	src := `Kysy nimi: Mikä on nimesi?
Sano: "Hei {nimi}"`
	out := runProgramWithInput(t, src, "Aku\n")
	if out != "Mikä on nimesi?Hei Aku\n" {
		t.Errorf("got %q, want %q", out, "Mikä on nimesi?Hei Aku\n")
	}
}

func TestKysyNumericInput(t *testing.T) {
	src := `Kysy luku: Anna luku:
Mu tulos: luku + 5
Sano: {tulos}`
	out := runProgramWithInput(t, src, "10\n")
	if out != "Anna luku:15\n" {
		t.Errorf("got %q, want %q", out, "Anna luku:15\n")
	}
}

func TestToisLoopVarShadowing(t *testing.T) {
	// Loop variable shouldn't leak after the loop
	out := runProgram(t, `Mu i: 99
Sano: "ennen {i}"
Tois 3 i:
  Sano: {i}
Loppu
Sano: "jälkeen {i}"`)
	if out != "ennen 99\n1\n2\n3\njälkeen 99\n" {
		t.Errorf("got %q, want %q", out, "ennen 99\n1\n2\n3\njälkeen 99\n")
	}
}

func TestJokaLoopVarShadowing(t *testing.T) {
	out := runProgram(t, `Mu alkio: vanha
Mu lista: [x y z]
Joka lista alkio:
  Sano: {alkio}
Loppu
Sano: {alkio}`)
	want := "x\ny\nz\nvanha\n"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestStringEscapingInSano(t *testing.T) {
	out := runProgram(t, `Sano: "rivi1\nrivi2\tvälissä\\peruutus\"laina"`)
	// \n, \t, \\, \" should be interpreted
	want := "rivi1\nrivi2\tvälissä\\peruutus\"laina\n"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestDivNonIntResult(t *testing.T) {
	out := runProgram(t, `Mu x: 10 / 4
Sano: {x}`)
	if out != "2.5\n" {
		t.Errorf("got %q, want %q", out, "2.5\n")
	}
}

func TestDivIntResult(t *testing.T) {
	out := runProgram(t, `Mu x: 8 / 4
Sano: {x}`)
	if out != "2\n" {
		t.Errorf("got %q, want %q", out, "2\n")
	}
}

func TestMultipleStatementsOnSeparateLines(t *testing.T) {
	out := runProgram(t, `Mu a: 1
Mu b: 2
Mu c: a + b
Sano: "tulos {c}"`)
	if out != "tulos 3\n" {
		t.Errorf("got %q, want %q", out, "tulos 3\n")
	}
}

func TestPituusOnString(t *testing.T) {
	out := runProgram(t, `Mu s: "moi"
Mu n: Pituus s
Sano: {n}`)
	if out != "3\n" {
		t.Errorf("got %q, want %q", out, "3\n")
	}
}

func TestBoolInConditionDirect(t *testing.T) {
	out := runProgram(t, `Jos kyllä:
  Sano: Toimii
Loppu`)
	if out != "Toimii\n" {
		t.Errorf("got %q, want %q", out, "Toimii\n")
	}
}

func TestNestedFunctionCall(t *testing.T) {
	out := runProgram(t, `Teko a:
  Anna 5
Loppu

Teko b:
  Mu t: Tee a
  Anna t + 3
Loppu

Mu x: Tee b
Sano: {x}`)
	if out != "8\n" {
		t.Errorf("got %q, want %q", out, "8\n")
	}
}

func TestFunctionParamRestore(t *testing.T) {
	// Global variable shadowed by function parameter should be restored
	out := runProgram(t, `Mu nimi: globaali
Teko testi nimi:
  Sano: {nimi}
Loppu

Tee testi paikallinen
Sano: {nimi}`)
	want := "paikallinen\nglobaali\n"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestToisLoopPreservesEnv(t *testing.T) {
	out := runProgram(t, `Tois 3 i:
  Mu x: i
Loppu
Mu result: x + 0
Sano: {result}`)
	if out != "3\n" {
		t.Errorf("got %q, want %q", out, "3\n")
	}
}
