package token

type TokenType int

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL TokenType = iota
	EOF

	IDENT
	INT

	ASSIGN
	PLUS
	MINUS
	BANG
	ASTERISK
	SLASH

	EQ
	NOT_EQ

	LT
	GT

	COMMA
	SEMICOLON

	LPAREN
	RPAREN

	LBRACE
	RBRACE

	FUNCTION
	LET
	RETURN
	TRUE
	FALSE
	IF
	ELSE
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
