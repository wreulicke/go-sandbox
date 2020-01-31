package object

import "fmt"

type ObjectType int

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	NUMBER ObjectType = iota
	BOOLEAN
	NULL
)

type Number struct {
	Value string
}

func (n *Number) Type() ObjectType {
	return NUMBER
}

func (n *Number) Inspect() string {
	return n.Value
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
