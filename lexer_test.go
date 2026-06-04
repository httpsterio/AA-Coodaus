package main

import "testing"

func collectTokens(input string) []Token {
	l := NewLexer(input)
	var tokens []Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == TOKEN_EOF {
			break
		}
	}
	return tokens
}

func TestLexerNumber(t *testing.T) {
	tokens := collectTokens("42")
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}
	if tokens[0].Type != TOKEN_NUMBER || tokens[0].Literal != "42" {
		t.Errorf("expected NUMBER 42, got %s %q", tokens[0].Type, tokens[0].Literal)
	}
	if tokens[1].Type != TOKEN_EOF {
		t.Errorf("expected EOF, got %s", tokens[1].Type)
	}
}

func TestLexerFloat(t *testing.T) {
	tokens := collectTokens("3.14")
	if tokens[0].Type != TOKEN_NUMBER || tokens[0].Literal != "3.14" {
		t.Errorf("expected NUMBER 3.14, got %s %q", tokens[0].Type, tokens[0].Literal)
	}
}

func TestLexerKeywords(t *testing.T) {
	tests := []struct {
		input string
		typ   TokenType
	}{
		{"Mu", TOKEN_MU},
		{"Sano", TOKEN_SANO},
		{"Kysy", TOKEN_KYSY},
		{"Jos", TOKEN_JOS},
		{"Muuten", TOKEN_MUUTEN},
		{"Loppu", TOKEN_LOPPU},
		{"Tois", TOKEN_TOIS},
		{"Loputon", TOKEN_LOPUTON},
		{"Kunnes", TOKEN_KUNNES},
		{"Teko", TOKEN_TEKO},
		{"Tee", TOKEN_TEE},
		{"Anna", TOKEN_ANNA},
		{"Lisää", TOKEN_LISAA},
		{"Pituus", TOKEN_PITUUS},
		{"Joka", TOKEN_JOKA},
	}
	for _, tt := range tests {
		tokens := collectTokens(tt.input)
		if len(tokens) < 2 || tokens[0].Type != tt.typ {
			t.Errorf("collectTokens(%q) = first token %s, want %s", tt.input, tokens[0].Type, tt.typ)
		}
	}
}

func TestLexerComparisonOps(t *testing.T) {
	tests := []struct {
		input string
		typ   TokenType
	}{
		{"on", TOKEN_ON},
		{"ei", TOKEN_EI},
		{"yli", TOKEN_YLI},
		{"alle", TOKEN_ALLE},
		{"yli-on", TOKEN_YLI_ON},
		{"alle-on", TOKEN_ALLE_ON},
	}
	for _, tt := range tests {
		tokens := collectTokens(tt.input)
		if len(tokens) < 2 || tokens[0].Type != tt.typ {
			t.Errorf("collectTokens(%q) = first token %s, want %s", tt.input, tokens[0].Type, tt.typ)
		}
	}
}

func TestLexerBoolKv(t *testing.T) {
	tokens := collectTokens("kyllä")
	if tokens[0].Type != TOKEN_KYLLA {
		t.Errorf("expected KYLLA, got %s", tokens[0].Type)
	}
}

func TestLexerString(t *testing.T) {
	tokens := collectTokens(`"moi maailma"`)
	if tokens[0].Type != TOKEN_STRING || tokens[0].Literal != "moi maailma" {
		t.Errorf("expected STRING \"moi maailma\", got %s %q", tokens[0].Type, tokens[0].Literal)
	}
}

func TestLexerStringEscapes(t *testing.T) {
	tokens := collectTokens(`"sano\n\t\\\"hei"`)
	want := "sano\n\t\\\"hei"
	if tokens[0].Literal != want {
		t.Errorf("expected %q, got %q", want, tokens[0].Literal)
	}
}

// TestLexerStringUnclosed omitted: ShowError calls os.Exit, cannot test in-process.

func TestLexerCommentSkipped(t *testing.T) {
	tokens := collectTokens("-- kommentti\nMoi")
	// Comment is replaced by newline (read until \n, break, then next NextToken produces NEWLINE)
	if len(tokens) < 3 {
		t.Fatalf("expected at least 3 tokens (NEWLINE, IDENT, EOF), got %d", len(tokens))
	}
	if tokens[1].Type != TOKEN_IDENT || tokens[1].Literal != "Moi" {
		t.Errorf("expected IDENT Moi at index 1, got %s %q", tokens[1].Type, tokens[1].Literal)
	}
}

func TestLexerCommentInline(t *testing.T) {
	tokens := collectTokens("Sano -- tämä on kommentti\n")
	if len(tokens) < 2 {
		t.Fatalf("expected at least 2 tokens, got %d", len(tokens))
	}
	if tokens[0].Type != TOKEN_SANO {
		t.Errorf("expected SANO, got %s", tokens[0].Type)
	}
	if tokens[1].Type != TOKEN_NEWLINE {
		t.Errorf("expected NEWLINE, got %s", tokens[1].Type)
	}
}

func TestLexerColon(t *testing.T) {
	tokens := collectTokens(":")
	if tokens[0].Type != TOKEN_COLON {
		t.Errorf("expected COLON, got %s", tokens[0].Type)
	}
}

func TestLexerBrackets(t *testing.T) {
	tokens := collectTokens("[1 2]")
	if tokens[0].Type != TOKEN_LBRACKET {
		t.Errorf("expected LBRACKET, got %s", tokens[0].Type)
	}
	// 1, 2 are numbers
	if tokens[3].Type != TOKEN_RBRACKET {
		t.Errorf("expected RBRACKET, got %s", tokens[3].Type)
	}
}

func TestLexerOperators(t *testing.T) {
	tokens := collectTokens("+ - * / %")
	ops := []TokenType{TOKEN_PLUS, TOKEN_MINUS, TOKEN_MUL, TOKEN_DIV, TOKEN_MOD}
	for i, op := range ops {
		if tokens[i].Type != op {
			t.Errorf("token %d: expected %s, got %s", i, op, tokens[i].Type)
		}
	}
}

func TestLexerNewline(t *testing.T) {
	tokens := collectTokens("a\nb")
	found := false
	for _, tok := range tokens {
		if tok.Type == TOKEN_NEWLINE {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected a NEWLINE token between a and b")
	}
}

func TestLexerCompoundIdent(t *testing.T) {
	tokens := collectTokens("yli-on alle-on")
	if tokens[0].Type != TOKEN_YLI_ON {
		t.Errorf("expected YLI_ON, got %s", tokens[0].Type)
	}
	if tokens[1].Type != TOKEN_ALLE_ON {
		t.Errorf("expected ALLE_ON, got %s", tokens[1].Type)
	}
}

func TestLexerIdentifier(t *testing.T) {
	tokens := collectTokens("nimi")
	if tokens[0].Type != TOKEN_IDENT || tokens[0].Literal != "nimi" {
		t.Errorf("expected IDENT nimi, got %s %q", tokens[0].Type, tokens[0].Literal)
	}
}

func TestLexerEmpty(t *testing.T) {
	tokens := collectTokens("")
	if len(tokens) != 1 || tokens[0].Type != TOKEN_EOF {
		t.Errorf("expected single EOF token, got %d tokens", len(tokens))
	}
}

func TestLexerMuStatement(t *testing.T) {
	tokens := collectTokens("Mu x: 5\n")
	if len(tokens) < 5 {
		t.Fatalf("expected at least 5 tokens, got %d", len(tokens))
	}
	if tokens[0].Type != TOKEN_MU {
		t.Errorf("token 0: expected MU, got %s", tokens[0].Type)
	}
	if tokens[1].Type != TOKEN_IDENT || tokens[1].Literal != "x" {
		t.Errorf("token 1: expected IDENT x, got %s %q", tokens[1].Type, tokens[1].Literal)
	}
	if tokens[2].Type != TOKEN_COLON {
		t.Errorf("token 2: expected COLON, got %s", tokens[2].Type)
	}
	if tokens[3].Type != TOKEN_NUMBER || tokens[3].Literal != "5" {
		t.Errorf("token 3: expected NUMBER 5, got %s %q", tokens[3].Type, tokens[3].Literal)
	}
	if tokens[4].Type != TOKEN_NEWLINE {
		t.Errorf("token 4: expected NEWLINE, got %s", tokens[4].Type)
	}
}

func TestLexerSanoKysyReadRest(t *testing.T) {
	// After Sano or Kysy, the parser signals the lexer to read the rest of line
	// as a single string via ReadRestOfLine. When called normally (without flag),
	// the rest of line is tokenized word by word.
	tokens := collectTokens("Sano:\nKysy")
	if len(tokens) < 3 {
		t.Fatalf("expected at least 3 tokens (SANO, COLON, NEWLINE, ...), got %d", len(tokens))
	}
	if tokens[0].Type != TOKEN_SANO {
		t.Errorf("expected SANO, got %s", tokens[0].Type)
	}
	if tokens[1].Type != TOKEN_COLON {
		t.Errorf("expected COLON, got %s", tokens[1].Type)
	}
}

func TestLexerFinnishLetters(t *testing.T) {
	tokens := collectTokens("å ä ö Å Ä Ö")
	for _, tok := range tokens[:6] {
		if tok.Type != TOKEN_IDENT {
			t.Errorf("expected IDENT for Finnish letter, got %s %q", tok.Type, tok.Literal)
		}
	}
}
