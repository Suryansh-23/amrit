package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Suryansh-23/amrit/token"

	"github.com/Suryansh-23/amrit/lexer"
)

const PROMPT = ">>> "

// scans the prompt input until eol and starts up the lexer for it
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()

		l := lexer.New(line)
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
