package interpreter

import (
	"fmt"
	"strconv"

	"github.com/wreulicke/go-sandbox/go-interpreter/monkey/ast"
	"github.com/wreulicke/go-sandbox/go-interpreter/monkey/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements, env)
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
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		return evalCallExpression(fn, node.Arguments, env)
	case *ast.ArrayLiteral:
		elements, err := evalExpressions(node.Elements, env)
		if err != nil {
			return err
		}
		return &object.Array{
			Elements: elements,
		}
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.NumberLiteral:
		i, err := strconv.ParseInt(node.Value, 10, 64)
		if err != nil {
			return newError("cannot convert int. %s", node.Value)
		}
		return &object.Integer{Value: i}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.FunctionLiteral:
		f := &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
		return f
	}
	return newError("Unsupported ast.Node got=%T", node)
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) ([]object.Object, object.Object) {
	var result []object.Object
	for _, e := range expressions {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return nil, evaluated
		}
		result = append(result, evaluated)
	}
	return result, nil
}

func evalCallExpression(fn object.Object, arguments []ast.Expression, env *object.Environment) object.Object {
	args, err := evalExpressions(arguments, env)
	if err != nil {
		return err
	}
	switch fn := fn.(type) {
	case *object.Function:
		functionEnv := extendFunctionEnv(fn, args)
		return unwrapReturnValue(Eval(fn.Body, functionEnv))
	case *object.Builtin:
		return fn.Fn(args...)
	}
	return newError("not a function: %s", fn.Type())

}

func extendFunctionEnv(function *object.Function, args []object.Object) *object.Environment {
	env := function.Env.NewEnclosedEnvironment()

	for paramIdx, param := range function.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(o object.Object) object.Object {
	switch o := o.(type) {
	case *object.ReturnValue:
		return o.Value
	default:
		return o
	}
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, s := range stmts {
		result = Eval(s, env)
		switch v := result.(type) {
		case *object.ReturnValue:
			return v.Value
		case *object.Error:
			return v
		}
	}

	return result

}

func evalBlockStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, s := range stmts {
		result = Eval(s, env)
		if result != nil {
			switch result.Type() {
			case object.RETURN:
				return result
			case object.ERROR:
				return result
			}
		}
	}

	return result
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	cond := Eval(ie.Condition, env)
	if isError(cond) {
		return cond
	}
	if isTruthy(cond) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}
	return NULL
}

func isTruthy(o object.Object) bool {
	switch o {
	case TRUE:
		return true
	case FALSE:
		return false
	case NULL:
		return false
	default:
		return true
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING && right.Type() == object.STRING:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value
	switch operator {
	case "+":
		return &object.String{Value: leftValue + rightValue}
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return FALSE
	default:
		return FALSE
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if v, ok := env.Get(node.Value); ok {
		return v
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("identifier is not found: %s", node.Value)
}

func nativeBoolToBooleanObject(b bool) object.Object {
	if b {
		return TRUE
	}
	return FALSE
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR
}
