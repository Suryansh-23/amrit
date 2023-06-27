package evaluator

import (
	"fmt"

	"github.com/Suryansh-23/amrit/ast"
	"github.com/Suryansh-23/amrit/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment, stdout *[]string) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, env, stdout)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env, stdout)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env, stdout)
	case *ast.IfExpression:
		return evalIfExpression(node, env, stdout)
	case *ast.WhileExpression:
		return evalWhileExpression(node, env, stdout)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env, stdout)
		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body

		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.CallExpression:
		function := Eval(node.Function, env, stdout)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env, stdout)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args, stdout)
	case *ast.LetStatement:
		val := Eval(node.Value, env, stdout)
		if isError(val) {
			return val
		}

		env.Set(node.Name.Value, val)
	case *ast.CompoundAssignmentStatement:
		val := Eval(node.Value, env, stdout)
		if isError(val) {
			return val
		}

		initVal, ok := env.Get(node.Name.Value)
		if !ok {
			return newError("identifier not found: " + node.Name.Value)
		}

		env.Set(node.Name.Value, computeOp(node.Operator, initVal, val))
	case *ast.Identifier:
		return evalIdentifier(node, env)

		// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env, stdout)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}

		return &object.Array{Elements: elements}
	case *ast.SliceExpression:
		left := Eval(node.Left, env, stdout)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env, stdout)
		if isError(right) {
			return right
		}
		return &object.Slice{Left: left, Right: right}

	case *ast.IndexExpression:
		left := Eval(node.Left, env, stdout)
		if isError(left) {
			return left
		}

		index := Eval(node.Index, env, stdout)
		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index)
	case *ast.SliceArrayExpression:
		left := Eval(node.Left, env, stdout)
		if isError(left) {
			return left
		}

		slice := Eval(&node.Slice, env, stdout)
		if isError(slice) {
			return slice
		}

		return evalSliceExpression(left, slice)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env, stdout)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env, stdout)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env, stdout)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env, stdout)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.Comment:
		return nil
	}

	return nil
}

func evalStatements(stmts []ast.Statement, env *object.Environment, stdout *[]string) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env, stdout)

		if returnVale, ok := result.(*object.ReturnValue); ok {
			return returnVale.Value
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s %s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
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

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		left := left.(*object.Boolean).Value
		right := right.(*object.Boolean).Value
		// fmt.Println(left, right, left == right, object.Boolean{Value: false} == object.Boolean{Value: false})
		return nativeBoolToBooleanObject(left == right)

	case operator == "!=":
		left := left.(*object.Boolean).Value
		right := right.(*object.Boolean).Value
		return nativeBoolToBooleanObject(left != right)

	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())

	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "<":
		return &object.Boolean{Value: leftVal < rightVal}
	case "<=":
		return &object.Boolean{Value: leftVal <= rightVal}
	case ">":
		return &object.Boolean{Value: leftVal > rightVal}
	case ">=":
		return &object.Boolean{Value: leftVal >= rightVal}
	case "==":
		return &object.Boolean{Value: leftVal == rightVal}
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return &object.Boolean{Value: leftVal == rightVal}
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment, stdout *[]string) object.Object {
	condition := Eval(ie.Condition, env, stdout)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env, stdout)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env, stdout)
	} else {
		// fmt.Println("returning NULL")
		return NULL
	}
}

func evalWhileExpression(we *ast.WhileExpression, env *object.Environment, stdout *[]string) object.Object {
	condition := Eval(we.Condition, env, stdout)
	if isError(condition) {
		return condition
	}

	for isTruthy(condition) {
		Eval(we.Body, env, stdout)
		condition = Eval(we.Condition, env, stdout)

		if isError(condition) {
			return condition
		}
	}

	return NULL
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Null:
		return false
	case *object.Boolean:
		return obj.Value
	default:
		return true
	}
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment, stdout *[]string) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env, stdout)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalProgram(program *ast.Program, env *object.Environment, stdout *[]string) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, env, stdout)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment, stdout *[]string) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env, stdout)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object, stdout *[]string) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendEnv, stdout)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(stdout, args...)

	default:
		return newError("not a function: %s", fn.Type())
	}

}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue
	}

	return obj
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported %s", left.Type())
	}
}

func evalArrayIndexExpression(arr, index object.Object) object.Object {
	arrObj := arr.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrObj.Elements) - 1)

	if idx < 0 || max < idx {
		return NULL
	}

	return arrObj.Elements[idx]
}

func evalSliceExpression(left, slice object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && slice.Type() == object.SLICE_OBJ:
		return evalArraySliceExpression(left, slice)
	default:
		return newError("slice operator not supported %s", left.Type())
	}
}

func evalArraySliceExpression(arr, slice object.Object) object.Object {
	arrObj := arr.(*object.Array)
	left := slice.(*object.Slice).Left.(*object.Integer).Value
	right := slice.(*object.Slice).Right.(*object.Integer).Value
	max := int64(len(arrObj.Elements))

	if left < 0 || max < left || right < 0 || max < right || left > right {
		return NULL
	}

	return &object.Array{Elements: arrObj.Elements[left:right]}
}

func computeOp(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return computeIntegerOp(operator, left, right)
	default:
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func computeIntegerOp(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+=":
		return &object.Integer{Value: leftVal + rightVal}
	case "-=":
		return &object.Integer{Value: leftVal - rightVal}
	case "*=":
		return &object.Integer{Value: leftVal * rightVal}
	case "/=":
		return &object.Integer{Value: leftVal / rightVal}
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment, stdout *[]string) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env, stdout)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env, stdout)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}
	return &object.Hash{Pairs: pairs}
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObj := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}
