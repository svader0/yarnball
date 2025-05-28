package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/svader0/yarnball/pkg/evaluator"
	"github.com/svader0/yarnball/pkg/lexer"
	"github.com/svader0/yarnball/pkg/parser"
	"github.com/svader0/yarnball/pkg/preprocessor"
)

// TODO:
/*
 - Make the preprocessor more robust and make it respect line numbers
 - Add a program trace / debug mode
 - Implement a more robust error handling system
 - Change language spec to look more like actual crochet
 - ADD SUPPORT FOR INPUT (e.g. reading from stdin)
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
	handler := log.New(os.Stderr)
	logger := slog.New(handler)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Yarnball REPL :) — type `\\q` to quit.")
	ev := evaluator.New(logger)

	var inputBuilder strings.Builder
	pre := preprocessor.New() // Create preprocessor once

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
			fmt.Println("Stack:", ev.Stack(), " <-- top ")
			continue
		}

		// Accumulate multi-line input
		inputBuilder.WriteString(line + "\n")
		if isCompleteInput(inputBuilder.String()) {
			processed, err := pre.Process(inputBuilder.String())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Preprocessing error: %v\n", err)
				inputBuilder.Reset()
				continue
			}

			// lex -> parse -> eval
			l := lexer.New(processed)
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
	handler := log.New(os.Stderr)
	// handler.SetLevel(log.DebugLevel)
	logger := slog.New(handler)

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	input := string(data)

	// Preprocess the input before lexing/parsing
	pre := preprocessor.New()
	input, err = pre.Process(input)
	if err != nil {
		return fmt.Errorf("Preprocessing error: %v", err)
	}

	// Process the entire file as a single program
	l := lexer.New(input)
	p := parser.New(l)
	prog, err := p.ParseProgram()
	if err != nil {
		return fmt.Errorf("Parse error: %v", err)
	}

	ev := evaluator.New(logger)
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
