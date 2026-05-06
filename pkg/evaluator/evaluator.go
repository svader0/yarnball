package evaluator

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/svader0/yarnball/pkg/parser"
	"github.com/svader0/yarnball/pkg/stack"
)

type Evaluator struct {
	log       *slog.Logger
	stack     *stack.Stack
	patterns  map[string]*parser.StitchDef
	stepLimit int
	steps     int
}

func New(logger *slog.Logger) *Evaluator {
	return &Evaluator{
		log:       logger,
		stack:     stack.New(),
		patterns:  make(map[string]*parser.StitchDef),
		stepLimit: 1_000_000,
	}
}

func (e *Evaluator) SetStepLimit(limit int) {
	if limit > 0 {
		e.stepLimit = limit
	}
}

// Passes the stack to the evaluator, allowing access to it from outside
// e.g. for debugging or inspection.
func (e *Evaluator) Stack() *stack.Stack {
	return e.stack
}

func (e *Evaluator) Eval(prog *parser.Program) error {
	e.log.Debug("Starting evaluation of program", "instructions", len(prog.Instructions))
	e.steps = 0
	for _, instr := range prog.Instructions {
		e.log.Debug("Evaluating instruction", "instruction", instr.TokenLiteral())
		// Execute the instruction based on its type
		if err := e.exec(instr); err != nil {
			e.log.Error("Error executing instruction", "instruction", instr.TokenLiteral(), "error", err)
			return fmt.Errorf("error executing instruction %s: %w", instr.TokenLiteral(), err)
		}
		e.log.Debug("Instruction executed successfully", "instruction", instr.TokenLiteral())
	}
	return nil
}

func (e *Evaluator) exec(instr parser.Instruction) error {
	if err := e.checkStep(); err != nil {
		return err
	}
	switch node := instr.(type) {
	case *parser.SimpleInstr:
		return e.execSimple(node)
	case *parser.RepeatInstr:
		return e.execRepeat(node)
	case *parser.CallInstr:
		return e.execCall(node)
	case *parser.StitchDef:
		e.patterns[node.Name] = node
		return nil
	case *parser.IfInstr:
		return e.execIf(node)
	default:
		return fmt.Errorf("unknown instruction type: %T", instr)
	}
}

func (e *Evaluator) execCall(ci *parser.CallInstr) error {
	e.log.Debug("Using stitch", "name", ci.Name)
	pat, exists := e.patterns[ci.Name]
	if !exists {
		return fmt.Errorf("undefined stitch %q", ci.Name)
	}

	for _, instr := range pat.Body {
		if err := e.exec(instr); err != nil {
			return fmt.Errorf("error executing stitch %s: %w", ci.Name, err)
		}
	}
	return nil
}

func (e *Evaluator) execIf(ii *parser.IfInstr) error {
	cond, err := e.stack.Pop()
	if err != nil {
		return fmt.Errorf("if: stack underflow")
	}
	if cond != 0 {
		for _, instr := range ii.IfBody {
			if err := e.exec(instr); err != nil {
				return fmt.Errorf("error executing if body: %w", err)
			}
		}
	} else {
		for _, instr := range ii.ElseBody {
			if err := e.exec(instr); err != nil {
				return fmt.Errorf("error executing else body: %w", err)
			}
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
		os.Exit(0)
	case "sc":
		// pop top value
		if e.stack.IsEmpty() {
			return fmt.Errorf("sc: stack underflow")
		}
		_, _ = e.stack.Pop()
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
		new := top * second
		e.stack.Push(new) // push product
	case "bob":
		// add top two values
		if e.stack.Size() < 2 {
			return fmt.Errorf("bob: stack underflow")
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
	case ">":
		if e.stack.Size() < 2 {
			return fmt.Errorf("gt: stack underflow")
		}
		top, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("gt: %w", err)
		}
		second, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("gt: %w", err)
		}
		if second > top {
			e.stack.Push(1) // push true
		} else {
			e.stack.Push(0) // push false
		}
	case "<":
		if e.stack.Size() < 2 {
			return fmt.Errorf("lt: stack underflow")
		}
		top, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("lt: %w", err)
		}
		second, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("lt: %w", err)
		}
		if second < top {
			e.stack.Push(1) // push true
		} else {
			e.stack.Push(0) // push false
		}
	case "eq":
		if e.stack.Size() < 2 {
			return fmt.Errorf("eq: stack underflow")
		}
		top, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("eq: %w", err)
		}
		second, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("eq: %w", err)
		}
		if second == top {
			e.stack.Push(1) // push true
		} else {
			e.stack.Push(0) // push false
		}
	case "neq":
		if e.stack.Size() < 2 {
			return fmt.Errorf("neq: stack underflow")
		}
		top, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("neq: %w", err)
		}
		second, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("neq: %w", err)
		}
		if second != top {
			e.stack.Push(1) // push true
		} else {
			e.stack.Push(0) // push false
		}
	case "turn":
		// Same function as 'rot' in FORTH ( n1 n2 n3 — n2 n3 n1 )
		if e.stack.Size() < 3 {
			return fmt.Errorf("turn: stack underflow")
		}
		top, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("turn: %w", err)
		}
		second, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("turn: %w", err)
		}
		third, err := e.stack.Pop()
		if err != nil {
			return fmt.Errorf("turn: %w", err)
		}
		e.stack.Push(second)
		e.stack.Push(top)
		e.stack.Push(third)
	case "over":
		if e.stack.Size() < 2 {
			return fmt.Errorf("over: stack underflow")
		}
		val, err := e.stack.PeekAt(1)
		if err != nil {
			return fmt.Errorf("over: %w", err)
		}
		e.stack.Push(val)
	case "pick":
		if len(si.Args) != 1 {
			return fmt.Errorf("pick: missing depth argument")
		}
		depth, err := strconv.Atoi(si.Args[0])
		if err != nil || depth < 0 {
			return fmt.Errorf("pick: invalid depth %q", si.Args[0])
		}
		val, err := e.stack.PeekAt(depth)
		if err != nil {
			return fmt.Errorf("pick: %w", err)
		}
		e.stack.Push(val)
	case "roll":
		if len(si.Args) != 1 {
			return fmt.Errorf("roll: missing depth argument")
		}
		depth, err := strconv.Atoi(si.Args[0])
		if err != nil || depth < 0 {
			return fmt.Errorf("roll: invalid depth %q", si.Args[0])
		}
		if err := e.stack.Roll(depth); err != nil {
			return fmt.Errorf("roll: %w", err)
		}
	default:
		return fmt.Errorf("unknown stitch %s", si.Token)
	}
	return nil
}

func (e *Evaluator) execRepeat(ri *parser.RepeatInstr) error {
	switch ri.Mode {
	case parser.RepeatCount:
		for i := 0; i < ri.Count; i++ {
			for _, instr := range ri.Body {
				if err := e.exec(instr); err != nil {
					return err
				}
			}
		}
	case parser.RepeatUntil:
		for {
			cond, ok := e.stack.Peek()
			if !ok {
				return fmt.Errorf("repeat until: stack underflow")
			}
			if cond != 0 {
				break
			}
			for _, instr := range ri.Body {
				if err := e.exec(instr); err != nil {
					return err
				}
			}
		}
	case parser.RepeatWhile:
		for {
			cond, ok := e.stack.Peek()
			if !ok {
				return fmt.Errorf("repeat while: stack underflow")
			}
			if cond == 0 {
				break
			}
			for _, instr := range ri.Body {
				if err := e.exec(instr); err != nil {
					return err
				}
			}
		}
	default:
		return fmt.Errorf("repeat: unknown mode")
	}
	return nil
}

func (e *Evaluator) checkStep() error {
	e.steps++
	if e.stepLimit > 0 && e.steps > e.stepLimit {
		return fmt.Errorf("step limit exceeded")
	}
	return nil
}
