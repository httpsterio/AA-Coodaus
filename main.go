package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Käyttö: %s <tiedosto.aa>\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]
	contentBytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Tiedostoa ei voitu lukea: %v\n", err)
		os.Exit(1)
	}

	src := string(contentBytes)
	// Handle shebang by converting it to a comment to preserve line numbering
	if strings.HasPrefix(src, "#!") {
		idx := strings.Index(src, "\n")
		if idx >= 0 {
			src = "--" + src[2:idx] + src[idx:]
		}
	}

	// Pre-check source code for forbidden meme structures (like "Ota yhteys")
	PreCheckSource(src)

	// Lex, parse, and evaluate
	lexer := NewLexer(src)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	env := NewEnv()
	evaluator := NewEvaluator(env)
	evaluator.Evaluate(program)
}
