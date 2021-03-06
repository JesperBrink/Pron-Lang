package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL" // Tokens/characters that we don't know about
	EOF     = "EOF"     // End of File

	// Identifiers + literals
	IDENT  = "IDENT"  //add, foobar, x, y, ...
	INT    = "INT"    // 42
	STRING = "STRING" // "Hello World!"
	REAL   = "REAL"   // 42.0, 4.5, 3.15, ...

	// Operators
	ASSIGN    = "="
	PLUS      = "+"
	MINUS     = "-"
	BANG      = "!"
	ASTERISK  = "*"
	SLASH     = "/"
	MODULO    = "%"
	EQ        = "=="
	NOT_EQ    = "!="
	INCREMENT = "++"
	DECREMENT = "--"

	LT = "<"
	GT = ">"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	DOT       = "."

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	VAR      = "VAR"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	ELIF     = "ELIF"
	RETURN   = "RETURN"
	FOR      = "FOR"
	FROM     = "FROM"
	TO       = "TO"
	IN       = "IN"
	CLASS    = "CLASS"
	INIT     = "INIT"
	THIS     = "THIS"
	NEW      = "NEW"

	// Comments
	STARTBLOCKCOMMENT = "/*"
	ENDBLOCKCOMMENT   = "*/"
)

var keywords = map[string]TokenType{
	"func":   FUNCTION,
	"var":    VAR,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"elif":   ELIF,
	"return": RETURN,
	"for":    FOR,
	"from":   FROM,
	"to":     TO,
	"in":     IN,
	"class":  CLASS,
	"Init":   INIT,
	"this":   THIS,
	"new":    NEW,
}

// Returns the TokenType that matches the ident given as argument.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
