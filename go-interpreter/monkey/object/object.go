package object

import "fmt"

var typeNames = []string{
	"INTEGER",
	"BOOLEAN",
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
