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
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	out := os.Stdout

	if len(os.Args) > 1 {
		filename := os.Args[1]

		docIndex := strings.Index(filename, ".")

		// Check for missing fileformat
		if docIndex == -1 {
			fmt.Print("ERROR: Missing file type, should be .pron\n")
			os.Exit(0)
		}

		fileType := string(filename[docIndex:])

		// Check if fileformat is .pron
		if fileType != ".pron" {
			fmt.Print("ERROR: Filetype is not .pron\n")
			os.Exit(0)
		}

		// Run Program
		input, err := ioutil.ReadFile(filename)
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
		// Start REPL
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
