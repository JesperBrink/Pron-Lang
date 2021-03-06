package repl

import (
	"Pron-Lang/evaluator"
	"Pron-Lang/lexer"
	"Pron-Lang/object"
	"Pron-Lang/parser"
	"bufio"
	"fmt"
	"io"
)

const PROMT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "quit" {
			return
		}
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			PrintParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated.Inspect() != "null" {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func PrintParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+"- "+msg+"\n")
	}
}
