package ast

import (
	"testing"

	"github.com/Suryansh-23/amrit/token"
)

func TestString(t *testing.T) {
	program := &Program{Statements: []Statement{
		&LetStatement{
			Token: token.Token{Type: token.LET_LATIN, Literal: "mana"},
			Name: &Identifier{
				Token: token.Token{Type: token.IDENT, Literal: "myVar"},
				Value: "myVar",
			},
			Value: &Identifier{
				Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
				Value: "anotherVar",
			},
		},
	},
	}

	if program.String() != "mana myVar = anotherVar|" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
