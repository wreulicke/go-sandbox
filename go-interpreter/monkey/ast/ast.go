package ast

import "github.com/wreulicke/go-sandbox/go-interpreter/monkey/token"

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

type statement struct{}

func (s *statement) statementNode() {}

type expression struct{}

func (e *expression) expressionNode() {}

type LetStatement struct {
	statement
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type ReturnStatement struct {
	statement
	Token       token.Token
	ReturnValue Expression
}

func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

type Identifier struct {
	expression
	Token token.Token
	Value string
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}
