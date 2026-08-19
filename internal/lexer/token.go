package lexer

import "fmt"

type Token struct {
	Type    TokenType
	Literal string // raw slice from source, e.g. "42" or "x"
	Line    int    // 1-based for user-facing messages
	Column  int    // optional: 1-based column
}

// String renders a token as TYPE("literal") at line:column, for error messages
// and test failures.
func (t Token) String() string {
	return fmt.Sprintf("%s(%q) %d:%d", t.Type, t.Literal, t.Line, t.Column)
}
