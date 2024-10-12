package parser

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
		} else if isDigit(l.ch) {
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
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' {
			break
		}
	}

	str := l.input[position:l.position]
	// l.readChar()
	return str
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
	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()
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
