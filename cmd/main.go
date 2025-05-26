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

// TODO:
/*
 - Add a proper preprocessor to handle comments and whitespace
 - Add a program trace / debug mode
 - Implement a more robust error handling system
 - Support for multi-line instructions
 - Add more built-in patterns
*/

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
	fmt.Println("Yarnball REPL :) — type `\\q` to quit.")
	ev := evaluator.New()

	var inputBuilder strings.Builder
	for {
		fmt.Print("=> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if strings.TrimSpace(line) == "\\q" {
			break
		}

		// handle print stack command
		if strings.TrimSpace(line) == ".s" {
			fmt.Println("Stack:", ev.Stack())
			continue
		}

		// Accumulate multi-line input
		inputBuilder.WriteString(line + "\n")
		if isCompleteInput(inputBuilder.String()) {
			// lex → parse → eval
			l := lexer.New(inputBuilder.String())
			p := parser.New(l)
			prog, err := p.ParseProgram()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
				inputBuilder.Reset()
				continue
			}

			if err := ev.Eval(prog); err != nil && err.Error() != "FO: halt" {
				fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
			}
			inputBuilder.Reset()
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

	// Process the entire file as a single program
	l := lexer.New(input)
	p := parser.New(l)
	prog, err := p.ParseProgram()
	if err != nil {
		return fmt.Errorf("Parse error: %v", err)
	}

	ev := evaluator.New()
	if err := ev.Eval(prog); err != nil && err.Error() != "FO: halt" {
		return fmt.Errorf("Runtime error: %v", err)
	}
	return nil
}

// Helper function to check if the input is complete
func isCompleteInput(input string) bool {
	openParens := strings.Count(input, "(")
	closeParens := strings.Count(input, ")")
	return openParens == closeParens
}
