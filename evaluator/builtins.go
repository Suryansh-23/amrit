package evaluator

import (
	"github.com/Suryansh-23/amrit/object"
)

var builtins = map[string]*object.Builtin{
	"lambai": {
		Fn: func(stdout *[]string, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `lambai` not supported, got %s",
					args[0].Type())
			}
		},
	},
	"print": {
		Fn: func(stdout *[]string, args ...object.Object) object.Object {
			for _, arg := range args {
				if arg.Type() != object.STRING_OBJ && arg.Type() != object.INTEGER_OBJ && arg.Type() != object.BOOLEAN_OBJ && arg.Type() != object.NULL_OBJ && arg.Type() != object.ARRAY_OBJ {
					return newError("argument `%s` of type %s not supported in `print`", arg, arg.Type())
				}
			}
			s := ""

			for _, arg := range args {
				s += arg.Inspect() + " "
			}
			s += "\n"
			*stdout = append(*stdout, s)

			return &object.Null{}
		},
	},
	"pehla": {
		Fn: func(stdout *[]string, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `pehla` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},
	"aakhri": {
		Fn: func(stdout *[]string, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `aakhri` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[len(arr.Elements)-1]
			}

			return NULL
		},
	},
	"baaki": {
		Fn: func(stdout *[]string, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `aakhri` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])

				return &object.Array{Elements: newElements}
			}

			return NULL
		},
	},
	"push": {
		Fn: func(stdout *[]string, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `push` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			newElements := make([]object.Object, length+1)

			copy(newElements, arr.Elements)
			newElements[length] = args[1]

			return &object.Array{Elements: newElements}
		},
	},
	"pop": {
		Fn: func(stdout *[]string, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `pop` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			newElements := make([]object.Object, length-1)
			copy(newElements, arr.Elements[:length-1])

			return &object.Array{Elements: newElements}
		},
	},
}
