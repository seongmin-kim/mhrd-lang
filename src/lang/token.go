package lang

type token int

// mhrd tokens
const (
	EOF = iota
	INVALID
	WHITESPACE // '\t' | ' '
	NEWLINE    // '\n'
	DIGIT      // '0' ... '9'
	LETTER     // 'a' ... 'z' | 'A' ... 'Z'

	INPUTS  // "Inputs"
	OUTPUTs // "Outputs"
	PARTS   // "Parts"
	WIRES   // "Wires"

	DOT          // ','
	COMMA        // '.'
	SEMICOLON    // ';'
	COLON        // ':'
	LEFTBRACKET  // '['
	RIGHTBRACKET // ']'
	ARROW        // "->"

	COMMENT // '//'
)
