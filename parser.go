package main

import (
	"strconv"
	"strings"
)

// PreCheckSource checks the source text for invalid commands or common syntax errors.
func PreCheckSource(input string) {
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		// Strip comments
		commentIdx := strings.Index(line, "--")
		cleanLine := line
		if commentIdx >= 0 {
			cleanLine = line[:commentIdx]
		}
		cleanLine = strings.TrimSpace(cleanLine)

		// Check for specific forbidden meme patterns/errors
		if strings.Contains(cleanLine, "Ota yhteys") {
			ShowError(i+1, "\"Ota yhteys\" ei kelpaa. Kirjoita Ot-yhteys.")
		}
	}
}

// Precedence levels
const (
	_ int = iota
	LOWEST
	COMPARE // on, ei, yli, alle, yli-on, alle-on
	SUM     // +, -
	PRODUCT // *, /, %
	INDEX   // lista[1]
)

var precedences = map[TokenType]int{
	TOKEN_ON:       COMPARE,
	TOKEN_EI:       COMPARE,
	TOKEN_YLI:      COMPARE,
	TOKEN_ALLE:     COMPARE,
	TOKEN_YLI_ON:   COMPARE,
	TOKEN_ALLE_ON:  COMPARE,
	TOKEN_PLUS:     SUM,
	TOKEN_MINUS:    SUM,
	TOKEN_MUL:      PRODUCT,
	TOKEN_DIV:      PRODUCT,
	TOKEN_MOD:      PRODUCT,
	TOKEN_LBRACKET: INDEX,
}

// AST Nodes
type Node interface {
	LineNumber() int
}

type Expr interface {
	Node
	exprNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type astNode struct {
	Line int
}

func (a astNode) LineNumber() int { return a.Line }

// Expressions

type NumberExpr struct {
	astNode
	Val   float64
	IsInt bool
}

func (e *NumberExpr) exprNode() {}

type StringExpr struct {
	astNode
	Val string
}

func (e *StringExpr) exprNode() {}

type BoolExpr struct {
	astNode
	Val bool
}

func (e *BoolExpr) exprNode() {}

type ListExpr struct {
	astNode
	Elements []Expr
}

func (e *ListExpr) exprNode() {}

type ListIndexExpr struct {
	astNode
	List  Expr
	Index Expr
}

func (e *ListIndexExpr) exprNode() {}

type VarExpr struct {
	astNode
	Name string
}

func (e *VarExpr) exprNode() {}

type BinaryExpr struct {
	astNode
	Op    TokenType
	Left  Expr
	Right Expr
}

func (e *BinaryExpr) exprNode() {}

type PituusExpr struct {
	astNode
	List Expr
}

func (e *PituusExpr) exprNode() {}

type TeeExpr struct {
	astNode
	Name string
	Args []Expr
}

func (e *TeeExpr) exprNode() {}

// Statements

type VarDeclStmt struct {
	astNode
	Name string
	Val  Expr
}

func (s *VarDeclStmt) stmtNode() {}

type SanoStmt struct {
	astNode
	Content string
}

func (s *SanoStmt) stmtNode() {}

type KysyStmt struct {
	astNode
	Name   string
	Prompt string
}

func (s *KysyStmt) stmtNode() {}

type JosStmt struct {
	astNode
	Cond     Expr
	ThenBody []Stmt
	ElseBody []Stmt
}

func (s *JosStmt) stmtNode() {}

type ToisStmt struct {
	astNode
	Count   Expr
	VarName string
	Body    []Stmt
}

func (s *ToisStmt) stmtNode() {}

type LoputonStmt struct {
	astNode
	Body []Stmt
}

func (s *LoputonStmt) stmtNode() {}

type KunnesStmt struct {
	astNode
	Cond Expr
	Body []Stmt
}

func (s *KunnesStmt) stmtNode() {}

type TekoStmt struct {
	astNode
	Name   string
	Params []string
	Body   []Stmt
}

func (s *TekoStmt) stmtNode() {}

type TeeStmt struct {
	astNode
	Call *TeeExpr
}

func (s *TeeStmt) stmtNode() {}

type AnnaStmt struct {
	astNode
	Val Expr
}

func (s *AnnaStmt) stmtNode() {}

type LisaaStmt struct {
	astNode
	ListName string
	Val      Expr
}

func (s *LisaaStmt) stmtNode() {}

type JokaStmt struct {
	astNode
	ListExpr Expr
	VarName  string
	Body     []Stmt
}

func (s *JokaStmt) stmtNode() {}

// Parser struct and logic
type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectToken(t TokenType) {
	if p.curTokenIs(t) {
		p.nextToken()
	} else {
		if t == TOKEN_LOPPU {
			ShowError(p.curToken.Line, "Loppu puuttuu. Tarkista lohko.")
		}
		ShowError(p.curToken.Line, "Odotettiin merkkiä %s, saatiin %s", t, p.curToken.Type)
	}
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) ParseProgram() []Stmt {
	var program []Stmt
	for !p.curTokenIs(TOKEN_EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program = append(program, stmt)
		}
		p.consumeNewline()
	}
	return program
}

func (p *Parser) consumeNewline() {
	for p.curTokenIs(TOKEN_NEWLINE) {
		p.nextToken()
	}
}

func (p *Parser) consumeStatementTerminator() {
	if p.curTokenIs(TOKEN_NEWLINE) {
		p.consumeNewline()
	} else if p.curTokenIs(TOKEN_EOF) {
		// do nothing
	} else {
		ShowError(p.curToken.Line, "Syntaksivirhe: odotettiin rivinvaihtoa, saatiin %s (%s)", p.curToken.Type, p.curToken.Literal)
	}
}

func (p *Parser) parseStatement() Stmt {
	p.consumeNewline()
	if p.curTokenIs(TOKEN_EOF) {
		return nil
	}

	tok := p.curToken
	switch tok.Type {
	case TOKEN_MU:
		stmt := p.parseVarDeclStmt()
		p.consumeStatementTerminator()
		return stmt
	case TOKEN_SANO:
		stmt := p.parseSanoStmt()
		p.consumeStatementTerminator()
		return stmt
	case TOKEN_KYSY:
		stmt := p.parseKysyStmt()
		p.consumeStatementTerminator()
		return stmt
	case TOKEN_JOS:
		return p.parseJosStmt()
	case TOKEN_TOIS:
		return p.parseToisStmt()
	case TOKEN_LOPUTON:
		return p.parseLoputonStmt()
	case TOKEN_KUNNES:
		return p.parseKunnesStmt()
	case TOKEN_TEKO:
		return p.parseTekoStmt()
	case TOKEN_TEE:
		expr := p.parseTeeExpr()
		stmt := &TeeStmt{astNode: astNode{Line: tok.Line}, Call: expr.(*TeeExpr)}
		p.consumeStatementTerminator()
		return stmt
	case TOKEN_ANNA:
		stmt := p.parseAnnaStmt()
		p.consumeStatementTerminator()
		return stmt
	case TOKEN_LISAA:
		stmt := p.parseLisaaStmt()
		p.consumeStatementTerminator()
		return stmt
	case TOKEN_JOKA:
		return p.parseJokaStmt()
	default:
		ShowError(tok.Line, "Tuntematon lauseke tai komento: %s", tok.Literal)
		return nil
	}
}

func (p *Parser) parseVarDeclStmt() Stmt {
	tok := p.curToken
	p.nextToken() // consume Mu

	if !p.curTokenIs(TOKEN_IDENT) {
		ShowError(tok.Line, "Muuttujan nimi puuttuu Mu-määrittelyn jälkeen.")
	}
	name := p.curToken.Literal
	p.nextToken() // consume variable name

	p.expectToken(TOKEN_COLON)

	val := p.parseExpression(LOWEST)
	return &VarDeclStmt{
		astNode: astNode{Line: tok.Line},
		Name:    name,
		Val:     val,
	}
}

func (p *Parser) parseSanoStmt() Stmt {
	tok := p.curToken
	p.l.readRestOfLineNext = true
	p.nextToken() // consume Sano. curToken becomes COLON, peekToken becomes the string.

	p.expectToken(TOKEN_COLON) // curToken becomes the string. peekToken becomes newline.

	if !p.curTokenIs(TOKEN_STRING) {
		ShowError(tok.Line, "Sano-komennon jälkeen odotettiin tekstiä.")
	}
	content := p.curToken.Literal
	p.nextToken() // consume the string

	return &SanoStmt{
		astNode: astNode{Line: tok.Line},
		Content: content,
	}
}

func (p *Parser) parseKysyStmt() Stmt {
	tok := p.curToken
	p.nextToken() // consume Kysy. curToken becomes var name, peekToken becomes COLON.

	if !p.curTokenIs(TOKEN_IDENT) {
		ShowError(tok.Line, "Muuttujan nimi puuttuu Kysy-komennon jälkeen.")
	}
	name := p.curToken.Literal

	p.l.readRestOfLineNext = true
	p.nextToken() // consume variable name. curToken becomes COLON, peekToken becomes the string.

	p.expectToken(TOKEN_COLON) // curToken becomes the string. peekToken becomes newline.

	if !p.curTokenIs(TOKEN_STRING) {
		ShowError(tok.Line, "Kysy-komennon jälkeen odotettiin kysymystekstiä.")
	}
	prompt := p.curToken.Literal
	p.nextToken() // consume the string

	return &KysyStmt{
		astNode: astNode{Line: tok.Line},
		Name:    name,
		Prompt:  prompt,
	}
}

func (p *Parser) parseBlock(startLine int) []Stmt {
	var body []Stmt
	p.consumeNewline()
	for !p.curTokenIs(TOKEN_LOPPU) && !p.curTokenIs(TOKEN_EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			body = append(body, stmt)
		}
		p.consumeNewline()
	}
	if p.curTokenIs(TOKEN_EOF) {
		ShowError(startLine, "Loppu puuttuu. Tarkista lohko.")
	}
	p.nextToken() // consume Loppu
	return body
}

func (p *Parser) parseJosStmt() Stmt {
	tok := p.curToken
	p.nextToken() // consume Jos

	cond := p.parseExpression(LOWEST)

	p.expectToken(TOKEN_COLON)
	p.consumeNewline()

	var thenBody []Stmt
	var elseBody []Stmt

	for !p.curTokenIs(TOKEN_MUUTEN) && !p.curTokenIs(TOKEN_LOPPU) && !p.curTokenIs(TOKEN_EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			thenBody = append(thenBody, stmt)
		}
		p.consumeNewline()
	}

	if p.curTokenIs(TOKEN_MUUTEN) {
		p.nextToken() // consume Muuten
		p.expectToken(TOKEN_COLON)
		p.consumeNewline()

		for !p.curTokenIs(TOKEN_LOPPU) && !p.curTokenIs(TOKEN_EOF) {
			stmt := p.parseStatement()
			if stmt != nil {
				elseBody = append(elseBody, stmt)
			}
			p.consumeNewline()
		}
	}

	if p.curTokenIs(TOKEN_EOF) {
		ShowError(tok.Line, "Loppu puuttuu. Tarkista lohko.")
	}
	p.expectToken(TOKEN_LOPPU)

	return &JosStmt{
		astNode:  astNode{Line: tok.Line},
		Cond:     cond,
		ThenBody: thenBody,
		ElseBody: elseBody,
	}
}

func (p *Parser) parseToisStmt() Stmt {
	tok := p.curToken
	p.nextToken() // consume Tois

	count := p.parseExpression(LOWEST)

	var varName string
	if p.curTokenIs(TOKEN_IDENT) {
		varName = p.curToken.Literal
		p.nextToken()
	}

	p.expectToken(TOKEN_COLON)

	body := p.parseBlock(tok.Line)

	return &ToisStmt{
		astNode: astNode{Line: tok.Line},
		Count:   count,
		VarName: varName,
		Body:    body,
	}
}

func (p *Parser) parseLoputonStmt() Stmt {
	tok := p.curToken
	p.nextToken() // consume Loputon

	p.expectToken(TOKEN_COLON)

	body := p.parseBlock(tok.Line)

	return &LoputonStmt{
		astNode: astNode{Line: tok.Line},
		Body:    body,
	}
}

func (p *Parser) parseKunnesStmt() Stmt {
	tok := p.curToken
	p.nextToken() // consume Kunnes

	cond := p.parseExpression(LOWEST)

	p.expectToken(TOKEN_COLON)

	body := p.parseBlock(tok.Line)

	return &KunnesStmt{
		astNode: astNode{Line: tok.Line},
		Cond:    cond,
		Body:    body,
	}
}

func (p *Parser) parseTekoStmt() Stmt {
	tok := p.curToken
	p.nextToken() // consume Teko

	if !p.curTokenIs(TOKEN_IDENT) {
		ShowError(tok.Line, "Funktion nimi puuttuu Teko-määrittelyn jälkeen.")
	}
	name := p.curToken.Literal
	p.nextToken() // consume name

	var params []string
	for p.curTokenIs(TOKEN_IDENT) {
		params = append(params, p.curToken.Literal)
		p.nextToken()
	}

	p.expectToken(TOKEN_COLON)

	body := p.parseBlock(tok.Line)

	return &TekoStmt{
		astNode: astNode{Line: tok.Line},
		Name:    name,
		Params:  params,
		Body:    body,
	}
}

func (p *Parser) parseAnnaStmt() Stmt {
	tok := p.curToken
	p.nextToken() // consume Anna

	val := p.parseExpression(LOWEST)

	return &AnnaStmt{
		astNode: astNode{Line: tok.Line},
		Val:     val,
	}
}

func (p *Parser) parseLisaaStmt() Stmt {
	tok := p.curToken
	p.nextToken() // consume Lisää

	if !p.curTokenIs(TOKEN_IDENT) {
		ShowError(tok.Line, "Listan nimi puuttuu Lisää-komennon jälkeen.")
	}
	listName := p.curToken.Literal
	p.nextToken() // consume name

	p.expectToken(TOKEN_COLON)

	val := p.parseExpression(LOWEST)

	return &LisaaStmt{
		astNode:  astNode{Line: tok.Line},
		ListName: listName,
		Val:      val,
	}
}

func (p *Parser) parseJokaStmt() Stmt {
	tok := p.curToken
	p.nextToken() // consume Joka

	listExpr := p.parseExpression(LOWEST)

	if !p.curTokenIs(TOKEN_IDENT) {
		ShowError(tok.Line, "Muuttujan nimi puuttuu Joka-komennosta.")
	}
	varName := p.curToken.Literal
	p.nextToken() // consume name

	p.expectToken(TOKEN_COLON)

	body := p.parseBlock(tok.Line)

	return &JokaStmt{
		astNode:  astNode{Line: tok.Line},
		ListExpr: listExpr,
		VarName:  varName,
		Body:     body,
	}
}

// Expression parsing with Pratt parser
func (p *Parser) parseExpression(precedence int) Expr {
	prefix := p.findPrefixHandler(p.curToken.Type)
	if prefix == nil {
		ShowError(p.curToken.Line, "Syntaksivirhe: ei voida aloittaa lauseketta merkillä %s (%s)", p.curToken.Type, p.curToken.Literal)
	}
	leftExp := prefix()

	for !p.curTokenIs(TOKEN_NEWLINE) && !p.curTokenIs(TOKEN_EOF) && precedence < p.curPrecedence() {
		infix := p.findInfixHandler(p.curToken.Type)
		if infix == nil {
			return leftExp
		}
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) findPrefixHandler(t TokenType) func() Expr {
	switch t {
	case TOKEN_IDENT:
		return p.parseVarExpr
	case TOKEN_NUMBER:
		return p.parseNumberExpr
	case TOKEN_STRING:
		return p.parseStringExpr
	case TOKEN_KYLLA:
		return func() Expr {
			tok := p.curToken
			p.nextToken()
			return &BoolExpr{astNode: astNode{Line: tok.Line}, Val: true}
		}
	case TOKEN_EI:
		// ei can be a boolean literal false in a prefix position
		return func() Expr {
			tok := p.curToken
			p.nextToken()
			return &BoolExpr{astNode: astNode{Line: tok.Line}, Val: false}
		}
	case TOKEN_LBRACKET:
		return p.parseListLiteralExpr
	case TOKEN_PITUUS:
		return p.parsePituusExpr
	case TOKEN_TEE:
		return p.parseTeeExpr
	case TOKEN_MINUS:
		// support unary minus
		return p.parsePrefixMinusExpr
	}
	return nil
}

func (p *Parser) findInfixHandler(t TokenType) func(Expr) Expr {
	switch t {
	case TOKEN_PLUS, TOKEN_MINUS, TOKEN_MUL, TOKEN_DIV, TOKEN_MOD, TOKEN_ON, TOKEN_EI, TOKEN_YLI, TOKEN_ALLE, TOKEN_YLI_ON, TOKEN_ALLE_ON:
		return p.parseBinaryExpr
	case TOKEN_LBRACKET:
		return p.parseListIndexExpr
	}
	return nil
}

func (p *Parser) parseVarExpr() Expr {
	tok := p.curToken
	p.nextToken()
	return &VarExpr{
		astNode: astNode{Line: tok.Line},
		Name:    tok.Literal,
	}
}

func (p *Parser) parseNumberExpr() Expr {
	tok := p.curToken
	val, err := strconv.ParseFloat(tok.Literal, 64)
	if err != nil {
		ShowError(tok.Line, "Laiton numero: %s", tok.Literal)
	}
	isInt := !strings.Contains(tok.Literal, ".")
	p.nextToken()
	return &NumberExpr{
		astNode: astNode{Line: tok.Line},
		Val:     val,
		IsInt:   isInt,
	}
}

func (p *Parser) parseStringExpr() Expr {
	tok := p.curToken
	p.nextToken()
	return &StringExpr{
		astNode: astNode{Line: tok.Line},
		Val:     tok.Literal,
	}
}

func (p *Parser) parsePrefixMinusExpr() Expr {
	tok := p.curToken
	p.nextToken() // consume minus
	// Parse operand with high precedence
	operand := p.parseExpression(PRODUCT)

	// If the operand is a number, we can directly negate it
	if num, ok := operand.(*NumberExpr); ok {
		num.Val = -num.Val
		return num
	}

	// Otherwise build a binary expression with a 0 LHS
	return &BinaryExpr{
		astNode: astNode{Line: tok.Line},
		Op:      TOKEN_MINUS,
		Left:    &NumberExpr{astNode: astNode{Line: tok.Line}, Val: 0, IsInt: true},
		Right:   operand,
	}
}

func (p *Parser) parseBinaryExpr(left Expr) Expr {
	tok := p.curToken
	precedence := p.curPrecedence()
	p.nextToken()
	right := p.parseExpression(precedence)

	return &BinaryExpr{
		astNode: astNode{Line: tok.Line},
		Op:      tok.Type,
		Left:    left,
		Right:   right,
	}
}

func (p *Parser) parseListLiteralExpr() Expr {
	tok := p.curToken
	p.nextToken() // consume [

	var elements []Expr
	for !p.curTokenIs(TOKEN_RBRACKET) && !p.curTokenIs(TOKEN_EOF) {
		expr := p.parseExpression(LOWEST)
		elements = append(elements, expr)

		p.consumeNewline()
	}

	p.expectToken(TOKEN_RBRACKET)

	return &ListExpr{
		astNode:  astNode{Line: tok.Line},
		Elements: elements,
	}
}

func (p *Parser) parseListIndexExpr(left Expr) Expr {
	tok := p.curToken
	p.nextToken() // consume [

	index := p.parseExpression(LOWEST)

	p.expectToken(TOKEN_RBRACKET)

	return &ListIndexExpr{
		astNode: astNode{Line: tok.Line},
		List:    left,
		Index:   index,
	}
}

func (p *Parser) parsePituusExpr() Expr {
	tok := p.curToken
	p.nextToken() // consume Pituus

	// Pituus binds tightly
	operand := p.parseExpression(PRODUCT)

	return &PituusExpr{
		astNode: astNode{Line: tok.Line},
		List:    operand,
	}
}

func (p *Parser) canStartExpression(tok Token) bool {
	switch tok.Type {
	case TOKEN_IDENT, TOKEN_NUMBER, TOKEN_STRING, TOKEN_LBRACKET, TOKEN_KYLLA, TOKEN_EI, TOKEN_PITUUS, TOKEN_TEE, TOKEN_MINUS:
		return true
	}
	return false
}

func (p *Parser) parseTeeExpr() Expr {
	tok := p.curToken
	p.nextToken() // consume Tee

	if !p.curTokenIs(TOKEN_IDENT) {
		ShowError(tok.Line, "Funktion nimi puuttuu Tee-komennon jälkeen.")
	}
	name := p.curToken.Literal
	p.nextToken() // consume function name

	var args []Expr
	for p.canStartExpression(p.curToken) {
		arg := p.parseExpression(COMPARE)
		args = append(args, arg)
	}

	return &TeeExpr{
		astNode: astNode{Line: tok.Line},
		Name:    name,
		Args:    args,
	}
}
