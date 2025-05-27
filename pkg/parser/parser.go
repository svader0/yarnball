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
		p.nextToken()
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
		// All zero-arg stitches
	case lexer.SC, lexer.SLST, lexer.SWAP,
		lexer.INC, lexer.DEC, lexer.BOB,
		lexer.HDC, lexer.DC, lexer.TR, lexer.CL,
		lexer.YO, lexer.PIC, lexer.FO:
		return &SimpleInstr{Token: p.cur.Literal}, nil
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
	} else {
		return nil, fmt.Errorf("expected INT after ch, got %s", p.cur.Literal)
	}
	return instr, nil
}

func (p *Parser) parseSubpattern() (Instruction, error) {
	// subpattern Name(opt params) = ( instr… )
	p.nextToken() // consume 'subpattern' token
	if p.cur.Type != lexer.IDENT {
		return nil, fmt.Errorf("expected subpattern name, got %s", p.cur.Literal)
	}
	pat := &SubpatternDef{Name: p.cur.Literal}

	if p.peek.Literal != "=" {
		return nil, fmt.Errorf("expected '=', got %s", p.peek.Literal)
	}
	p.nextToken() // consume '='

	// body
	if p.peek.Type != lexer.LPAREN {
		return nil, fmt.Errorf("expected '(', got %s", p.peek.Literal)
	}
	p.nextToken() // consume '('
	p.nextToken()
	for p.cur.Type != lexer.RPAREN && p.cur.Type != lexer.EOF {
		instr, err := p.parseInstruction()
		if err != nil {
			return nil, err
		}
		pat.Body = append(pat.Body, instr)
		p.nextToken()
	}
	return pat, nil
}

func (p *Parser) parseUse() (Instruction, error) {
	p.nextToken() // consume 'use'
	if p.cur.Type != lexer.IDENT {
		return nil, fmt.Errorf("expected subpattern name, got %s", p.cur.Literal)
	}
	use := &UseInstr{Name: p.cur.Literal}

	if p.peek.Type == lexer.LPAREN {
		p.nextToken() // consume '('
		p.nextToken()
		for p.cur.Type != lexer.RPAREN {
			if p.cur.Type != lexer.INT && p.cur.Type != lexer.IDENT {
				return nil, fmt.Errorf("expected INT or IDENT arg, got %s", p.cur.Literal)
			}
			use.Args = append(use.Args, p.cur.Literal)
			p.nextToken()
		}
	}
	return use, nil
}

// parseRep implements "*[<instrs>]; rep from * <count> [times]"
func (p *Parser) parseRep() (Instruction, error) {
	ri := &RepInstr{}

	// consume '*'
	p.nextToken()

	// parse body until ';'
	for p.cur.Type != lexer.SEMICOLON && p.cur.Type != lexer.EOF {
		instr, err := p.parseInstruction()
		if err != nil {
			return nil, err
		}
		ri.Body = append(ri.Body, instr)
		p.nextToken()
	}
	if p.cur.Type != lexer.SEMICOLON {
		return nil, fmt.Errorf("expected ';' after crochet block, got %s", p.cur.Literal)
	}
	// consume ';'
	p.nextToken()

	// expect 'rep'
	if p.cur.Type != lexer.REP {
		return nil, fmt.Errorf("expected 'rep' after ';', got %s", p.cur.Literal)
	}
	p.nextToken()

	// expect 'from'
	if p.cur.Literal != "from" {
		return nil, fmt.Errorf("expected 'from' after 'rep', got %s", p.cur.Literal)
	}
	p.nextToken()

	// expect '*'
	if p.cur.Type != lexer.ASTERISK {
		return nil, fmt.Errorf("expected '*' after 'from', got %s", p.cur.Literal)
	}
	p.nextToken()

	// optional INT count
	if p.cur.Type == lexer.INT {
		ri.CountExpr = p.cur.Literal
		p.nextToken()
	}
	// optional 'times'
	if p.cur.Type == lexer.IDENT && p.cur.Literal == "times" {
		p.nextToken() // consume "times" so it doesn’t become a UseInstr
	}

	return ri, nil
}
