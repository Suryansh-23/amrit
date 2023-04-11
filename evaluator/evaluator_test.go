package evaluator

import (
	"testing"

	"github.com/Suryansh-23/amrit/lexer"
	"github.com/Suryansh-23/amrit/object"
	"github.com/Suryansh-23/amrit/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
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

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(prog, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"satya", true},
		{"asatya", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"satya == satya", true},
		{"asatya == asatya", true},
		{"satya == asatya", false},
		{"satya != asatya", true},
		{"asatya != satya", true},
		{"(1 < 2) == satya", true},
		{"(1 < 2) == asatya", false},
		{"(1 > 2) == satya", false},
		{"(1 > 2) == asatya", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
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

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!satya", false},
		{"!asatya", true},
		{"!5", false},
		{"!!satya", true},
		{"!!asatya", false},
		{"!!5", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

// why is this test case failing ?
func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"agar (satya) { 10 }", 10},
		{"agar (asatya) { 10 }", nil},
		{"agar (1) { 10 }", 10},
		{"agar (1 < 2) { 10 }", 10},
		{"agar (1 > 2) { 10 }", nil},
		{"agar (1 > 2) { 10 } varna { 20 }", 20},
		{"agar (1 < 2) { 10 } varna { 20 }", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := (tt.expected).(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"labh 10 |", 10},
		{"labh 10 | 9 |", 10},
		{"labh 2 * 5 | 9 |", 10},
		{"9 | labh 2 * 5 | 9 |", 10},
		{`agar (10 > 1) {
			agar (10 > 1) {
				labh 10 |
			}

			labh 1 | 
		}
		`, 10},
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
			"5 + satya |",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + satya; 5 |",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-satya",
			"unknown operator: -BOOLEAN",
		},
		{
			"satya + asatya |",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5| satya + asatya| 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { satya + asatya| }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`agar (10 > 1) {
		agar (10 > 1) {
		labh satya + asatya |
		}
		labh 1 |
		}
		`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
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

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"mana a = 5 | a |", 5},
		{"mana a = 5 * 5 | a |", 25},
		{"mana a = 5 | mana b = a | b |", 5},
		{"mana a = 5 | mana b = a | mana c = a + b + 5 | c |", 15},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}
