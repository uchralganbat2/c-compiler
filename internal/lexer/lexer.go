package lexer

type Lexer struct {
	Input []byte
}

func NewLexer(input []byte) *Lexer {
	return &Lexer{Input: input}
}

func (l *Lexer) NextToken() Token {
	return Token{Type: EOF, Literal: "", Line: 0, Column: 0}
}
