package evaluator

import (
	"fmt"
	"strconv"

	"github.com/svader0/yarnball/pkg/parser"
)

type Evaluator struct {
	stack    []int
	patterns map[string]*parser.PatternDef
}

func New() *Evaluator {
	return &Evaluator{
		stack:    []int{},
		patterns: make(map[string]*parser.PatternDef),
	}
}

func (e *Evaluator) Stack() []int {
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

func (e *Evaluator) execUse(ui *parser.UseInstr) error {
	pat, exists := e.patterns[ui.Name]
	if !exists {
		return fmt.Errorf("undefined pattern %q", ui.Name)
	}

	// Push arguments onto the stack
	for _, arg := range ui.Args {
		n, _ := strconv.Atoi(arg)
		e.stack = append(e.stack, n)
	}

	// Execute the pattern body
	for _, instr := range pat.Body {
		switch node := instr.(type) {
		case *parser.SimpleInstr:
			if err := e.execSimple(node); err != nil {
				return err
			}
		case *parser.RepInstr:
			if err := e.execRep(node); err != nil {
				return err
			}
		case *parser.UseInstr:
			if err := e.execUse(node); err != nil {
				return err
			}
		case *parser.PatternDef:
			return fmt.Errorf("nested pattern definitions are not allowed")
		default:
			return fmt.Errorf("unknown instruction type in pattern %q", ui.Name)
		}
	}
	return nil
}

func (e *Evaluator) execSimple(si *parser.SimpleInstr) error {
	switch si.Token {
	case "ch":
		n, _ := strconv.Atoi(si.Args[0])
		e.stack = append(e.stack, n)
	case "pic":
		n := e.pop()
		fmt.Printf("%c", n)
	case "yo":
		n := e.pop()
		fmt.Println(n)
	case "FO":
		return fmt.Errorf("FO: halt")
	case "sc":
		// pop top value
		if len(e.stack) == 0 {
			return fmt.Errorf("sc: stack underflow")
		}
		e.pop()
	case "dc":
		// product of top two values
		if len(e.stack) < 2 {
			return fmt.Errorf("dc: stack underflow")
		}
		top := e.pop()
		second := e.pop()
		new := top * second
		e.stack = append(e.stack, new) // push product
	case "bob":
		// add top two values
		if len(e.stack) < 2 {
			return fmt.Errorf("dc: stack underflow")
		}
		top := e.pop()
		second := e.pop()
		new := top + second
		e.stack = append(e.stack, new)
	case "hdc":
		// subtract top two values
		if len(e.stack) < 2 {
			return fmt.Errorf("hdc: stack underflow")
		}
		top := e.pop()
		second := e.pop()
		new := second - top
		e.stack = append(e.stack, new)
	case "tr":
		// divide top two values
		if len(e.stack) < 2 {
			return fmt.Errorf("tr: stack underflow")
		}
		top := e.pop()
		second := e.pop()
		if top == 0 {
			return fmt.Errorf("tr: division by zero")
		}
		new := second / top
		e.stack = append(e.stack, new)
	case "cl":
		// modulo top two values
		if len(e.stack) < 2 {
			return fmt.Errorf("cl: stack underflow")
		}
		top := e.pop()
		second := e.pop()
		if top == 0 {
			return fmt.Errorf("cl: division by zero")
		}
		new := second % top
		e.stack = append(e.stack, new)
	case "slst":
		if len(e.stack) == 0 {
			return fmt.Errorf("slst: stack underflow")
		}

		top := e.pop()
		e.stack = append(e.stack, top) // push back the top element
		e.stack = append(e.stack, top) // duplicate it
	case "swap":
		if len(e.stack) < 2 {
			return fmt.Errorf("swap: stack underflow")
		}
		a := e.pop()
		b := e.pop()
		e.stack = append(e.stack, a, b) // push in reverse order
	case "dup":
		if len(e.stack) == 0 {
			return fmt.Errorf("dup: stack underflow")
		}
		top := e.stack[len(e.stack)-1]
		e.stack = append(e.stack, top) // duplicate top element
	case "inc":
		if len(e.stack) == 0 {
			return fmt.Errorf("inc: stack underflow")
		}
		top := e.pop()
		top++
		e.stack = append(e.stack, top) // increment top element
	case "dec":
		if len(e.stack) == 0 {
			return fmt.Errorf("dec: stack underflow")
		}
		top := e.pop()
		top--
		e.stack = append(e.stack, top) // decrement top element
	default:
		return fmt.Errorf("unknown stitch %s", si.Token)
	}
	return nil
}

func (e *Evaluator) execRep(ri *parser.RepInstr) error {
	count := 0
	if ri.CountExpr != "" {
		count, _ = strconv.Atoi(ri.CountExpr)
	} else {
		count = e.pop()
	}
	for i := 0; i < count; i++ {
		for _, instr := range ri.Body {
			if err := e.Eval(&parser.Program{Instructions: []parser.Instruction{instr}}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Evaluator) pop() int {
	if len(e.stack) == 0 {
		panic("stack underflow")
	}
	v := e.stack[len(e.stack)-1]
	e.stack = e.stack[:len(e.stack)-1]
	return v
}
