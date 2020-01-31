package object

import "fmt"

import "github.com/wreulicke/go-sandbox/go-interpreter/monkey/ast"

import "bytes"

import "strings"

var typeNames = []string{
	"INTEGER",
	"BOOLEAN",
	"FUNCTION",
	"NULL",
	"RETURN",
	"ERROR",
}

type ObjectType int

func (o ObjectType) String() string {
	return typeNames[o]
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	INTEGER ObjectType = iota
	BOOLEAN
	FUNCTION
	NULL
	RETURN
	ERROR
)

type Integer struct {
	Value int64
}

func (n *Integer) Type() ObjectType {
	return INTEGER
}

func (n *Integer) Inspect() string {
	return fmt.Sprintf("%d", n.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type Null struct {
}

func (b *Null) Type() ObjectType {
	return NULL
}

func (b *Null) Inspect() string {
	return "null"
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, v := range f.Parameters {
		params = append(params, v.String())
	}
	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}
