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

func NewParser(data string) (*Parser, error) {
	lexer := NewLexer(data)
	p := &Parser{lexer: lexer}
	err := p.nextToken()
	if err != nil {
		return nil, err
	}

	err = p.nextToken()

	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Parser) Parse() (interface{}, error) {
	var result interface{}
	var err error

	switch p.curToken.Type {
	case TknLeftBrace:
		result, err = p.parseObject()
	case TknLeftBracket:
		result, err = p.parseArray()
	default:
		err = fmt.Errorf("expected '{' or '[', found '%s'", p.curToken.Literal)
	}

	if err != nil {
		return nil, err
	}

	if p.curToken.Type != TknEOF {
		return nil, fmt.Errorf("extra data found after JSON value: %s", p.curToken.Literal)
	}

	return result, nil
}

func (p *Parser) parseObject() (map[string]interface{}, error) {
	object := make(map[string]interface{})

	err := p.nextToken()
	if err != nil {
		return nil, err
	}

	for p.curToken.Type != TknRightBrace {
		if p.curToken.Type != TknString {
			return nil, fmt.Errorf("expected string for key, got '%s'", p.curToken.Literal)
		}

		key := p.curToken.Literal
		err := p.nextToken()

		if p.curToken.Type != TknColon {
			return nil, fmt.Errorf("expected colon after key, got '%s'", p.curToken.Literal)
		}

		err = p.nextToken()
		if err != nil {
			return nil, err
		}

		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		object[key] = value

		if p.curToken.Type == TknComma {
			err := p.nextToken()
			if err != nil {
				return nil, err
			}
			if p.curToken.Type == TknRightBrace {
				return nil, fmt.Errorf("unexpected trailing comma before closing brace")
			}
		} else if p.curToken.Type != TknRightBrace {
			return nil, fmt.Errorf("expected comma or closing brace, got '%s'", p.curToken.Literal)
		}
	}

	err = p.nextToken()
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (p *Parser) nextToken() error {
	p.curToken = p.peekToken
	next, err := p.lexer.NextToken()
	if err != nil {
		return err
	}
	p.peekToken = next
	return nil
}

func (p *Parser) parseValue() (interface{}, error) {
	switch p.curToken.Type {
	case TknString:
		val := p.curToken.Literal
		err := p.nextToken()
		if err != nil {
			return nil, err
		}
		return val, nil

	case TknNumber:
		num, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse %q as a number: %v", p.curToken.Literal, err)
		}
		err = p.nextToken()
		if err != nil {
			return nil, err
		}
		return num, nil

	case TknLeftBrace:
		obj, err := p.parseObject()
		if err != nil {
			return nil, err
		}
		return obj, nil

	case TknLeftBracket: // Handle arrays
		arr, err := p.parseArray()
		if err != nil {
			return nil, err
		}
		return arr, nil

	case TknBoolean:
		boolean := p.curToken.Literal == "true"
		err := p.nextToken()
		if err != nil {
			return nil, err
		}
		return boolean, nil

	case TknNull:
		err := p.nextToken()
		if err != nil {
			return nil, err
		}
		return nil, nil

	default:
		return nil, fmt.Errorf("unexpected token type %v when expecting a value", p.curToken.Literal)
	}
}

func (p *Parser) parseArray() ([]interface{}, error) {
	var array []interface{}

	p.nextToken()

	for p.curToken.Type != TknRightBracket {
		if p.curToken.Type == TknRightBracket {
			break
		}

		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		array = append(array, value)
		if p.curToken.Type == TknComma {
			p.nextToken()
			if p.curToken.Type == TknRightBracket {
				return nil, fmt.Errorf("expected value, got '%s'", p.curToken.Literal)
			}
		} else if p.curToken.Type != TknRightBracket {
			return nil, fmt.Errorf("expected comma or closing bracket, got '%s'", p.curToken.Literal)
		}
	}

	p.nextToken()
	return array, nil
}
