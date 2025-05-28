# Yarnball

Yarnball is an esoteric, stack-based programming language where every instruction reads like a crochet pattern.

As of now, the language is in its early stages, with a minimal set of instructions. Currently, it's pretty obvious that what you're reading isn't actually a crochet pattern, but in the future, I would like for a Yarnball program to be indistinguishable from a real crochet pattern.

## What is Yarnball?

In Yarnball, operations are expressed using crochet terminology:
- **ch (chain):** Pushes a number onto the stack.
- **pic (picot stitch):** Pops a value and prints it as a character.
- **yo (yarn over):** Pops a value and prints it as a number.
- **fo (finish off):** Immediately halts program execution.
- **subpattern**: A reusable sequence of instructions that can be invoked with **use** (like a function call).
- **AND MORE!**

Other instructions manipulate the stack (e.g., **dc**, **bob**, **hdc**) or control the flow with loops (`rep`) and conditionals (`if`). You can even define reusable stitch patterns with subpattern definitions and invoke them using **use**.

## Getting Started

To run a Yarnball program, use the command line. For example, to run the [fib.yarn](examples/fib.yarn) example, execute:

```sh
make
./bin/yarnball examples/fib.yarn
```

If you prefer an interactive environment, start the REPL by running:

```sh
make repl
```

## Repository Structure

- [cmd/main.go](cmd/main.go) - The application entry point that initializes the Yarnball interpreter.
- [pkg/evaluator](pkg/evaluator/evaluator.go) - Implements the evaluator that processes Yarnball instructions.
- [pkg/lexer](pkg/lexer/lexer.go) - Responsible for lexing Yarnball source code into tokens.
- [pkg/preprocessor/preprocessor.go](pkg/preprocessor/preprocessor.go) - Preprocesses Yarnball source code, handling comments and whitespace and other aesthetic features of the language.
- [pkg/parser/parser.go](pkg/parser/parser.go) - Parses Yarnball source code into an abstract syntax tree (AST).
- [docs/specification.md](docs/specification.md) - Provides a detailed description of Yarnball’s instructions and behavior.
- [examples/](examples/) - Contains sample Yarnball programs.

## Learn More

For a complete list of instructions and their behaviors, see the [Yarnball Language Specification](docs/specification.md).
