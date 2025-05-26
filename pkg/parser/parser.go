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
	case lexer.PATTERN:
		return p.parsePattern()
	case lexer.USE:
		return p.parseUse()
	case lexer.REP:
		return p.parseRep()
	case lexer.CH:
		return p.parseCh()
	// Allow bare identifiers to be treated as pattern invocations
	case lexer.IDENT:
		return &UseInstr{Name: p.cur.Literal}, nil
	// All zero-arg stitches
	case lexer.SC, lexer.SLST, lexer.SWAP, lexer.DUP,
		lexer.INC, lexer.DEC,
		lexer.ADD, lexer.MUL, lexer.DIV, lexer.SUB, lexer.MOD,
		lexer.HDC, lexer.DC, lexer.TR, lexer.CL,
		lexer.YO, lexer.PIC, lexer.FO:
		return &SimpleInstr{Token: p.cur.Literal}, nil
	default:
		return nil, fmt.Errorf("unexpected token %q at line %d", p.cur.Literal, p.cur.Line)
	}
}

func (p *Parser) parseCh() (Instruction, error) {
	instr := &SimpleInstr{Token: p.cur.Literal}
	if p.peek.Type != lexer.INT {
		return nil, fmt.Errorf("expected INT after ch, got %s", p.peek.Literal)
	}
	p.nextToken() // consume INT
	instr.Args = append(instr.Args, p.cur.Literal)
	return instr, nil
}

func (p *Parser) parsePattern() (Instruction, error) {
	// pattern Name(opt params) = ( instr… )
	p.nextToken() // consume 'pattern'
	if p.cur.Type != lexer.IDENT {
		return nil, fmt.Errorf("expected pattern name, got %s", p.cur.Literal)
	}
	pat := &PatternDef{Name: p.cur.Literal}

	// optional params
	if p.peek.Type == lexer.LPAREN {
		p.nextToken() // consume '('
		p.nextToken()
		for p.cur.Type != lexer.RPAREN {
			if p.cur.Type != lexer.IDENT {
				return nil, fmt.Errorf("expected param name, got %s", p.cur.Literal)
			}
			pat.Params = append(pat.Params, p.cur.Literal)
			if p.peek.Type == lexer.COMMA {
				p.nextToken()
			}
			p.nextToken()
		}
	}

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
		if p.cur.Type == lexer.COMMA {
			p.nextToken()
		}
	}
	return pat, nil
}

func (p *Parser) parseUse() (Instruction, error) {
	// use Name(opt args)
	p.nextToken() // consume 'use'
	if p.cur.Type != lexer.IDENT {
		return nil, fmt.Errorf("expected pattern name, got %s", p.cur.Literal)
	}
	use := &UseInstr{Name: p.cur.Literal}

	if p.peek.Type == lexer.LPAREN {
		p.nextToken() // consume '('
		p.nextToken()
		for p.cur.Type != lexer.RPAREN {
			if p.cur.Type != lexer.INT {
				return nil, fmt.Errorf("expected INT arg, got %s", p.cur.Literal)
			}
			use.Args = append(use.Args, p.cur.Literal)
			if p.peek.Type == lexer.COMMA {
				p.nextToken()
			}
			p.nextToken()
		}
	}
	return use, nil
}

func (p *Parser) parseRep() (Instruction, error) {
	// assume rep [count]? ( … )
	ri := &RepInstr{}
	if p.peek.Type == lexer.INT {
		p.nextToken()
		ri.CountExpr = p.cur.Literal
	}
	// expect LPAREN
	p.nextToken() // should be "("
	p.nextToken()
	for p.cur.Type != lexer.RPAREN && p.cur.Type != lexer.EOF {
		instr, err := p.parseInstruction()
		if err != nil {
			return nil, err
		}
		ri.Body = append(ri.Body, instr)
		p.nextToken()
		if p.cur.Type == lexer.COMMA {
			p.nextToken()
		}
	}
	return ri, nil
}
