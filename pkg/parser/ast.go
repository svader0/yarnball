package parser

// Node is the base interface for all AST nodes.
type Node interface {
	TokenLiteral() string
}

// Program is the root node of every parsed file.
type Program struct {
	Instructions []Instruction
}

// Instruction represents one “stitch” or a rep‐block.
type Instruction interface {
	Node
	instructionNode()
}

// SimpleInstr is a plain token, like ch, pic, yo, FO…
type SimpleInstr struct {
	Token string // literal, e.g. "ch" or "pic"
	Args  []string
}

func (si *SimpleInstr) instructionNode()     {}
func (si *SimpleInstr) TokenLiteral() string { return si.Token }

// RepInstr represents: rep [<count>] (<instr…>)
type RepInstr struct {
	CountExpr string        // literal count or empty = dynamic
	Body      []Instruction // nested instructions
}

func (ri *RepInstr) instructionNode()     {}
func (ri *RepInstr) TokenLiteral() string { return "rep" }

type PatternDef struct {
	Name   string
	Params []string
	Body   []Instruction
}

func (*PatternDef) instructionNode()        {}
func (pd *PatternDef) TokenLiteral() string { return "pattern" }

type UseInstr struct {
	Name string
	Args []string
}

func (*UseInstr) instructionNode()        {}
func (ui *UseInstr) TokenLiteral() string { return "use" }
