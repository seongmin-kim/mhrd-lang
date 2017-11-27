package lang

import (
	"fmt"
	"io"
	"strconv"
)

type Parser struct {
	s   *Scanner
	buf struct {
		tok     token
		literal string
		n       int
	}
}

func NewParser(rd io.Reader) *Parser {
	return &Parser{s: NewScanner(rd)}
}

func (p *Parser) scan() (token, string) {
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.literal
	}

	tok, literal := p.s.scan()

	p.buf.tok = tok
	p.buf.literal = literal

	return tok, literal
}

func (p *Parser) unscan() {
	p.buf.n = 1
}

func (p *Parser) scanIgnoreWhitespace() (token, string) {
	for {
		tok, literal := p.scan()

		if tok == COMMENT {
			continue
		}

		if tok != WHITESPACE && tok != NEWLINE {
			return tok, literal
		}
	}
}

func (p *Parser) parsePort() (*Port, error) {
	port := Port{Pins: 1}

	tok, literal := p.scanIgnoreWhitespace()
	port.Id = literal

	tok, literal = p.scanIgnoreWhitespace()

	if tok == LEFTBRACKET {
		tok, literal := p.scanIgnoreWhitespace()
		if tok != DIGIT {
			err := fmt.Errorf("Found '%s', expected digit", literal)
			return nil, err
		}
		val, _ := strconv.Atoi(literal)
		port.Pins = val

		tok, literal = p.scanIgnoreWhitespace()
		if tok != RIGHTBRACKET {
			err := fmt.Errorf("Found '%s', expected ']'", literal)
			return nil, err
		}
	} else if tok == COMMA || tok == SEMICOLON {
		p.unscan()
	} else {
		err := fmt.Errorf("Found '%s', expected ';'", literal)
		return nil, err
	}

	return &port, nil
}

func (p *Parser) ParseInputs() (*InputStatement, error) {

	tok, literal := p.scanIgnoreWhitespace()
	if tok != INPUTS {
		err := fmt.Errorf("Found '%s', expected 'Inputs'", literal)
		return nil, err
	}

	tok, literal = p.scanIgnoreWhitespace()
	if tok != COLON {
		err := fmt.Errorf("Found '%s', expected ':'", literal)
		return nil, err
	}

	statement := NewInputStatement()

	for {
		tok, _ := p.scanIgnoreWhitespace()

		if tok == SEMICOLON {
			return statement, nil
		}

		if tok == LETTER {
			p.unscan()
			port, err := p.parsePort()
			if err != nil {
				return nil, err
			}
			statement.add(*port)
		}
	}
}

func (p *Parser) ParseOutputs() (*OutputStatement, error) {

	tok, literal := p.scanIgnoreWhitespace()
	if tok != OUTPUTS {
		err := fmt.Errorf("Found '%s', expected 'Outputs'", literal)
		return nil, err
	}

	tok, literal = p.scanIgnoreWhitespace()
	if tok != COLON {
		err := fmt.Errorf("Found '%s', expected ':'", literal)
		return nil, err
	}

	statement := NewOutputStatement()

	for {
		tok, _ := p.scanIgnoreWhitespace()

		if tok == SEMICOLON {
			return statement, nil
		}

		if tok == LETTER {
			p.unscan()
			port, err := p.parsePort()
			if err != nil {
				return nil, err
			}
			statement.add(*port)
		}
	}
}

func (p *Parser) parsePart() (*Part, error) {

	part := Part{}
	tok, literal := p.scanIgnoreWhitespace()

	part.id = literal

	tok, literal = p.scanIgnoreWhitespace()

	if tok == LETTER {
		part.module = literal
	} else {
		err := fmt.Errorf("Found '%s', expected letter", literal)
		return nil, err
	}

	return &part, nil
}

func (p *Parser) ParseParts() (*PartStatement, error) {

	tok, literal := p.scanIgnoreWhitespace()
	if tok != PARTS {
		err := fmt.Errorf("Found '%s', expected 'Outputs'", literal)
		return nil, err
	}

	tok, literal = p.scanIgnoreWhitespace()
	if tok != COLON {
		err := fmt.Errorf("Found '%s', expected ':'", literal)
		return nil, err
	}

	statement := NewPartStatement()

	for {
		tok, _ := p.scanIgnoreWhitespace()

		if tok == SEMICOLON {
			return statement, nil
		}

		if tok == LETTER {
			p.unscan()
			part, err := p.parsePart()
			if err != nil {
				return nil, err
			}
			statement.add(*part)
		}
	}
}
