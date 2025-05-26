package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/svader0/yarnball/pkg/evaluator"
	"github.com/svader0/yarnball/pkg/lexer"
	"github.com/svader0/yarnball/pkg/parser"
)

func main() {
	if len(os.Args) > 1 {
		if err := runFile(os.Args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		repl()
	}
}

func repl() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Yarnball REPL — type `exit` to quit.")
	ev := evaluator.New()

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if strings.TrimSpace(line) == "exit" {
			break
		}

		// handle print stack command
		if strings.TrimSpace(line) == "print stack" {
			fmt.Println("Stack:", ev.Stack())
			continue
		}

		// lex → parse → eval
		l := lexer.New(line)
		p := parser.New(l)
		prog, err := p.ParseProgram()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
			continue
		}

		if err := ev.Eval(prog); err != nil && err.Error() != "FO: halt" {
			fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		}
	}
	fmt.Println("Goodbye.")
}

func runFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	input := string(data)
	l := lexer.New(input)
	for {
		tok := l.NextToken()
		fmt.Printf("%-8s literal=%q  line=%d col=%d\n",
			tok.Type, tok.Literal, tok.Line, tok.Column)
		if tok.Type == lexer.EOF {
			break
		}
	}
	return nil
}
