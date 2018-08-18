package evaluator

import (
	"Pron-Lang/lexer"
	"Pron-Lang/object"
	"Pron-Lang/parser"
	"strconv"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
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

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
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

func TestElifExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 } elif (false) { 9 }", 10},
		{"if (false) { 10 } elif (true) { 9 }", 9},
		{"if (false) { 10 } elif (true) { 9 } elif (false) { 8 }", 9},
		{"if (false) { 10 } elif (false) { 9 } elif (true) { 8 }", 8},
		{"if (true) { 10 } elif (true) { 9 } elif (true) { 8 }", 10},
		{"if (false) { 10 } elif (false) { 9 }", nil},
		{"if (false) { 10 } elif (false) { 9 } else { 7 }", 7},
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

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return 10;
				}
			return 1;
			}
			`, 10,
		},
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
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
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
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
				return 1;
			}
		`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			`{"name": "Monkey"}[func(x) { x }];`,
			"unusable as hash key: FUNCTION",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestVarStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var a = 5; a;", 5},
		{"var a = 5 * 5; a;", 25},
		{"var a = 5; var b = a; b;", 5},
		{"var a = 5; var b = a; b; var c = a + b + 5; c", 15},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var identity = func(x) { x; }; identity(5);", 5},
		{"var identity = func(x) { return x; }; identity(5);", 5},
		{"var double = func(x) { x * 2; }; double(5);", 10},
		{"var add = func(x, y) { x + y; }; add(5, 5);", 10},
		{"var add = func(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestDirectFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"func identity(x) { x; }; (5)", 5},
		{"func identity(x) { return x; }; identity(5);", 5},
		{"func double(x) { x * 2; }; double(5);", 10},
		{"func add(x, y) { x + y; }; add(5, 5);", 10},
		{"func add(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not string. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("evaluated is not *object.String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("str.Value is not %s. got=%s", "Hello World!", str.Value)
	}
}

func TestStringComparison(t *testing.T) {
	input := []struct {
		input    string
		expected bool
	}{
		{`"hello" == "hello"`, true},
		{`"hello" == "hey"`, false},
		{`"hello" != "hello"`, false},
		{`"hello" != "hey"`, true},
	}

	for _, tt := range input {
		evaluated := testEval(tt.input)
		b, ok := evaluated.(*object.Boolean)
		if !ok {
			t.Fatalf("evaluated is not *object.Boolean. got=%T (%+v)", evaluated, evaluated)
		}

		if b.Value != tt.expected {
			t.Errorf("b.Value is not %t. got=%t", tt.expected, b.Value)
		}
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		//Test len of string
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
		//Test len of array
		{`var a = [1, 2]; len(a)`, 2},
		{`var a = []; len(a);`, 0},
		//Test `first`
		{`var a = [4, 2, 5]; first(a)`, 4},
		{`var a = []; first(a)`, nil},
		//Test `last`
		{`var a = [4, 2, 5]; last(a)`, 5},
		{`var a = []; last(a)`, nil},
		//Test `rest`
		{`var a = [4, 2, 5]; rest(a)`, []int64{2, 5}},
		{`var a = [1]; rest(a)`, []int64{}},
		{`var a = []; rest(a)`, nil},
		//Test `push`
		{`var a = [4, 2, 5]; push(a, 7)`, []int64{4, 2, 5, 7}},
		{`var a = []; push(a, 1)`, []int64{1}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		case []int64:
			arr, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("evaluated is not object.Array. got=%+v", evaluated)
			}

			for i, expectedElem := range arr.Elements {
				testIntegerObject(t, expectedElem, expected[i])
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"var i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"var myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"var myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"var myArray = [1, 2, 3]; var i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
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
	input := `var two = "two";
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
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]

		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`var key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
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

func TestAssignValueToExistingVariable(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			"var a = 5; a = 6; a;",
			6,
		},
		{
			"var a = 5; a = a + 1; a",
			6,
		},
		{
			"var a = 5; var b = a; b = a + a; b",
			10,
		},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestIncrementForloopExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"var x = 0; for (i from 0 to 10) { x = x + 1 }; return x;",
			10,
		},
		{
			"var x = 0; var y = 0; var z = 5; for (i from y to z) { x = x + 1 }; return x;",
			5,
		},
		{
			"var x = 0; for (i from 0 to 0) { x = x + 1 }; return x;",
			0,
		},
		{
			"var x = 0; for (i from 0 to 1) { x = x + 1 }; return x;",
			1,
		},
		{
			"var x = 0; for (i from 0 to -1) { x = x + 1 }; return x;",
			1,
		},
		{
			"var x = 0; for (i from -5 to -1) { x = x + 1 }; return x;",
			4,
		},
		{
			"var x = 0; for (i from 0 to 10) { x = i }; return x;",
			9,
		},
		{
			"for (i from 4 to 9) { return i };",
			4,
		},
		{
			"var x = 0; for (i from 3 to 0) { x = x + i }; return x",
			6,
		},
		{
			"var x = 0; for (i from 4 to -4) { x = x + i }; return x",
			4,
		},
		{
			"var x = 0; for (i from true to 10) { x = i; }; return x;",
			"'from' expression in forloop was not integer. got=*object.Boolean",
		},
		{
			"for (i from 0 to 0) { return i };",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else if _, ok := tt.expected.(string); ok {
			errObj, ok := evaluated.(*object.Error)
			expected := tt.expected.(string)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestArrayForloopExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"var x = [6]; for (i in x) { return i };",
			6,
		},
		{
			"var x = [1,2,5,3]; var sum = 0; for (i in x) { sum = sum + i }; return sum",
			11,
		},
		{
			"var x = []; for (i in x) { return i };",
			nil,
		},
		{
			"var x = [1,4,2]; for (i in x) { return i };",
			1,
		},
		{
			"var x = 0; for (i from true to 10) { x = i; }; return x;",
			"'from' expression in forloop was not integer. got=*object.Boolean",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else if _, ok := tt.expected.(string); ok {
			errObj, ok := evaluated.(*object.Error)
			expected := tt.expected.(string)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestClassObject(t *testing.T) {
	input := `
	class Person {
		var name = ""
		var age = 0

		init(name, this.age) {
			name = "Hans"
		}

		func getName(dummyParam) {
			return name
		}
	}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.ClassInstance)
	if !ok {
		t.Fatalf("Eval didn't return ClassInstance. got=%T (%+v)", evaluated, evaluated)
	}

	if result.Name != "Person" {
		t.Errorf("result.Name is not Person. got=%s", result.Name)
	}

	env := result.Env

	// Check fields
	nameField, ok := env.Get("name")
	if !ok {
		t.Errorf("'name' is not in the env of Person")
	}

	nameStr, ok := nameField.(*object.String)
	if !ok {
		t.Errorf("'name' is not of type string. got=%T", nameField)
	}

	if nameStr.Value != "" {
		t.Errorf("nameStr is not 'name'. got=%s", nameStr.Value)
	}

	ageField, ok := env.Get("age")
	if !ok {
		t.Errorf("'age' is not in the env of Person")
	}

	testIntegerObject(t, ageField, 0)

	// check function
	getName, ok := env.Get("getName")
	if !ok {
		t.Errorf("'getName' is not in the env of Person")
	}

	getNameFunc, ok := getName.(*object.Function)
	if !ok {
		t.Errorf("'getNameFunc' is not of type Function. got=%T", getName)
	}

	if getNameFunc.Parameters[0].Value != "dummyParam" {
		t.Errorf("getNameFunc.Parameters[0].Value is not 'dummyParam'. got=%s",
			getNameFunc.Parameters[0].Value)
	}

	// Check init
	init, ok := env.Get("init")
	if !ok {
		t.Errorf("'init' is not in the env of Person")
	}

	initFunc, ok := init.(*object.InitFunction)
	if !ok {
		t.Errorf("'initFunc' is not of type InitFunction. got=%T", initFunc)
	}

	if initFunc.Parameters[0].Parameter.Value != "name" {
		t.Errorf("initFunc.Parameters[0].Value is not 'name'. got=%s",
			initFunc.Parameters[0].Parameter.Value)
	}

	if initFunc.Parameters[0].IsThisParam != false {
		t.Errorf("initFunc.Parameters[0].IsThisParam is not 'false'. got=%t",
			initFunc.Parameters[0].IsThisParam)
	}

	if initFunc.Parameters[1].Parameter.Value != "age" {
		t.Errorf("initFunc.Parameters[1].Value is not 'age'. got=%s",
			initFunc.Parameters[1].Parameter.Value)
	}

	if initFunc.Parameters[1].IsThisParam != true {
		t.Errorf("initFunc.Parameters[1].IsThisParam is not 'true'. got=%t",
			initFunc.Parameters[1].IsThisParam)
	}
}

func TestObjectInitializationWithParameters(t *testing.T) {
	input := `
	class Person {
		var name
		var age

		init(n, this.age) {
			name = n
		}

		func getName() {
			return name
		}
	}
	var p = new Person("Hans", 10)
	return p
	`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.ClassInstance)
	if !ok {
		t.Fatalf("Eval didn't return ClassInstance. got=%T (%+v)", evaluated, evaluated)
	}

	if result.Name != "Person" {
		t.Errorf("result.Name is not Person. got=%s", result.Name)
	}

	env := result.Env

	// Check fields
	nameField, ok := env.Get("name")
	if !ok {
		t.Errorf("'name' is not in the env of Person")
	}

	nameStr, ok := nameField.(*object.String)
	if !ok {
		t.Errorf("'name' is not of type string. got=%T", nameField)
	}

	if nameStr.Value != "Hans" {
		t.Errorf("nameStr is not 'Hans'. got=%s", nameStr.Value)
	}

	ageField, ok := env.Get("age")
	if !ok {
		t.Errorf("'age' is not in the env of Person")
	}

	testIntegerObject(t, ageField, 10)
}

func TestObjectInitializationWithoutParameters(t *testing.T) {
	input := `
	class Person {
		var name = ""
		var age = 0

		func getName() {
			return name
		}
	}
	var p = new Person()
	return p
	`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.ClassInstance)
	if !ok {
		t.Fatalf("Eval didn't return ClassInstance. got=%T (%+v)", evaluated, evaluated)
	}

	if result.Name != "Person" {
		t.Errorf("result.Name is not Person. got=%s", result.Name)
	}

	env := result.Env

	// Check fields
	nameField, ok := env.Get("name")
	if !ok {
		t.Errorf("'name' is not in the env of Person")
	}

	nameStr, ok := nameField.(*object.String)
	if !ok {
		t.Errorf("'name' is not of type string. got=%T", nameField)
	}

	if nameStr.Value != "" {
		t.Errorf("nameStr is not 'name'. got=%s", nameStr.Value)
	}

	ageField, ok := env.Get("age")
	if !ok {
		t.Errorf("'age' is not in the env of Person")
	}

	testIntegerObject(t, ageField, 0)
}

func TestCallObjectFunction(t *testing.T) {
	input := `
	class Person {
		var name = "Hans"

		func getName() {
			return name
		}
	}
	var p = new Person()
	return p.getName()
	`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("Eval didn't return String. got=%T (%+v)", evaluated, evaluated)
	}

	if result.Value != "Hans" {
		t.Errorf("result.Value is not 'Hans'. got=%s", result.Value)
	}
}

func TestMultipleObjectInitializations(t *testing.T) {
	input := `
	class Person {
		var name
		var age

		init(n, this.age) {
			name = n
		}
	}
	var p = new Person("Hans", 10)
	var p2 = new Person("Ole", 15)
	var p3 = new Person("Jens", 20)
	var arr = [p, p2, p3]
	return arr
	`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("Eval didn't return Array. got=%T (%+v)", evaluated, evaluated)
	}

	names := []string{"Hans", "Ole", "Jens"}
	ages := []int{10, 15, 20}

	for i, elem := range result.Elements {
		person := elem.(*object.ClassInstance)

		name, _ := person.Env.Get("name")
		if name.Inspect() != names[i] {
			t.Errorf("person.Env.Get(name) is not %s. got=%s", names[i], name.Inspect())
		}

		age, _ := person.Env.Get("age")
		if age.Inspect() != strconv.Itoa(ages[i]) {
			t.Errorf("person.Env.Get(age) is not %d. got=%s", ages[i], age.Inspect())
		}
	}
}

func TestThisParametersIsSetBeforeExecutingInitBody(t *testing.T) {
	input := `
	class Person {
		var name

		init(this.name) {
			name = "OVERRIDDEN"
		}
	}
	var p = new Person("Hans")
	return p
	`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.ClassInstance)
	if !ok {
		t.Fatalf("Eval didn't return ClassInstance. got=%T (%+v)", evaluated, evaluated)
	}

	name, _ := result.Env.Get("name")
	if name.Inspect() != "OVERRIDDEN" {
		t.Errorf("person.Env.Get(name) is not 'OVERRIDEN'. got=%s", name.Inspect())
	}
}

func TestNullInitializationOfVariables(t *testing.T) {
	input := `
	var nullObj
	return nullObj`

	evaluated := testEval(input)
	_, ok := evaluated.(*object.Null)

	if !ok {
		t.Errorf("evaluated was not *object.Null. got=%T", evaluated)
	}
}

func TestLaterInitializationOfNullVariable(t *testing.T) {
	input := `
	var nullObj 
	nullObj = "Hello"
	return nullObj`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)

	if !ok {
		t.Errorf("evaluated was not *object.String. got=%T", evaluated)
	}

	if str.Value != "Hello" {
		t.Errorf("str.Value was not 'Hello'. got=%s", str.Value)
	}
}

func TestThisPrefixedIdentifier(t *testing.T) {
	input := `
	class Person {
		var name = "Hans"

		func getHans(name) {
			return this.name
		}

		func getArgument(name) {
			return name
		}
	}
	var p = new Person()
	return [p.getHans("Jens"), p.getArgument("Jens")]
	`

	evaluated := testEval(input)

	arr, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("evaluated is not *object.Array. got=%T (%+v)", evaluated, evaluated)
	}

	str, ok := arr.Elements[0].(*object.String)
	if !ok {
		t.Errorf("arr.Elements[0] is not *object.String. got=%T", arr.Elements[0])
	}

	if str.Value != "Hans" {
		t.Errorf("str.Value is not 'Hans'. got=%s", str.Value)
	}

	str2, ok := arr.Elements[1].(*object.String)
	if !ok {
		t.Errorf("arr.Elements[1] is not *object.String. got=%T", arr.Elements[1])
	}

	if str2.Value != "Jens" {
		t.Errorf("str2.Value is not 'Jens'. got=%s", str2.Value)
	}
}

//////////////////////////////
////// Helper functions //////
//////////////////////////////

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong Value. got=%d, expected=%d", result.Value, expected)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
