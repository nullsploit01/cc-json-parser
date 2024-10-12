package parser

import (
	"fmt"
	"strconv"
)

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

func (l *Lexer) NextToken() (Token, error) {
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
		lit, err := l.readString()
		if err != nil {
			return Token{}, err
		}

		tok.Literal = lit
	case 0:
		tok.Literal = ""
		tok.Type = TknEOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = l.lookupIdent(tok.Literal)
			return tok, nil
		} else if l.ch == '-' || isDigit(l.ch) {
			tok.Type = TknNumber
			number, err := l.readNumber()

			if err != nil {
				return Token{}, err
			}
			tok.Literal = number
			return tok, nil
		} else {
			tok = Token{TknInvalid, string(l.ch)}
		}
	}

	l.readChar()
	return tok, nil
}

func (l *Lexer) readString() (string, error) {
	var result []byte
	for {
		l.readChar()
		if l.ch == '"' {
			break
		}
		if l.ch == '\\' {
			l.readChar() // Escape sequence
			switch l.ch {
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
				handledChar, err := l.handleStandardEscape(l.ch)
				if err != nil {
					return "", err
				}
				result = append(result, handledChar)
			case 'u':
				hexValue := l.input[l.position+1 : l.position+5]
				r, err := strconv.ParseInt(hexValue, 16, 32)
				if err == nil {
					result = append(result, byte(r))
				}
				l.position += 4 // Move past the 4 hex digits
				l.readPosition += 4
			default:
				return "", fmt.Errorf("Illegal backslash escape: \\%c\n", l.ch)
			}
		} else if l.ch < ' ' { // Control characters should be escaped
			return "", fmt.Errorf("unescaped control character: %#U", l.ch)
		} else {
			result = append(result, l.ch)
		}
	}
	return string(result), nil
}

func (l *Lexer) handleStandardEscape(ch byte) (byte, error) {
	switch ch {
	case '"':
		return '"', nil
	case '\\':
		return '\\', nil
	case '/':
		return '/', nil
	case 'b':
		return '\b', nil
	case 'f':
		return '\f', nil
	case 'n':
		return '\n', nil
	case 'r':
		return '\r', nil
	case 't':
		return '\t', nil
	default:
		return 0, fmt.Errorf("illegal escape character: '\\%c'", ch)
	}
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

func (l *Lexer) readNumber() (string, error) {
	position := l.position

	if l.ch == '-' {
		l.readChar()
	}

	if l.ch == '0' {
		nextChar := l.peekChar()
		if isDigit(nextChar) && nextChar != '.' {
			return "", fmt.Errorf("invalid number with leading zero at position %d", l.position)
		}
	}

	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}

	if l.ch == 'e' || l.ch == 'E' {
		l.readChar()

		if l.ch == '-' || l.ch == '+' {
			l.readChar()
		}

		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position], nil
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
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
