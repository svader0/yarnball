package evaluator

import (
	"fmt"
	"strconv"

	"github.com/svader0/yarnball/pkg/parser"
	"github.com/svader0/yarnball/pkg/stack"
)

type Evaluator struct {
	stack    *stack.Stack
	patterns map[string]*parser.SubpatternDef
}

func New() *Evaluator {
	return &Evaluator{
		stack:    stack.New(),
		patterns: make(map[string]*parser.SubpatternDef),
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
		case *parser.SubpatternDef:
			e.patterns[node.Name] = node
		case *parser.UseInstr:
			if err := e.execUse(node); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Evaluator) exec(instr parser.Instruction) error {
	fmt.Println("Executing instruction:", instr.TokenLiteral())
	fmt.Println("Current stack:", e.stack)
	fmt.Println("Patterns:", e.patterns)
	fmt.Println("Instruction type:", fmt.Sprintf("%T", instr))
	switch node := instr.(type) {
	case *parser.SimpleInstr:
		return e.execSimple(node)
	case *parser.RepInstr:
		return e.execRep(node)
	case *parser.UseInstr:
		return e.execUse(node)
	case *parser.SubpatternDef:
		e.patterns[node.Name] = node
		return nil
	default:
		return fmt.Errorf("unknown instruction type: %T", instr)
	}
}

func (e *Evaluator) execUse(ui *parser.UseInstr) error {
	pat, exists := e.patterns[ui.Name]
	if !exists {
		return fmt.Errorf("undefined subpattern %q", ui.Name)
	}

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
		n, err := strconv.Atoi(si.Args[0])
		if err != nil {
			return fmt.Errorf("ch: invalid argument %q: %w", si.Args[0], err)
		}
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
	case "fo":
		// terminate program
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

		top, _ := e.stack.Peek()
		e.stack.Push(top)
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
	var err error

	if ri.CountExpr != "" {
		count, err = strconv.Atoi(ri.CountExpr)
		if err != nil {
			return fmt.Errorf("rep: invalid count %q: %w", ri.CountExpr, err)
		}
	} else {
		if e.stack.IsEmpty() {
			return fmt.Errorf("rep: stack underflow")
		}
		count, err = e.stack.Pop()
		if err != nil {
			return fmt.Errorf("rep: %w", err)
		}
	}

	for i := 0; i < count; i++ {
		for _, instr := range ri.Body {
			if err := e.exec(instr); err != nil {
				return err
			}
		}
	}
	return nil
}
