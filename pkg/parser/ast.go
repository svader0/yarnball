package parser

/*
	The contents of this file define the abstract syntax tree (AST) we are going
	to use for Yarnball. Each node represents a specific part of the Yarnball
	program structure, like instructions, stitch definitions, and blocks.
*/

// base interface for all AST nodes.
type Node interface {
	TokenLiteral() string
}

// root node of every parsed file (a program is just a sequence of instructions).
type Program struct {
	Instructions []Instruction
}

// represents one “stitch” or a repeat block.
type Instruction interface {
	Node
	instructionNode()
}

// Represents a plain token, like ch, pic, yo, FO…
type SimpleInstr struct {
	Token string // literal, e.g. "ch" or "pic"
	Args  []string
}

func (si *SimpleInstr) instructionNode()     {}
func (si *SimpleInstr) TokenLiteral() string { return si.Token }

type RepeatMode int

const (
	RepeatCount RepeatMode = iota
	RepeatUntil
	RepeatWhile
)

// RepeatInstr represents a repeat block.
type RepeatInstr struct {
	Mode  RepeatMode
	Count int
	Body  []Instruction
}

func (ri *RepeatInstr) instructionNode()     {}
func (ri *RepeatInstr) TokenLiteral() string { return "repeat" }

// StitchDef defines a reusable stitch pattern.
type StitchDef struct {
	Name string
	Body []Instruction
}

func (*StitchDef) instructionNode()        {}
func (sd *StitchDef) TokenLiteral() string { return "stitch" }

type CallInstr struct {
	Name string
}

func (*CallInstr) instructionNode()        {}
func (ci *CallInstr) TokenLiteral() string { return ci.Name }

type IfInstr struct {
	IfBody   []Instruction // instructions to execute if condition is true
	ElseBody []Instruction // instructions to execute if condition is false (if any)
}

func (*IfInstr) instructionNode()     {}
func (*IfInstr) TokenLiteral() string { return "if" }
