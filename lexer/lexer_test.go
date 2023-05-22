package lexer

import (
	"testing"

	"github.com/Suryansh-23/amrit/token"
)

func TestNextToken(t *testing.T) {
	input := `mana five = 5|
mana ten = 10|
mana add = karya(x, y) {
	x + y|
}|
mana result = add(five, ten)|
!-/*5|
5 < 10 > 5|

agar (5 < 10) {
	labh सत्य|
} varna {
	labh asatya|
}

10 == 10|
10 != 9|
"haanji"
"acha thik hai"
[1,2]|
{"chota bheem": "motu patlu"}
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET_LATIN, "mana"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.TERM, "|"},
		{token.LET_LATIN, "mana"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.TERM, "|"},
		{token.LET_LATIN, "mana"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FN_LATIN, "karya"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.TERM, "|"},
		{token.RBRACE, "}"},
		{token.TERM, "|"},
		{token.LET_LATIN, "mana"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.TERM, "|"},
		//!-/*5; 5 < 10 > 5;
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.TERM, "|"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.TERM, "|"},
		{token.IF_LATIN, "agar"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN_LATIN, "labh"},
		{token.TRUE_LATIN, "satya"},
		{token.TERM, "|"},
		{token.RBRACE, "}"},
		{token.ELSE_LATIN, "varna"},
		{token.LBRACE, "{"},
		{token.RETURN_LATIN, "labh"},
		{token.FALSE_LATIN, "asatya"},
		{token.TERM, "|"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.TERM, "|"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.TERM, "|"},
		{token.STRING, "haanji"},
		{token.STRING, "acha thik hai"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.TERM, "|"},
		{token.LBRACE, "{"},
		{token.STRING, "chota bheem"},
		{token.COLON, ":"},
		{token.STRING, "motu patlu"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
		{token.EOF, ""},
	}
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
