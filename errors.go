package main

import (
	"fmt"
	"os"
)

// ShowError prints a Finnish error message in the format "VIRHE rivillä N: viesti" and exits.
func ShowError(line int, format string, a ...interface{}) {
	fmt.Printf("VIRHE rivillä %d: %s\n", line, fmt.Sprintf(format, a...))
	os.Exit(1)
}
