package lexer

import (
	"regexp"
	"strings"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

const (
	// Special
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Literals
	IDENT = "IDENT" // names, subpattern names, or unknown mnemonics
	INT   = "INT"

	// Delimiters
	LPAREN    = "("
	RPAREN    = ")"
	SEMICOLON = ";"
	ASTERISK  = "*"

	// Keywords / Stitch mnemonics
	CH          = "CH"
	SC          = "SC"
	DC          = "DC"
	HDC         = "HDC"
	TR          = "TR"
	CL          = "CL"
	INC         = "INC"
	DEC         = "DEC"
	SWAP        = "SWAP"
	SLST        = "SLST" // slip‐stitch; parser may want to combine "sl" "st"
	YO          = "YO"
	PIC         = "PIC"
	REP         = "REP"
	FO          = "FO"
	SUBPATTERN  = "SUBPATTERN"
	USE         = "USE"
	BOB         = "BOB"
	IF          = "IF"
	ELSE        = "ELSE"
	END         = "END"
	GREATERTHAN = ">"
	LESSERTHAN  = "<"
	EQUALS      = "EQ"
	NOTEQUALS   = "NEQ"
	TURN        = "TURN"
)

var keywords = map[string]TokenType{
	"ch":         CH,
	"sc":         SC,
	"dc":         DC,
	"hdc":        HDC,
	"tr":         TR,
	"cl":         CL,
	"inc":        INC,
	"dec":        DEC,
	"bob":        BOB,
	"swap":       SWAP,
	"sl st":      SLST,
	"yo":         YO,
	"pic":        PIC,
	"rep":        REP,
	"fo":         FO,
	"subpattern": SUBPATTERN,
	"use":        USE,
	"if":         IF,
	"else":       ELSE,
	"end":        END,
	">":          GREATERTHAN,
	"<":          LESSERTHAN,
	"eq":         EQUALS,
	"neq":        NOTEQUALS,
	"turn":       TURN,
}

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // next position in input (after current char)
	ch           byte // current char under examination
	Line         int
	Column       int
}

// New initializes a lexer for the given input.
func New(input string) *Lexer {
	// remove UTF-8 BOM if present
	if strings.HasPrefix(input, "\uFEFF") {
		input = strings.TrimPrefix(input, "\uFEFF")
	}
	// normalize non-breaking spaces -> regular spaces
	input = strings.ReplaceAll(input, "\u00A0", " ")

	// remove comments
	reComment := regexp.MustCompile(`#.*`)
	input = reComment.ReplaceAllString(input, "")

	// remove leading “Row N:” or “Round N:” labels (per-line)
	reLabel := regexp.MustCompile(`(?m)^\s*(?:Row|Round)\s+\d+:\s*`)
	input = reLabel.ReplaceAllString(input, "")

	l := &Lexer{input: input, Line: 1}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	if l.ch == '\n' {
		l.Line++
		l.Column = 0
	} else {
		l.Column++
	}
}

// peekChar lets us look ahead without consuming.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// returns the next token from the input.
func (l *Lexer) NextToken() Token {
	var tok Token

	// skip whitespace and comments if they haven't been caught by the preprocessor
	for {
		if l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
			l.readChar()
		} else if l.ch == '#' {
			// skip comment until end-of-line
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
		} else {
			break
		}
	}
	// ignore commas
	if l.ch == ',' {
		l.readChar()
	}

	tok.Line = l.Line
	tok.Column = l.Column

	if len(l.input)-l.position >= 5 &&
		strings.EqualFold(l.input[l.position:l.position+5], "sl st") {
		tok.Type = SLST
		tok.Literal = "slst"
		// consume all 5 chars: 's','l',' ','s','t'
		for i := 0; i < 5; i++ {
			l.readChar()
		}
		return tok
	}

	switch l.ch {
	case '(':
		tok = newToken(LPAREN, l.ch, tok.Line, tok.Column)
	case ')':
		tok = newToken(RPAREN, l.ch, tok.Line, tok.Column)
	case ';':
		tok = newToken(SEMICOLON, l.ch, tok.Line, tok.Column)
	case '*':
		tok = newToken(ASTERISK, l.ch, tok.Line, tok.Column)
	case '>':
		tok = newToken(GREATERTHAN, l.ch, tok.Line, tok.Column)
	case '<':
		tok = newToken(LESSERTHAN, l.ch, tok.Line, tok.Column)
	case 0:
		tok.Type = EOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			lit := l.readIdentifier()
			tok.Type = lookupIdent(lit)
			tok.Literal = lit
			return tok
		} else if isDigit(l.ch) {
			lit := l.readNumber()
			tok.Type = INT
			tok.Literal = lit
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch, tok.Line, tok.Column)
		}
	}

	l.readChar()
	return tok
}

func newToken(tt TokenType, ch byte, line, col int) Token {
	return Token{Type: tt, Literal: string(ch), Line: line, Column: col}
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func lookupIdent(ident string) TokenType {
	if tok, ok := keywords[strings.ToLower(ident)]; ok {
		return tok
	}
	return IDENT
}
