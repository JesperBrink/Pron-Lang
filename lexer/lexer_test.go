package lexer

import (
	"Pron-Lang/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input :=
		`var five = 5;
		var ten = 10;
		var ten2 = 10;
		
		var add = func(x, y) {
			x + y;
		};
		
		var result = add(five, ten);
		!-/*5;
		5 < 10 > 5;

		if (5 < 10) {
			return true;
		} elif (5 > 10) {
			return "what";
		} else {
			return false;
		}

		10 == 10;
		10 != 9;
		"foobar"
		"foo bar"
		var myList = [1, 2];
		{"foo": "bar"}
		for (i from 0 to x) {
			puts(i)
		}
		for (i in myList) {
			puts(i)
		}

		class Person {
			init(this.name) {
				puts("Initialized")
			}
		}
		var p1 = new Person("Hans")
		i++
		i--
		`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.VAR, "var"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.VAR, "var"},
		{token.IDENT, "ten2"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.VAR, "var"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "func"},
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
		{token.VAR, "var"},
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
		{token.ELIF, "elif"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.GT, ">"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.STRING, "what"},
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
		{token.STRING, "foo bar"},
		{token.VAR, "var"},
		{token.IDENT, "myList"},
		{token.ASSIGN, "="},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
		{token.FOR, "for"},
		{token.LPAREN, "("},
		{token.IDENT, "i"},
		{token.FROM, "from"},
		{token.INT, "0"},
		{token.TO, "to"},
		{token.IDENT, "x"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "puts"},
		{token.LPAREN, "("},
		{token.IDENT, "i"},
		{token.RPAREN, ")"},
		{token.RBRACE, "}"},
		{token.FOR, "for"},
		{token.LPAREN, "("},
		{token.IDENT, "i"},
		{token.IN, "in"},
		{token.IDENT, "myList"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "puts"},
		{token.LPAREN, "("},
		{token.IDENT, "i"},
		{token.RPAREN, ")"},
		{token.RBRACE, "}"},
		{token.CLASS, "class"},
		{token.IDENT, "Person"},
		{token.LBRACE, "{"},
		{token.INIT, "init"},
		{token.LPAREN, "("},
		{token.THIS, "this"},
		{token.DOT, "."},
		{token.IDENT, "name"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "puts"},
		{token.LPAREN, "("},
		{token.STRING, "Initialized"},
		{token.RPAREN, ")"},
		{token.RBRACE, "}"},
		{token.RBRACE, "}"},
		{token.VAR, "var"},
		{token.IDENT, "p1"},
		{token.ASSIGN, "="},
		{token.NEW, "new"},
		{token.IDENT, "Person"},
		{token.LPAREN, "("},
		{token.STRING, "Hans"},
		{token.RPAREN, ")"},
		{token.IDENT, "i"},
		{token.INCREMENT, "++"},
		{token.IDENT, "i"},
		{token.DECREMENT, "--"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}

	}
}
