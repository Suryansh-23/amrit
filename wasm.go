package main

import (
	"fmt"
	"syscall/js"

	"github.com/Suryansh-23/amrit/evaluator"
	"github.com/Suryansh-23/amrit/lexer"
	"github.com/Suryansh-23/amrit/object"
	"github.com/Suryansh-23/amrit/parser"
)

func ScriptMode(prog string) string {
	env := object.NewEnvironment()
	l := lexer.New(prog)
	p := parser.New(l)
	s := ""

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			s += "\t" + msg + "\n"
		}
		return s
	}

	stdout := []string{}
	evaluator.Eval(program, env, &stdout)
	for _, output := range stdout {
		s += output
	}

	return s
}

func ReplMode(env *object.Environment, line string) string {
	l := lexer.New(line)
	p := parser.New(l)
	s := ""

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			s += "\t" + msg + "\n"
		}
		return s
	}

	stdout := []string{}
	evaluated := evaluator.Eval(program, env, &stdout)
	if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
		s += evaluated.Inspect() + "\n"
	}

	for _, output := range stdout {
		s += output
	}

	return s
}

func main() {
	c := make(chan struct{}, 0)

	fmt.Println("Namaste Duniya🙏🏻, This is the Amrit Programming Language!")
	env := object.NewEnvironment()

	js.Global().Set("ScriptMode", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return ScriptMode(args[0].String())
	}))
	js.Global().Set("ReplMode", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return ReplMode(env, args[0].String())
	}))

	<-c
}
