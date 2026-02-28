package lexer

type Token struct {
	Type    TokenType
	Literal string // raw slice from source, e.g. "42" or "x"
	Line    int    // 1-based for user-facing messages
	Column  int    // optional: 1-based column
}
