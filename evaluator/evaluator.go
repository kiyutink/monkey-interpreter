package evaluator

import (
	"fmt"

	"monkey-interpreter/ast"
	"monkey-interpreter/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// Statements

	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return env.Set(node.Name.Value, val)

	// Expressions

	case *ast.HashLiteral:
		pairs := make(map[object.HashKey]object.HashPair)
		for key, value := range node.Pairs {
			keyObj := Eval(key, env)
			if isError(keyObj) {
				return keyObj
			}
			hashableKey, ok := keyObj.(object.Hashable)
			if !ok {
				return newError("Can't use expression of type %v as hash key", keyObj.Type())
			}

			valObj := Eval(value, env)
			if isError(valObj) {
				return valObj
			}

			pairs[hashableKey.HashKey()] = object.HashPair{
				Value: valObj,
				Key:   keyObj,
			}
		}

		return &object.Hash{
			Pairs: pairs,
		}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.FunctionLiteral:
		return &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.BooleanExpression:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)
		if isError(left) {
			return left
		}
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	default:
		return nil
	}
}

func evalIndexExpression(left object.Object, index object.Object) object.Object {
	switch left := left.(type) {
	case *object.Array:
		idx := index.(*object.Integer).Value
		max := int64(len(left.Elements) - 1)
		if idx < 0 || idx > max {
			return NULL
		}

		return left.Elements[idx]

	case *object.Hash:
		key, ok := index.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %v", index.Type())
		}
		pair, ok := left.Pairs[key.HashKey()]
		if !ok {
			return NULL
		}
		return pair.Value

	default:
		return NULL
	}
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch function := fn.(type) {
	case *object.Function:
		callEnv := object.NewEnclosedEnvironment(function.Env)
		for i, arg := range args {
			callEnv.Set(function.Parameters[i].Value, arg)
		}
		return unwrapReturnValue(evalBlockStatement(function.Body, callEnv))
	case *object.Builtin:
		return function.Fn(args...)
	default:
		return newError("not a function: %v", fn.Type())
	}
}

func unwrapReturnValue(obj object.Object) object.Object {
	if retVal, ok := obj.(*object.ReturnValue); ok {
		return retVal.Value
	}
	return obj
}

func evalExpressions(nodes []ast.Expression, env *object.Environment) []object.Object {
	objects := []object.Object{}

	for _, node := range nodes {
		obj := Eval(node, env)
		if isError(obj) {
			return []object.Object{obj}
		}
		objects = append(objects, obj)
	}

	return objects
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	if !ok {
		return newError("identifier not found: " + node.Value)
	}

	return val
}

func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	}
	return NULL
}

func evalInfixExpression(op string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return newError("type mismatch: %v %v %v", left.Type(), op, right.Type())
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalInfixIntegerExpression(op, left, right)

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalInfixStringExpression(op, left, right)

	// After here at least one of the operands is a bool
	case op == "==":
		return nativeBoolToBooleanObject(left == right)
	case op == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return newError("unknown operator: %v %v %v", left.Type(), op, right.Type())
	}
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangPrefixOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %v %v", op, right.Type())
	}
}

func evalInfixStringExpression(op string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	if op == "+" {
		return &object.String{Value: leftVal + rightVal}
	}

	return newError("unknown operator: %v %v %v",
		left.Type(), op, right.Type())
}

func evalInfixIntegerExpression(op string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch op {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	default:
		return newError("unknown operator: %v %v %v",
			left.Type(), op, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%v", right.Type())
	}
	integer := right.(*object.Integer)
	return &object.Integer{Value: -integer.Value}
}

func evalBangPrefixOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}

		if result != nil && result.Type() == object.ERROR_OBJ {
			return result
		}
	}

	return result
}

func nativeBoolToBooleanObject(val bool) *object.Boolean {
	if val {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj object.Object) bool {
	if obj == FALSE || obj == NULL {
		return false
	}

	return true
}

func newError(format string, a ...interface{}) *object.Error {
	err := &object.Error{
		Message: fmt.Sprintf(format, a...),
	}
	return err
}

func isError(obj object.Object) bool {
	if obj != nil && obj.Type() == object.ERROR_OBJ {
		return true
	}
	return false
}
