package parser

import (
	"fmt"
	"monkey-interpreter/ast"
	"monkey-interpreter/lexer"
	"testing"
)

func TestReturnStatements(t *testing.T) {

	tests := []struct {
		input              string
		expectedExpression interface{}
	}{
		{"return 5;", 5},
		{"return x;", "x"},
		{"return foo;", "foo"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		if program == nil {
			t.Fatalf("ParseProgram returned nil")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("Expected 1 statements, instead got %v", len(program.Statements))
		}

		returnStatement, ok := program.Statements[0].(*ast.ReturnStatement)

		if !ok {
			t.Errorf("Expected type to be *ast.ReturnStatement, instead got %T", program.Statements[0])
			continue
		}

		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("Expected a RETURN statement, instead got %v", program.Statements[0].TokenLiteral())
		}

		if !testLiteralExpression(t, returnStatement.ReturnValue, tt.expectedExpression) {
			return
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let a = 3;", "a", 3},
		{"let b = false;", "b", false},
		{"let foo = bar;", "foo", "bar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if program == nil {
			t.Fatalf("ParseProgram returned nil")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("Program.Statements does not contain 3 statements. got = %v", len(program.Statements))
		}

		statement := program.Statements[0]
		if !testLetStatement(t, statement, tt.expectedIdentifier) {
			return
		}

		val := statement.(*ast.LetStatement).Value

		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	if len(p.Errors()) == 0 {
		return
	}

	t.Errorf("parser has %v errors", len(p.Errors()))

	for _, msg := range p.Errors() {
		t.Errorf("parser error: %v", msg)
	}

	t.FailNow()
}

func testLetStatement(t *testing.T, statement ast.Statement, expectedIdentifer string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("TokenLiteral should be let, got %v", statement.TokenLiteral())
		return false
	}

	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("Statement should be LetStatement, got %T", statement)
		return false
	}

	if letStatement.Name.Value != expectedIdentifer {
		t.Errorf("Expected statement to have name %v, got %v", expectedIdentifer, letStatement.Name.Value)
		return false
	}

	if letStatement.Name.TokenLiteral() != expectedIdentifer {
		t.Errorf("Expected letStatement.Name.TokenLiteral() to be %v, got %v", expectedIdentifer, letStatement.Name.TokenLiteral())
	}

	return true
}

func TestParsingStrings(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"foobar"`, "foobar"},
		{`"foo bar"`, "foo bar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("Expected program to have 1 statement, instead got %v", len(program.Statements))
		}

		expr := program.Statements[0].(*ast.ExpressionStatement)
		str, ok := expr.Expression.(*ast.StringLiteral)
		if !ok {
			t.Errorf("Expected to receive a StringLiteral, instead got %T", expr)
		}

		if str.Literal != tt.expected {
			t.Errorf("Expected string literal to be equal to %v, instead got %v", tt.expected, str.Literal)
		}
	}

}
func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program was expected to have 1 statement, instead got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Expected statement to be an expression statement, instead got %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("Expected expression to be an identifier, instead got %T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Fatalf("Expected identifier to be foobar, instead got %v", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("Expected TokenLiteral to be foobar, instead got %v", ident.TokenLiteral())
	}
}

func TestIntegerExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program was expected to have 1 statement, instead got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Expected statement to be an expression statement, instead got %T", program.Statements[0])
	}

	if !testIntegerLiteral(t, stmt.Expression, 5) {
		return
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue bool
	}{
		{"false;", false},
		{"true;", true},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Expected program to have 1 Statement, instead got %v", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected statement to be an ExpressionStatement, instead got %T", program.Statements[0])
		}

		if !testBooleanLiteral(t, stmt.Expression, tt.expectedValue) {
			return
		}
	}
}

func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input             string
		operator          string
		expectedIdentifer interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Expected program to have 1 statement, instead got %v", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("Expected statement to be an ExpressionStatement, instead got %T", program.Statements[0])
		}

		expr, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("Expected a PrefixExpression, instead got %T", stmt.Expression)
		}

		if !testLiteralExpression(t, expr.Right, tt.expectedIdentifer) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, v int64) bool {
	intl, ok := il.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("Expected expression to be an IntegerLiteral, instead got %T", il)
		return false
	}

	if intl.Value != v {
		t.Errorf("Expected value to be %v, instead got %v", v, intl.Value)
		return false
	}

	if intl.TokenLiteral() != fmt.Sprintf("%d", intl.Value) {
		t.Errorf("Expected TokenLiteral to be %v, instead got %v", intl.Value, intl.TokenLiteral())
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("Expected program to have 1 statement, instead got %v", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected statement to be an ExpressionStatement, instead got %T", program.Statements[0])
		}
		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.BooleanExpression)
	if !ok {
		t.Errorf("exp not *ast.BooleanExpression. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("Expexted expression to be an ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"true", "true"},
		{"false", "false"},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		}, {
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected program to have 1 statement, instead got %v", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Expected statement to be an ExpressionStatement, instead got %T", program.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.IfExpression)

	if !ok {
		t.Fatalf("Expected expression to be an IfExpression, instead got %T", stmt.Expression)
	}

	if !testInfixExpression(t, expr.Condition, "x", "<", "y") {
		return
	}

	if len(expr.Consequence.Statements) != 1 {
		t.Fatalf("Expected consequence to have 1 statement, instead got %v", len(expr.Consequence.Statements))
	}

	stmt, ok = expr.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Expected statement to be an ExpressionStatement, instead got %T", expr.Consequence.Statements[0])
	}

	if !testIdentifier(t, stmt.Expression, "x") {
		return
	}

	if expr.Alternative != nil {
		t.Fatalf("Expected Alternative to be nil, instead got %+v", expr.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteral(t *testing.T) {
	input := `fn(a, b) { a + b }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected program to have 1 statement, instead got %v", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Expected statement to be an ExpressionStatement, instead got %T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)

	if !ok {
		t.Fatalf("Expected expression to be a FunctionLiteral, instead got %T", stmt.Expression)
	}

	if !testLiteralExpression(t, function.Parameters[0], "a") || !testLiteralExpression(t, function.Parameters[1], "b") {
		return
	}

	if len(function.Body.Statements) != 1 {
		t.Fatalf("Expected function body to have 1 statement, instead got %v", len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected statement to be ExpressionStatement, instead got %T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "a", "+", "b")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"fn() {}", []string{}},
		{"fn(a) {}", []string{"a"}},
		{"fn(a, b, c) {}", []string{"a", "b", "c"}},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(test.expectedParams) {
			t.Fatalf("Expected %v parameters, but received %v", len(function.Parameters), len(test.expectedParams))
		}

		for i, param := range function.Parameters {
			testLiteralExpression(t, param, test.expectedParams[i])
		}
	}
}

func TestCallExpressions(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected program to have 1 statement, instead got %v", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected statement to be an ExpressionStatement, instead got %T", program.Statements[0])
	}

	call, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("Expected expression to be a CallExpression, instead got %T", stmt.Expression)
	}

	if !testLiteralExpression(t, call.Function, "add") {
		return
	}

	if len(call.Arguments) != 3 {
		t.Fatalf("Expected to have 3 arguments, instead got %v", len(call.Arguments))
	}

	testLiteralExpression(t, call.Arguments[0], 1)
	testInfixExpression(t, call.Arguments[1], 2, "*", 3)
	testInfixExpression(t, call.Arguments[2], 4, "+", 5)
}
