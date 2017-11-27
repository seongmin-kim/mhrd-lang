package lang

import (
	"bufio"
	"bytes"
	"io"
)

const (
	eof = rune(0)
)

func isWhitespace(r rune) bool {
	if r == ' ' || r == '\t' {
		return true
	}
	return false
}

func isNewline(r rune) bool {
	if r == '\n' {
		return true
	}
	return false
}

func isUppercase(r rune) bool {
	if 'A' <= r && r <= 'Z' {
		return true
	}
	return false
}

func isLowercase(r rune) bool {
	if 'a' <= r && r <= 'z' {
		return true
	}
	return false
}

func isLetter(r rune) bool {
	if isUppercase(r) || isLowercase(r) {
		return true
	}
	return false
}

func isDigit(r rune) bool {
	if '0' <= r && r <= '9' {
		return true
	}
	return false
}

func isComma(r rune) bool {
	return r == ','
}

func isDot(r rune) bool {
	return r == '.'
}

func isSemicolon(r rune) bool {
	return r == ';'
}

func isColon(r rune) bool {
	return r == ':'
}

func isLeftBracket(r rune) bool {
	return r == '['
}

func isRightBracket(r rune) bool {
	return r == ']'
}

func isEOF(r rune) bool {
	return r == eof
}

// mhrd token scanner
type Scanner struct {
	r *bufio.Reader
}

func NewScanner(rd io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(rd)}
}

// read a rune and returns
func (s *Scanner) read() rune {
	tok, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return tok
}

// unread a rune and returns error if exists
func (s *Scanner) unread() error {
	return s.r.UnreadRune()
}

func (s *Scanner) scan() (token, string) {
	tok := s.read()

	// character
	if isWhitespace(tok) {
		s.unread()
		return s.scanWhitespace()
	} else if isNewline(tok) {
		s.unread()
		return s.scanNewline()
	} else if isLetter(tok) {
		s.unread()
		return s.scanLetter()
	} else if isDigit(tok) {
		s.unread()
		return s.scanDigit()
	}

	// delemiter
	if isComma(tok) {
		s.unread()
		return s.scanComma()
	} else if isDot(tok) {
		s.unread()
		return s.scanDot()
	} else if isSemicolon(tok) {
		s.unread()
		return s.scanSemicolon()
	} else if isColon(tok) {
		s.unread()
		return s.scanColon()
	}

	// bracket
	if isLeftBracket(tok) {
		s.unread()
		return s.scanLeftBracket()
	} else if isRightBracket(tok) {
		s.unread()
		return s.scanRightBracket()
	}

	// arrow
	if s.read() == '>' && tok == '-' {
		return ARROW, "->"
	} else {
		s.unread()
	}

	// comment
	if s.read() == '/' && tok == '/' {
		return s.scanComment()
	} else {
		s.unread()
	}
	// eof, invalid
	if isEOF(tok) {
		return EOF, string(eof)
	}
	return INVALID, string(tok)
}

func (s *Scanner) scanWhitespace() (token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if tok := s.read(); tok == eof {
			break
		} else if !isWhitespace(tok) {
			s.unread()
			break
		} else {
			buf.WriteRune(tok)
		}
	}
	return WHITESPACE, buf.String()
}

func (s *Scanner) scanNewline() (token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if tok := s.read(); tok == eof {
			break
		} else if !isNewline(tok) {
			s.unread()
			break
		} else {
			buf.WriteRune(tok)
		}
	}
	return NEWLINE, buf.String()
}

func (s *Scanner) scanComment() (token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if tok := s.read(); tok == eof {
			break
		} else if isNewline(tok) {
			s.unread()
			break
		} else {
			buf.WriteRune(tok)
		}
	}
	return COMMENT, buf.String()
}

func (s *Scanner) scanLetter() (token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if tok := s.read(); tok == eof {
			break
		} else if isLetter(tok) {
			buf.WriteRune(tok)
		} else if isDigit(tok) {
			buf.WriteRune(tok)
		} else {
			s.unread()
			break
		}
	}

	if "Inputs" == buf.String() {
		return INPUTS, buf.String()
	} else if "Outputs" == buf.String() {
		return OUTPUTS, buf.String()
	} else if "Parts" == buf.String() {
		return PARTS, buf.String()
	} else if "Wires" == buf.String() {
		return WIRES, buf.String()
	}

	return LETTER, buf.String()
}

func (s *Scanner) scanDigit() (token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if tok := s.read(); tok == eof {
			break
		} else if !isDigit(tok) {
			s.unread()
			break
		} else {
			buf.WriteRune(tok)
		}
	}

	return DIGIT, buf.String()
}

func (s *Scanner) scanComma() (token, string) {
	tok := s.read()
	return COMMA, string(tok)
}

func (s *Scanner) scanDot() (token, string) {
	tok := s.read()
	return DOT, string(tok)
}

func (s *Scanner) scanSemicolon() (token, string) {
	tok := s.read()
	return SEMICOLON, string(tok)
}

func (s *Scanner) scanColon() (token, string) {
	tok := s.read()
	return COLON, string(tok)
}

func (s *Scanner) scanLeftBracket() (token, string) {
	tok := s.read()
	return LEFTBRACKET, string(tok)
}

func (s *Scanner) scanRightBracket() (token, string) {
	tok := s.read()
	return RIGHTBRACKET, string(tok)
}
