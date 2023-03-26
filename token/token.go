package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "ANK"

	// 1343456
	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"
	EQ       = "=="
	NOT_EQ   = "!="

	// Delimiters
	COMMA  = ","
	TERM   = "|"
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	//LATIN
	FN_LATIN     = "karya"
	LET_LATIN    = "mana"
	TRUE_LATIN   = "satya"
	FALSE_LATIN  = "asatya"
	IF_LATIN     = "agar"
	ELSE_LATIN   = "varna"
	RETURN_LATIN = "labh"

	// DEVANAGIRI
	// FN_DEVANAGIRI     = "कार्य"
	// LET_DEVANAGIRI    = "माना"
	// TRUE_DEVANAGIRI   = "सत्य"
	// FALSE_DEVANAGIRI  = "असत्य"
	// IF_DEVANAGIRI     = "अगर"
	// ELSE_DEVANAGIRI   = "वरना"
	// RETURN_DEVANAGIRI = "लाभ"
)

var keywords_latin = map[string]TokenType{
	"karya":  FN_LATIN,
	"mana":   LET_LATIN,
	"satya":  TRUE_LATIN,
	"asatya": FALSE_LATIN,
	"agar":   IF_LATIN,
	"varna":  ELSE_LATIN,
	"labh":   RETURN_LATIN,
}

var keywords_devanagiri = map[string]TokenType{
	"कार्य": FN_LATIN,
	"माना":  LET_LATIN,
	"सत्य":  TRUE_LATIN,
	"असत्य": FALSE_LATIN,
	"अगर":   IF_LATIN,
	"वरना":  ELSE_LATIN,
	"लाभ":   RETURN_LATIN,
}

var devanagiri_to_latin = map[string]string{
	"कार्य": "karya",
	"माना":  "mana",
	"सत्य":  "satya",
	"असत्य": "asatya",
	"अगर":   "agar",
	"वरना":  "varna",
	"लाभ":   "labh",
}

func LookupIdent(ident string) TokenType {
	// fmt.Println(ident)
	if tok, ok := keywords_latin[ident]; ok {
		return tok
	} else if tok, ok := keywords_devanagiri[ident]; ok {
		return tok
	} else {
		return IDENT
	}
}

func LookupIdentLatin(ident string) string {
	// fmt.Println(ident)
	if str, ok := devanagiri_to_latin[ident]; ok {
		return str
	} else {
		return ident
	}
}
