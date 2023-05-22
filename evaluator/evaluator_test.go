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
			"agar (10 > 1) { satya + asatya| }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
agar (10 > 1) {
	agar (10 > 1) {
		labh satya + asatya|
	}
	
	labh 1|
}
`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"acha" - "kaise ho?"`,
			"unknown operator: STRING - STRING",
		},
		{
			`{"name": "Monkey"}[karya(x) { x }]|`,
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

func TestFunctionObject(t *testing.T) {
	input := "karya (x) { x + 2 | } |"
	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"mana identity = karya(x) { x | } | identity(5) |", 5},
		{"mana identity = karya(x) { labh x | } | identity(5) |", 5},
		{"mana double = karya(x) { x * 2 | } | double(5) |", 10},
		{"mana add = karya(x, y) { x + y | } | add(5, 5) |", 10},
		{"mana add = karya(x, y) { x + y | } | add(5 + 5, add(5, 5)) |", 20},
		{"karya(x) { x | }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
mana newAdder = karya(x) {
	karya(y) { x + y }|
}|
mana addTwo = newAdder(2)|
addTwo(2)|`
	testIntegerObject(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"namaste duniya!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "namaste duniya!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Namaste" + " " + "Duniya!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Namaste Duniya!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`lambai("")`, 0},
		{`lambai("panch")`, 5},
		{`lambai("namaste")`, 7},
		{`lambai(1)`, "argument to `lambai` not supported, got INTEGER"},
		{`lambai("ek", "do")`, "wrong number of arguments. got=2, want=1"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)

			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
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
			"mana i = 0| [1][i]|",
			1,
		},
		{
			"[1, 2, 3][1 + 1]|",
			3,
		},
		{
			"mana myArray = [1, 2, 3]| myArray[2]|",
			3,
		},
		{
			"mana myArray = [1, 2, 3]| myArray[0] + myArray[1] + myArray[2]|",
			6,
		},
		{
			"mana myArray = [1, 2, 3]| mana i = myArray[0]| myArray[i]",
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
	input := `mana do = "do"|
	{
	"ek": 10 - 9,
	do: 1 + 1,
	"te" + "en": 6 / 2,
	4: 4,
	satya: 5,
	asatya: 6
	}`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}
	expected := map[object.HashKey]int64{
		(&object.String{Value: "ek"}).HashKey():   1,
		(&object.String{Value: "do"}).HashKey():   2,
		(&object.String{Value: "teen"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():     4,
		TRUE.HashKey():                            5,
		FALSE.HashKey():                           6,
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
			`mana key = "foo"| {"foo": 5}[key]`,
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
			`{satya: 5}[satya]`,
			5,
		},
		{
			`{asatya: 5}[asatya]`,
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
