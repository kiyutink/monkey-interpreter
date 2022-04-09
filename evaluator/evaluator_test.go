package evaluator

import (
	"monkey-interpreter/lexer"
	"monkey-interpreter/object"
	"monkey-interpreter/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-10", -10},
		{"-20", -20},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestStringConcatenation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"foo" + "bar";`, "foobar"},
		{`"foo" + "";`, "foo"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		str, ok := evaluated.(*object.String)
		if !ok {
			t.Errorf("Expected a String object, instead got %T (%+v)", evaluated, evaluated)
			continue
		}

		if str.Value != tt.expected {
			t.Errorf("Expected String value to be %v, instead got %v", tt.expected, str.Value)
		}
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"5 == 5", true},
		{"5 == 10", false},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"true != true", false},
		{"1 == 1", true},
		{"2 == 2", true},
		{"1 == 2", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 != 1", false},
		{"1 != 2", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"!5", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 5;", 5},
		{"return 9; 10;", 9},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`if (10 > 1) {
			if (10 > 1) {
				if (true) {
					return 9
				}
				return 10;
			}
			return 1; }`, 9},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		}, {
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
			return 1; }
			`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{`"foobar" - "bar";`, "unknown operator: STRING - STRING"},
		{`"foobar" * "bar";`, "unknown operator: STRING * STRING"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("Expected object to be Error, instead got %T (%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("Expected error message to be %v, instead got %v", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; }"

	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("Expected object to be a Function, instead received %T", evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("Expected function to have 1 parameter, instead got %v", len(fn.Parameters))
	}

	if fn.Parameters[0].Value != "x" {
		t.Fatalf("Expected function parameter to be x, instead got %v", fn.Parameters[0].Value)
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("Expected body of the function to be %v, instead got %v", expectedBody, fn.Body.String())
	}
}

func TestStringObject(t *testing.T) {
	input := `"foobar";`
	expected := "foobar"
	evaluated := testEval(input)

	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("Expected a String object, instead got %T", evaluated)
	}

	if str.Value != expected {
		t.Fatalf("Expected String object to have value %v, instead got %v", expected, str.Value)
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
		let newAdder = fn(x) {
			fn(y) { x + y };
		};
		let addTwo = newAdder(2);
		addTwo(2);`
	testIntegerObject(t, testEval(input), 4)
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("abc")`, 3},
		{`len("hello world!")`, 12},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1)"},
		{`len([1, 2, 3]);`, 3},
		{`len([]);`, 0},
		{`head([]);`, nil},
		{`head([1, 2, 3]);`, 1},
		{`tail([1, 2, 3]);`, []int64{2, 3}},
		{`tail([]);`, []int64{}},
		{`tail([1]);`, []int64{1}},
		{`last([]);`, nil},
		{`last([1]);`, 1},
		{`last([1, 2, 3]);`, 3},
		{`push([1, 2], 3);`, []int64{1, 2, 3}},
		{`push(5);`, "wrong number of arguments. got=1, want=2)"},
		{`push(5, 5);`, "argument to `push` not supported, got INTEGER"},
		{`push([], 5);`, []int64{5}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case nil:
			testNullObject(t, evaluated)
		case int:
			testIntegerObject(t, evaluated, int64(expected))

		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("Expected an Error object, instead got %T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("Expected Message to be %v, instead got %v", expected, errObj.Message)
			}
		case []int64:
			arr, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("Expected an array, instead got %T(%+v))", evaluated, evaluated)
			}

			for i, val := range arr.Elements {
				testIntegerObject(t, val, expected[i])
			}
		}

	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 + 3, 4 * 5];"
	evaluated := testEval(input)

	arr, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("Expected object to be Array, instead got %T", evaluated)
	}

	if len(arr.Elements) != 3 {
		t.Fatalf("Expected array to have 3 elements, instead got %v", len(arr.Elements))
	}

	testIntegerObject(t, arr.Elements[0], 1)
	testIntegerObject(t, arr.Elements[1], 5)
	testIntegerObject(t, arr.Elements[2], 20)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"let i = 0; [1][i];", 1},
		{"[1, 2, 3][1 + 1];", 3},
		{"let myArray = [1, 2, 3]; myArray[2];", 3},
		{"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];", 6},
		{"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]", 2},
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)

	if !ok {
		t.Fatalf("Expected a Hash, instead got %T(%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():     1,
		(&object.String{Value: "two"}).HashKey():     2,
		(&object.String{Value: "three"}).HashKey():   3,
		(&object.Integer{Value: int64(4)}).HashKey(): 4,
		(&object.Boolean{Value: true}).HashKey():     5,
		(&object.Boolean{Value: false}).HashKey():    6,
	}

	if len(expected) != len(result.Pairs) {
		t.Fatalf("The length are not equal, expected %v, received %v", len(expected), len(result.Pairs))
	}

	for key, value := range expected {
		evaluatedPair, ok := result.Pairs[key]
		if !ok {
			t.Errorf("Value for key is absent")
		}
		testIntegerObject(t, evaluatedPair.Value, value)
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("Expected object to be NULL, instead got %T (%+v)", obj, obj)
		return false
	}
	return true
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("Expected object to be Integer, instead got %T (%v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("Expected integer value to be %v, instead got %v", expected, result.Value)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("Expected object to be Boolean, instead got %T", obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("Expected object value to be %v, instead got %v", expected, result.Value)
		return false
	}

	return true
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	env := object.NewEnvironment()

	program := p.ParseProgram()

	return Eval(program, env)
}
