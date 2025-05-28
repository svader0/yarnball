package parser

import (
	"fmt"

	"github.com/svader0/yarnball/pkg/lexer"
)

type Parser struct {
	l         *lexer.Lexer
	cur, peek lexer.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.cur = p.peek
	p.peek = p.l.NextToken()
}

func (p *Parser) ParseProgram() (*Program, error) {
	prog := &Program{}
	for p.cur.Type != lexer.EOF {
		instr, err := p.parseInstruction()
		if err != nil {
			return nil, err
		}
		prog.Instructions = append(prog.Instructions, instr)

	}
	return prog, nil
}

func (p *Parser) parseInstruction() (Instruction, error) {
	switch p.cur.Type {
	case lexer.SUBPATTERN:
		return p.parseSubpattern()
	case lexer.USE:
		return p.parseUse()
	case lexer.CH:
		return p.parseCh()
	case lexer.ASTERISK:
		return p.parseRep()
	case lexer.IDENT:
		return &UseInstr{Name: p.cur.Literal}, nil
	case lexer.IF:
		return p.parseIf()
	case lexer.SC, lexer.SLST, lexer.SWAP,
		lexer.INC, lexer.DEC, lexer.BOB,
		lexer.HDC, lexer.DC, lexer.TR, lexer.CL,
		lexer.GREATERTHAN, lexer.LESSERTHAN, lexer.TURN,
		lexer.EQUALS, lexer.NOTEQUALS,
		lexer.YO, lexer.PIC, lexer.FO:
		instr := &SimpleInstr{Token: p.cur.Literal}
		p.nextToken() // consume the instruction token
		return instr, nil
	default:
		return nil, fmt.Errorf("unexpected token %q at line %d", p.cur.Literal, p.cur.Line)
	}
}

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

func (p *Parser) parseSubpattern() (Instruction, error) {
	p.nextToken() // consume 'subpattern' keyword
	if p.cur.Type != lexer.IDENT {
		return nil, fmt.Errorf("expected subpattern name, got %s", p.cur.Literal)
	}
	pat := &SubpatternDef{Name: p.cur.Literal}

	// Expect '='
	if p.peek.Literal != "=" {
		return nil, fmt.Errorf("expected '=', got %s", p.peek.Literal)
	}
	p.nextToken() // consume name
	p.nextToken() // consume '='

	// Now p.cur should be the '(' token.
	if p.cur.Type != lexer.LPAREN {
		return nil, fmt.Errorf("expected '(', got %s", p.cur.Literal)
	}
	p.nextToken() // consume '(' and move to the first token in the subpattern body

	// Parse instructions until closing ')'
	for p.cur.Type != lexer.RPAREN && p.cur.Type != lexer.EOF {
		instr, err := p.parseInstruction()
		if err != nil {
			return nil, err
		}
		pat.Body = append(pat.Body, instr)
	}
	p.nextToken() // consume ')'
	return pat, nil
}

func (p *Parser) parseUse() (Instruction, error) {
	p.nextToken() // consume 'use'
	if p.cur.Type != lexer.IDENT {
		return nil, fmt.Errorf("expected subpattern name, got %s", p.cur.Literal)
	}
	use := &UseInstr{Name: p.cur.Literal}
	p.nextToken() // advance past the IDENT token

	// If there is a parameter list, consume it.
	if p.cur.Type == lexer.LPAREN {
		p.nextToken() // consume '('
		for p.cur.Type != lexer.RPAREN && p.cur.Type != lexer.EOF {
			if p.cur.Type != lexer.INT && p.cur.Type != lexer.IDENT {
				return nil, fmt.Errorf("expected INT or IDENT arg, got %s", p.cur.Literal)
			}
			use.Args = append(use.Args, p.cur.Literal)
			p.nextToken()
		}
		if p.cur.Type != lexer.RPAREN {
			return nil, fmt.Errorf("expected RPAREN, got %s", p.cur.Literal)
		}
		p.nextToken() // consume ')'
	}
	return use, nil
}

// parseRep implements "*[<instrs>]; rep from * <count> [times]"
func (p *Parser) parseRep() (Instruction, error) {
	ri := &RepInstr{}

	// Consume the '*' that starts the rep block.
	p.nextToken()

	// Parse the rep block body until we hit a semicolon.
	for p.cur.Type != lexer.SEMICOLON && p.cur.Type != lexer.EOF {
		instr, err := p.parseInstruction()
		if err != nil {
			return nil, err
		}
		ri.Body = append(ri.Body, instr)
		// Removed extra p.nextToken() here.
	}

	if p.cur.Type != lexer.SEMICOLON {
		return nil, fmt.Errorf("expected ';' after rep block, got %q at line %d", p.cur.Literal, p.cur.Line)
	}
	// Consume the semicolon.
	p.nextToken()

	// Next, expect the 'rep' keyword.
	if p.cur.Type != lexer.REP {
		return nil, fmt.Errorf("expected 'rep' after ';', got %q at line %d", p.cur.Literal, p.cur.Line)
	}
	p.nextToken() // Consume 'rep'

	// Expect the literal "from"
	if p.cur.Literal != "from" {
		return nil, fmt.Errorf("expected 'from' after 'rep', got %q at line %d", p.cur.Literal, p.cur.Line)
	}
	p.nextToken() // Consume "from"

	// Expect an asterisk '*' for the count.
	if p.cur.Type != lexer.ASTERISK {
		return nil, fmt.Errorf("expected '*' after 'from', got %q at line %d", p.cur.Literal, p.cur.Line)
	}
	p.nextToken() // Consume '*'

	// Optionally, consume INT count.
	if p.cur.Type == lexer.INT {
		ri.CountExpr = p.cur.Literal
		p.nextToken()
	}

	// Optionally, consume a 'times' keyword.
	if p.cur.Type == lexer.IDENT && p.cur.Literal == "times" {
		p.nextToken()
	}

	return ri, nil
}

func (p *Parser) parseIf() (Instruction, error) {
	// Consume the 'if' token
	p.nextToken()

	var ifBody []Instruction
	// Parse IF branch until we hit ELSE or END
	for p.cur.Type != lexer.ELSE && p.cur.Type != lexer.END && p.cur.Type != lexer.EOF {
		instr, err := p.parseInstruction()
		if err != nil {
			return nil, err
		}
		ifBody = append(ifBody, instr)
		// Do not call p.nextToken() here; let each parse* function consume its tokens.
	}

	var elseBody []Instruction
	// If an ELSE is encountered, parse ELSE branch
	if p.cur.Type == lexer.ELSE {
		p.nextToken() // consume 'else'
		for p.cur.Type != lexer.END && p.cur.Type != lexer.EOF {
			instr, err := p.parseInstruction()
			if err != nil {
				return nil, err
			}
			elseBody = append(elseBody, instr)
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
