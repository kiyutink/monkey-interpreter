package evaluator

import (
	"fmt"
	"monkey-interpreter/object"
)

func length(args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Message: "wrong number of arguments. got=2, want=1"}
	}

	if args[0].Type() != object.STRING_OBJ {
		return &object.Error{Message: fmt.Sprintf("argument to `len` not supported, got %v", args[0].Type())}
	}

	str := args[0].(*object.String)
	return &object.Integer{
		Value: int64(len(str.Value)),
	}
}

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: length,
	},
}
