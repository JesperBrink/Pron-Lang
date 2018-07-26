package main

import (
	"Pron-Lang/evaluator"
	"Pron-Lang/lexer"
	"Pron-Lang/object"
	"Pron-Lang/parser"
	"Pron-Lang/repl"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	out := os.Stdout

	if len(os.Args) > 1 {
		// Check file format
		//docIndex :0 strings.Index(filename, ".")

		input, err := ioutil.ReadFile(os.Args[1])
		check(err)

		env := object.NewEnvironment()

		l := lexer.New(string(input))
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			PrintParserErrors(out, p.Errors())
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	} else {
		user, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Hello %s! Welcome to Pron-Lang \n", user.Username)
		repl.Start(os.Stdin, os.Stdout)
	}

}

func PrintParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+"- "+msg+"\n")
	}
}
