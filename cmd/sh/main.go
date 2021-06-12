package main

import (
	"os"

	"gitlab.com/yarbelk/slimbox/lib/sh"
)

// Lex that has error implemented

func main() {

	sh.ShParse(sh.NewLexer(os.Stdin))
}
