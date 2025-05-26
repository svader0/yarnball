package evaluator

import (
	"fmt"
	"strconv"

	"github.com/svader0/yarnball/pkg/parser"
	"github.com/svader0/yarnball/pkg/stack"
)

type Evaluator struct {
	stack    *stack.Stack
	patterns map[string]*parser.PatternDef
}

func New() *Evaluator {
	return &Evaluator{
		stack:    stack.New(),
		patterns: make(map[string]*parser.PatternDef),
	}
}

// Passes the stack to the evaluator, allowing access to it from outside
// e.g. for debugging or inspection.
func (e *Evaluator) Stack() *stack.Stack {
	return e.stack
}

func (e *Evaluator) Eval(prog *parser.Program) error {
	for _, instr := range prog.Instructions {
		switch node := instr.(type) {
		case *parser.SimpleInstr:
			if err := e.execSimple(node); err != nil {
				return err
			}
		case *parser.RepInstr:
			if err := e.execRep(node); err != nil {
				return err
			}
		case *parser.PatternDef:
			e.patterns[node.Name] = node // Store the pattern
		case *parser.UseInstr:
			if err := e.execUse(node); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Evaluator) exec(instr parser.Instruction) error {
	switch node := instr.(type) {
	case *parser.SimpleInstr:
		return e.execSimple(node)
	case *parser.RepInstr:
		return e.execRep(node)
	case *parser.UseInstr:
		return e.execUse(node)
	case *parser.PatternDef:
		e.patterns[node.Name] = node // Store the pattern definition
		return nil
	default:
		return fmt.Errorf("unknown instruction type: %T", instr)
	}
}

func (e *Evaluator) execUse(ui *parser.UseInstr) error {
	pat, exists := e.patterns[ui.Name]
	if !exists {
		return fmt.Errorf("undefined pattern %q", ui.Name)
	}

	// Push arguments onto the stack
	for _, arg := range ui.Args {
		n, _ := strconv.Atoi(arg)
		e.stack.Push(n)
	}
	// Execute the instructions in the pattern
	for _, instr := range pat.Body {
		if err := e.exec(instr); err != nil {
			return err
		}
	}

	return nil
}

func (e *Evaluator) execSimple(si *parser.SimpleInstr) error {
	switch si.Token {
	case "ch":
		n, _ := strconv.Atoi(si.Args[0])
		e.stack.Push(n)
	case "pic":
		n, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("pic: %w", err)
		}
		fmt.Printf("%c", n)
	case "yo":
		n, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("yo: %w", err)
		}
		fmt.Println(n)
	case "FO":
		return fmt.Errorf("FO: halt")
	case "sc":
		// pop top value
		if e.stack.IsEmpty() {
			return fmt.Errorf("sc: stack underflow")
		}
		e.stack.Pop()
	case "dc":
		// product of top two values
		if e.stack.Size() < 2 {
			return fmt.Errorf("dc: stack underflow")
		}
		top, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("dc: %w", err)
		}
		second, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("dc: %w", err)
		}
		if second == 0 {
			return fmt.Errorf("dc: division by zero")
		}
		new := top * second
		e.stack.Push(new) // push product
	case "bob":
		// add top two values
		if e.stack.Size() < 2 {
			return fmt.Errorf("dc: stack underflow")
		}
		top, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("bob: %w", err)
		}
		second, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("bob: %w", err)
		}
		new := top + second
		e.stack.Push(new)
	case "hdc":
		// subtract top two values
		if e.stack.Size() < 2 {
			return fmt.Errorf("hdc: stack underflow")
		}
		top, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("hdc: %w", err)
		}
		second, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("hdc: %w", err)
		}
		new := second - top
		e.stack.Push(new)
	case "tr":
		// divide top two values
		if e.stack.Size() < 2 {
			return fmt.Errorf("tr: stack underflow")
		}
		top, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("tr: %w", err)
		}
		second, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("tr: %w", err)
		}
		if top == 0 {
			return fmt.Errorf("tr: division by zero")
		}
		new := second / top
		e.stack.Push(new)
	case "cl":
		// modulo top two values
		if e.stack.Size() < 2 {
			return fmt.Errorf("cl: stack underflow")
		}
		top, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("cl: %w", err)
		}
		second, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("cl: %w", err)
		}
		if top == 0 {
			return fmt.Errorf("cl: division by zero")
		}
		new := second % top
		e.stack.Push(new)
	case "slst":
		if e.stack.IsEmpty() {
			return fmt.Errorf("slst: stack underflow")
		}

		top, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("slst: %w", err)
		}
		e.stack.Push(top) // push back the top element
		e.stack.Push(top) // duplicate it
	case "swap":
		if e.stack.Size() < 2 {
			return fmt.Errorf("swap: stack underflow")
		}
		a, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("swap: %w", err)
		}
		b, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("swap: %w", err)
		}
		e.stack.Push(a)
		e.stack.Push(b)
	case "dup":
		top, _ := e.stack.Peek()
		e.stack.Push(top) // duplicate top element
	case "inc":
		if e.stack.IsEmpty() {
			return fmt.Errorf("inc: stack underflow")
		}
		top, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("inc: %w", err)
		}
		top++
		e.stack.Push(top) // increment top element
	case "dec":
		if e.stack.IsEmpty() {
			return fmt.Errorf("dec: stack underflow")
		}
		top, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("dec: %w", err)
		}
		top--
		e.stack.Push(top) // decrement top element
	default:
		return fmt.Errorf("unknown stitch %s", si.Token)
	}
	return nil
}

func (e *Evaluator) execRep(ri *parser.RepInstr) error {
	var count int
	if ri.CountExpr != "" {
		// Evaluate the count expression
		var err error
		count, err = strconv.Atoi(ri.CountExpr)
		if err != nil {
			return fmt.Errorf("rep: invalid count expression: %w", err)
		}
	} else {
		// Pop the count from the stack
		if e.stack.IsEmpty() {
			return fmt.Errorf("rep: stack underflow")
		}
		var err error
		count, err = e.stack.Pop()
		if err != nil {
			return fmt.Errorf("rep: %w", err)
		}
	}

	// Execute the body of the loop `count` times
	for i := 0; i < count; i++ {
		for _, instr := range ri.Body {
			if err := e.exec(instr); err != nil {
				return err
			}
		}
	}
	return nil
}
