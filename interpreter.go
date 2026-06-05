package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type loopBreak struct{}

// Val represents a runtime value in the interpreter.
type Val interface {
	Type() string
	String() string
}

type NumberVal struct {
	Value float64
	IsInt bool
}

func (v *NumberVal) Type() string { return "Numero" }
func (v *NumberVal) String() string {
	if v.IsInt {
		return fmt.Sprintf("%.0f", v.Value)
	}
	return strconv.FormatFloat(v.Value, 'f', -1, 64)
}

type StringVal struct {
	Value string
}

func (v *StringVal) Type() string   { return "Merkkijono" }
func (v *StringVal) String() string { return v.Value }

type BoolVal struct {
	Value bool
}

func (v *BoolVal) Type() string { return "Totuusarvo" }
func (v *BoolVal) String() string {
	if v.Value {
		return "kyllä"
	}
	return "ei"
}

type ListVal struct {
	Elements []Val
}

func (v *ListVal) Type() string { return "Lista" }
func (v *ListVal) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, el := range v.Elements {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(el.String())
	}
	sb.WriteString("]")
	return sb.String()
}

type FunctionVal struct {
	Name   string
	Params []string
	Body   []Stmt
}

func (v *FunctionVal) Type() string { return "Funktio" }
func (v *FunctionVal) String() string {
	return fmt.Sprintf("<funktio %s>", v.Name)
}

// Env represents the single global environment.
type Env struct {
	vars map[string]Val
}

func NewEnv() *Env {
	return &Env{vars: make(map[string]Val)}
}

func (e *Env) Get(name string, line int) Val {
	v, ok := e.vars[name]
	if !ok {
		ShowError(line, "Muuttujaa \"%s\" ei löydy. Muistitko Mu %s: 0?", name, name)
	}
	return v
}

func (e *Env) Set(name string, val Val) {
	e.vars[name] = val
}

type returnValue struct {
	val  Val
	line int
}

// Evaluator walks the AST and executes it.
type Evaluator struct {
	env       *Env
	reader    *bufio.Reader
	loopDepth int
}

func NewEvaluator(env *Env) *Evaluator {
	return &Evaluator{
		env:    env,
		reader: bufio.NewReader(os.Stdin),
	}
}

// Interpolate parses a string template and replaces `{var}` with variable values.
func Interpolate(s string, env *Env, line int) string {
	var sb strings.Builder
	runes := []rune(s)
	i := 0
	for i < len(runes) {
		if runes[i] == '{' {
			start := i + 1
			end := start
			for end < len(runes) && runes[end] != '}' {
				end++
			}
			if end >= len(runes) {
				ShowError(line, "Sulkeva aaltosulku '}' puuttuu merkkijonon interpoloinnissa.")
			}
			varName := string(runes[start:end])
			val, exists := env.vars[varName]
			if !exists {
				ShowError(line, "Muuttujaa \"%s\" ei löydy. Muistitko Mu %s: 0?", varName, varName)
			}
			sb.WriteString(val.String())
			i = end + 1
		} else {
			sb.WriteRune(runes[i])
			i++
		}
	}
	return sb.String()
}

func (e *Evaluator) Evaluate(program []Stmt) {
	defer func() {
		if r := recover(); r != nil {
			if ret, ok := r.(returnValue); ok {
				ShowError(ret.line, "Anna toimii vain funktion sisällä.")
			}
			panic(r)
		}
	}()
	e.evalStatements(program)
}

func (e *Evaluator) evalStatements(stmts []Stmt) {
	for _, stmt := range stmts {
		e.evalStatement(stmt)
	}
}

func (e *Evaluator) evalStatement(stmt Stmt) {
	switch s := stmt.(type) {
	case *VarDeclStmt:
		val := e.evalExpr(s.Val, true)
		e.env.Set(s.Name, val)
	case *SanoStmt:
		e.evalSanoStmt(s)
	case *KysyStmt:
		e.evalKysyStmt(s)
	case *JosStmt:
		e.evalJosStmt(s)
	case *ToisStmt:
		e.evalToisStmt(s)
	case *LoputonStmt:
		e.evalLoputonStmt(s)
	case *KunnesStmt:
		e.evalKunnesStmt(s)
	case *TekoStmt:
		e.env.Set(s.Name, &FunctionVal{
			Name:   s.Name,
			Params: s.Params,
			Body:   s.Body,
		})
	case *TeeStmt:
		e.evalExpr(s.Call, false)
	case *AnnaStmt:
		val := e.evalExpr(s.Val, true)
		panic(returnValue{val: val, line: s.Line})
	case *LisaaStmt:
		e.evalLisaaStmt(s)
	case *JokaStmt:
		e.evalJokaStmt(s)
	case *KirjoitaStmt:
		e.evalKirjoitaStmt(s)
	case *YlikirjoitaStmt:
		e.evalYlikirjoitaStmt(s)
	case *LiitaStmt:
		e.evalLiitaStmt(s)
	case *LuoHakemistoStmt:
		e.evalLuoHakemistoStmt(s)
	case *LopetaStmt:
		if e.loopDepth == 0 {
			ShowError(s.Line, "Lopeta toimii vain silmukan sisällä.")
		}
		panic(loopBreak{})
	default:
		ShowError(stmt.LineNumber(), "Tuntematon lauseketyyppi suorituksessa.")
	}
}

func (e *Evaluator) evalSanoStmt(stmt *SanoStmt) {
	content := stmt.Content

	// 1. Quoted string
	if len(content) >= 2 && content[0] == '"' && content[len(content)-1] == '"' {
		strVal := content[1 : len(content)-1]
		fmt.Println(Interpolate(strVal, e.env, stmt.Line))
		return
	}

	// 2. Single word matching variable name
	if !strings.ContainsAny(content, " \t\r\n") {
		if val, exists := e.env.vars[content]; exists {
			fmt.Println(val.String())
			return
		}
	}

	// 3. General string with interpolation
	fmt.Println(Interpolate(content, e.env, stmt.Line))
}

func (e *Evaluator) evalKysyStmt(stmt *KysyStmt) {
	prompt := stmt.Prompt

	var finalPrompt string
	if len(prompt) >= 2 && prompt[0] == '"' && prompt[len(prompt)-1] == '"' {
		finalPrompt = Interpolate(prompt[1:len(prompt)-1], e.env, stmt.Line)
	} else if !strings.ContainsAny(prompt, " \t\r\n") {
		if val, exists := e.env.vars[prompt]; exists {
			finalPrompt = val.String()
		} else {
			finalPrompt = Interpolate(prompt, e.env, stmt.Line)
		}
	} else {
		finalPrompt = Interpolate(prompt, e.env, stmt.Line)
	}

	fmt.Print(finalPrompt)

	lineText, err := e.reader.ReadString('\n')
	if err != nil {
		lineText = ""
	}
	lineText = strings.TrimSuffix(lineText, "\n")
	lineText = strings.TrimSuffix(lineText, "\r")

	// Try auto-parsing numeric inputs or boolean inputs to avoid type errors in comparison/arithmetic
	trimmed := strings.TrimSpace(lineText)
	if floatVal, err := strconv.ParseFloat(trimmed, 64); err == nil {
		isInt := !strings.Contains(trimmed, ".")
		e.env.Set(stmt.Name, &NumberVal{Value: floatVal, IsInt: isInt})
	} else if trimmed == "kyllä" {
		e.env.Set(stmt.Name, &BoolVal{Value: true})
	} else if trimmed == "ei" {
		e.env.Set(stmt.Name, &BoolVal{Value: false})
	} else {
		e.env.Set(stmt.Name, &StringVal{Value: lineText})
	}
}

func (e *Evaluator) evalJosStmt(stmt *JosStmt) {
	val := e.evalExpr(stmt.Cond, false)
	boolVal, ok := val.(*BoolVal)
	if !ok {
		ShowError(stmt.Line, "Jos-lauseen ehto pitää olla totuusarvo, saatiin %s.", val.Type())
	}
	if boolVal.Value {
		e.evalStatements(stmt.ThenBody)
	} else if stmt.ElseBody != nil {
		e.evalStatements(stmt.ElseBody)
	}
}

func (e *Evaluator) evalToisStmt(stmt *ToisStmt) {
	val := e.evalExpr(stmt.Count, false)
	numVal, ok := val.(*NumberVal)
	if !ok {
		ShowError(stmt.Line, "Tois-silmukka vaatii lukuarvon, saatiin %s.", val.Type())
	}

	count := int(numVal.Value)
	if count <= 0 {
		return
	}

	e.loopDepth++
	defer func() { e.loopDepth-- }()
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(loopBreak); ok {
				return
			}
			panic(r)
		}
	}()

	if stmt.VarName != "" {
		oldVal, exists := e.env.vars[stmt.VarName]
		for i := 1; i <= count; i++ {
			e.env.Set(stmt.VarName, &NumberVal{Value: float64(i), IsInt: true})
			e.evalStatements(stmt.Body)
		}
		if exists {
			e.env.Set(stmt.VarName, oldVal)
		} else {
			delete(e.env.vars, stmt.VarName)
		}
	} else {
		for i := 0; i < count; i++ {
			e.evalStatements(stmt.Body)
		}
	}
}

func (e *Evaluator) evalLoputonStmt(stmt *LoputonStmt) {
	e.loopDepth++
	defer func() { e.loopDepth-- }()
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(loopBreak); ok {
				return
			}
			panic(r)
		}
	}()
	for {
		e.evalStatements(stmt.Body)
	}
}

func (e *Evaluator) evalKunnesStmt(stmt *KunnesStmt) {
	e.loopDepth++
	defer func() { e.loopDepth-- }()
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(loopBreak); ok {
				return
			}
			panic(r)
		}
	}()
	for {
		val := e.evalExpr(stmt.Cond, false)
		boolVal, ok := val.(*BoolVal)
		if !ok {
			ShowError(stmt.Line, "Kunnes-silmukan ehto pitää olla totuusarvo, saatiin %s.", val.Type())
		}
		if boolVal.Value {
			break
		}
		e.evalStatements(stmt.Body)
	}
}

func (e *Evaluator) evalLisaaStmt(stmt *LisaaStmt) {
	val, exists := e.env.vars[stmt.ListName]
	if !exists {
		ShowError(stmt.Line, "Muuttujaa \"%s\" ei löydy. Muistitko Mu %s: 0?", stmt.ListName, stmt.ListName)
	}
	listVal, ok := val.(*ListVal)
	if !ok {
		ShowError(stmt.Line, "Yritettiin lisätä tyyppiin %s, vaaditaan lista.", val.Type())
	}
	appendVal := e.evalExpr(stmt.Val, true)
	listVal.Elements = append(listVal.Elements, appendVal)
}

func (e *Evaluator) evalJokaStmt(stmt *JokaStmt) {
	val := e.evalExpr(stmt.ListExpr, false)
	listVal, ok := val.(*ListVal)
	if !ok {
		ShowError(stmt.Line, "Joka-silmukka vaatii listan, saatiin %s.", val.Type())
	}

	e.loopDepth++
	defer func() { e.loopDepth-- }()
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(loopBreak); ok {
				return
			}
			panic(r)
		}
	}()

	oldVal, exists := e.env.vars[stmt.VarName]

	for _, el := range listVal.Elements {
		e.env.Set(stmt.VarName, el)
		e.evalStatements(stmt.Body)
	}

	if exists {
		e.env.Set(stmt.VarName, oldVal)
	} else {
		delete(e.env.vars, stmt.VarName)
	}
}

func (e *Evaluator) evalExpr(expr Expr, allowBareString bool) Val {
	switch ex := expr.(type) {
	case *NumberExpr:
		return &NumberVal{Value: ex.Val, IsInt: ex.IsInt}
	case *StringExpr:
		return &StringVal{Value: ex.Val}
	case *BoolExpr:
		return &BoolVal{Value: ex.Val}
	case *VarExpr:
		val, exists := e.env.vars[ex.Name]
		if !exists {
			if allowBareString {
				return &StringVal{Value: ex.Name}
			}
			ShowError(ex.Line, "Muuttujaa \"%s\" ei löydy. Muistitko Mu %s: 0?", ex.Name, ex.Name)
		}
		return val
	case *ListExpr:
		var elements []Val
		for _, elemExpr := range ex.Elements {
			elements = append(elements, e.evalExpr(elemExpr, true))
		}
		return &ListVal{Elements: elements}
	case *ListIndexExpr:
		return e.evalListIndexExpr(ex)
	case *PituusExpr:
		return e.evalPituusExpr(ex)
	case *BinaryExpr:
		return e.evalBinaryExpr(ex)
	case *TeeExpr:
		return e.evalTeeExpr(ex)
	case *PilkoExpr:
		return e.evalPilkoExpr(ex)
	case *KorvaaExpr:
		return e.evalKorvaaExpr(ex)
	case *SisaltaaExpr:
		return e.evalSisaltaaExpr(ex)
	case *TrimmaaExpr:
		return e.evalTrimmaaExpr(ex)
	case *IsotExpr:
		return e.evalIsotExpr(ex)
	case *PienetExpr:
		return e.evalPienetExpr(ex)
	case *LueExpr:
		return e.evalLueExpr(ex)
	case *ListaaExpr:
		return e.evalListaaExpr(ex)
	case *OnkoExpr:
		return e.evalOnkoExpr(ex)
	case *SatunnainenExpr:
		return e.evalSatunnainenExpr(ex)
	default:
		ShowError(expr.LineNumber(), "Tuntematon lauseketyyppi laskennassa.")
		return nil
	}
}

func (e *Evaluator) evalListIndexExpr(expr *ListIndexExpr) Val {
	lVal := e.evalExpr(expr.List, false)
	listVal, ok := lVal.(*ListVal)
	if !ok {
		ShowError(expr.Line, "Yritettiin hakea indeksiä tyypistä %s, vaaditaan lista.", lVal.Type())
	}

	iVal := e.evalExpr(expr.Index, false)
	idxVal, ok := iVal.(*NumberVal)
	if !ok {
		ShowError(expr.Line, "Listan indeksin pitää olla numero, saatiin %s.", iVal.Type())
	}

	idx := int(idxVal.Value)
	sliceIdx := idx - 1

	if sliceIdx < 0 || sliceIdx >= len(listVal.Elements) {
		ShowError(expr.Line, "Listan indeksi %d on rajojen ulkopuolella (listan pituus %d).", idx, len(listVal.Elements))
	}

	return listVal.Elements[sliceIdx]
}

func (e *Evaluator) evalPituusExpr(expr *PituusExpr) Val {
	val := e.evalExpr(expr.List, false)
	switch v := val.(type) {
	case *ListVal:
		return &NumberVal{Value: float64(len(v.Elements)), IsInt: true}
	case *StringVal:
		runes := []rune(v.Value)
		return &NumberVal{Value: float64(len(runes)), IsInt: true}
	default:
		ShowError(expr.Line, "Pituus-komento vaatii listan tai merkkijonon, saatiin %s.", val.Type())
		return nil
	}
}

func (e *Evaluator) evalPilkoExpr(expr *PilkoExpr) Val {
	strVal := e.evalExpr(expr.Str, false)
	sepVal := e.evalExpr(expr.Sep, false)

	str, ok := strVal.(*StringVal)
	if !ok {
		ShowError(expr.Line, "Pilko vaatii merkkijonon, saatiin %s.", strVal.Type())
	}
	sep, ok := sepVal.(*StringVal)
	if !ok {
		ShowError(expr.Line, "Pilko vaatii erottimen merkkijonona, saatiin %s.", sepVal.Type())
	}

	parts := strings.Split(str.Value, sep.Value)
	elements := make([]Val, len(parts))
	for i, p := range parts {
		elements[i] = &StringVal{Value: p}
	}
	return &ListVal{Elements: elements}
}

func (e *Evaluator) evalKorvaaExpr(expr *KorvaaExpr) Val {
	strVal := e.evalExpr(expr.Str, false)
	oldVal := e.evalExpr(expr.Old, false)
	newVal := e.evalExpr(expr.New, false)

	str, ok := strVal.(*StringVal)
	if !ok {
		ShowError(expr.Line, "Korvaa vaatii merkkijonon, saatiin %s.", strVal.Type())
	}
	old, ok := oldVal.(*StringVal)
	if !ok {
		ShowError(expr.Line, "Korvaa vaatii korvattavan merkkijonona, saatiin %s.", oldVal.Type())
	}
	repl, ok := newVal.(*StringVal)
	if !ok {
		ShowError(expr.Line, "Korvaa vaatii korvaavan merkkijonona, saatiin %s.", newVal.Type())
	}

	return &StringVal{Value: strings.ReplaceAll(str.Value, old.Value, repl.Value)}
}

func (e *Evaluator) evalSisaltaaExpr(expr *SisaltaaExpr) Val {
	strVal := e.evalExpr(expr.Str, false)
	subVal := e.evalExpr(expr.Sub, false)

	str, ok := strVal.(*StringVal)
	if !ok {
		ShowError(expr.Line, "Sisaltaa vaatii merkkijonon, saatiin %s.", strVal.Type())
	}
	sub, ok := subVal.(*StringVal)
	if !ok {
		ShowError(expr.Line, "Sisaltaa vaatii osamerkkijonon merkkijonona, saatiin %s.", subVal.Type())
	}

	return &BoolVal{Value: strings.Contains(str.Value, sub.Value)}
}

func (e *Evaluator) evalTrimmaaExpr(expr *TrimmaaExpr) Val {
	strVal := e.evalExpr(expr.Str, false)
	str, ok := strVal.(*StringVal)
	if !ok {
		ShowError(expr.Line, "Trimmaa vaatii merkkijonon, saatiin %s.", strVal.Type())
	}
	return &StringVal{Value: strings.TrimSpace(str.Value)}
}

func (e *Evaluator) evalIsotExpr(expr *IsotExpr) Val {
	strVal := e.evalExpr(expr.Str, false)
	str, ok := strVal.(*StringVal)
	if !ok {
		ShowError(expr.Line, "Isot vaatii merkkijonon, saatiin %s.", strVal.Type())
	}
	return &StringVal{Value: strings.ToUpper(str.Value)}
}

func (e *Evaluator) evalPienetExpr(expr *PienetExpr) Val {
	strVal := e.evalExpr(expr.Str, false)
	str, ok := strVal.(*StringVal)
	if !ok {
		ShowError(expr.Line, "Pienet vaatii merkkijonon, saatiin %s.", strVal.Type())
	}
	return &StringVal{Value: strings.ToLower(str.Value)}
}

func (e *Evaluator) evalLueExpr(expr *LueExpr) Val {
	pathVal := e.evalExpr(expr.Path, false)
	path, ok := pathVal.(*StringVal)
	if !ok {
		ShowError(expr.Line, "Lue vaatii tiedostopolun merkkijonona, saatiin %s.", pathVal.Type())
	}
	data, err := os.ReadFile(path.Value)
	if err != nil {
		ShowError(expr.Line, "Tiedostoa \"%s\" ei voitu lukea: %s", path.Value, err)
	}
	content := strings.TrimRight(string(data), "\n\r")
	return &StringVal{Value: content}
}

func (e *Evaluator) evalListaaExpr(expr *ListaaExpr) Val {
	pathVal := e.evalExpr(expr.Path, false)
	path, ok := pathVal.(*StringVal)
	if !ok {
		ShowError(expr.Line, "Listaa vaatii hakemistopolun merkkijonona, saatiin %s.", pathVal.Type())
	}
	entries, err := os.ReadDir(path.Value)
	if err != nil {
		ShowError(expr.Line, "Hakemistoa \"%s\" ei voitu lukea: %s", path.Value, err)
	}
	elements := make([]Val, len(entries))
	for i, entry := range entries {
		elements[i] = &StringVal{Value: entry.Name()}
	}
	return &ListVal{Elements: elements}
}

func (e *Evaluator) evalOnkoExpr(expr *OnkoExpr) Val {
	pathVal := e.evalExpr(expr.Path, false)
	path, ok := pathVal.(*StringVal)
	if !ok {
		ShowError(expr.Line, "Onko vaatii polun merkkijonona, saatiin %s.", pathVal.Type())
	}
	info, err := os.Stat(path.Value)
	if err != nil {
		return &BoolVal{Value: false}
	}
	if expr.IsFile {
		return &BoolVal{Value: !info.IsDir()}
	}
	return &BoolVal{Value: info.IsDir()}
}

func (e *Evaluator) evalSatunnainenExpr(expr *SatunnainenExpr) Val {
	minVal := e.evalExpr(expr.Min, false)
	maxVal := e.evalExpr(expr.Max, false)

	minNum, ok := minVal.(*NumberVal)
	if !ok {
		ShowError(expr.Line, "Satunnainen vaatii lukuarvon minimiksi, saatiin %s.", minVal.Type())
	}
	maxNum, ok := maxVal.(*NumberVal)
	if !ok {
		ShowError(expr.Line, "Satunnainen vaatii lukuarvon maksimiksi, saatiin %s.", maxVal.Type())
	}

	if minNum.Value > maxNum.Value {
		ShowError(expr.Line, "Satunnainen: minimi ei voi olla suurempi kuin maksimi.")
	}

	min := int(minNum.Value)
	max := int(maxNum.Value)

	if expr.Desimaalit == nil {
		return &NumberVal{Value: float64(rand.Intn(max-min+1) + min), IsInt: true}
	}

	desVal := e.evalExpr(expr.Desimaalit, false)
	desNum, ok := desVal.(*NumberVal)
	if !ok {
		ShowError(expr.Line, "Satunnainen: desimaalien määrän tulee olla luku, saatiin %s.", desVal.Type())
	}

	des := int(desNum.Value)
	if des < 0 || des > 10 {
		ShowError(expr.Line, "Satunnainen: desimaalien määrän tulee olla 0–10.")
	}

	val := minNum.Value + rand.Float64()*(maxNum.Value-minNum.Value)
	pow := math.Pow(10, float64(des))
	rounded := math.Round(val*pow) / pow
	return &NumberVal{Value: rounded, IsInt: des == 0}
}

func (e *Evaluator) evalKirjoitaStmt(stmt *KirjoitaStmt) {
	pathVal := e.evalExpr(stmt.Path, false)
	contentVal := e.evalExpr(stmt.Content, false)
	path, ok := pathVal.(*StringVal)
	if !ok {
		ShowError(stmt.Line, "Kirjoita vaatii tiedostopolun merkkijonona, saatiin %s.", pathVal.Type())
	}
	content, ok := contentVal.(*StringVal)
	if !ok {
		ShowError(stmt.Line, "Kirjoita vaatii sisällön merkkijonona, saatiin %s.", contentVal.Type())
	}
	if _, err := os.Stat(path.Value); err == nil {
		ShowError(stmt.Line, "Tiedosto \"%s\" on jo olemassa. Käytä Ylikirjoita jos haluat korvata.", path.Value)
	}
	if err := os.WriteFile(path.Value, []byte(content.Value), 0644); err != nil {
		ShowError(stmt.Line, "Tiedostoon \"%s\" kirjoittaminen epäonnistui: %s", path.Value, err)
	}
}

func (e *Evaluator) evalYlikirjoitaStmt(stmt *YlikirjoitaStmt) {
	pathVal := e.evalExpr(stmt.Path, false)
	contentVal := e.evalExpr(stmt.Content, false)
	path, ok := pathVal.(*StringVal)
	if !ok {
		ShowError(stmt.Line, "Ylikirjoita vaatii tiedostopolun merkkijonona, saatiin %s.", pathVal.Type())
	}
	content, ok := contentVal.(*StringVal)
	if !ok {
		ShowError(stmt.Line, "Ylikirjoita vaatii sisällön merkkijonona, saatiin %s.", contentVal.Type())
	}
	if err := os.WriteFile(path.Value, []byte(content.Value), 0644); err != nil {
		ShowError(stmt.Line, "Tiedostoon \"%s\" kirjoittaminen epäonnistui: %s", path.Value, err)
	}
}

func (e *Evaluator) evalLiitaStmt(stmt *LiitaStmt) {
	pathVal := e.evalExpr(stmt.Path, false)
	contentVal := e.evalExpr(stmt.Content, false)
	path, ok := pathVal.(*StringVal)
	if !ok {
		ShowError(stmt.Line, "Liitä vaatii tiedostopolun merkkijonona, saatiin %s.", pathVal.Type())
	}
	content, ok := contentVal.(*StringVal)
	if !ok {
		ShowError(stmt.Line, "Liitä vaatii sisällön merkkijonona, saatiin %s.", contentVal.Type())
	}
	f, err := os.OpenFile(path.Value, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		ShowError(stmt.Line, "Tiedostoon \"%s\" liittäminen epäonnistui: %s", path.Value, err)
	}
	defer f.Close()
	if _, err := f.WriteString(content.Value); err != nil {
		ShowError(stmt.Line, "Tiedostoon \"%s\" kirjoittaminen epäonnistui: %s", path.Value, err)
	}
}

func (e *Evaluator) evalLuoHakemistoStmt(stmt *LuoHakemistoStmt) {
	pathVal := e.evalExpr(stmt.Path, false)
	path, ok := pathVal.(*StringVal)
	if !ok {
		ShowError(stmt.Line, "Luo-hakemisto vaatii polun merkkijonona, saatiin %s.", pathVal.Type())
	}
	if err := os.MkdirAll(path.Value, 0755); err != nil {
		ShowError(stmt.Line, "Hakemiston \"%s\" luominen epäonnistui: %s", path.Value, err)
	}
}

func (e *Evaluator) evalTeeExpr(expr *TeeExpr) (result Val) {
	val, exists := e.env.vars[expr.Name]
	if !exists {
		ShowError(expr.Line, "Funktiota \"%s\" ei ole. Muistitko Teko %s?", expr.Name, expr.Name)
	}
	fn, ok := val.(*FunctionVal)
	if !ok {
		ShowError(expr.Line, "Muuttuja \"%s\" ei ole funktio.", expr.Name)
	}

	var args []Val
	for _, argExpr := range expr.Args {
		args = append(args, e.evalExpr(argExpr, true))
	}

	if len(args) != len(fn.Params) {
		ShowError(expr.Line, "Funktio \"%s\" odotti %d parametria, saatiin %d.", fn.Name, len(fn.Params), len(args))
	}

	oldVals := make(map[string]Val)
	for i, param := range fn.Params {
		if v, exists := e.env.vars[param]; exists {
			oldVals[param] = v
		}
		e.env.Set(param, args[i])
	}

	defer func() {
		for _, param := range fn.Params {
			if oldVal, exists := oldVals[param]; exists {
				e.env.Set(param, oldVal)
			} else {
				delete(e.env.vars, param)
			}
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			if ret, ok := r.(returnValue); ok {
				result = ret.val
			} else {
				panic(r)
			}
		}
	}()

	e.evalStatements(fn.Body)
	return &BoolVal{Value: false}
}

func (e *Evaluator) evalBinaryExpr(expr *BinaryExpr) Val {
	leftVal := e.evalExpr(expr.Left, false)
	rightVal := e.evalExpr(expr.Right, false)

	switch expr.Op {
	case TOKEN_PLUS:
		return e.evalPlus(leftVal, rightVal, expr.Line)
	case TOKEN_MINUS:
		return e.evalMinus(leftVal, rightVal, expr.Line)
	case TOKEN_MUL:
		return e.evalMul(leftVal, rightVal, expr.Line)
	case TOKEN_DIV:
		return e.evalDiv(leftVal, rightVal, expr.Line)
	case TOKEN_MOD:
		return e.evalMod(leftVal, rightVal, expr.Line)
	case TOKEN_ON:
		return &BoolVal{Value: e.evalOn(leftVal, rightVal)}
	case TOKEN_EI:
		return &BoolVal{Value: !e.evalOn(leftVal, rightVal)}
	case TOKEN_YLI:
		return &BoolVal{Value: e.evalYli(leftVal, rightVal, expr.Line)}
	case TOKEN_ALLE:
		return &BoolVal{Value: e.evalAlle(leftVal, rightVal, expr.Line)}
	case TOKEN_YLI_ON:
		return &BoolVal{Value: e.evalYliOn(leftVal, rightVal, expr.Line)}
	case TOKEN_ALLE_ON:
		return &BoolVal{Value: e.evalAlleOn(leftVal, rightVal, expr.Line)}
	default:
		ShowError(expr.Line, "Tuntematon operaattori: %s", expr.Op)
		return nil
	}
}

func (e *Evaluator) evalPlus(left, right Val, line int) Val {
	_, leftIsStr := left.(*StringVal)
	_, rightIsStr := right.(*StringVal)

	if leftIsStr || rightIsStr {
		if _, ok := left.(*BoolVal); ok {
			ShowError(line, "Totuusarvoa ei voi yhdistää merkkijonoon.")
		}
		if _, ok := right.(*BoolVal); ok {
			ShowError(line, "Totuusarvoa ei voi yhdistää merkkijonoon.")
		}
		return &StringVal{Value: left.String() + right.String()}
	}

	lNum, lOk := left.(*NumberVal)
	rNum, rOk := right.(*NumberVal)
	if !lOk || !rOk {
		ShowError(line, "Virheelliset tyypit plus-laskutoimituksessa: %s ja %s.", left.Type(), right.Type())
	}

	isInt := lNum.IsInt && rNum.IsInt
	return &NumberVal{Value: lNum.Value + rNum.Value, IsInt: isInt}
}

func (e *Evaluator) evalMinus(left, right Val, line int) Val {
	lNum, lOk := left.(*NumberVal)
	rNum, rOk := right.(*NumberVal)
	if !lOk || !rOk {
		ShowError(line, "Virheelliset tyypit miinus-laskutoimituksessa: %s ja %s.", left.Type(), right.Type())
	}
	isInt := lNum.IsInt && rNum.IsInt
	return &NumberVal{Value: lNum.Value - rNum.Value, IsInt: isInt}
}

func (e *Evaluator) evalMul(left, right Val, line int) Val {
	lNum, lOk := left.(*NumberVal)
	rNum, rOk := right.(*NumberVal)
	if !lOk || !rOk {
		ShowError(line, "Virheelliset tyypit kerto-laskutoimituksessa: %s ja %s.", left.Type(), right.Type())
	}
	isInt := lNum.IsInt && rNum.IsInt
	return &NumberVal{Value: lNum.Value * rNum.Value, IsInt: isInt}
}

func (e *Evaluator) evalDiv(left, right Val, line int) Val {
	lNum, lOk := left.(*NumberVal)
	rNum, rOk := right.(*NumberVal)
	if !lOk || !rOk {
		ShowError(line, "Virheelliset tyypit jako-laskutoimituksessa: %s ja %s.", left.Type(), right.Type())
	}
	if rNum.Value == 0 {
		ShowError(line, "Jako nollalla.")
	}
	val := lNum.Value / rNum.Value
	isInt := false
	if lNum.IsInt && rNum.IsInt && math.Mod(lNum.Value, rNum.Value) == 0 {
		isInt = true
	}
	return &NumberVal{Value: val, IsInt: isInt}
}

func (e *Evaluator) evalMod(left, right Val, line int) Val {
	lNum, lOk := left.(*NumberVal)
	rNum, rOk := right.(*NumberVal)
	if !lOk || !rOk {
		ShowError(line, "Virheelliset tyypit modulo-laskutoimituksessa: %s ja %s.", left.Type(), right.Type())
	}
	if rNum.Value == 0 {
		ShowError(line, "Jako nollalla.")
	}
	val := math.Mod(lNum.Value, rNum.Value)
	isInt := lNum.IsInt && rNum.IsInt
	return &NumberVal{Value: val, IsInt: isInt}
}

func (e *Evaluator) evalOn(left, right Val) bool {
	lBool, lIsBool := left.(*BoolVal)
	rBool, rIsBool := right.(*BoolVal)
	if lIsBool || rIsBool {
		if lIsBool && rIsBool {
			return lBool.Value == rBool.Value
		}
		return false
	}

	lNum, lIsNum := left.(*NumberVal)
	rNum, rIsNum := right.(*NumberVal)
	if lIsNum && rIsNum {
		return lNum.Value == rNum.Value
	}

	lStr, lIsStr := left.(*StringVal)
	rStr, rIsStr := right.(*StringVal)
	if lIsStr && rIsStr {
		return lStr.Value == rStr.Value
	}

	lList, lIsList := left.(*ListVal)
	rList, rIsList := right.(*ListVal)
	if lIsList && rIsList {
		if len(lList.Elements) != len(rList.Elements) {
			return false
		}
		for i := range lList.Elements {
			if !e.evalOn(lList.Elements[i], rList.Elements[i]) {
				return false
			}
		}
		return true
	}

	return false
}

func (e *Evaluator) evalYli(left, right Val, line int) bool {
	lNum, lIsNum := left.(*NumberVal)
	rNum, rIsNum := right.(*NumberVal)
	if lIsNum && rIsNum {
		return lNum.Value > rNum.Value
	}
	lStr, lIsStr := left.(*StringVal)
	rStr, rIsStr := right.(*StringVal)
	if lIsStr && rIsStr {
		return lStr.Value > rStr.Value
	}
	ShowError(line, "Vertailu yli (>) ei ole sallittu tyypeille %s ja %s.", left.Type(), right.Type())
	return false
}

func (e *Evaluator) evalAlle(left, right Val, line int) bool {
	lNum, lIsNum := left.(*NumberVal)
	rNum, rIsNum := right.(*NumberVal)
	if lIsNum && rIsNum {
		return lNum.Value < rNum.Value
	}
	lStr, lIsStr := left.(*StringVal)
	rStr, rIsStr := right.(*StringVal)
	if lIsStr && rIsStr {
		return lStr.Value < rStr.Value
	}
	ShowError(line, "Vertailu alle (<) ei ole sallittu tyypeille %s ja %s.", left.Type(), right.Type())
	return false
}

func (e *Evaluator) evalYliOn(left, right Val, line int) bool {
	lNum, lIsNum := left.(*NumberVal)
	rNum, rIsNum := right.(*NumberVal)
	if lIsNum && rIsNum {
		return lNum.Value >= rNum.Value
	}
	lStr, lIsStr := left.(*StringVal)
	rStr, rIsStr := right.(*StringVal)
	if lIsStr && rIsStr {
		return lStr.Value >= rStr.Value
	}
	ShowError(line, "Vertailu yli-on (>=) ei ole sallittu tyypeille %s ja %s.", left.Type(), right.Type())
	return false
}

func (e *Evaluator) evalAlleOn(left, right Val, line int) bool {
	lNum, lIsNum := left.(*NumberVal)
	rNum, rIsNum := right.(*NumberVal)
	if lIsNum && rIsNum {
		return lNum.Value <= rNum.Value
	}
	lStr, lIsStr := left.(*StringVal)
	rStr, rIsStr := right.(*StringVal)
	if lIsStr && rIsStr {
		return lStr.Value <= rStr.Value
	}
	ShowError(line, "Vertailu alle-on (<=) ei ole sallittu tyypeille %s ja %s.", left.Type(), right.Type())
	return false
}
