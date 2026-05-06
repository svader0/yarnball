package parser

import (
	"fmt"
	"strconv"

	"github.com/svader0/yarnball/pkg/lexer"
)

type Parser struct {
	l         *lexer.Lexer
	cur, peek lexer.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken() // Initialize current token
	p.nextToken() // Initialize peek token
	return p
}

// nextToken advances the parser to the next token, updating current and peek tokens.
func (p *Parser) nextToken() {
	p.cur = p.peek
	p.peek = p.l.NextToken()
}

// Parses the entire program, which consists of a sequence of instructions.
func (p *Parser) ParseProgram() (*Program, error) {
	prog := &Program{}
	for p.cur.Type != lexer.EOF {
		p.skipFillers()
		if p.cur.Type == lexer.EOF {
			break
		}
		instr, err := p.parseInstruction()
		if err != nil {
			return nil, err
		}
		if instr != nil {
			prog.Instructions = append(prog.Instructions, instr)
		}

	}
	return prog, nil
}

// Parses an instruction based on the current token type.
func (p *Parser) parseInstruction() (Instruction, error) {
	switch p.cur.Type {
	case lexer.STITCHDEF:
		return p.parseStitchDef()
	case lexer.USE:
		return p.parseUse()
	case lexer.CH:
		return p.parseCh()
	case lexer.PICK, lexer.ROLL:
		return p.parsePickRoll()
	case lexer.ASTERISK, lexer.LBRACKET:
		return p.parseRepeatBlock()
	case lexer.INT:
		return p.parsePrefixedCount()
	case lexer.IDENT:
		return p.parseCall()
	case lexer.IF:
		return p.parseIf()
	case lexer.SC, lexer.SLST, lexer.SWAP,
		lexer.INC, lexer.DEC, lexer.BOB,
		lexer.HDC, lexer.DC, lexer.TR, lexer.CL,
		lexer.GREATERTHAN, lexer.LESSERTHAN, lexer.TURN,
		lexer.EQ, lexer.NEQ,
		lexer.OVER, lexer.YO, lexer.PIC, lexer.FO:
		return p.parseSimpleWithOptionalCount()
	case lexer.FILLER:
		p.nextToken()
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected token %q at line %d", p.cur.Literal, p.cur.Line)
	}
}

// parseCh parses the 'ch' instruction, which expects an INT argument.
func (p *Parser) parseCh() (Instruction, error) {
	instr := &SimpleInstr{Token: p.cur.Literal}
	p.nextToken() // consume 'ch'
	// If next token is INT, just store it
	if p.cur.Type == lexer.INT {
		instr.Args = append(instr.Args, p.cur.Literal)
		p.nextToken() // consume INT
	} else {
		return nil, fmt.Errorf("expected INT after ch, got %s", p.cur.Literal)
	}
	return instr, nil
}

func (p *Parser) parsePickRoll() (Instruction, error) {
	instr := &SimpleInstr{Token: p.cur.Literal}
	p.nextToken() // consume 'pick' or 'roll'
	if p.cur.Type == lexer.INT {
		instr.Args = append(instr.Args, p.cur.Literal)
		p.nextToken() // consume INT
	} else {
		return nil, fmt.Errorf("expected INT after %s, got %s", instr.Token, p.cur.Literal)
	}
	return instr, nil
}

func (p *Parser) parseStitchDef() (Instruction, error) {
	p.nextToken() // consume 'stitch' keyword
	if p.cur.Type != lexer.IDENT {
		return nil, fmt.Errorf("expected stitch name, got %s", p.cur.Literal)
	}
	def := &StitchDef{Name: p.cur.Literal}

	// Expect '='
	if p.peek.Type != lexer.ASSIGN {
		return nil, fmt.Errorf("expected '=', got %s", p.peek.Literal)
	}
	p.nextToken() // consume name
	p.nextToken() // consume '='

	if p.cur.Type != lexer.LPAREN {
		return nil, fmt.Errorf("expected '(', got %s", p.cur.Literal)
	}
	p.nextToken() // consume '('

	// Parse instructions until closing ')'
	for p.cur.Type != lexer.RPAREN && p.cur.Type != lexer.EOF {
		p.skipFillers()
		if p.cur.Type == lexer.RPAREN || p.cur.Type == lexer.EOF {
			break
		}
		instr, err := p.parseInstruction()
		if err != nil {
			return nil, err
		}
		if instr != nil {
			def.Body = append(def.Body, instr)
		}
	}
	p.nextToken() // consume ')'
	return def, nil
}

func (p *Parser) parseUse() (Instruction, error) {
	p.nextToken() // consume 'use'
	if p.cur.Type != lexer.IDENT {
		return nil, fmt.Errorf("expected stitch name, got %s", p.cur.Literal)
	}
	call := &CallInstr{Name: p.cur.Literal}
	p.nextToken() // advance past the IDENT token
	return p.wrapPostfixCount(call)
}

func (p *Parser) parseCall() (Instruction, error) {
	call := &CallInstr{Name: p.cur.Literal}
	p.nextToken()
	return p.wrapPostfixCount(call)
}

func (p *Parser) parseIf() (Instruction, error) {
	// Consume the 'if' token
	p.nextToken()

	var ifBody []Instruction
	// Parse IF branch until we hit ELSE or END
	for p.cur.Type != lexer.ELSE && p.cur.Type != lexer.END && p.cur.Type != lexer.EOF {
		p.skipFillers()
		if p.cur.Type == lexer.ELSE || p.cur.Type == lexer.END || p.cur.Type == lexer.EOF {
			break
		}
		instr, err := p.parseInstruction()
		if err != nil {
			return nil, err
		}
		if instr != nil {
			ifBody = append(ifBody, instr)
		}
	}

	var elseBody []Instruction
	// If an ELSE is encountered, parse ELSE branch
	if p.cur.Type == lexer.ELSE {
		p.nextToken() // consume 'else'
		for p.cur.Type != lexer.END && p.cur.Type != lexer.EOF {
			p.skipFillers()
			if p.cur.Type == lexer.END || p.cur.Type == lexer.EOF {
				break
			}
			instr, err := p.parseInstruction()
			if err != nil {
				return nil, err
			}
			if instr != nil {
				elseBody = append(elseBody, instr)
			}
		}
	}

	// Make sure we have an END token
	if p.cur.Type != lexer.END {
		return nil, fmt.Errorf("expected 'end' token, got %q at line %d", p.cur.Literal, p.cur.Line)
	}
	// Consume the END token
	p.nextToken()

	return &IfInstr{IfBody: ifBody, ElseBody: elseBody}, nil
}

// parseRepeatBlock handles both * ... * and [ ... ] repeat blocks.
func (p *Parser) parseRepeatBlock() (Instruction, error) {
	startToken := p.cur.Type
	endToken := startToken
	if startToken == lexer.LBRACKET {
		endToken = lexer.RBRACKET
	}

	ri := &RepeatInstr{}
	p.nextToken() // consume '*' or '['

	for p.cur.Type != endToken && p.cur.Type != lexer.EOF {
		p.skipFillers()
		if p.cur.Type == endToken || p.cur.Type == lexer.EOF {
			break
		}
		instr, err := p.parseInstruction()
		if err != nil {
			return nil, err
		}
		if instr != nil {
			ri.Body = append(ri.Body, instr)
		}
	}

	if p.cur.Type != endToken {
		return nil, fmt.Errorf("expected closing %q, got %q at line %d", endToken, p.cur.Literal, p.cur.Line)
	}
	p.nextToken() // consume closing

	p.skipFillers()
	if p.cur.Type != lexer.REPEAT {
		return nil, fmt.Errorf("expected 'repeat' after block, got %q at line %d", p.cur.Literal, p.cur.Line)
	}
	p.nextToken() // consume 'repeat'

	p.skipFillers()
	if p.cur.Type == lexer.ASTERISK || p.cur.Type == lexer.LBRACKET {
		p.nextToken() // allow "repeat from *"
		p.skipFillers()
	}

	switch p.cur.Type {
	case lexer.INT:
		count, err := strconv.Atoi(p.cur.Literal)
		if err != nil || count < 0 {
			return nil, fmt.Errorf("invalid repeat count %q", p.cur.Literal)
		}
		ri.Mode = RepeatCount
		ri.Count = count
		p.nextToken()
		p.skipFillers()
	case lexer.UNTIL:
		ri.Mode = RepeatUntil
		p.nextToken()
	case lexer.WHILE:
		ri.Mode = RepeatWhile
		p.nextToken()
	default:
		return nil, fmt.Errorf("expected repeat count, 'until', or 'while', got %q at line %d", p.cur.Literal, p.cur.Line)
	}

	return ri, nil
}

func (p *Parser) parsePrefixedCount() (Instruction, error) {
	count, err := strconv.Atoi(p.cur.Literal)
	if err != nil || count < 0 {
		return nil, fmt.Errorf("invalid count %q", p.cur.Literal)
	}
	p.nextToken()
	p.skipFillers()

	instr, err := p.parseInstruction()
	if err != nil {
		return nil, err
	}
	if instr == nil || !countableInstr(instr) {
		return nil, fmt.Errorf("count prefix must apply to a stitch or stitch call")
	}
	return &RepeatInstr{Mode: RepeatCount, Count: count, Body: []Instruction{instr}}, nil
}

func (p *Parser) parseSimpleWithOptionalCount() (Instruction, error) {
	instr := &SimpleInstr{Token: p.cur.Literal}
	p.nextToken() // consume instruction token
	return p.wrapPostfixCount(instr)
}

func (p *Parser) wrapPostfixCount(instr Instruction) (Instruction, error) {
	if p.cur.Type == lexer.INT && countableInstr(instr) {
		count, err := strconv.Atoi(p.cur.Literal)
		if err != nil || count < 0 {
			return nil, fmt.Errorf("invalid count %q", p.cur.Literal)
		}
		p.nextToken()
		return &RepeatInstr{Mode: RepeatCount, Count: count, Body: []Instruction{instr}}, nil
	}
	return instr, nil
}

func (p *Parser) skipFillers() {
	for p.cur.Type == lexer.FILLER {
		p.nextToken()
	}
}

func countableInstr(instr Instruction) bool {
	switch node := instr.(type) {
	case *SimpleInstr:
		return isCountableOp(node.Token)
	case *CallInstr:
		return true
	default:
		return false
	}
}

func isCountableOp(op string) bool {
	switch op {
	case "ch", "pick", "roll":
		return false
	default:
		return true
	}
}
