package evaluator

import (
	"fmt"
	"monkey-interpreter/object"
)

func length(args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%v, want=1)", len(args))}
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{
			Value: int64(len(arg.Value)),
		}

	case *object.Array:
		return &object.Integer{
			Value: int64(len(arg.Elements)),
		}
	default:
		return &object.Error{Message: fmt.Sprintf("argument to `len` not supported, got %v", args[0].Type())}
	}
}

func head(args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%v, want=1)", len(args))}
	}

	switch arg := args[0].(type) {
	case *object.Array:
		if len(arg.Elements) == 0 {
			return NULL
		}
		return arg.Elements[0]
	default:
		return &object.Error{Message: fmt.Sprintf("argument to `len` not supported, got %v", args[0].Type())}
	}
}

func tail(args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%v, want=1)", len(args))}
	}

	switch arg := args[0].(type) {
	case *object.Array:
		if len(arg.Elements) < 2 {
			return &object.Array{
				Elements: []object.Object{},
			}
		}
		return &object.Array{
			Elements: arg.Elements[1:],
		}
	default:
		return &object.Error{Message: fmt.Sprintf("argument to `len` not supported, got %v", args[0].Type())}
	}
}

func last(args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%v, want=1)", len(args))}
	}

	switch arg := args[0].(type) {
	case *object.Array:
		if len(arg.Elements) == 0 {
			return NULL
		}
		return arg.Elements[len(arg.Elements)-1]
	default:
		return &object.Error{Message: fmt.Sprintf("argument to `len` not supported, got %v", args[0].Type())}
	}
}

func push(args ...object.Object) object.Object {
	if len(args) != 2 {
		return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%v, want=2)", len(args))}
	}

	arr, ok := args[0].(*object.Array)
	if !ok {
		return &object.Error{Message: fmt.Sprintf("argument to `push` not supported, got %v", args[0].Type())}
	}

	return &object.Array{
		Elements: append(arr.Elements, args[1]),
	}
}

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: length,
	},
	"head": {
		Fn: head,
	},
	"tail": {
		Fn: tail,
	},
	"last": {
		Fn: last,
	},
	"push": {
		Fn: push,
	},
}
