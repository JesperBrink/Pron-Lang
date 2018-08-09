package evaluator

import (
	"Pron-Lang/ast"
	"Pron-Lang/object"
	"fmt"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}

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

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

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

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.VarStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.DirectFunctionStatement:
		var newNode ast.Node
		newNode = &node.Function
		val := Eval(newNode, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.ClassStatement:
		return evalClassStatement(node, env)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		if node.Operator == "=" {
			return evalAssignValueToExistingVariable(node.Left, node.Right, env)
		}

		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.ElseIfExpression:
		return evalElseIfExpression(node, env)

	case *ast.IncrementForloopExpression:
		return evalIncrementForloopExpression(node, env)

	case *ast.ArrayForloopExpression:
		return evalArrayForloopExpression(node, env)

	case *ast.ObjectInitialization:
		return evalObjectInitialization(node, env)

	case *ast.CallObjectFunction:
		return evalCallObejctFunction(node, env)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalCallObejctFunction(node *ast.CallObjectFunction, env *object.Environment) object.Object {
	objObject, ok := env.Get(node.ObjectName.Value)
	if !ok {
		return newError("%s is not defined", node.ObjectName.Value)
	}

	obj, ok := objObject.(*object.ClassInstance)
	if !ok {
		return newError("%s is not an object. It's a %T", obj.Name, objObject)
	}

	functionObject, ok := obj.Env.Get(node.FunctionName.Value)
	if !ok {
		return newError("%s is not a defined method", node.FunctionName.Value)
	}
	function := functionObject.(*object.Function)

	//evalExpressions returns []object.Object
	arguments := evalExpressions(node.Arguments, obj.Env)
	return applyFunction(function, arguments)
}

func evalObjectInitialization(node *ast.ObjectInitialization, env *object.Environment) object.Object {
	classInstanceObject, ok := env.Get(node.Name.Value)
	if !ok {
		return newError("There is no Class called: " + node.Name.Value)
	}
	classInstance := classInstanceObject.(*object.ClassInstance)

	initFunctionObject, ok := classInstance.Env.Get("init")
	if !ok {
		return classInstance
	}
	initFunction := initFunctionObject.(*object.InitFunction)

	args := node.Arguments

	// Create env with all arguments that isn't a 'this.' argument
	newEnv := object.NewEnclosedEnvironment(initFunction.Env)

	for paramIdx, param := range initFunction.Parameters {
		if param.IsThisParam {
			val := Eval(args[paramIdx], classInstance.Env)
			classInstance.Env.Update(param.Parameter.Value, val)
		} else {
			val := Eval(args[paramIdx], newEnv)
			newEnv.Set(param.Parameter.Value, val)
		}

	}

	Eval(initFunction.Body, newEnv)
	return classInstance
}

func evalClassStatement(node *ast.ClassStatement, env *object.Environment) object.Object {
	// Create local env
	classEnv := object.NewEnvironment()

	// Eval fields
	for _, field := range node.Fields {
		val := Eval(field.Value, classEnv)
		if isError(val) {
			return val
		}
		classEnv.Set(field.Name.Value, val)
	}

	// Eval functions
	for _, function := range node.Functions {
		var newNode ast.Node
		newNode = &function.Function
		val := Eval(newNode, classEnv)
		if isError(val) {
			return val
		}
		classEnv.Set(function.Name.Value, val)
	}

	var initFunction object.Object

	// Eval init
	if node.InitBody != nil {
		initFunction = &object.InitFunction{Parameters: node.InitParams, Body: node.InitBody, Env: classEnv}
		classEnv.Set("init", initFunction)
	}

	// Put class into global env
	result := &object.ClassInstance{Name: node.Name.Value, Env: classEnv}
	env.Set(result.Name, result)
	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
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
		return newError("unknown operator: %s%s", operator, right.Type())
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

func evalInfixExpression(operator string, left, right object.Object, env *object.Environment) object.Object {
	switch {
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
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
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalElseIfExpression(ei *ast.ElseIfExpression, env *object.Environment) object.Object {
	for _, conditionAndBlockstatement := range ei.ConditionAndBlockstatementList {
		condition := Eval(conditionAndBlockstatement.Condition, env)
		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			return Eval(conditionAndBlockstatement.Consequence, env)
		}
	}

	if ei.Alternative != nil {
		return Eval(ei.Alternative, env)
	} else {
		return NULL
	}

}

func evalIncrementForloopExpression(incForloopExp *ast.IncrementForloopExpression, env *object.Environment) object.Object {
	from := Eval(incForloopExp.From, env)
	if from.Type() != object.INTEGER_OBJ {
		return newError("'from' expression in forloop was not integer. got=%T", from)
	}

	to := Eval(incForloopExp.To, env)
	if to.Type() != object.INTEGER_OBJ {
		return newError("'to' expression in forloop was not integer. got=%T", to)
	}

	// create new extended env with local var
	newEnv := object.NewEnclosedEnvironment(env)
	newEnv.Set(incForloopExp.LocalVar.String(), NULL)

	var result object.Object = NULL
	fromValue := from.(*object.Integer).Value
	toValue := to.(*object.Integer).Value

	if fromValue < toValue {
		for i := fromValue; i < toValue; i++ {
			newEnv.Update(incForloopExp.LocalVar.String(), &object.Integer{Value: i})
			result = evalBlockStatement(incForloopExp.Body, newEnv)

			switch result := result.(type) {
			case *object.ReturnValue:
				return result.Value
			case *object.Error:
				return result
			}
		}
	} else {
		for i := fromValue; i > toValue; i-- {
			newEnv.Update(incForloopExp.LocalVar.String(), &object.Integer{Value: i})
			result = evalBlockStatement(incForloopExp.Body, newEnv)

			switch result := result.(type) {
			case *object.ReturnValue:
				return result.Value
			case *object.Error:
				return result
			}
		}
	}
	return result
}

func evalArrayForloopExpression(arrayForloopExp *ast.ArrayForloopExpression, env *object.Environment) object.Object {
	array, ok := env.Get(arrayForloopExp.ArrayName.String())
	if !ok {
		return newError("%s is not defined", arrayForloopExp.ArrayName.String())
	}
	arrayObject := array.(*object.Array)

	// create new extended env with local var
	newEnv := object.NewEnclosedEnvironment(env)
	newEnv.Set(arrayForloopExp.LocalVar.String(), NULL)

	var result object.Object = NULL

	for _, elem := range arrayObject.Elements {
		newEnv.Update(arrayForloopExp.LocalVar.String(), elem)
		result = evalBlockStatement(arrayForloopExp.Body, newEnv)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
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

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendedFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function %s", fn.Type())
	}
}

func extendedFunctionEnv(function *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(function.Env)

	for paramIdx, param := range function.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env

}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalAssignValueToExistingVariable(left, right ast.Expression, env *object.Environment) object.Object {
	leftIdentifier, ok := left.(*ast.Identifier)
	if !ok {
		return newError("leftside of assignment is not a string. got=%T (%+v)", left, left)
	}

	val := Eval(right, env)
	// check existens of variables
	if !env.Update(leftIdentifier.Value, val) {
		return newError("%s is not defined", leftIdentifier.Value)
	}
	//env.Set(leftIdentifier.Value, val)
	return val
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	switch operator {
	case "+":
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		return &object.String{Value: leftVal + rightVal}
	case "==":
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		return &object.Boolean{Value: leftVal == rightVal}
	case "!=":
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		return &object.Boolean{Value: leftVal != rightVal}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}
