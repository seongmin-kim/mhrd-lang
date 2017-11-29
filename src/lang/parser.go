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
	port := Port{Pins: newRange(1, 1)}

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
		port.Pins = newRange(1, val)

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
		err := fmt.Errorf("Found '%s', expected 'Parts'", literal)
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

func (p *Parser) parseWire() (*Wire, error) {

	wire := Wire{
		SourceRange: newRange(0, 0),
		TargetRange: newRange(0, 0),
	}

	// left part

	tok, literal := p.scanIgnoreWhitespace()
	if tok == DIGIT {
		wire.Source = literal
	} else if tok == LETTER {
		wire.Source = literal

		if tok, _ := p.scanIgnoreWhitespace(); tok == DOT {
			if tok, literal := p.scanIgnoreWhitespace(); tok == LETTER {
				wire.SourceOut = literal
			} else {
				err := fmt.Errorf("Found '%s', expected letter", literal)
				return nil, err
			}
		} else {
			p.unscan()
		}

		if tok, _ := p.scanIgnoreWhitespace(); tok == LEFTBRACKET {
			if tok, literal := p.scanIgnoreWhitespace(); tok == DIGIT {
				// wire.SourceStartPin, _ = strconv.Atoi(literal)
				val, _ := strconv.Atoi(literal)
				wire.SourceRange = newRange(val, val)
				if tok, _ := p.scanIgnoreWhitespace(); tok == COLON {
					if tok, literal := p.scanIgnoreWhitespace(); tok == DIGIT {
						// wire.SourceEndPin, _ = strconv.Atoi(literal)
						val, _ := strconv.Atoi(literal)
						wire.SourceRange.End = val
						if tok, _ := p.scanIgnoreWhitespace(); tok != RIGHTBRACKET {
							err := fmt.Errorf("Found '%s', expected ']'", literal)
							return nil, err
						}
					} else {
						err := fmt.Errorf("Found '%s', expected digit", literal)
						return nil, err
					}
				} else if tok == RIGHTBRACKET {
					val, _ := strconv.Atoi(literal)
					wire.SourceRange.End = val
				} else {
					err := fmt.Errorf("Found '%s', expected ']'", literal)
					return nil, err
				}
			} else {
				err := fmt.Errorf("Found '%s', expected digit", literal)
				return nil, err
			}
		} else if tok == ARROW {
			p.unscan()
		} else {
			err := fmt.Errorf("Found '%s', expected '[' or '->'", literal)
			return nil, err
		}
	} else {
		err := fmt.Errorf("Found '%s', expected letter or digit", literal)
		return nil, err
	}

	// Arrow
	if tok, literal = p.scanIgnoreWhitespace(); tok != ARROW {
		err := fmt.Errorf("Found '%s', expected ->", literal)
		return nil, err
	}

	// right part

	tok, literal = p.scanIgnoreWhitespace()
	wire.Target = literal

	if tok, _ := p.scanIgnoreWhitespace(); tok == DOT {
		if tok, literal := p.scanIgnoreWhitespace(); tok == LETTER {
			wire.TargetIn = literal
		} else {
			err := fmt.Errorf("Found '%s', expected letter", literal)
			return nil, err
		}
	} else {
		p.unscan()
	}

	if tok, _ := p.scanIgnoreWhitespace(); tok == LEFTBRACKET {
		if tok, literal := p.scanIgnoreWhitespace(); tok == DIGIT {
			// wire.TargetStartPin, _ = strconv.Atoi(literal)
			val, _ := strconv.Atoi(literal)
			wire.TargetRange = newRange(val, val)
			if tok, _ := p.scanIgnoreWhitespace(); tok == COLON {
				if tok, literal := p.scanIgnoreWhitespace(); tok == DIGIT {
					// wire.TargetEndPin, _ = strconv.Atoi(literal)
					val, _ := strconv.Atoi(literal)
					wire.TargetRange.End = val
					if tok, _ := p.scanIgnoreWhitespace(); tok != RIGHTBRACKET {
						err := fmt.Errorf("Found '%s', expected ']'", literal)
						return nil, err
					}
				} else {
					err := fmt.Errorf("Found '%s', expected digit", literal)
					return nil, err
				}
			} else if tok == RIGHTBRACKET {
				// wire.TargetEndPin = wire.TargetStartPin
				val, _ := strconv.Atoi(literal)
				wire.TargetRange.End = val
			} else {
				err := fmt.Errorf("Found '%s', expected ']'", literal)
				return nil, err
			}
		} else {
			err := fmt.Errorf("Found '%s', expected digit", literal)
			return nil, err
		}
	} else if tok == COMMA || tok == SEMICOLON {
		p.unscan()
	} else {
		err := fmt.Errorf("Found '%s', expected ';'", literal)
		return nil, err
	}

	return &wire, nil
}

func (p *Parser) ParseWires() (*WireStatement, error) {

	tok, literal := p.scanIgnoreWhitespace()
	if tok != WIRES {
		err := fmt.Errorf("Found '%s', expected 'Wires'", literal)
		return nil, err
	}

	tok, literal = p.scanIgnoreWhitespace()
	if tok != COLON {
		err := fmt.Errorf("Found '%s', expected ':'", literal)
		return nil, err
	}

	statement := NewWireStatement()

	for {
		tok, _ := p.scanIgnoreWhitespace()

		if tok == SEMICOLON {
			return statement, nil
		}

		if tok == LETTER || tok == DIGIT {
			p.unscan()
			wire, err := p.parseWire()
			if err != nil {
				return nil, err
			}
			statement.add(*wire)
		}
	}
}

func (p *Parser) Parse() (*InputStatement, *OutputStatement, *PartStatement, *WireStatement, error) {
	inputs, err := p.ParseInputs()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	outputs, err := p.ParseOutputs()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	parts, err := p.ParseParts()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	wires, err := p.ParseWires()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return inputs, outputs, parts, wires, nil
}
