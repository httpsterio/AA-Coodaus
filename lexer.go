package main

import (
	"unicode"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

type Lexer struct {
	input              []rune
	pos                int
	line               int
	readRestOfLineNext bool
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:              []rune(input),
		pos:                0,
		line:               1,
		readRestOfLineNext: false,
	}
}

func (l *Lexer) NextToken() Token {
	if l.readRestOfLineNext {
		l.readRestOfLineNext = false
		return l.ReadRestOfLine()
	}

	l.skipWhitespaceExceptNewline()

	if l.pos >= len(l.input) {
		return Token{Type: TOKEN_EOF, Literal: "", Line: l.line}
	}

	r := l.input[l.pos]

	// Handle comments starting with --
	if r == '-' && l.peek() == '-' {
		l.skipComment()
		return l.NextToken()
	}

	// Handle newline
	if r == '\n' {
		tok := Token{Type: TOKEN_NEWLINE, Literal: "\n", Line: l.line}
		l.pos++
		l.line++
		return tok
	}

	// Handle colon
	if r == ':' {
		tok := Token{Type: TOKEN_COLON, Literal: ":", Line: l.line}
		l.pos++
		return tok
	}

	// Handle brackets
	if r == '[' {
		tok := Token{Type: TOKEN_LBRACKET, Literal: "[", Line: l.line}
		l.pos++
		return tok
	}
	if r == ']' {
		tok := Token{Type: TOKEN_RBRACKET, Literal: "]", Line: l.line}
		l.pos++
		return tok
	}

	// Handle standard single-rune operators
	switch r {
	case '+':
		l.pos++
		return Token{Type: TOKEN_PLUS, Literal: "+", Line: l.line}
	case '-':
		// Since we already checked for '--', this is a single minus operator
		l.pos++
		return Token{Type: TOKEN_MINUS, Literal: "-", Line: l.line}
	case '*':
		l.pos++
		return Token{Type: TOKEN_MUL, Literal: "*", Line: l.line}
	case '/':
		l.pos++
		return Token{Type: TOKEN_DIV, Literal: "/", Line: l.line}
	case '%':
		l.pos++
		return Token{Type: TOKEN_MOD, Literal: "%", Line: l.line}
	}

	// Handle quoted string
	if r == '"' {
		return l.readString()
	}

	// Handle numbers (integers and floats)
	if unicode.IsDigit(r) {
		return l.readNumber()
	}

	// Handle identifiers / keywords (letters, digits, underscore, and hyphens if surrounded by letters/digits)
	if isLetter(r) {
		return l.readIdentifier()
	}

	// Illegal character
	errChar := string(r)
	l.pos++
	ShowError(l.line, "Laiton merkki: %s", errChar)
	return Token{Type: TOKEN_EOF, Literal: "", Line: l.line}
}

func (l *Lexer) ReadRestOfLine() Token {
	var content []rune
	for l.pos < len(l.input) {
		if l.input[l.pos] == '\n' {
			break
		}
		// If we see comment --, stop reading content
		if l.input[l.pos] == '-' && l.peek() == '-' {
			// Skip comment until newline
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
			break
		}
		content = append(content, l.input[l.pos])
		l.pos++
	}

	// Trim spaces from content
	raw := string(content)
	// We want to trim leading/trailing spaces
	trimmed := trimSpaces(raw)

	// Process escape sequences if this is a quoted string
	if len(trimmed) >= 2 && trimmed[0] == '"' && trimmed[len(trimmed)-1] == '"' {
		inner := trimmed[1 : len(trimmed)-1]
		trimmed = "\"" + processEscapeSequences(inner) + "\""
	}

	return Token{Type: TOKEN_STRING, Literal: trimmed, Line: l.line}
}

func processEscapeSequences(s string) string {
	var result []rune
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		if runes[i] == '\\' && i+1 < len(runes) {
			switch runes[i+1] {
			case 'n':
				result = append(result, '\n')
				i++
			case 't':
				result = append(result, '\t')
				i++
			case '\\':
				result = append(result, '\\')
				i++
			case '"':
				result = append(result, '"')
				i++
			default:
				result = append(result, runes[i])
			}
		} else {
			result = append(result, runes[i])
		}
	}
	return string(result)
}

func trimSpaces(s string) string {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t' || s[start] == '\r') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

func (l *Lexer) peek() rune {
	if l.pos+1 >= len(l.input) {
		return 0
	}
	return l.input[l.pos+1]
}

func (l *Lexer) skipWhitespaceExceptNewline() {
	for l.pos < len(l.input) {
		r := l.input[l.pos]
		if r == ' ' || r == '\t' || r == '\r' {
			l.pos++
		} else {
			break
		}
	}
}

func (l *Lexer) skipComment() {
	// Skip --
	l.pos += 2
	for l.pos < len(l.input) {
		if l.input[l.pos] == '\n' {
			break
		}
		l.pos++
	}
}

func (l *Lexer) readString() Token {
	l.pos++ // Skip opening quote
	start := l.pos
	var content []rune

	for l.pos < len(l.input) {
		r := l.input[l.pos]
		if r == '"' {
			l.pos++ // Skip closing quote
			return Token{Type: TOKEN_STRING, Literal: string(content), Line: l.line}
		}
		if r == '\n' {
			ShowError(l.line, "Lainausmerkki puuttuu rivin lopusta.")
		}
		if r == '\\' {
			if l.pos+1 >= len(l.input) {
				ShowError(l.line, "Laiton perääntymisjono merkkijonossa.")
			}
			next := l.input[l.pos+1]
			switch next {
			case 'n':
				content = append(content, '\n')
			case 't':
				content = append(content, '\t')
			case '\\':
				content = append(content, '\\')
			case '"':
				content = append(content, '"')
			default:
				ShowError(l.line, "Laiton perääntymisjono: \\%c", next)
			}
			l.pos += 2
		} else {
			content = append(content, r)
			l.pos++
		}
	}

	ShowError(l.line, "Merkkijonoa ei suljettu. Alkoi kohdasta %d", start)
	return Token{}
}

func (l *Lexer) readNumber() Token {
	start := l.pos
	isFloat := false
	for l.pos < len(l.input) {
		r := l.input[l.pos]
		if unicode.IsDigit(r) {
			l.pos++
		} else if r == '.' && unicode.IsDigit(l.peek()) {
			if isFloat {
				break // Cannot have two decimal points
			}
			isFloat = true
			l.pos++
		} else {
			break
		}
	}
	return Token{Type: TOKEN_NUMBER, Literal: string(l.input[start:l.pos]), Line: l.line}
}

func (l *Lexer) readIdentifier() Token {
	start := l.pos
	for l.pos < len(l.input) {
		r := l.input[l.pos]
		if isLetter(r) || unicode.IsDigit(r) || r == '_' {
			l.pos++
		} else if r == '-' && isLetter(l.peek()) {
			// Hyphen is allowed in compound words if followed by a letter (e.g. yli-on, Ot-yhteys)
			l.pos++
		} else {
			break
		}
	}
	lit := string(l.input[start:l.pos])
	tokType := LookupIdent(lit)
	return Token{Type: tokType, Literal: lit, Line: l.line}
}

func isLetter(r rune) bool {
	return unicode.IsLetter(r) || r == '_' || r == 'å' || r == 'ä' || r == 'ö' || r == 'Å' || r == 'Ä' || r == 'Ö'
}

type TokenType string

const (
	TOKEN_EOF      TokenType = "EOF"
	TOKEN_NEWLINE  TokenType = "NEWLINE"
	TOKEN_IDENT    TokenType = "IDENT"
	TOKEN_NUMBER   TokenType = "NUMBER"
	TOKEN_STRING   TokenType = "STRING"
	TOKEN_COLON    TokenType = "COLON"
	TOKEN_LBRACKET TokenType = "LBRACKET"
	TOKEN_RBRACKET TokenType = "RBRACKET"

	// Keywords
	TOKEN_MU      TokenType = "Mu"
	TOKEN_SANO    TokenType = "Sano"
	TOKEN_KYSY    TokenType = "Kysy"
	TOKEN_JOS     TokenType = "Jos"
	TOKEN_MUUTEN  TokenType = "Muuten"
	TOKEN_LOPPU   TokenType = "Loppu"
	TOKEN_TOIS    TokenType = "Tois"
	TOKEN_LOPUTON TokenType = "Loputon"
	TOKEN_KUNNES  TokenType = "Kunnes"
	TOKEN_TEKO    TokenType = "Teko"
	TOKEN_TEE     TokenType = "Tee"
	TOKEN_ANNA    TokenType = "Anna"
	TOKEN_LISAA   TokenType = "Lisää"
	TOKEN_PITUUS  TokenType = "Pituus"
	TOKEN_JOKA    TokenType = "Joka"

	// String operations (Milestone 3)
	TOKEN_PILKO    TokenType = "Pilko"
	TOKEN_KORVAA   TokenType = "Korvaa"
	TOKEN_SISALTAA TokenType = "Sisaltaa"
	TOKEN_TRIMMAA  TokenType = "Trimmaa"
	TOKEN_ISOT     TokenType = "Isot"
	TOKEN_PIENET   TokenType = "Pienet"

	// File I/O (Milestone 2)
	TOKEN_LUE        TokenType = "Lue"
	TOKEN_KIRJOITA   TokenType = "Kirjoita"
	TOKEN_YLIKIRJOITA TokenType = "Ylikirjoita"
	TOKEN_LIITA      TokenType = "Liitä"
	TOKEN_LISTAA     TokenType = "Listaa"
	TOKEN_HAK        TokenType = "Hak"
	TOKEN_LUO        TokenType = "Luo"
	TOKEN_LUOHAK     TokenType = "Luohak"
	TOKEN_ONKO       TokenType = "Onko"

	// Milestone 4
	TOKEN_SATUNNAINEN TokenType = "Satunnainen"
	TOKEN_SATU        TokenType = "Satu"
	TOKEN_LOPETA      TokenType = "Lopeta"
	TOKEN_LOPT        TokenType = "Lopt"

	// Operators
	TOKEN_PLUS    TokenType = "+"
	TOKEN_MINUS   TokenType = "-"
	TOKEN_MUL     TokenType = "*"
	TOKEN_DIV     TokenType = "/"
	TOKEN_MOD     TokenType = "%"
	TOKEN_ON      TokenType = "on"
	TOKEN_EI      TokenType = "ei"
	TOKEN_YLI     TokenType = "yli"
	TOKEN_ALLE    TokenType = "alle"
	TOKEN_YLI_ON  TokenType = "yli-on"
	TOKEN_ALLE_ON TokenType = "alle-on"
	TOKEN_KYLLA   TokenType = "kyllä"
)

var keywords = map[string]TokenType{
	"Mu":      TOKEN_MU,
	"Sano":    TOKEN_SANO,
	"Kysy":    TOKEN_KYSY,
	"Jos":     TOKEN_JOS,
	"Muuten":  TOKEN_MUUTEN,
	"Loppu":   TOKEN_LOPPU,
	"Tois":    TOKEN_TOIS,
	"Loputon": TOKEN_LOPUTON,
	"Kunnes":  TOKEN_KUNNES,
	"Teko":    TOKEN_TEKO,
	"Tee":     TOKEN_TEE,
	"Anna":    TOKEN_ANNA,
	"Lisää":   TOKEN_LISAA,
	"Pituus":  TOKEN_PITUUS,
	"Joka":    TOKEN_JOKA,

	"Pilko":    TOKEN_PILKO,
	"Pilk":     TOKEN_PILKO,
	"Korvaa":   TOKEN_KORVAA,
	"Korv":     TOKEN_KORVAA,
	"Sisaltaa": TOKEN_SISALTAA,
	"Sis":      TOKEN_SISALTAA,
	"Trimmaa":  TOKEN_TRIMMAA,
	"Trim":     TOKEN_TRIMMAA,
	"Isot":     TOKEN_ISOT,
	"Pienet":   TOKEN_PIENET,

	"Lue":       TOKEN_LUE,
	"Kirjoita":  TOKEN_KIRJOITA,
	"Kirj":      TOKEN_KIRJOITA,
	"Ylikirjoita": TOKEN_YLIKIRJOITA,
	"Ylikirj":   TOKEN_YLIKIRJOITA,
	"Liitä":     TOKEN_LIITA,
	"Liit":      TOKEN_LIITA,
	"Listaa":    TOKEN_LISTAA,
	"Hak":       TOKEN_HAK,
	"Luo":       TOKEN_LUO,
	"Luohak":    TOKEN_LUOHAK,
	"Onko":      TOKEN_ONKO,
	"Satunnainen": TOKEN_SATUNNAINEN,
	"Satu":        TOKEN_SATU,
	"Lopeta":      TOKEN_LOPETA,
	"Lopt":        TOKEN_LOPT,
	"on":       TOKEN_ON,
	"ei":       TOKEN_EI,
	"yli":      TOKEN_YLI,
	"alle":     TOKEN_ALLE,
	"yli-on":   TOKEN_YLI_ON,
	"alle-on":  TOKEN_ALLE_ON,
	"kyllä":    TOKEN_KYLLA,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return TOKEN_IDENT
}
