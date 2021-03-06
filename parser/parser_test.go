package parser

import (
	"Pron-Lang/ast"
	"Pron-Lang/lexer"
	"fmt"
	"testing"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error %q", msg)
	}
	t.FailNow()
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	testLiteralExpression(t, stmt.Expression, 5)
}

func TestRealLiteralExpression(t *testing.T) {
	input := "5.4"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	testLiteralExpression(t, stmt.Expression, 5.4)
}

func TestVarStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"var x = 5;", "x", 5},
		{"var y = true;", "y", true},
		{"var foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testVarStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.VarStatement).Value

		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return y;", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}
		stmt := program.Statements[0]

		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}

		expression := returnStmt.ReturnValue
		if !testLiteralExpression(t, expression, tt.expectedValue) {
			return
		}
	}
}
func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Staments does not contain %d statement. got=%d",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.PrefixExpression. got=%T",
				stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}

		intValue, ok := tt.value.(int64)
		if ok && !testIntegerLiteral(t, exp.Right, intValue) {
			return
		}
	}
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] was not ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.InfixExpression. got=%T",
				stmt.Expression)
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}

		if !testLiteralExpression(t, exp.Left, tt.leftValue) {
			return
		}

		if !testLiteralExpression(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == false",
			"((3 < 5) == false)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			"true",
			true,
		},
		{
			"false",
			false,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] was not Boolean. got=%T",
				program.Statements[0])
		}

		testLiteralExpression(t, stmt.Expression, tt.expected)
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("exp.Consequence.Statements is not 1 statement. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("exp.Consequence.Statements is not 1 statement. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements is not 1 statement. got=%d\n",
			len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Alternative.Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestElifExpression(t *testing.T) {
	input := `if (x < y) { x } elif (x > y) { y } elif (x == y) { 10 } else { z } `

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.ElseIfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ElseIfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.ConditionAndBlockstatementList[0].Condition, "x", "<", "y") {
		return
	}

	if len(exp.ConditionAndBlockstatementList[0].Consequence.Statements) != 1 {
		t.Errorf("exp.ConditionAndBlockstatementList[0].Consequence.Statements is not 1 statement. got=%d\n",
			len(exp.ConditionAndBlockstatementList[0].Consequence.Statements))
	}

	consequence, ok := exp.ConditionAndBlockstatementList[0].Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.ConditionAndBlockstatementList[0].Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if !testInfixExpression(t, exp.ConditionAndBlockstatementList[1].Condition, "x", ">", "y") {
		return
	}

	if len(exp.ConditionAndBlockstatementList[1].Consequence.Statements) != 1 {
		t.Errorf("exp.ConditionAndBlockstatementList[1].Consequence.Statements is not 1 statement. got=%d\n",
			len(exp.ConditionAndBlockstatementList[1].Consequence.Statements))
	}

	consequence, ok = exp.ConditionAndBlockstatementList[1].Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[1] is not ast.ExpressionStatement. got=%T",
			exp.ConditionAndBlockstatementList[1].Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "y") {
		return
	}

	if !testInfixExpression(t, exp.ConditionAndBlockstatementList[2].Condition, "x", "==", "y") {
		return
	}

	if len(exp.ConditionAndBlockstatementList[2].Consequence.Statements) != 1 {
		t.Errorf("exp.ConditionAndBlockstatementList[2].Consequence.Statements is not 1 statement. got=%d\n",
			len(exp.ConditionAndBlockstatementList[2].Consequence.Statements))
	}

	consequence, ok = exp.ConditionAndBlockstatementList[2].Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[2] is not ast.ExpressionStatement. got=%T",
			exp.ConditionAndBlockstatementList[2].Consequence.Statements[0])
	}

	if !testIntegerLiteral(t, consequence.Expression, 10) {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements is not 1 statement. got=%d\n",
			len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Alternative.Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "z") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `var add = func(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.VarStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.VarStatement. got=%T",
			program.Statements[0])
	}

	function, ok := stmt.Value.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Value)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "var doNoting = func() {};", expectedParams: []string{}},
		{input: "var doNoting = func(x) {};", expectedParams: []string{"x"}},
		{input: "var doNoting = func(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.VarStatement)
		function := stmt.Value.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestDirectFunctionStatementParsing(t *testing.T) {
	input := `func add(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.DirectFunctionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.DirectFunctionStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Name.String() != "add" {
		t.Errorf("stmt.Name is not add. got=%s", stmt.Name)
	}
}

func TestDirectFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "func doNothing() {};", expectedParams: []string{}},
		{input: "func doNothing(x) {};", expectedParams: []string{"x"}},
		{input: "func doNothing(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.DirectFunctionStatement)

		if len(stmt.Function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(stmt.Function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, stmt.Function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "add();", expectedParams: []string{}},
		{input: "add(x);", expectedParams: []string{"x"}},
		{input: "add(x, y, z);", expectedParams: []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.CallExpression)

		if len(function.Arguments) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Arguments))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Arguments[i], ident)
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"Hello World";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "Hello World" {
		t.Errorf("literal.Value is not %q. got=%q", "Hello World", literal.Value)
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}

		expectedValue := expected[literal.String()]

		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsIntegerKeys(t *testing.T) {
	input := `{1: 1, 2: 2, 3: 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	expected := map[int64]int64{
		1: 1,
		2: 2,
		3: 3,
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("key is not ast.IntegerLiteral. got=%T", key)
		}

		expectedValue := expected[literal.Value]

		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsBooleanKeys(t *testing.T) {
	input := `{true: 1, false: 2}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 2 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	expected := map[bool]int64{
		true:  1,
		false: 2,
	}

	for key, value := range hash.Pairs {
		boolLiteral, ok := key.(*ast.Boolean)
		if !ok {
			t.Errorf("key is not ast.Boolean. got=%T", key)
		}

		expectedValue := expected[boolLiteral.Value]

		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}
	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}
		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}
		testFunc(value)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("stmt.Expression was not *ast.HastLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("hast.Pairs does not contain 0 elements. contains=%d", len(hash.Pairs))
	}
}

func TestIncrementForloopExpression(t *testing.T) {
	input := `var x = 10; for (i from 0 to x) { i };`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[1])
	}

	exp, ok := stmt.Expression.(*ast.IncrementForloopExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IncrementForloopExpression. got=%T", stmt.Expression)
	}

	if exp.LocalVar.String() != "i" {
		t.Fatalf("exp.LocalVar is not i. got=%s", exp.LocalVar.String())
	}

	if exp.From.TokenLiteral() != "0" {
		t.Fatalf("exp.From is not 0. got=%s", exp.From.TokenLiteral())
	}

	if exp.To.TokenLiteral() != "x" {
		t.Fatalf("exp.From is not x. got=%s", exp.To.TokenLiteral())
	}

	if len(exp.Body.Statements) != 1 {
		t.Errorf("exp.Body.Statements is not 1 statement. got=%d\n",
			len(exp.Body.Statements))
	}

	bodyExp, ok := exp.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Body.Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Body.Statements[0])
	}

	if !testIdentifier(t, bodyExp.Expression, "i") {
		return
	}
}

func TestArrayForloopExpression(t *testing.T) {
	input := `var myArr = [1,4,2]; for (i in myArr) { i };`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[1])
	}

	exp, ok := stmt.Expression.(*ast.ArrayForloopExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ArrayForloopExpression. got=%T", stmt.Expression)
	}

	if exp.LocalVar.String() != "i" {
		t.Fatalf("exp.LocalVar is not i. got=%s", exp.LocalVar.String())
	}

	if exp.ArrayName.String() != "myArr" {
		t.Fatalf("exp.ArrayName.String() is not myArr. got=%s", exp.ArrayName.String())
	}

	if len(exp.Body.Statements) != 1 {
		t.Errorf("exp.Body.Statements is not 1 statement. got=%d\n",
			len(exp.Body.Statements))
	}

	bodyExp, ok := exp.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Body.Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Body.Statements[0])
	}

	if !testIdentifier(t, bodyExp.Expression, "i") {
		return
	}
}

func TestChangeValueOfExistingVariable(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"var x = 5; x = 6;", "x", 6},
		{"var y = true; y = false", "y", false},
		{"var foobar = y; foobar = x", "foobar", "x"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 2 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testVarStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		expStmt, ok := program.Statements[1].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[1] was not ExpressionStatement. got=%T", program.Statements[1])
		}

		infixExp, ok := expStmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("expStmt.Expression is not *ast.InfixExpression. got=%T", expStmt.Expression)
		}

		if !testInfixExpression(t, infixExp, tt.expectedIdentifier, "=", tt.expectedValue) {
			return
		}
	}
}

func TestClassStatementParsing(t *testing.T) {
	input := `
	class Person {
		var name = ""
		var age = 0

		Init(name, this.age) {
			name = "Hans"
		}

		func GetName() {
			return name
		}

		func GetAge() {
			return age
		}
	}
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ClassStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ClassStatement. got=%T",
			program.Statements[0])
	}

	// Name
	if stmt.Name.Value != "Person" {
		t.Errorf("Name of class is not Person. got=%s", stmt.Name.Value)
	}

	// Fields
	if stmt.Fields[0].Name.Value != "name" {
		t.Errorf("Name of first Varstatement is not name. got=%s", stmt.Fields[0].Name.Value)
	}

	if stmt.Fields[0].Value.String() != "" {
		t.Errorf("Value of first Varstatement is not empty string. got=%s", stmt.Fields[0].Value.String())
	}

	if stmt.Fields[1].Name.Value != "age" {
		t.Errorf("Name of second Varstatement is not age. got=%s", stmt.Fields[1].Name.Value)
	}

	if stmt.Fields[1].Value.String() != "0" {
		t.Errorf("Value of second Varstatement is not 0. got=%s", stmt.Fields[1].Value.String())
	}

	// Init function
	if stmt.InitParams[0].Parameter.Value != "name" {
		t.Errorf("Name of first InitParameter is not name. got=%s", stmt.InitParams[0].Parameter.Value)
	}

	if stmt.InitParams[0].IsThisParam != false {
		t.Errorf("First InitParameter does not have 'false' isThisParam. got=%t",
			stmt.InitParams[0].IsThisParam)
	}

	if stmt.InitParams[1].Parameter.Value != "age" {
		t.Errorf("Name of second InitParameter is not name. got=%s", stmt.InitParams[1].Parameter.Value)
	}

	if stmt.InitParams[1].IsThisParam != true {
		t.Errorf("Second InitParameter does not have 'true' isThisParam. got=%t",
			stmt.InitParams[1].IsThisParam)
	}

	if stmt.Functions[0].Name.String() != "GetName" {
		t.Errorf("stmt.Functions[0].Name is not GetName. got=%s", stmt.Functions[0].Name.String())
	}

	if stmt.Functions[1].Name.String() != "GetAge" {
		t.Errorf("stmt.Functions[1].Name is not GetAge. got=%s", stmt.Functions[1].Name.String())
	}
}

func TestClassInitializationParsing(t *testing.T) {
	input := `new Person("Hans", 10)`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	objectInitExp, ok := stmt.Expression.(*ast.ObjectInitialization)
	if !ok {
		t.Fatalf("exp not *ast.ObjectInitialization. got=%T", stmt.Expression)
	}

	if objectInitExp.Name.Value != "Person" {
		t.Errorf("objectInitExp.Name is not 'Person'. got=%s", objectInitExp.Name.Value)
	}

	strArg, ok := objectInitExp.Arguments[0].(*ast.StringLiteral)
	if !ok {
		t.Errorf("strArg is not *ast.Identifier. got=%T", objectInitExp.Arguments[0])
	}

	if strArg.Value != "Hans" {
		t.Errorf("strArg.Value is not 'Person'. got=%s", strArg.Value)
	}

	intArg, ok := objectInitExp.Arguments[1].(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("objectInitExp.Arguments[1] is not *ast.IntegerLiteral. got=%T",
			objectInitExp.Arguments[1])
	}

	if intArg.Value != 10 {
		t.Errorf("intArg.Value is not 10. got=%d", intArg.Value)
	}
}

func TestCallObjectFunction(t *testing.T) {
	input := `p.changeName("Ole")`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallObjectFunction)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.CallObjectFunction. got=%T", stmt.Expression)
	}

	if exp.ObjectName.Value != "p" {
		t.Errorf("exp.ObjectName.Value is not 'p'. got=%s", exp.ObjectName.Value)
	}

	if exp.FunctionName.Value != "changeName" {
		t.Errorf("exp.FunctionName.Value is not 'changeName'. got=%s", exp.FunctionName.Value)
	}

	if len(exp.Arguments) != 1 {
		t.Errorf("len(exp.Arguments) is not 1. got=%d", len(exp.Arguments))
	}

	if exp.Arguments[0].String() != "Ole" {
		t.Errorf("exp.Arguments[0].String() is not 'Ole'. got=%s", exp.Arguments[0].String())
	}
}

func TestNullInitializaitonOfVarStatements(t *testing.T) {
	input := `
	var x
	var y
	var z`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	identifiers := []string{"x", "y", "z"}

	for i, ident := range identifiers {
		stmt := program.Statements[i]
		if !testVarStatement(t, stmt, ident) {
			return
		}

		val := stmt.(*ast.VarStatement).Value
		_, ok := val.(*ast.Null)

		if !ok {
			t.Errorf("VarStatement.Value is not *ast.Null. got=%T", val)
		}
	}
}

func TestThisKeywordGivesTheOuterMostVariable(t *testing.T) {
	input := `this.varName`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.Identifier. got=%T", program.Statements[0])
	}

	if ident.TokenLiteral() != "varName" {
		t.Errorf("ident.Token.Literal is not IDENT. got=%s", ident.Token.Literal)
	}

	if ident.Value != "varName" {
		t.Errorf("ident.Value is not 'varName'. got=%s", ident.Value)
	}

	if !ident.HasThisPrefix {
		t.Errorf("ident.HasThisPrefix is not 'true'. got=%t", ident.HasThisPrefix)
	}
}

func TestIncrementIdentifier(t *testing.T) {
	input := `i++`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T", program.Statements[0])
	}

	increment, ok := stmt.Expression.(*ast.Increment)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.Increment. got=%T", stmt.Expression)
	}

	if increment.TokenLiteral() != "++" {
		t.Errorf("increment.TokenLiteral() is not '++'. got=%s", increment.TokenLiteral())
	}

	if increment.Name.Value != "i" {
		t.Errorf("increment.Name.Value is not 'i'. got=%s", increment.Name.Value)
	}
}

func TestDecrementIdentifier(t *testing.T) {
	input := `i--`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T", program.Statements[0])
	}

	decrement, ok := stmt.Expression.(*ast.Decrement)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.Decrement. got=%T", stmt.Expression)
	}

	if decrement.TokenLiteral() != "--" {
		t.Errorf("decrement.TokenLiteral() is not '--'. got=%s", decrement.TokenLiteral())
	}

	if decrement.Name.Value != "i" {
		t.Errorf("decrement.Name.Value is not 'i'. got=%s", decrement.Name.Value)
	}
}

func TestBlockComment(t *testing.T) {
	input := `
	var i = 0
	i++
	i--
	/*
	i++
	i++
	i++
	*/
	i++`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 5 {
		t.Fatalf("program.Statements does not contain 4 statement. got=%d",
			len(program.Statements))
	}

	if _, ok := program.Statements[0].(*ast.VarStatement); !ok {
		t.Errorf("program.Statements[0] is not *ast.VarStatement. got=%T", program.Statements[0])
	}

	stmt := program.Statements[1].(*ast.ExpressionStatement)
	if _, ok := stmt.Expression.(*ast.Increment); !ok {
		t.Errorf("program.Statements[1] is not *ast.Increment. got=%T", stmt.Expression)
	}

	stmt = program.Statements[2].(*ast.ExpressionStatement)
	if _, ok := stmt.Expression.(*ast.Decrement); !ok {
		t.Errorf("program.Statements[2] is not *ast.Decrement. got=%T", stmt.Expression)
	}

	stmt = program.Statements[3].(*ast.ExpressionStatement)
	if _, ok := stmt.Expression.(*ast.Null); !ok {
		t.Errorf("program.Statements[3] is not *ast.Null. got=%T", stmt.Expression)
	}

	stmt = program.Statements[4].(*ast.ExpressionStatement)
	if _, ok := stmt.Expression.(*ast.Increment); !ok {
		t.Errorf("program.Statements[4] is not *ast.Increment. got=%T", stmt.Expression)
	}
}

func TestPublicAndPrivateMethods(t *testing.T) {
	input := `class Person {
		var name = ""

		Init(this.name) {}

		func GetName() {
			return name
		}

		func doPrivateStuff() {}
	}
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ClassStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ClassStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Functions[0].Name.Value != "GetName" {
		t.Errorf("Name of first function is not GetName. got=%s", stmt.Functions[0].Name.Value)
	}

	if stmt.Functions[0].IsPublic != true {
		t.Errorf("GetName's IsPublic is not true. got=%t", stmt.Functions[0].IsPublic)
	}

	if stmt.Functions[1].Name.Value != "doPrivateStuff" {
		t.Errorf("Name of first function is not doPrivateStuff. got=%s", stmt.Functions[1].Name.Value)
	}

	if stmt.Functions[1].IsPublic != false {
		t.Errorf("doPrivateStuff's IsPublic is not false. got=%t", stmt.Functions[1].IsPublic)
	}
}

///////////////////////////////////////////
//////////// Helper functions /////////////
///////////////////////////////////////////
func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	case float64:
		return testRealLiteral(t, exp, float64(v))
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "var" {
		t.Errorf("s.TokenLiteral not 'var'. got=%q", s.TokenLiteral())
		return false
	}

	varStmt, ok := s.(*ast.VarStatement)
	if !ok {
		t.Errorf("s not *ast.VarStatement. got=%T", s)
		return false
	}

	if varStmt.Name.Value != name {
		t.Errorf("varStmt.Name.Value not '%s'. got=%s", name, varStmt.Value)
		return false
	}

	if varStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, varStmt.Name)
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
	}
	return true
}

func testRealLiteral(t *testing.T, il ast.Expression, value float64) bool {
	real, ok := il.(*ast.RealLiteral)
	if !ok {
		t.Errorf("il not *ast.RealLiteral. got=%T", il)
		return false
	}

	if real.Value != value {
		t.Errorf("real.Value not %f. got=%f", value, real.Value)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}
	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}
	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}
	return true
}
