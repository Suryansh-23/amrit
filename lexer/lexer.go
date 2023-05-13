package lexer

import (
	"unicode/utf8"

	"github.com/Suryansh-23/amrit/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           rune
}

// creates a lexer struct for parsing
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// reads a single character from the program (Note - it only supports ASCII - change ch's type from byte to rune for UTF support)
func (l *Lexer) readChar() {
	size := 0
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		// l.ch = rune(l.input[l.readPosition])
		// fmt.Println(l.input[l.readPosition : l.readPosition+3])
		// if !('a' <= l.ch && l.ch <= 'z' || 'A' <= l.ch && l.ch <= 'Z' || l.ch == '_') && 0x900 <= l.ch && l.ch <= 0x97F {
		// 	l.readDevanagiri()
		// }
		l.ch, size = utf8.DecodeRuneInString(l.input[l.readPosition:])
		// fmt.Printf("# %c\n", l.ch)
	}

	l.position = l.readPosition
	l.readPosition += size
}

// returns the subsequent token from the program string
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	// Main logic for parsing through the input string and thus generating resp. tokens
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(l.ch) + string(ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}

	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '|':
		tok = newToken(token.TERM, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = token.LookupIdentLatin(l.readIdentifier())
			tok.Type = token.LookupIdent(tok.Literal)

			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok

}

// returns a Token generated from tokenType and its Literal string
func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// reads through a valid continuous var name or keyword until it ends
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}

	// fmt.Printf("Read identifier: %s\n", l.input[position:l.position])
	return l.input[position:l.position]
}

// reads continuous digits i.e. numbers
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	return l.input[position:l.position]
}

// Add chars here to allow them in var names or keywords
func isLetter(ch rune) bool {
	// fmt.Println('a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || 0x900 <= ch && ch <= 0x97F)
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || 0x900 <= ch && ch <= 0x97F
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// eats up all the whitespace b/w the tokens (cuz they are just dividers)
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// peeks the next character but doesn't update the position as well as thre readPosition
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		ch, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
		return ch
	}
}
