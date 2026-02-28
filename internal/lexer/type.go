package lexer

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
