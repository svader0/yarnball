# Yarnball

Yarnball is an esoteric, stack-based programming language where every instruction reads sort of like a crochet pattern.

This is just a fun little side project and doesn't implement anything particularly novel. Yarnball programs are also fairly distinguishable from real crochet patterns. Otherwise, it is technically totally turing complete and pretty fun! 

## What is Yarnball?

In Yarnball, operations are expressed using crochet terminology:
- **ch (chain):** Pushes a number onto the stack.
- **pic (picot stitch):** Pops a value and prints it as a character.
- **yo (yarn over):** Pops a value and prints it as a number.
- **fo (finish off):** Immediately halts program execution.
- **stitch**: Defines a reusable stitch pattern with `stitch name = (...)`, called by writing the name (or `use name`).
- **repeat**: Uses crochet-style blocks like `* ... * repeat while/until` or `* ... * repeat 3`.
- **AND MORE!**

Other instructions manipulate the stack (e.g., **dc**, **bob**, **hdc**) or control the flow with loops (`repeat`) and conditionals (`if`).

## Getting Started

To run a Yarnball program, use the command line. For example, to run the [loop_counter.yarn](examples/loop_counter.yarn) example, execute:

```sh
make
./bin/yarnball examples/fib.yarn
```

If you prefer an interactive environment, start the REPL by running:

```sh
make repl
```

or 

```sh
make # to build
./yarnball
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
