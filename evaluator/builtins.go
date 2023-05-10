package evaluator

import (
	"fmt"

	"github.com/Suryansh-23/amrit/object"
)

var builtins = map[string]*object.Builtin{
	"lambai": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `lambai` not supported, got %s",
					args[0].Type())
			}
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				if arg.Type() != object.STRING_OBJ && arg.Type() != object.INTEGER_OBJ && arg.Type() != object.BOOLEAN_OBJ && arg.Type() != object.NULL_OBJ {
					return newError("argument `%s` of type %s not supported in `print`", arg, arg.Type())
				}
			}

			for _, arg := range args {
				fmt.Printf("%s ", arg.Inspect())
			}
			fmt.Println()

			return &object.Null{}
		},
	},
}
