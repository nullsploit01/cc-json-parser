package parser

import "strconv"

// Define token types
type TokenType int

const (
	TknInvalid TokenType = iota
	TknEOF
	TknWhiteSpace
	TknLeftBrace    // {
	TknRightBrace   // }
	TknLeftBracket  // [
	TknRightBracket // ]
	TknColon        // :
	TknComma        // ,
	TknString
	TknNumber
	TknBoolean
	TknNull
)

type Token struct {
	Type    TokenType
	Literal string
}

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII code for "NUL"
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()
	switch l.ch {
	case '{':
		tok = Token{TknLeftBrace, "{"}
	case '}':
		tok = Token{TknRightBrace, "}"}
	case '[':
		tok = Token{TknLeftBracket, "["}
	case ']':
		tok = Token{TknRightBracket, "]"}
	case ':':
		tok = Token{TknColon, ":"}
	case ',':
		tok = Token{TknComma, ","}
	case '"':
		tok.Type = TknString
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = TknEOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = l.lookupIdent(tok.Literal)
			return tok
		} else if l.ch == '-' || isDigit(l.ch) {
			tok.Type = TknNumber
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = Token{TknInvalid, string(l.ch)}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readString() string {
	var result []byte
	for {
		l.readChar()
		if l.ch == '"' {
			break
		}
		if l.ch == '\\' {
			l.readChar() // Escape sequence
			switch l.ch {
			case '"':
				result = append(result, '"')
			case '\\':
				result = append(result, '\\')
			case 'b':
				result = append(result, '\b')
			case 'f':
				result = append(result, '\f')
			case 'n':
				result = append(result, '\n')
			case 'r':
				result = append(result, '\r')
			case 't':
				result = append(result, '\t')
			case 'u':
				// Handle Unicode sequence
				hexValue := l.input[l.position+1 : l.position+5]
				r, err := strconv.ParseInt(hexValue, 16, 32)
				if err == nil {
					result = append(result, byte(r))
				}
				l.position += 4 // Move past the 4 hex digits
				l.readPosition += 4
			default:
				// Handle invalid escape sequences
				result = append(result, '\\', l.ch)
			}
		} else {
			result = append(result, l.ch)
		}
	}
	return string(result)
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() string {
	position := l.position

	// Check if the number starts with a '-' sign
	if l.ch == '-' {
		l.readChar() // Move past the '-' sign
	}

	// Read the main part of the number (digits and optional decimal point)
	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}

	// Handle scientific notation (e.g., 1.23e+10)
	if l.ch == 'e' || l.ch == 'E' {
		l.readChar() // Move past the 'e' or 'E'

		// Handle the optional '+' or '-' sign in the exponent
		if l.ch == '-' || l.ch == '+' {
			l.readChar()
		}

		// Read the exponent part (must be digits)
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) lookupIdent(ident string) TokenType {
	switch ident {
	case "true", "false":
		return TknBoolean
	case "null":
		return TknNull
	default:
		return TknInvalid
	}
}
