package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Suryansh-23/amrit/evaluator"
	"github.com/Suryansh-23/amrit/lexer"
	"github.com/Suryansh-23/amrit/object"
	"github.com/Suryansh-23/amrit/parser"
)

const PROMPT = ">>> "

// scans the prompt input until eol and starts up the lexer for it
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			PrintParserErrors(out, p.Errors())
			continue
		}

		stdout := []string{}
		evaluated := evaluator.Eval(program, env, &stdout)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

		for _, s := range stdout {
			io.WriteString(out, s)
		}
	}
}

func PrintParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
