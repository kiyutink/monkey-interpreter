package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"monkey-interpreter/ast"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: INTEGER_OBJ, Value: uint64(i.Value)}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) HashKey() HashKey {
	value := 0
	if b.Value {
		value = 1
	}
	return HashKey{Type: BOOLEAN_OBJ, Value: uint64(value)}
}

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
	buf := bytes.Buffer{}

	params := []string{}

	for _, p := range f.Parameters {
		params = append(params, p.Value)
	}

	buf.WriteString("fn")
	buf.WriteString("(")
	buf.WriteString(strings.Join(params, ", "))
	buf.WriteString(")")
	buf.WriteString("{\n")
	buf.WriteString(f.Body.String())
	buf.WriteString("\n}")

	return buf.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return `"` + s.Value + `"` }
func (s *String) HashKey() HashKey {
	hk := HashKey{Type: STRING_OBJ}

	h := fnv.New64a()
	h.Write([]byte(s.Value))
	hk.Value = h.Sum64()
	return hk
}

type (
	BuiltinFn func(args ...Object) Object
	Builtin   struct {
		Fn BuiltinFn
	}
)

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}
type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	buf := bytes.Buffer{}

	pairs := []string{}

	for _, pair := range h.Pairs {
		pairs = append(pairs, strings.Join([]string{pair.Key.Inspect(), pair.Value.Inspect()}, " : "))
	}

	buf.WriteString("{\n")
	buf.WriteString(strings.Join(pairs, ",\n"))
	buf.WriteString("\n}")

	return buf.String()
}

type Hashable interface {
	HashKey() HashKey
}
