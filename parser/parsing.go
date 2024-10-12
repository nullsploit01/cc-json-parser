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

func NewParser(data string) *Parser {
	lexer := NewLexer(data)
	p := &Parser{lexer: lexer}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Parse() error {
	if p.curToken.Type != TknLeftBrace {
		return fmt.Errorf("expected left brace, found %s", p.curToken.Literal)
	}

	_, err := p.parseObject()
	return err
}

func (p *Parser) parseObject() (map[string]interface{}, error) {
	object := make(map[string]interface{})

	p.nextToken()

	for p.curToken.Type != TknRightBrace {
		if p.curToken.Type != TknString {
			return nil, fmt.Errorf("expected string for key, got '%s'", p.curToken.Literal)
		}

		key := p.curToken.Literal
		p.nextToken()

		if p.curToken.Type != TknColon {
			return nil, fmt.Errorf("expected colon after key, got '%s'", p.curToken.Literal)
		}

		p.nextToken()

		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		object[key] = value

		if p.curToken.Type == TknComma {
			p.nextToken()
			if p.curToken.Type == TknRightBrace {
				return nil, fmt.Errorf("unexpected trailing comma before closing brace")
			}
		} else if p.curToken.Type != TknRightBrace {
			return nil, fmt.Errorf("expected comma or closing brace, got '%s'", p.curToken.Literal)
		}
	}

	p.nextToken()
	return object, nil
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) parseValue() (interface{}, error) {
	switch p.curToken.Type {
	case TknString:
		val := p.curToken.Literal
		p.nextToken()
		return val, nil

	case TknNumber:
		num, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse %q as a number: %v", p.curToken.Literal, err)
		}
		p.nextToken()
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
		p.nextToken()
		return boolean, nil

	case TknNull:
		p.nextToken()
		return nil, nil

	default:
		return nil, fmt.Errorf("unexpected token type %v when expecting a value", p.curToken.Literal)
	}
}

func (p *Parser) parseArray() ([]interface{}, error) {
	var array []interface{}

	p.nextToken()

	for p.curToken.Type != TknRightBracket {
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		array = append(array, value)

		if p.curToken.Type == TknComma {
			p.nextToken()
		} else if p.curToken.Type != TknRightBracket {
			return nil, fmt.Errorf("expected comma or closing bracket, got '%s'", p.curToken.Literal)
		}
	}

	p.nextToken()

	return array, nil
}
