package sh_test

import (
	"strings"
	"testing"

	"gitlab.com/yarbelk/slimbox/lib/sh"
)

// Testing the tokenizing lexer

func TestCommands(t *testing.T) {
	var tests = []struct {
		name     string
		expected int
		given    string
	}{
		{"AND_IF", sh.AND_IF, "&&"},
		{"OR_IF", sh.OR_IF, "||"},
		{"DSEMI", sh.DSEMI, ";;"},
		{"DLESS", sh.DLESS, "<<"},
		{"LESSAND", sh.LESSAND, "<&"},
		{"DGREAT", sh.DGREAT, ">>"},
		{"GREATAND", sh.GREATAND, ">&"},
		{"LESSGREAT", sh.LESSGREAT, "<>"},
		{"DLESSDASH", sh.DLESSDASH, "<<-"},
		{"CLOBBER", sh.CLOBBER, ">|"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := sh.NewLexer(strings.NewReader(tt.given))
			actual := tokenizer.Lex(&sh.ShSymType{})
			if actual != tt.expected {
				t.Errorf("given(%s): expected %d, actual %d", tt.given, tt.expected, actual)
			}
		})
	}
}

func TestReservedWords(t *testing.T) {
	var tests = []struct {
		name     string
		expected int
		given    string
	}{
		{"If", sh.If, "if"},
		{"Then", sh.Then, "then"},
		{"Else", sh.Else, "else"},
		{"Elif", sh.Elif, "elif"},
		{"Fi", sh.Fi, "fi"},
		{"Do", sh.Do, "do"},
		{"Done", sh.Done, "done"},
		{"Case", sh.Case, "case"},
		{"Esac", sh.Esac, "esac"},
		{"While", sh.While, "while"},
		{"Until", sh.Until, "until"},
		{"For", sh.For, "for"},
		{"{", sh.Lbrace, "{"},
		{"}", sh.Rbrace, "}"},
		{"!", sh.Bang, "!"},
		{"in", sh.In, "in"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := sh.NewLexer(strings.NewReader(tt.given))
			actual := tokenizer.Lex(&sh.ShSymType{})
			if actual != tt.expected {
				t.Errorf("given(%s): expected %d, actual %d", tt.given, tt.expected, actual)
			}
		})
	}
}
