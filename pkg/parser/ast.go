package parser

import (
	"strings"

	"github.com/svader0/yarnball/pkg/lexer"
)

/*
	The contents of this file define the abstract syntax tree (AST) we are going
	to use for Yarnball. Each node represents a specific part of the Yarnball
	program structure, like instructions, subpatterns, and their definitions.
	The AST is used to represent the parsed program in a structure that makes it
	easy to analyze and evaluate.
*/

// base interface for all AST nodes.
type Node interface {
	TokenLiteral() string
}

// root node of every parsed file (a program is just a sequence of instructions).
type Program struct {
	Instructions []Instruction
}

// represents one “stitch” or a rep‐block.
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

// rep [<count>] (<instr…>)
type RepInstr struct {
	CountExpr string        // literal count or empty = dynamic
	Body      []Instruction // nested instructions
}

func (ri *RepInstr) instructionNode()     {}
func (ri *RepInstr) TokenLiteral() string { return strings.ToLower(lexer.REP) }

// Represents a subpattern definition, which stores the name of the subpattern,
// its parameters (it just pushes that number), and the body of instructions that make up the subpattern.
type SubpatternDef struct {
	Name   string
	Params []string
	Body   []Instruction
}

func (*SubpatternDef) instructionNode()        {}
func (pd *SubpatternDef) TokenLiteral() string { return strings.ToLower(lexer.SUBPATTERN) }

type UseInstr struct {
	Name string
	Args []string
}

func (*UseInstr) instructionNode()        {}
func (ui *UseInstr) TokenLiteral() string { return strings.ToLower(lexer.USE) }
