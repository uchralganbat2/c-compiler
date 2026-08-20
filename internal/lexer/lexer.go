package lexer

// Lexer performs a single pass over the source, turning it into tokens.
// It tracks line and column so every token can point back at its origin.
type Lexer struct {
	Input []byte

	pos  int  // index of ch in Input; == len(Input) once exhausted
	ch   byte // char under examination, 0 at end of input
	line int  // 1-based line of ch
	col  int  // 1-based column of ch
}

func NewLexer(input []byte) *Lexer {
	l := &Lexer{Input: input, line: 1, col: 1}
	if len(input) > 0 {
		l.ch = input[0]
	}
	return l
}

// Tokens lexes the whole input and returns every token, ending with EOF.
// Illegal characters are reported as ILLEGAL tokens; lexing continues past them.
func (l *Lexer) Tokens() []Token {
	var toks []Token
	for {
		tok := l.NextToken()
		toks = append(toks, tok)
		if tok.Type == EOF {
			return toks
		}
	}
}

// NextToken returns the next token in the input, or EOF once it is exhausted.
func (l *Lexer) NextToken() Token {
	if bad, isBad := l.skipSpaceAndComments(); isBad {
		return bad
	}

	line, col := l.line, l.col

	switch l.ch {
	case 0:
		return Token{Type: EOF, Literal: "", Line: line, Column: col}
	case '=':
		return l.oneOrTwo('=', EQ, ASSIGN, line, col)
	case '!':
		// '!' only exists as part of "!=" in the current C subset.
		if l.peekChar() == '=' {
			return l.twoChar(NE, line, col)
		}
		return l.illegalChar(line, col)
	case '<':
		return l.oneOrTwo('=', LE, LT, line, col)
	case '>':
		return l.oneOrTwo('=', GE, GT, line, col)
	case '+':
		return l.singleChar(PLUS, line, col)
	case '-':
		return l.singleChar(MINUS, line, col)
	case '*':
		return l.singleChar(ASTERISK, line, col)
	case '/':
		// Comments were already consumed above, so this is division.
		return l.singleChar(SLASH, line, col)
	case ';':
		return l.singleChar(SEMICOLON, line, col)
	case ',':
		return l.singleChar(COMMA, line, col)
	case '(':
		return l.singleChar(LPAREN, line, col)
	case ')':
		return l.singleChar(RPAREN, line, col)
	case '{':
		return l.singleChar(LBRACE, line, col)
	case '}':
		return l.singleChar(RBRACE, line, col)
	}

	switch {
	case isLetter(l.ch):
		lit := l.readWhile(isIdentChar)
		return Token{Type: LookupIdent(lit), Literal: lit, Line: line, Column: col}
	case isDigit(l.ch):
		lit := l.readWhile(isDigit)
		// "42abc" is not a number followed by an identifier in C; report the
		// whole run as one illegal token so the message names the real problem.
		if isIdentChar(l.ch) {
			lit += l.readWhile(isIdentChar)
			return Token{Type: ILLEGAL, Literal: lit, Line: line, Column: col}
		}
		return Token{Type: INT_LIT, Literal: lit, Line: line, Column: col}
	default:
		return l.illegalChar(line, col)
	}
}

// skipSpaceAndComments advances past whitespace, // line comments and /* block
// comments. It reports an ILLEGAL token for an unterminated block comment.
func (l *Lexer) skipSpaceAndComments() (Token, bool) {
	for {
		switch {
		case isSpace(l.ch):
			l.readChar()
		case l.ch == '/' && l.peekChar() == '/':
			for l.ch != 0 && l.ch != '\n' {
				l.readChar()
			}
		case l.ch == '/' && l.peekChar() == '*':
			line, col := l.line, l.col
			l.readChar() // on '*'
			l.readChar() // past "/*"
			for l.ch != 0 && !(l.ch == '*' && l.peekChar() == '/') {
				l.readChar()
			}
			if l.ch == 0 {
				return Token{Type: ILLEGAL, Literal: "/*", Line: line, Column: col}, true
			}
			l.readChar() // on '/'
			l.readChar() // past "*/"
		default:
			return Token{}, false
		}
	}
}

func (l *Lexer) singleChar(t TokenType, line, col int) Token {
	lit := string(l.ch)
	l.readChar()
	return Token{Type: t, Literal: lit, Line: line, Column: col}
}

// twoChar consumes the current char and the one after it as a single token.
func (l *Lexer) twoChar(t TokenType, line, col int) Token {
	lit := string([]byte{l.ch, l.peekChar()})
	l.readChar()
	l.readChar()
	return Token{Type: t, Literal: lit, Line: line, Column: col}
}

// oneOrTwo emits two if the next char is next, otherwise one.
func (l *Lexer) oneOrTwo(next byte, two, one TokenType, line, col int) Token {
	if l.peekChar() == next {
		return l.twoChar(two, line, col)
	}
	return l.singleChar(one, line, col)
}

func (l *Lexer) illegalChar(line, col int) Token {
	return l.singleChar(ILLEGAL, line, col)
}

// readWhile consumes chars as long as pred holds and returns them.
func (l *Lexer) readWhile(pred func(byte) bool) string {
	start := l.pos
	for l.ch != 0 && pred(l.ch) {
		l.readChar()
	}
	return string(l.Input[start:l.pos])
}

func (l *Lexer) readChar() {
	if l.pos >= len(l.Input) {
		l.ch = 0
		return
	}
	if l.ch == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	l.pos++
	if l.pos < len(l.Input) {
		l.ch = l.Input[l.pos]
	} else {
		l.ch = 0
	}
}

func (l *Lexer) peekChar() byte {
	if l.pos+1 >= len(l.Input) {
		return 0
	}
	return l.Input[l.pos+1]
}

func isSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' || ch == '\v' || ch == '\f'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isIdentChar(ch byte) bool {
	return isLetter(ch) || isDigit(ch)
}
