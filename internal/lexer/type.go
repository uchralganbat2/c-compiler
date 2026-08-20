package lexer

import "strconv"

type TokenType int

const (
	// Special
	ILLEGAL TokenType = iota
	EOF

	// Literals
	INT_LIT // 0, 42, ...
	IDENT   // x, main, ...

	// Keywords (keep together so you can range-check if needed)
	INT
	RETURN
	IF
	ELSE
	WHILE
	// add more as you expand the C subset

	// Operators
	ASSIGN   // =
	PLUS     // +
	MINUS    // -
	ASTERISK // *
	SLASH    // /
	EQ       // ==
	NE       // !=
	LT       // <
	LE       // <=
	GT       // >
	GE       // >=

	// Delimiters
	SEMICOLON // ;
	COMMA     // ,
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
)

var tokenNames = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	INT_LIT: "INT_LIT",
	IDENT:   "IDENT",

	INT:    "INT",
	RETURN: "RETURN",
	IF:     "IF",
	ELSE:   "ELSE",
	WHILE:  "WHILE",

	ASSIGN:   "ASSIGN",
	PLUS:     "PLUS",
	MINUS:    "MINUS",
	ASTERISK: "ASTERISK",
	SLASH:    "SLASH",
	EQ:       "EQ",
	NE:       "NE",
	LT:       "LT",
	LE:       "LE",
	GT:       "GT",
	GE:       "GE",

	SEMICOLON: "SEMICOLON",
	COMMA:     "COMMA",
	LPAREN:    "LPAREN",
	RPAREN:    "RPAREN",
	LBRACE:    "LBRACE",
	RBRACE:    "RBRACE",
}

func (t TokenType) String() string {
	if int(t) < len(tokenNames) && tokenNames[t] != "" {
		return tokenNames[t]
	}
	return "TokenType(" + strconv.Itoa(int(t)) + ")"
}

// keywords maps reserved words to their token type; anything else is an IDENT.
var keywords = map[string]TokenType{
	"int":    INT,
	"return": RETURN,
	"if":     IF,
	"else":   ELSE,
	"while":  WHILE,
}

// LookupIdent returns the keyword token type for ident, or IDENT if it is not
// a reserved word.
func LookupIdent(ident string) TokenType {
	if t, ok := keywords[ident]; ok {
		return t
	}
	return IDENT
}
