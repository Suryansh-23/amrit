package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"unicode/utf8"

	"github.com/Suryansh-23/amrit/evaluator"
	"github.com/Suryansh-23/amrit/lexer"
	"github.com/Suryansh-23/amrit/object"
	"github.com/Suryansh-23/amrit/parser"
	"github.com/Suryansh-23/amrit/repl"
)

const EXT = ".amr"

func main() {
	out := os.Stdout
	in := os.Stdin
	if len(os.Args) == 1 {

		user, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Hello %s! This is the Amrit Programming Language!\n",
			user.Username)
		fmt.Printf("Feel free to type in commands\n")
		repl.Start(in, out)
	} else {
		fpath := os.Args[1]
		if filepath.Ext(fpath) != EXT {
			io.WriteString(out, fmt.Sprintf("the file type is invalid: `%s`, must be of `.amr` filetype\n", fpath))
			return
		}

		progScript, err := ioutil.ReadFile(fpath)
		if err != nil {
			io.WriteString(out, fmt.Sprintf("the following error occured while opening %s:\n\t%s", fpath, err.Error()))
			return
		}

		if !utf8.Valid(progScript) {
			io.WriteString(out, fmt.Sprintf("the file `%s` is not encoded in utf-8 format, try resaving the file again.\n", fpath))
			return
		}

		env := object.NewEnvironment()
		l := lexer.New(string(progScript))
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			repl.PrintParserErrors(out, p.Errors())
			return
		}

		stdout := []string{}
		evaluator.Eval(program, env, &stdout)
		for _, s := range stdout {
			io.WriteString(out, s)
		}
	}
}
