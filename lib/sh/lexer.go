package sh

import (
	"fmt"
)

const eof = 0

func getOperator(s string) int {
	switch s {
	case "&&":
		return AND_IF
	case "||":
		return OR_IF
	case ";;":
		return DSEMI
	case "<<":
		return DLESS
	case "<<-":
		return DLESSDASH
	case "<&":
		return LESSAND
	case "<>":
		return LESSGREAT
	case ">>":
		return DGREAT
	case ">&":
		return GREATAND
	case ">|":
		return CLOBBER
	default:
		panic(fmt.Sprintf("Syntax Error operator %s", s))
	}
}

func getReserved(s string) int {
	switch s {
	case "if":
		return If
	case "then":
		return Then
	case "else":
		return Else
	case "elif":
		return Elif
	case "fi":
		return Fi
	case "do":
		return Do
	case "done":
		return Done
	case "case":
		return Case
	case "esac":
		return Esac
	case "while":
		return While
	case "until":
		return Until
	case "for":
		return For
	case "{":
		return Lbrace
	case "}":
		return Rbrace
	case "!":
		return Bang
	case "in":
		return In
	default:
		panic(fmt.Sprintf("Syntax Error reserved word %s", s))
	}
}
