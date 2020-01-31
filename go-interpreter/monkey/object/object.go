package object

import "fmt"

type ObjectType int

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	INTEGER ObjectType = iota
	BOOLEAN
	NULL
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
