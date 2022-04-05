//lint:file-ignore U1000 fake methods for separating interfacess

package ast

import (
	"bytes"
	"monkey-interpreter/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	buf := bytes.Buffer{}
	buf.WriteString(ls.TokenLiteral() + " " + ls.Name.Value)
	buf.WriteString(" = ")

	if ls.Value != nil {
		buf.WriteString(ls.Value.String())
	}

	buf.WriteString(";")

	return buf.String()
}

type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
	buf := bytes.Buffer{}

	buf.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		buf.WriteString(rs.ReturnValue.String())
	}

	buf.WriteString(";")
	return buf.String()
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es != nil {
		return es.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PrefixExpression) String() string {
	buf := bytes.Buffer{}
	buf.WriteString("(")
	buf.WriteString(pe.Operator)
	buf.WriteString(pe.Right.String())
	buf.WriteString(")")
	return buf.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *InfixExpression) String() string {
	buf := bytes.Buffer{}
	buf.WriteString("(")
	buf.WriteString(ie.Left.String())
	buf.WriteString(" " + ie.Operator + " ")
	buf.WriteString(ie.Right.String())
	buf.WriteString(")")
	return buf.String()
}

type BooleanExpression struct {
	Token token.Token
	Value bool
}

func (b *BooleanExpression) expressionNode() {}
func (b *BooleanExpression) TokenLiteral() string {
	return b.Token.Literal
}
func (b *BooleanExpression) String() string {
	return b.Token.Literal
}

type IfExpression struct {
	Token       token.Token // The "if" token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var buf bytes.Buffer
	buf.WriteString("if")
	buf.WriteString(ie.Condition.String())
	buf.WriteString(" ")
	buf.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		buf.WriteString("else ")
		buf.WriteString(ie.Alternative.String())
	}
	return buf.String()
}

type BlockStatement struct {
	Token      token.Token // The "{" token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *BlockStatement) String() string {
	var buf bytes.Buffer
	for _, s := range bs.Statements {
		buf.WriteString(s.String())
	}
	return buf.String()
}

type FunctionLiteral struct {
	Token      token.Token // The "fn" token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}
func (fl *FunctionLiteral) String() string {
	var buf bytes.Buffer
	params := []string{}
	for _, param := range fl.Parameters {
		params = append(params, param.TokenLiteral())
	}
	buf.WriteString(fl.TokenLiteral())
	buf.WriteString("(")
	buf.WriteString(strings.Join(params, ", "))
	buf.WriteString(")")
	buf.WriteString("{")
	buf.WriteString(fl.Body.String())
	buf.WriteString("}")

	return buf.String()
}

type CallExpression struct {
	Token     token.Token // The "(" token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	buf := bytes.Buffer{}

	args := []string{}

	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}

	buf.WriteString(ce.Function.String())
	buf.WriteString("(")
	buf.WriteString(strings.Join(args, ", "))
	buf.WriteString(")")
	return buf.String()
}

type StringLiteral struct {
	Token token.Token // The String token
	Value string
}

func (s *StringLiteral) expressionNode()      {}
func (s *StringLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *StringLiteral) String() string {
	buf := bytes.Buffer{}
	buf.WriteByte('"')
	buf.WriteString(s.Value)
	buf.WriteByte('"')

	return buf.String()
}

type ArrayLiteral struct {
	Token    token.Token // the "[" token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	buf := bytes.Buffer{}
	buf.WriteByte('[')
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	buf.WriteString(strings.Join(elements, ", "))
	buf.WriteByte(']')
	return buf.String()
}
