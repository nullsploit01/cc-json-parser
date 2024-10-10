package parser

import (
	"fmt"
	"strconv"
)

type Parser struct {
	lexer     *Lexer
	curToken  Token
	peekToken Token
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer}
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) Parse() (interface{}, error) {
	return p.parseValue()
}

func (p *Parser) parseValue() (interface{}, error) {
	switch p.curToken.Type {
	case TknLeftBrace:
		return p.parseObject()
	case TknLeftBracket:
		return p.parseArray()
	case TknString:
		return p.parseString()
	case TknNumber:
		return p.parseNumber()
	case TknBoolean:
		return p.parseBoolean()
	case TknNull:
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected token: %s", p.curToken.Literal)
	}
}

func (p *Parser) parseObject() (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	if p.curToken.Type != TknLeftBrace {
		return nil, fmt.Errorf("expected '{', got %s", p.curToken.Literal)
	}

	p.nextToken()

	for p.curToken.Type != TknRightBrace {
		key, err := p.parseString()
		if err != nil {
			return nil, err
		}

		if p.nextToken(); p.curToken.Type != TknColon {
			return nil, fmt.Errorf("expected ':', got %s", p.curToken.Literal)
		}

		p.nextToken()
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		obj[key] = value

		p.nextToken()

		if p.curToken.Type == TknComma {
			p.nextToken()
		} else if p.curToken.Type != TknRightBrace {
			return nil, fmt.Errorf("expected ',' or '}', got %s", p.curToken.Literal)
		}
	}

	p.nextToken() // Move past the '}' to the next token
	return obj, nil
}

func (p *Parser) parseArray() ([]interface{}, error) {
	array := []interface{}{}

	if p.curToken.Type != TknLeftBracket {
		return nil, fmt.Errorf("expected '[', got %s", p.curToken.Literal)
	}

	p.nextToken()

	for p.curToken.Type != TknRightBracket {
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		array = append(array, value)

		p.nextToken()

		if p.curToken.Type == TknComma {
			p.nextToken()
		} else if p.curToken.Type != TknRightBracket {
			return nil, fmt.Errorf("expected ',' or ']', got %s", p.curToken.Literal)
		}
	}

	p.nextToken() // Move past the ']' to the next token
	return array, nil
}

func (p *Parser) parseString() (string, error) {
	if p.curToken.Type != TknString {
		return "", fmt.Errorf("expected string, got %s", p.curToken.Literal)
	}
	val := p.curToken.Literal
	p.nextToken() // Move to the next token
	return val, nil
}

func (p *Parser) parseNumber() (float64, error) {
	if p.curToken.Type != TknNumber {
		return 0, fmt.Errorf("expected number, got %s", p.curToken.Literal)
	}
	val, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		return 0, err
	}
	p.nextToken() // Move to the next token
	return val, nil
}

func (p *Parser) parseBoolean() (bool, error) {
	if p.curToken.Type != TknBoolean {
		return false, fmt.Errorf("expected boolean, got %s", p.curToken.Literal)
	}
	val := p.curToken.Literal == "true"
	p.nextToken() // Move to the next token
	return val, nil
}

func (p *Parser) parseNull() (interface{}, error) {
	if p.curToken.Type != TknNull {
		return nil, fmt.Errorf("expected null, got %s", p.curToken.Literal)
	}
	p.nextToken() // Move to the next token
	return nil, nil
}

// func (p *Parser) parseString() (string, error) {
// 	var out strings.Builder
// 	if p.curToken.Type != TknString {
// 		return "", fmt.Errorf("expected string, got %s", p.curToken.Literal)
// 	}

// 	// Start after the opening quote
// 	for p.lexer.pos < len(p.lexer.input) {
// 		ch := p.lexer.input[p.lexer.pos]
// 		if ch == '"' {
// 			p.lexer.pos++ // Consume the closing quote
// 			p.nextToken() // Move to the next token
// 			return out.String(), nil
// 		}
// 		if ch == '\\' { // Handle escapes
// 			p.lexer.pos++
// 			if p.lexer.pos >= len(p.lexer.input) {
// 				return "", fmt.Errorf("unexpected end of string")
// 			}
// 			escapeChar := p.lexer.input[p.lexer.pos]
// 			switch escapeChar {
// 			case '"':
// 				out.WriteByte('"')
// 			case '\\':
// 				out.WriteByte('\\')
// 			case '/':
// 				out.WriteByte('/')
// 			case 'b':
// 				out.WriteByte('\b')
// 			case 'f':
// 				out.WriteByte('\f')
// 			case 'n':
// 				out.WriteByte('\n')
// 			case 'r':
// 				out.WriteByte('\r')
// 			case 't':
// 				out.WriteByte('\t')
// 			case 'u':
// 				if p.lexer.pos+4 >= len(p.lexer.input) {
// 					return "", fmt.Errorf("invalid Unicode escape sequence")
// 				}
// 				// Parse the 4 hex digits
// 				hex := p.lexer.input[p.lexer.pos+1 : p.lexer.pos+5]
// 				runeValue, err := strconv.ParseInt(hex, 16, 64)
// 				if err != nil {
// 					return "", fmt.Errorf("invalid Unicode escape sequence: %v", err)
// 				}
// 				if utf16.IsSurrogate(rune(runeValue)) {
// 					// If surrogate, expect a second \u escape
// 					// This part can be expanded based on UTF-16 parsing requirements
// 				}
// 				out.WriteRune(rune(runeValue))
// 				p.lexer.pos += 4 // Move past the hex digits
// 			default:
// 				return "", fmt.Errorf("invalid escape character: \\%c", escapeChar)
// 			}
// 		} else {
// 			out.WriteByte(ch)
// 		}
// 		p.lexer.pos++
// 	}
// 	return "", fmt.Errorf("unterminated string")
// }
