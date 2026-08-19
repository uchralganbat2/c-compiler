package lexer

import "testing"

// want is a token expectation without line/column, for tests that only care
// about the type/literal sequence.
type want struct {
	typ TokenType
	lit string
}

func lex(t *testing.T, input string) []Token {
	t.Helper()
	return NewLexer([]byte(input)).Tokens()
}

func assertTokens(t *testing.T, input string, wants []want) {
	t.Helper()
	got := lex(t, input)
	if len(got) != len(wants) {
		t.Fatalf("input %q: got %d tokens %v, want %d", input, len(got), got, len(wants))
	}
	for i, w := range wants {
		if got[i].Type != w.typ || got[i].Literal != w.lit {
			t.Errorf("input %q: token %d = %v, want %s(%q)", input, i, got[i], w.typ, w.lit)
		}
	}
}

func TestNextTokenDeclaration(t *testing.T) {
	// The README's own example.
	assertTokens(t, "int x = 42;", []want{
		{INT, "int"},
		{IDENT, "x"},
		{ASSIGN, "="},
		{INT_LIT, "42"},
		{SEMICOLON, ";"},
		{EOF, ""},
	})
}

func TestNextTokenOperators(t *testing.T) {
	assertTokens(t, "+ - * / = == != < <= > >=", []want{
		{PLUS, "+"},
		{MINUS, "-"},
		{ASTERISK, "*"},
		{SLASH, "/"},
		{ASSIGN, "="},
		{EQ, "=="},
		{NE, "!="},
		{LT, "<"},
		{LE, "<="},
		{GT, ">"},
		{GE, ">="},
		{EOF, ""},
	})
}

func TestNextTokenOperatorsUnspaced(t *testing.T) {
	// Two-char operators must win over their single-char prefixes even with no
	// separating whitespace.
	assertTokens(t, "a==b<=c>=d!=e=f", []want{
		{IDENT, "a"}, {EQ, "=="},
		{IDENT, "b"}, {LE, "<="},
		{IDENT, "c"}, {GE, ">="},
		{IDENT, "d"}, {NE, "!="},
		{IDENT, "e"}, {ASSIGN, "="},
		{IDENT, "f"},
		{EOF, ""},
	})
}

func TestNextTokenDelimiters(t *testing.T) {
	assertTokens(t, "; , ( ) { }", []want{
		{SEMICOLON, ";"},
		{COMMA, ","},
		{LPAREN, "("},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{RBRACE, "}"},
		{EOF, ""},
	})
}

func TestNextTokenKeywordsAndIdents(t *testing.T) {
	// Keywords are exact matches: "integer" and "iff" are identifiers.
	assertTokens(t, "int return if else while integer iff returned _x x1 X", []want{
		{INT, "int"},
		{RETURN, "return"},
		{IF, "if"},
		{ELSE, "else"},
		{WHILE, "while"},
		{IDENT, "integer"},
		{IDENT, "iff"},
		{IDENT, "returned"},
		{IDENT, "_x"},
		{IDENT, "x1"},
		{IDENT, "X"},
		{EOF, ""},
	})
}

func TestNextTokenProgram(t *testing.T) {
	input := `int main() {
	int x = 3;
	while (x < 10) {
		x = x + 1;
	}
	if (x == 10) {
		return x / 2;
	} else {
		return 0;
	}
}`
	assertTokens(t, input, []want{
		{INT, "int"}, {IDENT, "main"}, {LPAREN, "("}, {RPAREN, ")"}, {LBRACE, "{"},
		{INT, "int"}, {IDENT, "x"}, {ASSIGN, "="}, {INT_LIT, "3"}, {SEMICOLON, ";"},
		{WHILE, "while"}, {LPAREN, "("}, {IDENT, "x"}, {LT, "<"}, {INT_LIT, "10"}, {RPAREN, ")"}, {LBRACE, "{"},
		{IDENT, "x"}, {ASSIGN, "="}, {IDENT, "x"}, {PLUS, "+"}, {INT_LIT, "1"}, {SEMICOLON, ";"},
		{RBRACE, "}"},
		{IF, "if"}, {LPAREN, "("}, {IDENT, "x"}, {EQ, "=="}, {INT_LIT, "10"}, {RPAREN, ")"}, {LBRACE, "{"},
		{RETURN, "return"}, {IDENT, "x"}, {SLASH, "/"}, {INT_LIT, "2"}, {SEMICOLON, ";"},
		{RBRACE, "}"},
		{ELSE, "else"}, {LBRACE, "{"},
		{RETURN, "return"}, {INT_LIT, "0"}, {SEMICOLON, ";"},
		{RBRACE, "}"},
		{RBRACE, "}"},
		{EOF, ""},
	})
}

func TestNextTokenComments(t *testing.T) {
	input := `// leading line comment
int /* inline */ x; /* multi
line
comment */ x = 1; // trailing`
	assertTokens(t, input, []want{
		{INT, "int"},
		{IDENT, "x"},
		{SEMICOLON, ";"},
		{IDENT, "x"},
		{ASSIGN, "="},
		{INT_LIT, "1"},
		{SEMICOLON, ";"},
		{EOF, ""},
	})
}

func TestNextTokenDivisionIsNotAComment(t *testing.T) {
	assertTokens(t, "a / b", []want{
		{IDENT, "a"}, {SLASH, "/"}, {IDENT, "b"}, {EOF, ""},
	})
}

func TestNextTokenEmptyAndBlankInput(t *testing.T) {
	for _, input := range []string{"", "   ", "\n\n\t", "// only a comment"} {
		got := lex(t, input)
		if len(got) != 1 || got[0].Type != EOF {
			t.Errorf("input %q: got %v, want a single EOF", input, got)
		}
	}
}

func TestNextTokenIsIdempotentAtEOF(t *testing.T) {
	l := NewLexer([]byte("x"))
	if tok := l.NextToken(); tok.Type != IDENT {
		t.Fatalf("first token = %v, want IDENT", tok)
	}
	for i := 0; i < 3; i++ {
		if tok := l.NextToken(); tok.Type != EOF {
			t.Fatalf("call %d past the end = %v, want EOF", i, tok)
		}
	}
}

func TestNextTokenIllegal(t *testing.T) {
	tests := []struct {
		input string
		want  []want
	}{
		// Unknown character, then lexing recovers.
		{"@x", []want{{ILLEGAL, "@"}, {IDENT, "x"}, {EOF, ""}}},
		// '!' is only valid as part of "!=" in this subset.
		{"!x", []want{{ILLEGAL, "!"}, {IDENT, "x"}, {EOF, ""}}},
		// A number running into letters is one bad token, not INT_LIT + IDENT.
		{"42abc;", []want{{ILLEGAL, "42abc"}, {SEMICOLON, ";"}, {EOF, ""}}},
		// Unterminated block comment.
		{"int /* nope", []want{{INT, "int"}, {ILLEGAL, "/*"}, {EOF, ""}}},
	}
	for _, tt := range tests {
		assertTokens(t, tt.input, tt.want)
	}
}

func TestNextTokenPositions(t *testing.T) {
	input := "int x;\n  x = 42;\r\nreturn x;"
	wants := []Token{
		{INT, "int", 1, 1},
		{IDENT, "x", 1, 5},
		{SEMICOLON, ";", 1, 6},
		{IDENT, "x", 2, 3},
		{ASSIGN, "=", 2, 5},
		{INT_LIT, "42", 2, 7},
		{SEMICOLON, ";", 2, 9},
		{RETURN, "return", 3, 1},
		{IDENT, "x", 3, 8},
		{SEMICOLON, ";", 3, 9},
		{EOF, "", 3, 10},
	}
	got := lex(t, input)
	if len(got) != len(wants) {
		t.Fatalf("got %d tokens %v, want %d", len(got), got, len(wants))
	}
	for i, w := range wants {
		if got[i] != w {
			t.Errorf("token %d = %v, want %v", i, got[i], w)
		}
	}
}

func TestPositionAfterMultiLineComment(t *testing.T) {
	// The token after a block comment spanning lines reports the real line.
	got := lex(t, "/* one\ntwo */ x")
	if len(got) != 2 {
		t.Fatalf("got %v, want IDENT then EOF", got)
	}
	if want := (Token{IDENT, "x", 2, 8}); got[0] != want {
		t.Errorf("got %v, want %v", got[0], want)
	}
}

func TestTokenTypeString(t *testing.T) {
	if got := INT_LIT.String(); got != "INT_LIT" {
		t.Errorf("INT_LIT.String() = %q, want %q", got, "INT_LIT")
	}
	if got := TokenType(9999).String(); got != "TokenType(9999)" {
		t.Errorf("unknown type String() = %q, want %q", got, "TokenType(9999)")
	}
}

func TestLookupIdent(t *testing.T) {
	if got := LookupIdent("while"); got != WHILE {
		t.Errorf("LookupIdent(\"while\") = %v, want WHILE", got)
	}
	if got := LookupIdent("whilst"); got != IDENT {
		t.Errorf("LookupIdent(\"whilst\") = %v, want IDENT", got)
	}
}

func TestTokenString(t *testing.T) {
	tok := Token{Type: IDENT, Literal: "x", Line: 2, Column: 5}
	if got, want := tok.String(), `IDENT("x") 2:5`; got != want {
		t.Errorf("Token.String() = %q, want %q", got, want)
	}
}
