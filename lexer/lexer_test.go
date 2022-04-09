package lexer

import (
	"testing"

	"monkey-interpreter/token"
)

func TestNextToken(t *testing.T) {
	input := `
	let five = 5;
	let ten = 10;
	let add = fn(x, y) {
	 x + y;
	};
	let result = add(five, ten);
	!-/*5;
	5 < 10 > 5;
	if (5 < 10) {
		return true;
	} else {
		return false;
	}

	10 == 10; 10 != 9;

	"foobar";
	"foo bar";
	test[a];
	arr[3];
	[1, 2, 3];
	{ "a": "b" };
`

	tests := []struct {
		expectedToken   token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.STRING, "foobar"},
		{token.SEMICOLON, ";"},
		{token.STRING, "foo bar"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "test"},
		{token.LBRACKET, "["},
		{token.IDENT, "a"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "arr"},
		{token.LBRACKET, "["},
		{token.INT, "3"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.COMMA, ","},
		{token.INT, "3"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.LBRACE, "{"},
		{token.STRING, "a"},
		{token.COLON, ":"},
		{token.STRING, "b"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
	}

	l := New(input)

	for _, test := range tests {
		tok := l.NextToken()

		if tok.Type != test.expectedToken {
			t.Fatalf("Expected token type %v but received %v", test.expectedToken, tok.Type)
		}

		if tok.Literal != test.expectedLiteral {
			t.Fatalf("Expected literal %v but received %v", test.expectedLiteral, tok.Literal)
		}
	}
}
